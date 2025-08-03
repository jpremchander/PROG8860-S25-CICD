pipeline {
    agent any
    
    environment {
        // Azure credentials - these should be configured in Jenkins credentials
        AZURE_CLIENT_ID = credentials('azure-client-id')
        AZURE_CLIENT_SECRET = credentials('azure-client-secret')
        AZURE_TENANT_ID = credentials('azure-tenant-id')
        AZURE_SUBSCRIPTION_ID = credentials('azure-subscription-id')
        
        // Azure Function App details - update these with your actual values
        RESOURCE_GROUP = 'your-resource-group'
        FUNCTION_APP_NAME = 'your-function-app-name'
        
        // Node.js version
        NODE_VERSION = '18'
    }
    
    stages {
        stage('Checkout') {
            steps {
                script {
                    echo 'Checking out source code...'
                    checkout scm
                }
            }
        }
        
        stage('Setup Node.js') {
            steps {
                script {
                    echo 'Setting up Node.js environment...'
                    // Install Node.js if not available
                    bat '''
                        node --version || echo "Node.js not found"
                        npm --version || echo "npm not found"
                    '''
                }
            }
        }
        
        stage('Build') {
            steps {
                script {
                    echo 'Building the application...'
                    bat '''
                        echo "Installing dependencies..."
                        npm install
                        echo "Build completed successfully!"
                    '''
                }
            }
        }
        
        stage('Test') {
            steps {
                script {
                    echo 'Running tests...'
                    bat '''
                        echo "Running Jest test suite..."
                        npm test
                        echo "All tests completed!"
                    '''
                }
            }
            post {
                always {
                    // Publish test results
                    publishTestResults testResultsPattern: 'coverage/lcov.info'
                    
                    // Archive test coverage reports
                    publishHTML([
                        allowMissing: false,
                        alwaysLinkToLastBuild: true,
                        keepAll: true,
                        reportDir: 'coverage',
                        reportFiles: 'index.html',
                        reportName: 'Coverage Report'
                    ])
                }
            }
        }
        
        stage('Package') {
            steps {
                script {
                    echo 'Packaging Azure Function...'
                    bat '''
                        echo "Creating deployment package..."
                        
                        REM Create a deployment directory
                        if exist "deployment" rmdir /s /q deployment
                        mkdir deployment
                        
                        REM Copy necessary files
                        copy package.json deployment\\
                        copy host.json deployment\\
                        xcopy /E /I HelloWorld deployment\\HelloWorld
                        
                        REM Install production dependencies
                        cd deployment
                        npm install --production
                        cd ..
                        
                        REM Create ZIP file for deployment
                        powershell -Command "Compress-Archive -Path 'deployment\\*' -DestinationPath 'function-app.zip' -Force"
                        
                        echo "Package created: function-app.zip"
                    '''
                }
            }
        }
        
        stage('Deploy') {
            steps {
                script {
                    echo 'Deploying to Azure Functions...'
                    bat '''
                        echo "Logging in to Azure..."
                        az login --service-principal -u %AZURE_CLIENT_ID% -p %AZURE_CLIENT_SECRET% --tenant %AZURE_TENANT_ID%
                        
                        echo "Setting Azure subscription..."
                        az account set --subscription %AZURE_SUBSCRIPTION_ID%
                        
                        echo "Deploying function app..."
                        az functionapp deployment source config-zip ^
                            --resource-group %RESOURCE_GROUP% ^
                            --name %FUNCTION_APP_NAME% ^
                            --src function-app.zip
                        
                        echo "Deployment completed successfully!"
                        
                        echo "Function URL:"
                        az functionapp function show ^
                            --resource-group %RESOURCE_GROUP% ^
                            --name %FUNCTION_APP_NAME% ^
                            --function-name HelloWorld ^
                            --query "invokeUrlTemplate" ^
                            --output tsv
                    '''
                }
            }
        }
        
        stage('Verify Deployment') {
            steps {
                script {
                    echo 'Verifying deployment...'
                    bat '''
                        echo "Getting function URL..."
                        for /f %%i in ('az functionapp function show --resource-group %RESOURCE_GROUP% --name %FUNCTION_APP_NAME% --function-name HelloWorld --query "invokeUrlTemplate" --output tsv') do set FUNCTION_URL=%%i
                        
                        echo "Function URL: %FUNCTION_URL%"
                        
                        echo "Testing the deployed function..."
                        powershell -Command "try { $response = Invoke-RestMethod -Uri '%FUNCTION_URL%' -Method Get; Write-Host 'Response: ' $response; if ($response -match 'Hello') { Write-Host 'Deployment verification: SUCCESS' } else { Write-Host 'Deployment verification: FAILED'; exit 1 } } catch { Write-Host 'Deployment verification: FAILED - ' $_.Exception.Message; exit 1 }"
                    '''
                }
            }
        }
    }
    
    post {
        always {
            // Clean up workspace
            script {
                echo 'Cleaning up...'
                bat '''
                    if exist "deployment" rmdir /s /q deployment
                    if exist "function-app.zip" del function-app.zip
                    if exist "node_modules" rmdir /s /q node_modules
                '''
            }
        }
        
        success {
            echo 'Pipeline completed successfully! ✅'
            echo 'Your Azure Function has been deployed and verified.'
        }
        
        failure {
            echo 'Pipeline failed! ❌'
            echo 'Please check the logs for more details.'
        }
    }
}
