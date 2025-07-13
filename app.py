#!/usr/bin/env python3
import aws_cdk as cdk
import sys
from cdk_project_9015480.cdk_project_9015480_stack import CdkProject9015480Stack

try:
    app = cdk.App()
    CdkProject9015480Stack(app, "CdkProject9015480Stack")
    app.synth()
except Exception as e:
    print("Error during CDK synth:", e)
    sys.exit(1)
