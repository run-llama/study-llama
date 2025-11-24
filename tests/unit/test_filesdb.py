import pytest
import os
from study_llama.filesdb.query_files import AsyncQuerier
from study_llama.filesdb.models import File
from study_llama.classify_and_extract.resources import get_db_conn


@pytest.mark.skipif(
    condition=(os.getenv("POSTGRES_CONNECTION_STRING") is None),
    reason="Postgres connection string is not available",
)
@pytest.mark.asyncio
async def test_file_upload() -> None:
    async with get_db_conn() as db_conn:
        querier = AsyncQuerier(conn=db_conn)
        fl = await querier.create_file(
            username="testuser", file_name="testfile.pdf", file_category="test"
        )
        assert fl is not None
        assert isinstance(fl, File)
        assert fl.file_category == "test"
        assert fl.username == "testuser"
        assert fl.file_name == "testfile.pdf"
