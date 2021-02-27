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
)


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
        rule_name="RunETLTask",
        subnet_selection=ec2.SubnetSelection(
            subnet_type=ec2.SubnetType.PUBLIC
        ),
    )

    bucket.grant_read(task.task_definition.task_role)

    results_table.grant_read_write_data(task.task_definition.task_role)
    archive_table.grant_read_write_data(task.task_definition.task_role)

    return task


def LAMBDA_REFRESH_TRADEABLES(
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
            nat_gateways=0, # $1/day is too damn high
        )
        cluster = ecs.Cluster(self, "PhishFood-EcsCluster", vpc=vpc)

        tradeables_bucket = S3_TRADEABLES(self)
        refresh_tradeables_func = LAMBDA_REFRESH_TRADEABLES(
            self, tradeables_bucket
        )

        main_table = DYNAMO_SCRAPERESULTS(self)
        rarchive_table = DYNAMO_REDDITARCHIVE(self)

        elt_task = FARGATE_ETL(
            self, cluster, vpc, tradeables_bucket, main_table, rarchive_table
        )
