pipeline {
    agent any
    
    environment {
        AZURE_SUBSCRIPTION_ID = credentials('azure-subscription-id')
        AZURE_TENANT_ID = credentials('azure-tenant-id')
        AZURE_CLIENT_ID = credentials('azure-client-id')
        AZURE_CLIENT_SECRET = credentials('azure-client-secret')
        RESOURCE_GROUP = 'premfunc8860_group'
        FUNCTION_APP_NAME = 'premfunc8860'
        AZURE_REGION = 'East US'
    }
    
    stages {
        stage('Build') {
            steps {
                echo '🔧 Building Azure Functions project...'
                
                script {
                    // Clean any previous builds
                    if (fileExists('node_modules')) {
                        bat 'rmdir /s /q node_modules'
                    }
                    
                    // Install dependencies
                    bat 'npm install'
                    
                    echo '✅ Build completed successfully'
                }
            }
        }
        
        stage('Test') {
            steps {
                echo '🧪 Running tests...'
                
                script {
                    try {
                        // Run Jest tests
                        bat 'npm test'
                        echo '✅ All tests passed successfully'
                    } catch (Exception e) {
                        echo "❌ Tests failed: ${e.message}"
                        currentBuild.result = 'FAILURE'
                        error("Test stage failed")
                    }
                }
            }
        }
        
        stage('Deploy') {
            steps {
                echo '🚀 Deploying to Azure Functions...'
                
                script {
                    try {
                        // Login to Azure using service principal
                        bat """
                            az login --service-principal -u %AZURE_CLIENT_ID% -p %AZURE_CLIENT_SECRET% --tenant %AZURE_TENANT_ID%
                        """
                        
                        // Set the subscription
                        bat """
                            az account set --subscription %AZURE_SUBSCRIPTION_ID%
                        """
                        
                        // Verify resource group exists
                        def rgExists = bat(
                            script: "az group show --name %RESOURCE_GROUP% --output table",
                            returnStatus: true
                        )
                        
                        if (rgExists != 0) {
                            echo "Creating resource group: ${RESOURCE_GROUP}"
                            bat """
                                az group create --name %RESOURCE_GROUP% --location "%AZURE_REGION%"
                            """
                        }
                        
                        // Check if function app exists
                        def appExists = bat(
                            script: "az functionapp show --name %FUNCTION_APP_NAME% --resource-group %RESOURCE_GROUP% --output table",
                            returnStatus: true
                        )
                        
                        if (appExists != 0) {
                            echo "Creating function app: ${FUNCTION_APP_NAME}"
                            bat """
                                az functionapp create --resource-group %RESOURCE_GROUP% --consumption-plan-location "%AZURE_REGION%" --runtime node --runtime-version 18 --functions-version 4 --name %FUNCTION_APP_NAME% --storage-account premfunc8860storage
                            """
                        }
                        
                        // Create deployment package
                        echo "Creating deployment package..."
                        
                        // Create a clean directory for deployment
                        if (fileExists('deploy')) {
                            bat 'rmdir /s /q deploy'
                        }
                        bat 'mkdir deploy'
                        
                        // Copy necessary files
                        bat 'copy package.json deploy\\'
                        bat 'copy host.json deploy\\'
                        bat 'xcopy src deploy\\src\\ /E /I'
                        
                        // Install production dependencies in deploy folder
                        bat 'cd deploy && npm install --production'
                        
                        // Create zip package
                        bat 'cd deploy && powershell "Compress-Archive -Path * -DestinationPath ..\\deployment.zip -Force"'
                        
                        // Deploy the function app
                        echo "Deploying function app..."
                        bat """
                            az functionapp deployment source config-zip --resource-group %RESOURCE_GROUP% --name %FUNCTION_APP_NAME% --src deployment.zip
                        """
                        
                        // Wait for deployment to complete
                        sleep 30
                        
                        echo '✅ Deployment completed successfully'
                        echo "🌐 Function App URL: https://${FUNCTION_APP_NAME}.azurewebsites.net"
                        
                    } catch (Exception e) {
                        echo "❌ Deployment failed: ${e.message}"
                        currentBuild.result = 'FAILURE'
                        error("Deploy stage failed")
                    } finally {
                        // Logout from Azure
                        bat 'az logout || echo "Already logged out"'
                    }
                }
            }
        }
    }
    
    post {
        always {
            echo '🧹 Cleaning up workspace...'
            
            // Clean up deployment artifacts
            script {
                if (fileExists('deployment.zip')) {
                    bat 'del deployment.zip'
                }
                if (fileExists('deploy')) {
                    bat 'rmdir /s /q deploy'
                }
            }
        }
        
        success {
            echo '🎉 Pipeline completed successfully!'
            echo "✅ Build: Successful"
            echo "✅ Test: All tests passed"
            echo "✅ Deploy: Function deployed to Azure"
            echo "🌐 Access your function at: https://${FUNCTION_APP_NAME}.azurewebsites.net/api/HelloWorld"
        }
        
        failure {
            echo '💥 Pipeline failed!'
            echo "❌ Check the logs above for details"
        }
    }
}