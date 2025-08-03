# Azure Functions CI/CD Pipeline with Jenkins

This project demonstrates a complete CI/CD pipeline using Jenkins to deploy Azure Functions. The pipeline includes automated testing, building, and deployment to Azure.

## Project Structure

```
PROG8860-S25-CICD/
├── HelloWorld/
│   ├── function.json          # Azure Function configuration
│   └── index.js              # Main function code
├── tests/
│   ├── hello-world.test.js   # Test cases
│   └── setup.js              # Jest setup
├── package.json              # Node.js dependencies
├── host.json                 # Azure Functions host configuration
├── jest.config.js            # Jest testing configuration
├── Jenkinsfile               # Jenkins pipeline definition
└── README.md                 # This file
```

## Features

- **HTTP-triggered Azure Function** that returns a personalized greeting
- **Comprehensive test suite** with 5 test cases covering various scenarios
- **Jenkins CI/CD pipeline** with Build, Test, and Deploy stages
- **Automated deployment** to Azure Functions
- **Deployment verification** to ensure the function is working correctly

## Prerequisites

### Azure Setup
1. Azure subscription with access to Azure Functions
2. Azure Function App created in your Azure portal
3. Azure Service Principal for Jenkins authentication

### Jenkins Setup
1. Jenkins server (local or cloud-based)
2. Required Jenkins plugins:
   - GitHub Plugin
   - Azure CLI Plugin
   - Pipeline Plugin
   - HTML Publisher Plugin (for test reports)

### Local Development
1. Node.js 18+ installed
2. Azure Functions Core Tools
3. Azure CLI

## Quick Start

### 1. Clone the Repository
```bash
git clone <your-repo-url>
cd PROG8860-S25-CICD
```

### 2. Install Dependencies
```bash
npm install
```

### 3. Run Tests Locally
```bash
npm test
```

### 4. Run Function Locally (Optional)
```bash
func start
```

## Jenkins Configuration

### 1. Create Jenkins Credentials
In Jenkins, go to "Manage Jenkins" > "Manage Credentials" and add:

- `azure-client-id`: Your Azure Service Principal client ID
- `azure-client-secret`: Your Azure Service Principal client secret
- `azure-tenant-id`: Your Azure tenant ID
- `azure-subscription-id`: Your Azure subscription ID

### 2. Update Jenkinsfile Variables
Edit the `Jenkinsfile` and update these environment variables:

```groovy
RESOURCE_GROUP = 'your-actual-resource-group'
FUNCTION_APP_NAME = 'your-actual-function-app-name'
```

### 3. Create Jenkins Pipeline Job
1. In Jenkins, create a new "Pipeline" job
2. Under "Pipeline" section, select "Pipeline script from SCM"
3. Set SCM to "Git" and provide your repository URL
4. Set the script path to `Jenkinsfile`

## Azure Function Details

### Function Endpoint
The function responds to both GET and POST requests:
- **GET**: `https://your-function-app.azurewebsites.net/api/HelloWorld`
- **GET with name**: `https://your-function-app.azurewebsites.net/api/HelloWorld?name=YourName`
- **POST**: Send JSON body with `{"name": "YourName"}`

### Response Format
```json
{
  "status": 200,
  "body": "Hello, World!",
  "headers": {
    "Content-Type": "text/plain"
  }
}
```

## Test Cases

The project includes 5 comprehensive test cases:

1. **Basic Response Test**: Verifies the function returns "Hello, World!" by default
2. **Status Code Test**: Ensures the function returns HTTP 200 status
3. **Query Parameter Test**: Tests personalized greeting with query parameter
4. **POST Body Test**: Tests personalized greeting from POST request body
5. **Edge Case Test**: Handles empty name parameter gracefully

## CI/CD Pipeline Stages

### 1. Checkout
- Pulls the latest code from GitHub

### 2. Setup Node.js
- Verifies Node.js and npm installation

### 3. Build
- Installs npm dependencies
- Prepares the application for testing and deployment

### 4. Test
- Runs the Jest test suite
- Generates code coverage reports
- Publishes test results and coverage reports

### 5. Package
- Creates a deployment package
- Installs production dependencies
- Creates a ZIP file for Azure deployment

### 6. Deploy
- Authenticates with Azure using Service Principal
- Deploys the function using Azure CLI
- Uploads the deployment package to Azure Functions

### 7. Verify Deployment
- Tests the deployed function endpoint
- Verifies the function is responding correctly

## Troubleshooting

### Common Issues

1. **Azure CLI not found in Jenkins**
   - Install Azure CLI on the Jenkins server
   - Add Azure CLI to the system PATH

2. **Authentication failures**
   - Verify Service Principal credentials in Jenkins
   - Ensure Service Principal has Contributor role on the Resource Group

3. **Test failures**
   - Run tests locally first: `npm test`
   - Check Jest configuration and dependencies

4. **Deployment failures**
   - Verify Resource Group and Function App names
   - Check Azure CLI authentication
   - Ensure Function App allows deployment from external sources

### Viewing Logs
- **Jenkins**: Check the pipeline console output
- **Azure**: Use Azure Portal > Function App > Functions > Monitor
- **Local**: Use `func start` and check terminal output

## Security Considerations

- Never commit Azure credentials to the repository
- Use Jenkins credentials store for sensitive information
- Regularly rotate Service Principal secrets
- Use Azure Key Vault for production environments

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests locally
5. Submit a pull request

## Assignment Requirements Checklist

- ✅ Jenkins Setup (3%): Complete Jenkins configuration with GitHub integration
- ✅ Pipeline Stages (3%): Build, Test, and Deploy stages implemented
- ✅ Test Cases (2%): 5 comprehensive test cases (exceeds requirement of 3)
- ✅ Azure Deployment (2%): Automated deployment with verification

## License

This project is created for educational purposes as part of the PROG8860 course assignment.
