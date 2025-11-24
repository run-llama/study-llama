import pytest
import os
from llama_cloud.types.classifier_rule import ClassifierRule
from study_llama.rulesdb.query_rules import AsyncQuerier
from study_llama.rulesdb.models import Rule
from study_llama.classify_and_extract.resources import get_db_conn
from study_llama.classify_and_extract.utils import rules_to_classify_rules


@pytest.mark.skipif(
    condition=(
        os.getenv("POSTGRES_CONNECTION_STRING") is None
        or os.getenv("TEST_USER") is None
    ),
    reason="Postgres connection string or test user name are not available",
)
@pytest.mark.asyncio
async def test_file_upload() -> None:
    async with get_db_conn() as db_conn:
        querier = AsyncQuerier(conn=db_conn)
        rules: list[Rule] = []
        response = querier.get_rules(username=os.getenv("TEST_USER", ""))
        async for rule in response:
            assert isinstance(rule, Rule)
            assert rule.username == os.getenv("TEST_USER", "")
            rules.append(rule)
        assert len(rules) > 0
        class_rules = rules_to_classify_rules(rules)
        assert len(class_rules) == len(rules)
        for i, rule in enumerate(class_rules):
            assert isinstance(rule, ClassifierRule)
            assert rule.description == rules[i].rule_description
            assert rule.type == rules[i].rule_type
