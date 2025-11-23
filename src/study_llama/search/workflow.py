from workflows import Workflow, Context, step
from workflows.resource import Resource
from typing import Annotated, TYPE_CHECKING
from .resources import get_vector_db_faqs, get_vector_db_summaries
from .events import SearchInputEvent, SearchOutputEvent
if TYPE_CHECKING:
    from study_llama.vectordb.vectordb import SummaryVectorDB, FaqsVectorDB

class SearchWorkflow(Workflow):
    @step
    async def search(self, ev: SearchInputEvent, summaries_vdb: Annotated[SummaryVectorDB, Resource(get_vector_db_summaries)], faqs_vdb: Annotated[FaqsVectorDB, Resource(get_vector_db_faqs)]) -> SearchOutputEvent:
        if ev.search_type == "faqs":
            results = await faqs_vdb.search(ev.search_input, ev.username, ev.category, ev.file_name)
            return SearchOutputEvent(results=results)
        else:
            results = await summaries_vdb.search(ev.search_input, ev.username, ev.category, ev.file_name)
            return SearchOutputEvent(results=results)

workflow = SearchWorkflow(timeout=600)
