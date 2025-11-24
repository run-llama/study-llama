import uuid
import os
from pydantic import BaseModel
from qdrant_client import AsyncQdrantClient
from qdrant_client.models import PointStruct, Filter, FieldCondition, MatchValue
from typing import cast, Literal
from openai import AsyncOpenAI
from .embeddings import OpenAIEmbedder


class Result(BaseModel):
    result_type: Literal["answer", "summary"]
    text: str
    file_name: str
    category: str
    similarity: float


class SummaryVectorDB:
    def __init__(self, client: AsyncQdrantClient, collection_name: str):
        self._client = client
        self.collection_name = collection_name
        openai_client = AsyncOpenAI(api_key=os.getenv("OPENAI_API_KEY"))
        self._embedder = OpenAIEmbedder(client=openai_client)

    async def upload(
        self, summary: str, username: str, category: str, file_name: str
    ) -> None:
        vec = await self._embedder.embed([summary])
        point = PointStruct(
            id=uuid.uuid4(),
            vector=vec[0],
            payload={
                "category": category,
                "file_name": file_name,
                "summary": summary,
                "username": username,
            },
        )
        self._client.upload_points(self.collection_name, points=[point])
        return None

    async def search(
        self,
        text: str,
        username: str,
        category: str | None = None,
        file_name: str | None = None,
    ) -> list[Result]:
        filters = Filter(
            must=[FieldCondition(key="username", match=MatchValue(value=username))]
        )
        if category is not None:
            (cast(list[FieldCondition], filters.must)).append(
                FieldCondition(key="category", match=MatchValue(value=category))
            )
        if file_name is not None:
            (cast(list[FieldCondition], filters.must)).append(
                FieldCondition(key="file_name", match=MatchValue(value=file_name))
            )
        vec = await self._embedder.embed([text])
        results = await self._client.query_points(
            self.collection_name,
            query=vec[0],
            query_filter=filters,
            score_threshold=0.75,
        )
        points = results.points
        return [
            Result(
                text=point.payload["summary"],
                similarity=point.score,
                file_name=point.payload["file_name"],
                category=point.payload["category"],
                result_type="summary",
            )
            for point in points
            if point.payload is not None
        ]


class FaqsVectorDB:
    def __init__(self, client: AsyncQdrantClient, collection_name: str):
        self._client = client
        self.collection_name = collection_name
        openai_client = AsyncOpenAI(api_key=os.getenv("OPENAI_API_KEY"))
        self._embedder = OpenAIEmbedder(client=openai_client)

    async def upload(
        self,
        questions: list[str],
        answers: list[str],
        username: str,
        category: str,
        file_name: str,
    ) -> None:
        vecs = await self._embedder.embed(questions)
        points = []
        for i, vec in enumerate(vecs):
            points.append(
                PointStruct(
                    id=uuid.uuid4(),
                    vector=vec,
                    payload={
                        "category": category,
                        "file_name": file_name,
                        "question": questions[i],
                        "answer": answers[i],
                        "username": username,
                    },
                )
            )
        self._client.upload_points(self.collection_name, points=points)
        return None

    async def search(
        self,
        text: str,
        username: str,
        category: str | None = None,
        file_name: str | None = None,
    ) -> list[Result]:
        filters = Filter(
            must=[FieldCondition(key="username", match=MatchValue(value=username))]
        )
        if category is not None:
            (cast(list[FieldCondition], filters.must)).append(
                FieldCondition(key="category", match=MatchValue(value=category))
            )
        if file_name is not None:
            (cast(list[FieldCondition], filters.must)).append(
                FieldCondition(key="file_name", match=MatchValue(value=file_name))
            )
        vec = await self._embedder.embed([text])
        results = await self._client.query_points(
            self.collection_name,
            query=vec[0],
            query_filter=filters,
            score_threshold=0.75,
        )
        points = results.points
        return [
            Result(
                text=point.payload["answer"],
                similarity=point.score,
                file_name=point.payload["file_name"],
                category=point.payload["category"],
                result_type="answer",
            )
            for point in points
            if point.payload is not None
        ]
