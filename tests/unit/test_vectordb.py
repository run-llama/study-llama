import pytest
import os
from qdrant_client import AsyncQdrantClient
from openai import AsyncOpenAI
from study_llama.vectordb.vectordb import FaqsVectorDB, SummaryVectorDB, Result
from study_llama.vectordb.embeddings import OpenAIEmbedder


@pytest.mark.asyncio
@pytest.mark.skipif(
    condition=(os.getenv("OPENAI_API_KEY") is None),
    reason="OpenAI API key not available",
)
async def test_embed_text() -> None:
    texts = ["hello world", "this is a test"]
    client = AsyncOpenAI(api_key=os.getenv("OPENAI_API_KEY"))
    vecs = await OpenAIEmbedder(client=client).embed(texts=texts)
    assert len(vecs) == 2
    assert len(vecs[0]) == 768
    assert isinstance(vecs[0][0], float)
    assert len(vecs[1]) == 768
    assert isinstance(vecs[1][0], float)


@pytest.mark.asyncio
@pytest.mark.skipif(
    condition=(
        os.getenv("OPENAI_API_KEY") is None
        or os.getenv("QDRANT_HOST") is None
        or os.getenv("QDRANT_API_KEY") is None
    ),
    reason="OpenAI API key, Qdrant host server or Qdrant API key not available",
)
async def test_upload_summary() -> None:
    client = AsyncQdrantClient(
        api_key=os.getenv("QDRANT_API_KEY"),
        https=True,
        port=443,
        host=os.getenv("QDRANT_HOST"),
        check_compatibility=False,
    )
    db = SummaryVectorDB(client=client, collection_name="summaries")
    try:
        await db.upload(
            summary="This is a test summary",
            username="testuser",
            category="test",
            file_name="test_summary.pdf",
        )
        success = True
    except Exception:
        success = False
    assert success


@pytest.mark.asyncio
@pytest.mark.skipif(
    condition=(
        os.getenv("OPENAI_API_KEY") is None
        or os.getenv("QDRANT_HOST") is None
        or os.getenv("QDRANT_API_KEY") is None
    ),
    reason="OpenAI API key, Qdrant host server or Qdrant API key not available",
)
async def test_search_summary() -> None:
    client = AsyncQdrantClient(
        api_key=os.getenv("QDRANT_API_KEY"),
        https=True,
        port=443,
        host=os.getenv("QDRANT_HOST"),
        check_compatibility=False,
    )
    db = SummaryVectorDB(client=client, collection_name="summaries")
    try:
        results = await db.search(
            text="This is a test summary",
            username="testuser",
            category="test",
            file_name="test_summary.pdf",
        )
        success = True
    except Exception:
        results = None
        success = False
    assert success
    assert results is not None
    for result in results:
        assert isinstance(result, Result)
        assert result.category == "test"
        assert result.file_name == "test_summary.pdf"
        if result.text == "This is a test summary":
            assert 0.9 <= result.similarity <= 1.1


@pytest.mark.asyncio
@pytest.mark.skipif(
    condition=(
        os.getenv("OPENAI_API_KEY") is None
        or os.getenv("QDRANT_HOST") is None
        or os.getenv("QDRANT_API_KEY") is None
    ),
    reason="OpenAI API key, Qdrant host server or Qdrant API key not available",
)
async def test_upload_faq() -> None:
    client = AsyncQdrantClient(
        api_key=os.getenv("QDRANT_API_KEY"),
        https=True,
        port=443,
        host=os.getenv("QDRANT_HOST"),
        check_compatibility=False,
    )
    db = FaqsVectorDB(client=client, collection_name="faqs")
    try:
        await db.upload(
            questions=["What is this summary?"],
            answers=["This is a test summary"],
            username="testuser",
            category="test",
            file_name="test_summary.pdf",
        )
        success = True
    except Exception:
        success = False
    assert success


@pytest.mark.asyncio
@pytest.mark.skipif(
    condition=(
        os.getenv("OPENAI_API_KEY") is None
        or os.getenv("QDRANT_HOST") is None
        or os.getenv("QDRANT_API_KEY") is None
    ),
    reason="OpenAI API key, Qdrant host server or Qdrant API key not available",
)
async def test_search_faq() -> None:
    client = AsyncQdrantClient(
        api_key=os.getenv("QDRANT_API_KEY"),
        https=True,
        port=443,
        host=os.getenv("QDRANT_HOST"),
        check_compatibility=False,
    )
    db = FaqsVectorDB(client=client, collection_name="faqs")
    try:
        results = await db.search(
            text="What is this summary?",
            username="testuser",
            category="test",
            file_name="test_summary.pdf",
        )
        success = True
    except Exception:
        results = None
        success = False
    assert success
    assert results is not None
    for result in results:
        assert isinstance(result, Result)
        assert result.category == "test"
        assert result.file_name == "test_summary.pdf"
        if result.text == "This is a test summary":
            assert 0.9 <= result.similarity <= 1.1
