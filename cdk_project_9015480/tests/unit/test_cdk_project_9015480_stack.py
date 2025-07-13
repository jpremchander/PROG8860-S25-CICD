import aws_cdk as core
import aws_cdk.assertions as assertions

from cdk_project_9015480.cdk_project_9015480_stack import CdkProject9015480Stack

# example tests. To run these tests, uncomment this file along with the example
# resource in cdk_project_9015480/cdk_project_9015480_stack.py
def test_sqs_queue_created():
    app = core.App()
    stack = CdkProject9015480Stack(app, "cdk-project-9015480")
    template = assertions.Template.from_stack(stack)

#     template.has_resource_properties("AWS::SQS::Queue", {
#         "VisibilityTimeout": 300
#     })
