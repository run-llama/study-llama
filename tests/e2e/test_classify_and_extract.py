import pytest
import os

from workflows.testing import WorkflowTestRunner
from study_llama.classify_and_extract.workflow import workflow
from study_llama.classify_and_extract.events import InputFileEvent, IngestedFileEvent


condition = (
    os.getenv("LLAMA_CLOUD_API_KEY") is None
    or os.getenv("LLAMA_CLOUD_PROJECT_ID") is None
    or os.getenv("QDRANT_API_KEY") is None
    or os.getenv("OPENAI_API_KEY") is None
    or os.getenv("POSTGRES_CONNECTION_STRING") is None
    or os.getenv("QDRANT_HOST") is None
    or os.getenv("UPLOADED_FILE_ID") is None
    or os.getenv("TEST_USER") is None
)


@pytest.mark.skipif(
    condition=condition, reason="Needed environment variables are not available"
)
@pytest.mark.asyncio
async def test_workflow() -> None:
    test_runner = WorkflowTestRunner(workflow=workflow)
    try:
        result = await test_runner.run(
            start_event=InputFileEvent(
                file_id=os.getenv("UPLOADED_FILE_ID", ""),
                file_name="observability.pdf",
                username=os.getenv("TEST_USER", ""),
            ),
            expose_internal=False,
        )
    except Exception as e:
        result = None
    assert result is not None
    assert isinstance(result.result, IngestedFileEvent)
    assert result.result.error is None
    assert result.result.success
