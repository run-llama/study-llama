import re
from llama_cloud.types.classifier_rule import ClassifierRule
from study_llama.rulesdb.models import Rule


def rules_to_classify_rules(rules: list[Rule]) -> list[ClassifierRule]:
    class_rules: list[ClassifierRule] = []
    for rule in rules:
        class_rules.append(
            ClassifierRule(
                type=re.sub(r"\s+", "_", rule.rule_type.lower().strip()),
                description=rule.rule_description,
            )
        )
    return class_rules
