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


class EtlStack(core.NestedStack):
    def __init__(
        self,
        scope: core.Construct,
        construct_id: str,
        vpc: ec2.Vpc,
        cluster: ecs.Cluster,
        count_results_table: dynamodb.Table,
        reddit_archive_bucket: s3.Bucket,
        **kwargs
    ) -> None:

        super().__init__(scope, construct_id, **kwargs)
        """
        Should alread have all tables created

        Create tradeables bucket, function
        Create Fargate task
        """
        tradeables_bucket, tradeables_update_func = self.tradeables()
        distdb_bucket, distdb_update_func = self.distdb(
            tradeables_bucket,
            count_results_table,
            reddit_archive_bucket,
        )

        etl_task = self.fargate_etl(
            vpc,
            cluster,
            tradeables_bucket,
            count_results_table,
            reddit_archive_bucket,
        )

    def tradeables(self) -> (s3.Bucket, lambda_.Function):

        bucket = s3.Bucket(
            self,
            "TradeableSecurities",
            removal_policy=core.RemovalPolicy.DESTROY,
        )

        handler = lambda_.Function(
            self,
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
            self,
            "RefreshTrabeablesSchedule",
            schedule=aws_events.Schedule.cron(
                minute="0", hour="0", day="*", month="*", year="*"
            ),
        )

        rule.add_target(targets.LambdaFunction(handler))

        return bucket, handler

    def distdb(
        self,
        tradeables_bucket: s3.Bucket,
        count_results_table: dynamodb.Table,
        reddit_archive_bucket: s3.Bucket,
    ) -> (s3.Bucket, lambda_.Function):

        bucket = s3.Bucket(
            self,
            "DistributedDatabase",
            removal_policy=core.RemovalPolicy.DESTROY,
            versioned=True,
            public_read_access=True,
        )

        handler = lambda_.Function(
            self,
            "UpdateDistDbFunction",
            runtime=lambda_.Runtime.GO_1_X,
            code=lambda_.Code.from_asset(
                ".",
                bundling=core.BundlingOptions(
                    user="root",
                    image=lambda_.Runtime.GO_1_X.bundling_docker_image,
                    command=[
                        "bash",
                        "-c",
                        "GOOS=linux go build -o /asset-output/main cmd/lambda/update-distdb/main.go",
                    ],
                ),
            ),
            handler="main",
            environment={
                "DIST_BUCKET": bucket.bucket_name,
                "TRADEABLES_BUCKET": tradeables_bucket.bucket_name,
                "ETL_RESULTS_TABLE": count_results_table.table_name,
                "REDDIT_ARCHIVE_BUCKET": reddit_archive_bucket.bucket_name,
            },
            timeout=core.Duration.minutes(5),
        )

        bucket.grant_read_write(handler)
        tradeables_bucket.grant_read_write(handler)
        count_results_table.grant_read_data(handler)

        # Update db ever day at 11:30
        # This should ensure that all of that day's ETL results are available
        rule = aws_events.Rule(
            self,
            "UpdateDistDbSchedule",
            schedule=aws_events.Schedule.cron(
                minute="30", hour="23", day="*", month="*", year="*"
            ),
        )

        rule.add_target(targets.LambdaFunction(handler))

        return bucket, handler

    def fargate_etl(
        self,
        vpc: ec2.Vpc,
        cluster: ecs.Cluster,
        tradeables_bucket: s3.Bucket,
        count_results_table: dynamodb.Table,
        reddit_archive_bucket: s3.Bucket,
    ) -> ecs.TaskDefinition:

        app_secret = os.environ["APP_SECRET"]
        app_id = os.environ["APP_ID"]
        if app_secret == "" or app_id == "":
            raise ValueError("Could not find reddit app secrets")

        task = ecs_patterns.ScheduledFargateTask(
            self,
            "ETLTask",
            cluster=cluster,
            scheduled_fargate_task_image_options=ecs_patterns.ScheduledFargateTaskImageOptions(
                image=ecs.ContainerImage.from_asset(
                    ".",
                    file="Dockerfile.etl",
                ),
                environment={
                    "APP_ID": app_id,
                    "APP_SECRET": app_secret,
                    "TRADEABLES_BUCKET": tradeables_bucket.bucket_name,
                    "ETL_RESULTS_TABLE": count_results_table.table_name,
                    "REDDIT_ARCHIVE_BUCKET": reddit_archive_bucket.bucket_name,
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

        tradeables_bucket.grant_read(task.task_definition.task_role)

        count_results_table.grant_read_write_data(
            task.task_definition.task_role
        )

        reddit_archive_bucket.grant_read_write(task.task_definition.task_role)

        return task
