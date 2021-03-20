import os

from aws_cdk import (
    core,
    aws_lambda as lambda_,
    aws_s3 as s3,
    aws_dynamodb as dynamodb,
    aws_apigateway as apigateway,
    aws_certificatemanager as certificates,
)


class ApiStack(core.NestedStack):
    def __init__(
        self,
        scope: core.Construct,
        construct_id: str,
        count_results_table: str,
        **kwargs
    ) -> None:

        super().__init__(scope, construct_id, **kwargs)

        # TODO: This needs to be changed to a dynamic cert request
        # instead of relying on a cert already existing
        self.api = apigateway.RestApi(
            self,
            "RedditTrendsAPI",
            domain_name=apigateway.DomainNameOptions(
                domain_name="api.thekettle.org",
                certificate=certificates.Certificate.from_certificate_arn(
                    self,
                    "DomainCertificateEast2",
                    "arn:aws:acm:us-east-2:261392311630:certificate/8509c657-9ad9-4c9a-80e2-f11d9535b13d",
                ),
                security_policy=apigateway.SecurityPolicy.TLS_1_2,
            ),
            deploy_options=apigateway.StageOptions(
                stage_name="prod",
            ),
        )

        self.api.root.add_method("ANY")

        """
        Reddit data access API
        """
        reddit_resource = self.api.root.add_resource("reddit")

        get_cr_integration = apigateway.LambdaIntegration(
            self.lambda_get_countresults(count_results_table)
        )

        reddit_resource.add_method(
            "GET",
            get_cr_integration,
            api_key_required=True,
            request_validator_options={
                "validate_request_parameters": True,
            },
            request_parameters={
                "method.request.querystring.subreddit": True,
                "method.request.querystring.date": True,
            },
        )

    def lambda_get_countresults(self, count_results_table) -> lambda_.Function:
        handler = lambda_.Function(
            self,
            "GetCountResultsFunction",
            runtime=lambda_.Runtime.GO_1_X,
            code=lambda_.Code.from_asset(
                ".",
                bundling=core.BundlingOptions(
                    user="root",
                    image=lambda_.Runtime.GO_1_X.bundling_docker_image,
                    command=[
                        "bash",
                        "-c",
                        "GOOS=linux go build -o /asset-output/main cmd/lambda/get-count-results/main.go",
                    ],
                ),
            ),
            handler="main",
            environment={
                "TABLE": count_results_table.table_name,
            },
        )

        count_results_table.grant_read_data(handler)

        return handler
