from aws_cdk import (
    Stack,
    RemovalPolicy,
    aws_s3 as s3,
    aws_lambda as _lambda,
    aws_dynamodb as ddb,
)
from constructs import Construct

class CdkProject9015480Stack(Stack):
    def __init__(self, scope: Construct, construct_id: str, **kwargs) -> None:
        super().__init__(scope, construct_id, **kwargs)

        # ✅ S3 Bucket
        my_bucket = s3.Bucket(self, "S3Bucket9015480",
            versioned=True,
            removal_policy=RemovalPolicy.DESTROY
        )

        # ✅ Lambda Function
        my_lambda = _lambda.Function(self, "LambdaFunction9015480",
            runtime=_lambda.Runtime.PYTHON_3_9,
            handler="index.handler",
            code=_lambda.InlineCode("""
def handler(event, context):
    print("Hello from Lambda!")
    return {
        'statusCode': 200,
        'body': 'Hello from 9015480!'
    }
"""),
            environment={
                "BUCKET_NAME": my_bucket.bucket_name
            }
        )

        # ✅ DynamoDB Table
        my_table = ddb.Table(self, "DynamoDBTable9015480",
            partition_key={"name": "id", "type": ddb.AttributeType.STRING},
            table_name="DynamoTable9015480",
            removal_policy=RemovalPolicy.DESTROY
        )
