from workflows.events import StartEvent, StopEvent
from typing import Literal
from pydantic import ConfigDict
from study_llama.vectordb.vectordb import Result


class SearchInputEvent(StartEvent):
    search_type: Literal["summary", "faqs"]
    search_input: str
    username: str
    file_name: str | None = None
    category: str | None = None


class SearchOutputEvent(StopEvent):
    results: list[Result]

    model_config = ConfigDict(arbitrary_types_allowed=True)
