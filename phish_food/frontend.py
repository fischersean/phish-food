import os

from aws_cdk import (
    core,
    aws_s3 as s3,
    aws_certificatemanager as certificates,
    aws_s3_deployment as s3deploy,
    aws_cloudfront as cloudfront,
    aws_cloudfront_origins as origins,
    aws_route53 as route53,
)


class FrontendStack(core.NestedStack):
    def __init__(
        self, scope: core.Construct, construct_id: str, **kwargs
    ) -> None:

        super().__init__(scope, construct_id, **kwargs)

        frontend_bucket = s3.Bucket(
            self,
            "FrontendBucket",
            removal_policy=core.RemovalPolicy.DESTROY,
            versioned=True,
        )

        distribution = cloudfront.Distribution(
            self,
            "FrontendDistribution",
            # TODO: The domain and cert info should be env vars
            domain_names=["www.thekettle.org"],
            certificate=certificates.Certificate.from_certificate_arn(
                self,
                "DomainCertificateEast1",
                "arn:aws:acm:us-east-1:261392311630:certificate/02a75969-25ce-47d3-acf6-d93408b2eed1",
            ),
            default_root_object="index.html",
            default_behavior=cloudfront.BehaviorOptions(
                origin=origins.S3Origin(
                    frontend_bucket,
                ),
                viewer_protocol_policy=cloudfront.ViewerProtocolPolicy.REDIRECT_TO_HTTPS,
            ),
        )

        s3deploy.BucketDeployment(
            self,
            "FrontendS3Deployment",
            sources=[s3deploy.Source.asset("./web/dist")],
            destination_bucket=frontend_bucket,
            distribution=distribution,
        )

        # Below is currently broken
        # Domain name setup
        # cname_record = route53.CnameRecord(
            # self,
            # "CloudFrontFrontendCnameRecord",
            # # TODO: This needs to be an env variable
            # zone=route53.HostedZone.from_hosted_zone_attributes(
                # self,
                # "DomainHostedZoneId",
                # hosted_zone_id="Z0864562JJDZ4PZXPZGZ",
                # zone_name="thekettle.org",
            # ),
            # domain_name=distribution.domain_name,
        # )
