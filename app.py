#!/usr/bin/env python3
from aws_cdk import core

from phish_food.phish_food_stack import PhishFoodBackendStack
from phish_food.frontend_stack import PhishFoodFrontendStack


app = core.App()
PhishFoodBackendStack(app, "PhishFood")
PhishFoodFrontendStack(app, "PhishFoodFrontend")

app.synth()
