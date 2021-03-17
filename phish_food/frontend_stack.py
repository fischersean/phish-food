import os

from aws_cdk import (
    core,
    aws_s3 as s3,
    aws_certificatemanager as certificates,
    aws_s3_deployment as s3deploy,
    aws_cloudfront as cloudfront,
    aws_cloudfront_origins as origins,
)


def S3_FRONTENDDEPLOY(stack: core.Construct) -> s3.Bucket:
    bucket = s3.Bucket(
        stack,
        "FrontendBucket",
        removal_policy=core.RemovalPolicy.DESTROY,
        versioned=True,
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
        default_root_object="index.html",
        default_behavior=cloudfront.BehaviorOptions(
            origin=origins.S3Origin(
                bucket,
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


class PhishFoodFrontendStack(core.Stack):
    def __init__(
        self, scope: core.Construct, construct_id: str, **kwargs
    ) -> None:

        super().__init__(scope, construct_id, **kwargs)

        os.system("cd web && npm run build")
        bucket = s3.Bucket(
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
                    bucket,
                ),
                viewer_protocol_policy=cloudfront.ViewerProtocolPolicy.REDIRECT_TO_HTTPS,
            ),
        )

        s3deploy.BucketDeployment(
            self,
            "FrontendS3Deployment",
            sources=[s3deploy.Source.asset("./web/dist")],
            destination_bucket=bucket,
            distribution=distribution,
        )
