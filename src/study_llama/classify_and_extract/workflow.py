import os
from workflows import Workflow, Context, step
from workflows.resource import Resource
from llama_cloud.types.file import File
from llama_cloud.types.extract_run import ExtractRun
from typing import Annotated, cast
from study_llama.filesdb.query_files import AsyncQuerier as AsyncFilesQuerier
from study_llama.rulesdb.query_rules import AsyncQuerier as AsyncRulesQuerier
from .events import InputFileEvent, IngestedFileEvent, ClassifiedFileEvent, ExtractedFileEvent
from .resources import get_db_conn, get_llama_classify, get_llama_extract, get_vector_db_faqs, get_vector_db_summaries
from .models import StudyNotes, WorkflowState, EXTRACT_CONFIG
from .utils import rules_to_classify_rules
from llama_cloud_services.beta.classifier.client import LlamaClassify
from llama_cloud_services.extract import LlamaExtract
from study_llama.vectordb.vectordb import SummaryVectorDB, FaqsVectorDB

class ClassifyExtractWorkflow(Workflow):
    @step
    async def classify_file(self, ev: InputFileEvent, ctx: Context[WorkflowState], classifier: Annotated[LlamaClassify, Resource(get_llama_classify)]) -> ClassifiedFileEvent | IngestedFileEvent:
        async with get_db_conn() as db_conn:
            querier = AsyncRulesQuerier(conn=db_conn)
            response = querier.get_rules(username=ev.username)
            rules = []
            async for rule in response:
                rules.append(rule)
            class_rules = rules_to_classify_rules(rules=rules)
            result = await classifier.aclassify_file_ids(rules=class_rules, file_ids=[ev.file_id])
            file_type: str | None = None
            for item in result.items:
                if (class_res := item.result) is not None:
                    if class_res.type is not None:
                        file_type = class_res.type
                        break
            if file_type is not None:
                async with ctx.store.edit_state() as state:
                    state.file_name = ev.file_name
                    state.username = ev.username
                    state.file_id = ev.file_id
                    state.file_type = file_type
                querier_files = AsyncFilesQuerier(conn=db_conn)
                await querier_files.create_file(username=ev.username, file_name=ev.file_name, file_category=file_type)
                return ClassifiedFileEvent(file_type=file_type)
            else:
                return IngestedFileEvent(success=False, error="It was not possible to classify the provided file based on the existing categories")
    

    @step
    async def extract_file_details(self, ev: ClassifiedFileEvent, extractor: Annotated[LlamaExtract, Resource(get_llama_extract)], ctx: Context[WorkflowState]) -> ExtractedFileEvent | IngestedFileEvent:
        state = await ctx.store.get_state()
        result = await extractor.aextract(
            data_schema=StudyNotes,
            config=EXTRACT_CONFIG,
            files=File(
                id=state.file_id, 
                name=state.file_name, 
                project_id=os.getenv("LLAMA_CLOUD_PROJECT_ID"), 
                data_source_id=None,
                created_at=None,
                external_file_id=None,
                file_size=None,
                file_type=None,
                last_modified_at=None,
                permission_info=None,
                resource_info=None,
                updated_at=None
            ),
        )
        if (data := cast(ExtractRun, result).data) is not None:
            extraction_result = StudyNotes.model_validate(data)
            return ExtractedFileEvent(summary=extraction_result.summary, faqs=extraction_result.faqs)
        else:
            return IngestedFileEvent(success=False, error="It was not possible to extract details from the provided file")
    
    @step
    async def ingest_file_details(self, ev: ExtractedFileEvent, summaries_vdb: Annotated[SummaryVectorDB, Resource(get_vector_db_summaries)], faqs_vdb: Annotated[FaqsVectorDB, Resource(get_vector_db_faqs)], ctx: Context[WorkflowState]) -> IngestedFileEvent:
        state = await ctx.store.get_state()
        questions = [faq.question for faq in ev.faqs]
        answers = [faq.answer for faq in ev.faqs]
        await faqs_vdb.upload(questions, answers, state.username, state.file_type, state.file_name)
        await summaries_vdb.upload(ev.summary, state.username, state.file_type, state.file_name)
        return IngestedFileEvent(success=True)


workflow = ClassifyExtractWorkflow(timeout=1000)