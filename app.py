#!/usr/bin/env python3
import os

from aws_cdk import core

from phish_food.main_stack import PhishFood

env_USA = core.Environment(
    account=os.getenv("AWS_ACCOUNT"), region=os.getenv("AWS_REGION")
)

app = core.App()
PhishFood(app, "PhishFood", env=env_USA)

app.synth()
