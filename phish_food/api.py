import os

from aws_cdk import (
    core,
    aws_lambda as lambda_,
    aws_s3 as s3,
    aws_dynamodb as dynamodb,
    aws_events as aws_events,
    aws_apigateway as apigateway,
    aws_certificatemanager as certificates,
)


def API_MAIN(
    stack: core.Construct, count_results_table: dynamodb.Table
) -> apigateway.RestApi:

    # TODO: This needs to be changed to a dynamic cert request
    # instead of relying on a cert already existing
    api = apigateway.RestApi(
        stack,
        "RedditTrendsAPI",
        domain_name=apigateway.DomainNameOptions(
            domain_name="api.thekettle.org",
            certificate=certificates.Certificate.from_certificate_arn(
                stack,
                "DomainCertificateEast2",
                "arn:aws:acm:us-east-2:261392311630:certificate/8509c657-9ad9-4c9a-80e2-f11d9535b13d",
            ),
            security_policy=apigateway.SecurityPolicy.TLS_1_2,
        ),
        deploy_options=apigateway.StageOptions(
            stage_name="v0",
        ),
    )

    api.root.add_method("ANY")

    """
    Reddit data access API
    """
    reddit = api.root.add_resource("reddit")

    get_cr_integration = apigateway.LambdaIntegration(
        LAMBDA_GET_COUNTRESULTS(stack, count_results_table)
    )
    reddit.add_method(
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

    """
    Primary Application API
    """

    return api


def LAMBDA_GET_COUNTRESULTS(
    stack: core.Construct, table: dynamodb.Table
) -> lambda_.Function:

    handler = lambda_.Function(
        stack,
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
            "TABLE": table.table_name,
        },
    )

    table.grant_read_data(handler)

    return handler


def DYNAMO_USERDATA(stack: core.Construct) -> dynamodb.Table:

    table = dynamodb.Table(
        stack,
        "UserDataTable",
        partition_key=dynamodb.Attribute(
            name="id", type=dynamodb.AttributeType.STRING
        ),
    )

    return table
