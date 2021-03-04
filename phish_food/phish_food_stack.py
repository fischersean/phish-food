import os

from aws_cdk import (
    core,
    aws_ec2 as ec2,
    aws_ecs as ecs,
    aws_applicationautoscaling as events,
    aws_lambda as lambda_,
    aws_s3 as s3,
    aws_dynamodb as dynamodb,
    aws_ecs_patterns as ecs_patterns,
    aws_events as aws_events,
    aws_events_targets as targets,
    aws_apigateway as apigateway,
    aws_certificatemanager as certificates,
    aws_s3_deployment as s3deploy,
    aws_cloudfront as cloudfront,
    aws_cloudfront_origins as origins,
)


def S3_FRONTENDDEPLOY(stack: core.Construct) -> s3.Bucket:
    bucket = s3.Bucket(
        stack,
        "FrontendBucket",
        website_index_document="index.html",
        # public_read_access=True,
    )

    distribution = cloudfront.Distribution(
        stack,
        "FrontendDistribution",
        # TODO: The domain and cert info should be env vars
        domain_names=["www.thekettle.org"],
        certificate=certificates.Certificate.from_certificate_arn(
            stack,
            "DomainCertificateEast1",
            "arn:aws:acm:us-east-1:261392311630:certificate/02a75969-25ce-47d3-acf6-d93408b2eed1",
        ),
        default_behavior=cloudfront.BehaviorOptions(
            origin=origins.S3Origin(
                bucket,
                origin_access_identity=cloudfront.OriginAccessIdentity(
                    stack,
                    "FrontendDeplotmentIdentity",
                ),
            ),
            viewer_protocol_policy=cloudfront.ViewerProtocolPolicy.REDIRECT_TO_HTTPS,
        ),
    )

    s3deploy.BucketDeployment(
        stack,
        "FrontendS3Deployment",
        sources=[s3deploy.Source.asset("./web/dist")],
        destination_bucket=bucket,
        distribution=distribution,
    )

    return bucket


def S3_TRADEABLES(stack: core.Construct) -> s3.Bucket:
    bucket = s3.Bucket(
        stack,
        "TradeableSecurities",
        removal_policy=core.RemovalPolicy.DESTROY,
    )

    return bucket


def DYNAMO_SCRAPERESULTS(stack: core.Construct) -> dynamodb.Table:
    # parition is sub+YYYY+MM+DD
    table = dynamodb.Table(
        stack,
        "RedditTrendingStocks",
        partition_key=dynamodb.Attribute(
            name="id", type=dynamodb.AttributeType.STRING
        ),
        sort_key=dynamodb.Attribute(
            name="hour", type=dynamodb.AttributeType.NUMBER
        ),
    )

    return table


def DYNAMO_REDDITARCHIVE(stack: core.Construct) -> dynamodb.Table:
    # parition is sub+YYYY+MM+DD
    table = dynamodb.Table(
        stack,
        "RedditPermalinkArchive",
        partition_key=dynamodb.Attribute(
            name="id", type=dynamodb.AttributeType.STRING
        ),
        sort_key=dynamodb.Attribute(
            name="hour", type=dynamodb.AttributeType.NUMBER
        ),
    )

    return table


def FARGATE_ETL(
    stack: core.Construct,
    cluster: ecs.Cluster,
    vpc: ec2.Vpc,
    bucket: s3.Bucket,
    results_table: dynamodb.Table,
    archive_table: dynamodb.Table,
) -> ecs_patterns.ScheduledFargateTask:

    app_secret = os.environ["APP_SECRET"]
    app_id = os.environ["APP_ID"]

    task = ecs_patterns.ScheduledFargateTask(
        stack,
        "ETLTask",
        cluster=cluster,
        vpc=vpc,
        scheduled_fargate_task_image_options=ecs_patterns.ScheduledFargateTaskImageOptions(
            image=ecs.ContainerImage.from_asset(
                ".",
                file="Dockerfile.etl",
            ),
            environment={
                "name": "TRIGGER",
                "value": "CloudWatch Events",
                "BUCKET": bucket.bucket_name,
                "APP_ID": app_id,
                "APP_SECRET": app_secret,
                "TABLE": results_table.table_name,
                "ARCHIVE_TABLE": archive_table.table_name,
            },
            cpu=2048,
            memory_limit_mib=4096,
        ),
        enabled=True,
        schedule=events.Schedule.expression(
            "cron(0 * * * ? *)",  # Run at beginning of every hour
        ),
        subnet_selection=ec2.SubnetSelection(
            subnet_type=ec2.SubnetType.PUBLIC
        ),
    )

    bucket.grant_read(task.task_definition.task_role)

    results_table.grant_read_write_data(task.task_definition.task_role)
    archive_table.grant_read_write_data(task.task_definition.task_role)

    return task


def API_MAIN(
    stack: core.Construct, get_countresults_func: lambda_.Function
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
    reddit = api.root.add_resource("reddit")

    get_cr_integration = apigateway.LambdaIntegration(get_countresults_func)
    reddit.add_method(
        "GET",
        get_cr_integration,
        request_validator_options={
            "validate_request_parameters": True,
        },
        request_parameters={
            "method.request.querystring.subreddit": True,
            "method.request.querystring.date": True,
        },
    )

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


def LAMBDA_REFRESHTRADEABLES(
    stack: core.Construct, bucket: s3.Bucket
) -> lambda_.Function:

    handler = lambda_.Function(
        stack,
        "RefreshTrabeablesFunction",
        runtime=lambda_.Runtime.GO_1_X,
        code=lambda_.Code.from_asset(
            ".",
            bundling=core.BundlingOptions(
                user="root",
                image=lambda_.Runtime.GO_1_X.bundling_docker_image,
                command=[
                    "bash",
                    "-c",
                    "GOOS=linux go build -o /asset-output/main cmd/lambda/refresh-tradeables/main.go",
                ],
            ),
        ),
        handler="main",
        environment=dict(BUCKET=bucket.bucket_name),
    )

    bucket.grant_read_write(handler)

    rule = aws_events.Rule(
        stack,
        "RefreshTrabeablesSchedule",
        schedule=aws_events.Schedule.cron(
            minute="0", hour="0", day="*", month="*", year="*"
        ),
    )

    rule.add_target(targets.LambdaFunction(handler))

    return handler


class PhishFoodStack(core.Stack):
    def __init__(
        self, scope: core.Construct, construct_id: str, **kwargs
    ) -> None:

        super().__init__(scope, construct_id, stack_name="PhishFood", **kwargs)

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
            self, cluster, vpc, tradeables_bucket, main_table, rarchive_table
        )

        # API
        get_count_results_func = LAMBDA_GET_COUNTRESULTS(self, main_table)
        api = API_MAIN(self, get_count_results_func)

        # Front end deployment
        S3_FRONTENDDEPLOY(self)
