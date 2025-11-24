import os
from qdrant_client import AsyncQdrantClient
from study_llama.vectordb.vectordb import SummaryVectorDB, FaqsVectorDB


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
