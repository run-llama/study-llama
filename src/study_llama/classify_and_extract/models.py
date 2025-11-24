from pydantic import BaseModel, Field, ConfigDict
from llama_cloud_services.extract import ExtractConfig
from llama_cloud.types.extract_mode import ExtractMode


class QuestionAndAnswer(BaseModel):
    question: str = Field(description="Question related to the main document")
    answer: str = Field(description="Answer to the question")


class StudyNotes(BaseModel):
    summary: str = Field(description="Summary of the study notes in the document")
    faqs: list[QuestionAndAnswer] = Field(
        description="List of potential 'Frequently Asked Questions' (with associated answer) to help a student review and prepare with the study notes"
    )

    model_config = ConfigDict(arbitrary_types_allowed=True)


class WorkflowState(BaseModel):
    username: str = ""
    file_name: str = ""
    file_type: str = ""
    file_id: str = ""


EXTRACT_CONFIG = ExtractConfig(
    extraction_mode=ExtractMode.MULTIMODAL,
)
