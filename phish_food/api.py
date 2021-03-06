import os

from aws_cdk import (
    core,
    aws_s3 as s3,
    aws_dynamodb as dynamodb,
    aws_certificatemanager as certificates,
    aws_route53 as route53,
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
        count_results_table: dynamodb.Table,
        reddit_archive_bucket: s3.Bucket,
        api_key_table: dynamodb.Table,
        hosted_zone: route53.HostedZone,
        **kwargs,
    ) -> None:

        super().__init__(scope, construct_id, **kwargs)

        """
        Reddit data access API
        """
        ecs_service = self.ecs_get_countresults(
            cluster,
            count_results_table,
            reddit_archive_bucket,
            api_key_table,
            hosted_zone,
        )

        return

    def ecs_get_countresults(
        self,
        cluster: ecs.Cluster,
        count_results_table: dynamodb.Table,
        reddit_archive_bucket: s3.Bucket,
        api_key_table: dynamodb.Table,
        hosted_zone: route53.HostedZone,
    ):

        # TODO: Make this a env variable at build time
        port = os.getenv("API_PORT")
        if port == "":
            raise ValueError("Could not find API_PORT env variable")

        MAX_CAPACITY = 5

        # autoscale_group = cluster.add_capacity(
        # "DefaultAutoScalingGroup",
        # instance_type=ec2.InstanceType("t3a.micro"),
        # vpc_subnets=ec2.SubnetSelection(subnet_type=ec2.SubnetType.PUBLIC),
        # can_containers_access_instance_role=True,
        # max_capacity=MAX_CAPACITY,
        # )

        # autoscale_group.scale_on_cpu_utilization(
        # "CpuUtilizationScaler",
        # target_utilization_percent=95,
        # )

        ecs_service = ecs_patterns.ApplicationLoadBalancedFargateService(
            self,
            "ApiEcs",
            cluster=cluster,
            cpu=256,
            memory_limit_mib=512,
            task_image_options=ecs_patterns.ApplicationLoadBalancedTaskImageOptions(
                image=ecs.ContainerImage.from_asset(
                    ".",
                    file="Dockerfile.api",
                ),
                environment={
                    "API_PORT": port,
                    "ETL_RESULTS_TABLE": count_results_table.table_name,
                    "REDDIT_ARCHIVE_BUCKET": reddit_archive_bucket.bucket_name,
                    "API_KEY_TABLE": api_key_table.table_name,
                    "AWS_REGION": os.getenv("AWS_REGION"),
                },
                container_port=int(port),
                enable_logging=True,
            ),
            domain_name="api",
            domain_zone=hosted_zone,
            certificate=certificates.Certificate.from_certificate_arn(
                self,
                "DomainCertificateEast2",
                "arn:aws:acm:us-east-2:261392311630:certificate/8509c657-9ad9-4c9a-80e2-f11d9535b13d",
            ),
            redirect_http=True,
        )

        ecs_service.service.auto_scale_task_count(
            max_capacity=MAX_CAPACITY,
        )

        count_results_table.grant_read_write_data(
            ecs_service.task_definition.task_role
        )
        reddit_archive_bucket.grant_read_write(
            ecs_service.task_definition.task_role
        )
        api_key_table.grant_read_write_data(
            ecs_service.task_definition.task_role
        )

        return ecs_service
