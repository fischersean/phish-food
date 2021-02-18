#!/usr/bin/env python3

from aws_cdk import core

from phish_food.phish_food_stack import PhishFoodStack


app = core.App()
PhishFoodStack(app, "phish-food")

app.synth()
