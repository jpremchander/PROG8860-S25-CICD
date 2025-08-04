pipeline {
    agent any

    environment {
        // Azure credentials from Jenkins credentials store
        AZURE_CLIENT_ID = credentials('azure-client-id')
        AZURE_CLIENT_SECRET = credentials('azure-client-secret')
        AZURE_TENANT_ID = credentials('azure-tenant-id')
        AZURE_SUBSCRIPTION_ID = credentials('azure-subscription-id')

        // Azure Function App settings — update these
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
                '''
            }
        }

/*
        stage('Package') {
            steps {
                echo 'Packaging Azure Function app...'
                sh '''
                    rm -rf deployment function-app.zip
                    mkdir deployment
                    cp package.json host.json deployment/
                    cp -r src deployment/src
                    cd deployment
                    npm install --production
                    cd ..
                    zip -r function-app.zip deployment/*
                '''
            }
        }

*/
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
                sh '''
                    FUNCTION_URL=$(az functionapp function show \
                        --resource-group "$RESOURCE_GROUP" \
                        --name "$FUNCTION_APP_NAME" \
                        --function-name HttpExample \
                        --query "invokeUrlTemplate" \
                        --output tsv)

                    echo "Function URL: $FUNCTION_URL"

                    RESPONSE=$(curl -s -w "%{http_code}" $FUNCTION_URL?name=TestUser)
                    echo "Response: $RESPONSE"

                    if echo "$RESPONSE" | grep -q "Hello"; then
                        echo "✅ Deployment verification succeeded."
                    else
                        echo "❌ Deployment verification failed."
                        exit 1
                    fi
                '''
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
