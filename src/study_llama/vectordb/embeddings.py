import os
from openai import AsyncOpenAI


class OpenAIEmbedder:
    def __init__(self, client: AsyncOpenAI) -> None:
        self._client = client

    async def embed(self, texts: list[str]) -> list[list[float]]:
        response = await self._client.embeddings.create(
            input=texts, model="text-embedding-3-small", dimensions=768
        )
        return [d.embedding for d in response.data]
