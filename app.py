#!/usr/bin/env python3
from aws_cdk import core

from phish_food.main_stack import PhishFood

app = core.App()
PhishFood(app, "PhishFood")

app.synth()
