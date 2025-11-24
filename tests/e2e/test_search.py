import pytest
import os

from workflows.testing import WorkflowTestRunner
from study_llama.search.workflow import workflow
from study_llama.search.events import SearchInputEvent, SearchOutputEvent
from study_llama.vectordb.vectordb import Result

condition = (
    os.getenv("QDRANT_API_KEY") is None
    or os.getenv("OPENAI_API_KEY") is None
    or os.getenv("QDRANT_HOST") is None
)


@pytest.mark.skipif(
    condition=condition, reason="Needed environment variables are not available"
)
@pytest.mark.asyncio
async def test_workflow() -> None:
    test_runner = WorkflowTestRunner(workflow=workflow)
    try:
        result = await test_runner.run(
            start_event=SearchInputEvent(
                search_type="summary",
                search_input="This is a test summary",
                username="testuser",
                file_name="test_summary.pdf",
                category="test",
            )
        )
    except Exception as e:
        result = None
    assert result is not None
    assert isinstance(result.result, SearchOutputEvent)
    assert len(result.result.results) > 0
    assert all(isinstance(res, Result) for res in result.result.results)
    try:
        result = await test_runner.run(
            start_event=SearchInputEvent(
                search_type="faqs",
                search_input="What is this summary?",
                username="testuser",
                file_name="test_summary.pdf",
                category="test",
            )
        )
    except Exception as e:
        result = None
    assert result is not None
    assert isinstance(result.result, SearchOutputEvent)
    assert len(result.result.results) > 0
    assert all(isinstance(res, Result) for res in result.result.results)
