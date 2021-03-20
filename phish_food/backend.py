import os

from aws_cdk import (
    core,
    aws_ecs as ecs,
    aws_ec2 as ec2,
    aws_dynamodb as dynamodb,
)

from phish_food.etl import EtlStack
from phish_food.api import ApiStack


class BackendStack(core.NestedStack):
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

        count_results_table = self.dynamo_scraperesults()
        reddit_archive_table = self.dynamo_redditarchive()

        etl = EtlStack(
            self,
            "ETL",
            vpc=vpc,
            cluster=cluster,
            count_results_table=count_results_table,
            reddit_archive_table=reddit_archive_table,
        )

        api = ApiStack(self, "API", count_results_table=count_results_table)

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

        return table

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
