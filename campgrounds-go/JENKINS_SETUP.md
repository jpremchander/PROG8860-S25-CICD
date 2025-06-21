# Jenkins Configuration Guide

## Required Jenkins Credentials

You need to configure the following credentials in Jenkins:

### 1. AWS Credentials
- **Credential ID**: `AWSCREDENTIALS`
- **Type**: AWS Credentials
- **Access Key ID**: Your AWS Access Key
- **Secret Access Key**: Your AWS Secret Key
- **Description**: AWS credentials for S3 access

### 2. Session Secrets
- **Credential ID**: `SESSION_SECRET_DEV`
- **Type**: Secret text
- **Secret**: Generate with `openssl rand -hex 32`
- **Description**: Development session secret

- **Credential ID**: `SESSION_SECRET_PROD`
- **Type**: Secret text
- **Secret**: Generate with `openssl rand -hex 32` (different from dev)
- **Description**: Production session secret

### 3. MongoDB URIs
- **Credential ID**: `MONGODB_URI_DEV`
- **Type**: Secret text
- **Secret**: `mongodb://localhost:27017/campgrounds_dev`
- **Description**: Development MongoDB connection string

- **Credential ID**: `MONGODB_URI_PROD`
- **Type**: Secret text
- **Secret**: `mongodb://localhost:27017/campgrounds_prod`
- **Description**: Production MongoDB connection string

## Setup Steps

### Step 1: Generate Session Secrets
\`\`\`bash
# Generate development secret
echo "Dev Session Secret: $(openssl rand -hex 32)"

# Generate production secret  
echo "Prod Session Secret: $(openssl rand -hex 32)"
\`\`\`

### Step 2: Add Credentials to Jenkins
1. Go to Jenkins Dashboard
2. Click "Manage Jenkins"
3. Click "Manage Credentials"
4. Click "System" → "Global credentials (unrestricted)"
5. Click "Add Credentials"

### Step 3: Configure Each Credential

#### AWS Credentials
- Kind: AWS Credentials
- ID: `AWSCREDENTIALS`
- Access Key ID: [Your AWS Access Key]
- Secret Access Key: [Your AWS Secret Key]
- Description: AWS credentials for S3 access

#### Session Secret - Development
- Kind: Secret text
- Secret: [Generated dev secret from Step 1]
- ID: `SESSION_SECRET_DEV`
- Description: Development session secret

#### Session Secret - Production
- Kind: Secret text
- Secret: [Generated prod secret from Step 1]
- ID: `SESSION_SECRET_PROD`
- Description: Production session secret

#### MongoDB URI - Development
- Kind: Secret text
- Secret: `mongodb://localhost:27017/campgrounds_dev`
- ID: `MONGODB_URI_DEV`
- Description: Development MongoDB connection

#### MongoDB URI - Production
- Kind: Secret text
- Secret: `mongodb://localhost:27017/campgrounds_prod`
- ID: `MONGODB_URI_PROD`
- Description: Production MongoDB connection

### Step 4: Verify Credentials
After adding all credentials, you should see:
- AWSCREDENTIALS (AWS Credentials)
- SESSION_SECRET_DEV (Secret text)
- SESSION_SECRET_PROD (Secret text)
- MONGODB_URI_DEV (Secret text)
- MONGODB_URI_PROD (Secret text)

### Step 5: Test Pipeline
1. Push your code to the repository
2. Trigger a Jenkins build
3. Verify all credentials are loaded correctly
4. Check deployment logs for any credential issues

## Security Best Practices

✅ **DO:**
- Use different secrets for dev/prod
- Rotate secrets regularly
- Use Jenkins credentials for all sensitive data
- Monitor credential access logs

❌ **DON'T:**
- Commit real secrets to version control
- Share credentials in plain text
- Use the same secret across environments
- Log credential values in build output

## Troubleshooting

### Common Issues:
1. **Credential not found**: Check credential ID matches exactly
2. **Permission denied**: Ensure Jenkins user has access to credentials
3. **Invalid secret format**: Verify secret generation and format
4. **AWS access denied**: Check AWS credentials and S3 bucket permissions

### Debug Commands:
\`\`\`bash
# Test AWS credentials
aws sts get-caller-identity

# Test S3 access
aws s3 ls s3://yelp-camp-project/

# Test MongoDB connection
mongo mongodb://localhost:27017/campgrounds_dev --eval "db.stats()"
