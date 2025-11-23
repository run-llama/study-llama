import os
from contextlib import asynccontextmanager
from typing import AsyncIterator
from qdrant_client import AsyncQdrantClient
from llama_cloud_services.beta.classifier import LlamaClassify
from llama_cloud_services.extract import LlamaExtract
from sqlalchemy.ext.asyncio import create_async_engine, AsyncConnection
from study_llama.vectordb.vectordb import SummaryVectorDB, FaqsVectorDB

async def get_llama_classify(*args, **kwargs):
    return LlamaClassify.from_api_key(api_key=os.getenv("LLAMA_CLOUD_API_KEY", ""))

async def get_llama_extract(*args, **kwargs):
    return LlamaExtract(
        api_key=os.getenv("LLAMA_CLOUD_API_KEY", ""),
    )

@asynccontextmanager
async def get_db_conn() -> AsyncIterator[AsyncConnection]:
    base_url = os.getenv("POSTGRES_CONNECTION_STRING", "").replace("postgresql://", "postgresql+asyncpg://").split("?")[0]
    eng = create_async_engine(url=base_url)
    async with eng.connect() as db_conn:
        yield db_conn

async def get_vector_db_summaries(*args, **kwargs):
    client = AsyncQdrantClient(
        api_key=os.getenv("QDRANT_API_KEY"),
        https=True,
        port=443,
        host=os.getenv("QDRANT_HOST"),
        check_compatibility=False,
    )
    return SummaryVectorDB(client=client, collection_name="summaries")

async def get_vector_db_faqs(*args, **kwargs):
    client = AsyncQdrantClient(
        api_key=os.getenv("QDRANT_API_KEY"),
        https=True,
        port=443,
        host=os.getenv("QDRANT_HOST"),
        check_compatibility=False,
    )
    return FaqsVectorDB(client=client, collection_name="faqs")