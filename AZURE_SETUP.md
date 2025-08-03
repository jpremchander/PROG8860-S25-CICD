# Azure Setup Instructions

This document provides step-by-step instructions for setting up Azure resources for the CI/CD pipeline.

## 1. Create Azure Function App

### Using Azure Portal

1. Log in to the [Azure Portal](https://portal.azure.com)
2. Click "Create a resource" > "Function App"
3. Fill in the details:
   - **Subscription**: Your subscription
   - **Resource Group**: Create new or use existing (e.g., `rg-cicd-assignment`)
   - **Function App Name**: Choose a unique name (e.g., `fa-cicd-yourname`)
   - **Runtime Stack**: Node.js
   - **Version**: 18 LTS
   - **Region**: Choose closest region
4. Click "Review + Create" and then "Create"

### Using Azure CLI

```bash
# Create resource group
az group create --name rg-cicd-assignment --location eastus

# Create storage account (required for Function App)
az storage account create --name stcicdyourname --resource-group rg-cicd-assignment --location eastus --sku Standard_LRS

# Create Function App
az functionapp create --resource-group rg-cicd-assignment --consumption-plan-location eastus --runtime node --runtime-version 18 --functions-version 4 --name fa-cicd-yourname --storage-account stcicdyourname
```

## 2. Create Service Principal for Jenkins

### Using Azure CLI

```bash
# Create service principal
az ad sp create-for-rbac --name "jenkins-cicd-sp" --role contributor --scopes /subscriptions/{subscription-id}/resourceGroups/rg-cicd-assignment

# Note down the output:
# - appId (this is your client ID)
# - password (this is your client secret)
# - tenant (this is your tenant ID)
```

### Get Required Information

```bash
# Get subscription ID
az account show --query id --output tsv

# Get tenant ID
az account show --query tenantId --output tsv
```

## 3. Test Manual Deployment

Before setting up the CI/CD pipeline, test manual deployment:

```bash
# Login to Azure
az login

# Deploy the function (after creating the ZIP package)
az functionapp deployment source config-zip --resource-group rg-cicd-assignment --name fa-cicd-yourname --src function-app.zip
```

## 4. Update Jenkinsfile

Replace the placeholder values in the Jenkinsfile:

```groovy
RESOURCE_GROUP = 'rg-cicd-assignment'
FUNCTION_APP_NAME = 'fa-cicd-yourname'
```

## 5. Configure Jenkins Credentials

In Jenkins:
1. Go to "Manage Jenkins" > "Manage Credentials"
2. Add these credentials:
   - `azure-client-id`: The appId from service principal
   - `azure-client-secret`: The password from service principal
   - `azure-tenant-id`: The tenant from service principal
   - `azure-subscription-id`: Your subscription ID
