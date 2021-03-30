import os

from aws_cdk import (
    core,
    aws_ecs as ecs,
    aws_ec2 as ec2,
    aws_s3 as s3,
    aws_dynamodb as dynamodb,
    aws_route53 as route53,
)

from phish_food.etl import EtlStack
from phish_food.api import ApiStack


class BackendStack(core.NestedStack):
    def __init__(
        self,
        scope: core.Construct,
        construct_id: str,
        vpc: ec2.Vpc,
        hosted_zone: route53.HostedZone,
        **kwargs
    ) -> None:
        super().__init__(scope, construct_id, **kwargs)

        cluster = ecs.Cluster(self, "PhishFood-EcsCluster", vpc=vpc)

        count_results_table = self.dynamo_scraperesults()
        # reddit_archive_table = self.dynamo_redditarchive()
        api_key_table = self.dynamo_apikeys()
        reddit_archive_bucket = self.s3_redditarchive()

        etl = EtlStack(
            self,
            "ETL",
            vpc=vpc,
            cluster=cluster,
            count_results_table=count_results_table,
            reddit_archive_bucket=reddit_archive_bucket,
            api_key_table=api_key_table,
        )

        api = ApiStack(
            self,
            "API",
            vpc=vpc,
            cluster=cluster,
            count_results_table=count_results_table,
            reddit_archive_bucket=reddit_archive_bucket,
            api_key_table=api_key_table,
            hosted_zone=hosted_zone,
        )

    def dynamo_scraperesults(self: core.Construct) -> dynamodb.Table:
        # parition is sub+YYYY+MM+DD
        table = dynamodb.Table(
            self,
            "RedditTrendingStocks",
            partition_key=dynamodb.Attribute(
                name="id", type=dynamodb.AttributeType.STRING
            ),
            sort_key=dynamodb.Attribute(
                name="hour", type=dynamodb.AttributeType.NUMBER
            ),
        )

        table.auto_scale_write_capacity(
            min_capacity=5,
            max_capacity=50,
        )

        table.auto_scale_read_capacity(
            min_capacity=5,
            max_capacity=50,
        )

        return table

    def s3_redditarchive(self) -> s3.Bucket:
        bucket = s3.Bucket(
            self,
            "RedditPostsArchive",
        )

        return bucket

    def dynamo_redditarchive(self: core.Construct) -> dynamodb.Table:
        # parition is sub+YYYY+MM+DD
        table = dynamodb.Table(
            self,
            "RedditPermalinkArchive",
            partition_key=dynamodb.Attribute(
                name="id", type=dynamodb.AttributeType.STRING
            ),
            sort_key=dynamodb.Attribute(
                name="hour", type=dynamodb.AttributeType.NUMBER
            ),
        )

        return table

    def dynamo_apikeys(self: core.Construct) -> dynamodb.Table:
        # parition is sub+YYYY+MM+DD
        table = dynamodb.Table(
            self,
            "ApiKeys",
            partition_key=dynamodb.Attribute(
                name="key_hash", type=dynamodb.AttributeType.STRING
            ),
            write_capacity=1,
        )

        table.auto_scale_read_capacity(
            min_capacity=5,
            max_capacity=50,
        )

        return table
