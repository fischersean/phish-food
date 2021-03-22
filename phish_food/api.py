import os

from aws_cdk import (
    core,
    aws_lambda as lambda_,
    aws_s3 as s3,
    aws_dynamodb as dynamodb,
    aws_apigateway as apigateway,
    aws_certificatemanager as certificates,
    aws_route53 as route53,
    aws_route53_targets as alias,
    aws_ecs_patterns as ecs_patterns,
    aws_ecs as ecs,
    aws_ec2 as ec2,
    aws_elasticloadbalancingv2 as loadbalancing,
)


class ApiStack(core.NestedStack):
    def __init__(
        self,
        scope: core.Construct,
        construct_id: str,
        vpc: ec2.Vpc,
        cluster: ecs.Cluster,
        count_results_table: str,
        reddit_archive_table: str,
        hosted_zone: route53.HostedZone,
        **kwargs
    ) -> None:

        super().__init__(scope, construct_id, **kwargs)

        # TODO: This needs to be changed to a dynamic cert request
        # instead of relying on a cert already existing
        # self.api = apigateway.RestApi(
        # self,
        # "RedditTrendsAPI",
        # domain_name=apigateway.DomainNameOptions(
        # domain_name="api.thekettle.org",
        # certificate=certificates.Certificate.from_certificate_arn(
        # self,
        # "DomainCertificateEast2",
        # "arn:aws:acm:us-east-2:261392311630:certificate/8509c657-9ad9-4c9a-80e2-f11d9535b13d",
        # ),
        # security_policy=apigateway.SecurityPolicy.TLS_1_2,
        # ),
        # deploy_options=apigateway.StageOptions(
        # stage_name="prod",
        # ),
        # )

        # self.api.root.add_method("ANY")

        """
        Reddit data access API
        """
        ecs_service = self.ecs_get_countresults(
            cluster, count_results_table, reddit_archive_table
        )

        return

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

    def ecs_get_countresults(
        self,
        cluster: ecs.Cluster,
        count_results_table: dynamodb.Table,
        reddit_archive_table: dynamodb.Table,
    ):

        # TODO: Make this a env variable at build time
        port = os.getenv("API_PORT")
        if port == "":
            raise ValueError("Could not find API_PORT env variable")

        cluster.add_capacity(
            "DefaultAutoScalingGroup",
            instance_type=ec2.InstanceType("t2.small"),
            vpc_subnets=ec2.SubnetSelection(subnet_type=ec2.SubnetType.PUBLIC),
            can_containers_access_instance_role=True,
        )

        ecs_service = ecs_patterns.ApplicationLoadBalancedEc2Service(
            self,
            "ApiEcs",
            cluster=cluster,
            memory_limit_mib=512,
            task_image_options=ecs_patterns.ApplicationLoadBalancedTaskImageOptions(
                image=ecs.ContainerImage.from_asset(
                    ".",
                    file="Dockerfile.api",
                ),
                environment={
                    "API_PORT": port,
                    "ETL_RESULTS_TABLE": count_results_table.table_name,
                    "REDDIT_ARCHIVE_TABLE": reddit_archive_table.table_name,
                    "AWS_REGION": os.getenv("AWS_REGION"),
                },
                container_port=8000,
                enable_logging=True,
            ),
            # public_load_balancer=True,
        )

        count_results_table.grant_read_write_data(
            ecs_service.task_definition.task_role
        )
        reddit_archive_table.grant_read_write_data(
            ecs_service.task_definition.task_role
        )

        # # Configure health check
        # ecs_service.target_group.configure_health_check(
        # protocol=loadbalancing.Protocol.HTTP,
        # )

        return ecs_service
