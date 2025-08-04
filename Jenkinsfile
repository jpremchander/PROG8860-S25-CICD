pipeline {
    agent any

    environment {
        // Azure credentials from Jenkins credentials store
        AZURE_CLIENT_ID = credentials('azure-client-id')
        AZURE_CLIENT_SECRET = credentials('azure-client-secret')
        AZURE_TENANT_ID = credentials('azure-tenant-id')
        AZURE_SUBSCRIPTION_ID = credentials('azure-subscription-id')

        // Azure Function App settings
        RESOURCE_GROUP = 'premfunc8860_group'
        FUNCTION_APP_NAME = 'premfunc8860'
        NODE_VERSION = '18'
    }

    stages {
        stage('Checkout') {
            steps {
                echo 'Checking out source code...'
                checkout scm
            }
        }

        stage('Setup Node.js') {
            steps {
                echo 'Verifying Node.js...'
                sh 'node --version || echo "Node.js not found"'
                sh 'npm --version || echo "npm not found"'
            }
        }

        stage('Install Dependencies') {
            steps {
                echo 'Installing npm dependencies...'
                sh 'npm install'
            }
        }

        stage('Test') {
            steps {
                echo 'Running tests...'
                sh 'npm test'
            }
        }

        stage('Package') {
            steps {
                echo "\nPackaging Azure Function app..."
                sh '''
                    rm -rf deployment function-app.zip
                    mkdir -p deployment/HelloWorld
                    cp package.json host.json deployment/
                    cp -r HelloWorld deployment/
                    cd deployment
                    zip -r ../function-app.zip HelloWorld host.json package.json
                    cd ..
                '''
            }
        }

        stage('Azure Login & Deploy') {
            steps {
                echo 'Logging into Azure and deploying...'
                sh '''
                    az login --service-principal \
                        -u "$AZURE_CLIENT_ID" \
                        -p "$AZURE_CLIENT_SECRET" \
                        --tenant "$AZURE_TENANT_ID"

                    az account set --subscription "$AZURE_SUBSCRIPTION_ID"

                    az functionapp deployment source config-zip \
                        --resource-group "$RESOURCE_GROUP" \
                        --name "$FUNCTION_APP_NAME" \
                        --src function-app.zip
                '''
            }
        }

        stage('Verify Deployment') {
            steps {
                echo 'Verifying deployed function...'
                script {
                    def functionUrl = sh(
                        script: """az functionapp function show \
                            --resource-group "$RESOURCE_GROUP" \
                            --name "$FUNCTION_APP_NAME" \
                            --function-name HelloWorld \
                            --query invokeUrlTemplate \
                            --output tsv""",
                        returnStdout: true
                    ).trim()

                    echo "Function URL: ${functionUrl}"

                    def maxRetries = 5
                    def sleepSeconds = 6
                    def success = false

                    for (int i = 1; i <= maxRetries; i++) {
                        def response = sh(
                            script: "curl -s '${functionUrl}?name=TestUser'",
                            returnStdout: true
                        ).trim()

                        echo "Attempt ${i} response: ${response}"

                        if (response && response.toLowerCase().contains("hello")) {
                            echo "✅ Deployment verification succeeded."
                            success = true
                            break
                        } else {
                            echo "Function not ready yet, retrying in ${sleepSeconds} seconds..."
                            sleep sleepSeconds
                        }
                    }

                    if (!success) {
                        error("❌ Deployment verification failed after ${maxRetries} attempts.")
                    }
                }
            }
        }
    }

    post {
        always {
            echo 'Cleaning up workspace...'
            sh 'rm -rf deployment function-app.zip node_modules'
        }

        success {
            echo '✅ Pipeline completed successfully!'
        }

        failure {
            echo '❌ Pipeline failed. Check logs for details.'
        }
    }
}
