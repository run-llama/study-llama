import os
from qdrant_client import AsyncQdrantClient
from qdrant_client.models import VectorParams, Distance


async def create_collections():
    client = AsyncQdrantClient(
        api_key=os.getenv("QDRANT_API_KEY"),
        https=True,
        port=443,
        host=os.getenv("QDRANT_HOST"),
        check_compatibility=False,
    )
    for coll in ("summaries", "faqs"):
        if not (await client.collection_exists(coll)):
            succ = await client.create_collection(
                collection_name=coll,
                vectors_config=VectorParams(size=768, distance=Distance.COSINE),
            )
            print(
                f"Successfully created {coll}"
                if succ
                else f"Something went wrong while creating {coll}"
            )


if __name__ == "__main__":
    import asyncio

    asyncio.run(create_collections())
