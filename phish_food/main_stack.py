from aws_cdk import (
    core,
)

from phish_food.backend import BackendStack
from phish_food.frontend import FrontendStack


class PhishFood(core.Stack):
    def __init__(
        self, scope: core.Construct, construct_id: str, **kwargs
    ) -> None:

        super().__init__(scope, construct_id, **kwargs)

        frontend = FrontendStack(self, "Frontend")
        backend = BackendStack(self, "Backend")
