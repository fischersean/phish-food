from aws_cdk import (
    core,
    aws_route53 as route53,
    aws_ec2 as ec2,
)

from phish_food.backend import BackendStack
from phish_food.frontend import FrontendStack


class PhishFood(core.Stack):
    def __init__(
        self, scope: core.Construct, construct_id: str, **kwargs
    ) -> None:

        super().__init__(scope, construct_id, **kwargs)

        vpc = ec2.Vpc(
            self,
            "PhishFood-VPC",
            nat_gateways=0,  # $1/day is too damn high
        )

        # Route53 hosted zone
        zone = route53.HostedZone.from_lookup(
            self,
            "DomainHostedZone",
            domain_name="thekettle.org",
        )

        # TODO: Create both a SSL certificate for us-east-1 and us-east-2
        # us-east-1 will go to the frontend while us-east-2 will go to the backend

        frontend = FrontendStack(self, "Frontend", hosted_zone=zone)
        backend = BackendStack(self, "Backend", vpc=vpc, hosted_zone=zone)
