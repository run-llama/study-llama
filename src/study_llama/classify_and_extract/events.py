from workflows.events import StartEvent, StopEvent, Event
from pydantic import ConfigDict
from .models import QuestionAndAnswer


class InputFileEvent(StartEvent):
    file_id: str
    file_name: str
    username: str


class ClassifiedFileEvent(Event):
    file_type: str


class ExtractedFileEvent(Event):
    summary: str
    faqs: list[QuestionAndAnswer]

    model_config = ConfigDict(arbitrary_types_allowed=True)


class IngestedFileEvent(StopEvent):
    success: bool
    error: str | None = None
