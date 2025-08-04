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

                    # Ensure proper runtime settings
                    echo "Setting Function App configuration..."
                    az functionapp config appsettings set \
                        --resource-group "$RESOURCE_GROUP" \
                        --name "$FUNCTION_APP_NAME" \
                        --settings \
                            FUNCTIONS_WORKER_RUNTIME=node \
                            WEBSITE_NODE_DEFAULT_VERSION=~18 \
                            FUNCTIONS_EXTENSION_VERSION=~4

                    # Deploy the function
                    echo "Deploying function package..."
                    az functionapp deployment source config-zip \
                        --resource-group "$RESOURCE_GROUP" \
                        --name "$FUNCTION_APP_NAME" \
                        --src function-app.zip

                    # Wait a moment for the deployment to stabilize
                    echo "Waiting for deployment to stabilize..."
                    sleep 10
                '''
            }
        }

        stage('Verify Deployment') {
            steps {
                echo 'Verifying deployed function...'
                script {
                    // Get the function URL
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

                    // Also get the base app URL as a fallback
                    def appUrl = sh(
                        script: """az functionapp show \
                            --resource-group "$RESOURCE_GROUP" \
                            --name "$FUNCTION_APP_NAME" \
                            --query defaultHostName \
                            --output tsv""",
                        returnStdout: true
                    ).trim()
                    
                    def baseUrl = "https://${appUrl}/api/HelloWorld"
                    echo "Base URL (fallback): ${baseUrl}"

                    def maxRetries = 8
                    def sleepSeconds = 10
                    def success = false
                    def urlsToTry = [functionUrl, baseUrl]

                    for (String testUrl : urlsToTry) {
                        echo "Testing URL: ${testUrl}"
                        
                        for (int i = 1; i <= maxRetries; i++) {
                            try {
                                // First check if the function endpoint is reachable
                                def statusCheck = sh(
                                    script: "curl -s -o /dev/null -w '%{http_code}' '${testUrl}'",
                                    returnStdout: true
                                ).trim()
                                
                                echo "Attempt ${i} - HTTP Status: ${statusCheck}"
                                
                                if (statusCheck == "200") {
                                    // If we get 200, test with parameter
                                    def response = sh(
                                        script: "curl -s '${testUrl}?name=TestUser'",
                                        returnStdout: true
                                    ).trim()

                                    echo "Attempt ${i} response: '${response}'"

                                    if (response && response.toLowerCase().contains("hello")) {
                                        echo "✅ Deployment verification succeeded with URL: ${testUrl}"
                                        success = true
                                        break
                                    }
                                } else if (statusCheck == "404") {
                                    echo "Function not found (404), trying different case..."
                                    // Try lowercase version
                                    def lowercaseUrl = testUrl.replace("/HelloWorld", "/helloworld")
                                    def lowercaseResponse = sh(
                                        script: "curl -s '${lowercaseUrl}?name=TestUser'",
                                        returnStdout: true
                                    ).trim()
                                    
                                    if (lowercaseResponse && lowercaseResponse.toLowerCase().contains("hello")) {
                                        echo "✅ Deployment verification succeeded with lowercase URL: ${lowercaseUrl}"
                                        success = true
                                        break
                                    }
                                }
                                
                            } catch (Exception e) {
                                echo "Error during verification attempt ${i}: ${e.getMessage()}"
                            }
                            
                            if (!success) {
                                echo "Function not ready yet, retrying in ${sleepSeconds} seconds..."
                                sleep sleepSeconds
                            }
                        }
                        
                        if (success) break
                    }

                    if (!success) {
                        // Try to get more diagnostic information
                        echo "Getting diagnostic information..."
                        sh """
                            echo "Function App Status:"
                            az functionapp show --resource-group "$RESOURCE_GROUP" --name "$FUNCTION_APP_NAME" --query "{state: state, hostNames: hostNames}" --output table
                            
                            echo "Function List:"
                            az functionapp function list --resource-group "$RESOURCE_GROUP" --name "$FUNCTION_APP_NAME" --query "[].{name: name, state: properties.config.bindings[0].type}" --output table
                            
                            echo "Recent logs (if available):"
                            az functionapp log tail --resource-group "$RESOURCE_GROUP" --name "$FUNCTION_APP_NAME" --timeout 5 || echo "Logs not available"
                        """
                        
                        error("❌ Deployment verification failed after ${maxRetries} attempts on all URLs.")
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
