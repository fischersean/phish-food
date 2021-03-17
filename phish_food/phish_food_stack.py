import os

from aws_cdk import (
    core,
    aws_ecs as ecs,
    aws_cognito as cognito,
)

from .api import *
from .etl import *


def USERAUTH(stack: core.Construct) -> cognito.UserPool:
    pool = cognito.UserPool(
        stack,
        "PhishFoodUserPool",
        self_sign_up_enabled=True,
        sign_in_case_sensitive=False,
        sign_in_aliases=cognito.SignInAliases(
            email=True,
            username=False,
        ),
    )

    pool.register_identity_provider(
        provider=cognito.UserPoolIdentityProviderGoogle(
            stack,
            "GoogleIdentityProvider",
            user_pool=pool,
            client_id="74805178401-r3qn337n3avo2sj8fmnjg3r8knhtrjk4.apps.googleusercontent.com",
            client_secret="NNjSYYylIzCG-DVmttkDTQep",
        )
    )

    app_client = cognito.UserPoolClient(
        stack,
        "FrontendClient",
        user_pool=pool,
        prevent_user_existence_errors=True,
    )

    return pool


class PhishFoodBackendStack(core.Stack):
    def __init__(
        self, scope: core.Construct, construct_id: str, **kwargs
    ) -> None:

        super().__init__(scope, construct_id, **kwargs)

        vpc = ec2.Vpc(
            self,
            "PhishFood-VPC",
            nat_gateways=0,  # $1/day is too damn high
        )
        cluster = ecs.Cluster(self, "PhishFood-EcsCluster", vpc=vpc)

        tradeables_bucket = S3_TRADEABLES(self)
        refresh_tradeables_func = LAMBDA_REFRESHTRADEABLES(
            self, tradeables_bucket
        )

        main_table = DYNAMO_SCRAPERESULTS(self)
        rarchive_table = DYNAMO_REDDITARCHIVE(self)

        etl_task = FARGATE_ETL(
            self,
            cluster,
            vpc,
            tradeables_bucket,
            main_table,
            rarchive_table,
        )

        # Front end deployment
        # pool = USERAUTH(self)
        # user_data = DYNAMO_USERDATA(self)
        # API
        api = API_MAIN(self, main_table)
