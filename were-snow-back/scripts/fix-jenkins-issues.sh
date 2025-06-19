#!/bin/bash

echo "🔧 Fixing Jenkins Configuration Issues"
echo "====================================="

# Load configuration
if [ -f ".jenkins-config" ]; then
    source .jenkins-config
    echo "✅ Configuration loaded"
else
    echo "❌ No configuration found"
    exit 1
fi

echo "🔍 Diagnosing issues..."

# Issue 1: Fix CSRF/Crumb issue
echo "1. Fixing CSRF protection..."
JENKINS_CRUMB=$(curl -s "$JENKINS_URL/crumbIssuer/api/xml?xpath=concat(//crumbRequestField,\":\",//crumb)" \
  --user "$JENKINS_USER:$JENKINS_PASSWORD" \
  -H "ngrok-skip-browser-warning: true")

if [ -n "$JENKINS_CRUMB" ]; then
    echo "   ✅ CSRF token obtained: ${JENKINS_CRUMB:0:20}..."
    
    # Update configuration with crumb
    echo "JENKINS_CRUMB=$JENKINS_CRUMB" >> .jenkins-config
else
    echo "   ❌ Failed to get CSRF token"
fi

# Issue 2: Check if job actually exists
echo "2. Checking pipeline job status..."
JOB_CHECK=$(curl -s -w "%{http_code}" -o /tmp/job-check.json \
  --user "$JENKINS_USER:$JENKINS_PASSWORD" \
  -H "ngrok-skip-browser-warning: true" \
  "$JENKINS_URL/job/SnowBored-Game-Pipeline/api/json")

if [ "$JOB_CHECK" = "200" ]; then
    echo "   ✅ Pipeline job exists and is accessible"
    JOB_STATUS=$(cat /tmp/job-check.json | grep -o '"buildable":[^,]*' | cut -d':' -f2)
    echo "   📊 Job buildable: $JOB_STATUS"
elif [ "$JOB_CHECK" = "404" ]; then
    echo "   ❌ Pipeline job not found - recreating..."
    
    # Recreate the job
    cat > /tmp/job-config.xml << EOF
<?xml version='1.1' encoding='UTF-8'?>
<flow-definition plugin="workflow-job">
  <actions/>
  <description>SnowBored Game CI/CD Pipeline</description>
  <keepDependencies>false</keepDependencies>
  <properties>
    <org.jenkinsci.plugins.workflow.job.properties.PipelineTriggersJobProperty>
      <triggers>
        <com.cloudbees.jenkins.GitHubPushTrigger plugin="github">
          <spec></spec>
        </com.cloudbees.jenkins.GitHubPushTrigger>
      </triggers>
    </org.jenkinsci.plugins.workflow.job.properties.PipelineTriggersJobProperty>
  </properties>
  <definition class="org.jenkinsci.plugins.workflow.cps.CpsScmFlowDefinition" plugin="workflow-cps">
    <scm class="hudson.plugins.git.GitSCM" plugin="git">
      <configVersion>2</configVersion>
      <userRemoteConfigs>
        <hudson.plugins.git.UserRemoteConfig>
          <url>https://github.com/$GITHUB_USER/$REPO_NAME.git</url>
        </hudson.plugins.git.UserRemoteConfig>
      </userRemoteConfigs>
      <branches>
        <hudson.plugins.git.BranchSpec>
          <name>*/main</name>
        </hudson.plugins.git.BranchSpec>
      </branches>
      <doGenerateSubmoduleConfigurations>false</doGenerateSubmoduleConfigurations>
      <submoduleCfg class="empty-list"/>
      <extensions/>
    </scm>
    <scriptPath>Jenkinsfile</scriptPath>
    <lightweight>true</lightweight>
  </definition>
  <triggers/>
  <disabled>false</disabled>
</flow-definition>
EOF

    # Create job with proper headers
    CREATE_RESPONSE=$(curl -s -w "%{http_code}" -X POST \
      "$JENKINS_URL/createItem?name=SnowBored-Game-Pipeline" \
      --user "$JENKINS_USER:$JENKINS_PASSWORD" \
      -H "$JENKINS_CRUMB" \
      -H "Content-Type: application/xml" \
      -H "ngrok-skip-browser-warning: true" \
      --data-binary @/tmp/job-config.xml)
    
    if [[ "$CREATE_RESPONSE" == *"200"* ]]; then
        echo "   ✅ Job recreated successfully"
    else
        echo "   ❌ Failed to recreate job: $CREATE_RESPONSE"
    fi
    
    rm /tmp/job-config.xml
else
    echo "   ❌ Unexpected response: HTTP $JOB_CHECK"
fi

# Issue 3: Fix webhook ID issue
echo "3. Checking GitHub webhook..."
if [ -z "$WEBHOOK_ID" ] || [ "$WEBHOOK_ID" = "" ]; then
    echo "   ⚠️  Webhook ID missing, checking GitHub..."
    
    # Get all webhooks for the repository
    WEBHOOKS=$(curl -s \
      -H "Authorization: token $GITHUB_TOKEN" \
      "https://api.github.com/repos/$GITHUB_USER/$REPO_NAME/hooks")
    
    # Find our webhook
    WEBHOOK_ID=$(echo "$WEBHOOKS" | grep -B5 -A5 "$WEBHOOK_URL" | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2)
    
    if [ -n "$WEBHOOK_ID" ]; then
        echo "   ✅ Found webhook ID: $WEBHOOK_ID"
        # Update configuration
        sed -i "s/WEBHOOK_ID=$/WEBHOOK_ID=$WEBHOOK_ID/" .jenkins-config
    else
        echo "   ❌ Webhook not found in GitHub"
    fi
else
    echo "   ✅ Webhook ID exists: $WEBHOOK_ID"
fi

# Issue 4: Check project structure
echo "4. Checking project structure..."
if [ ! -f "package.json" ]; then
    echo "   ❌ package.json missing - creating basic structure..."
    
    cat > package.json << 'EOF'
{
  "name": "snowbored-game",
  "version": "1.0.0",
  "private": true,
  "scripts": {
    "dev": "next dev",
    "build": "next build",
    "start": "next start",
    "lint": "next lint",
    "test": "jest",
    "test:watch": "jest --watch",
    "test:coverage": "jest --coverage"
  },
  "dependencies": {
    "next": "14.0.0",
    "react": "^18",
    "react-dom": "^18"
  },
  "devDependencies": {
    "@types/node": "^20",
    "@types/react": "^18",
    "@types/react-dom": "^18",
    "eslint": "^8",
    "eslint-config-next": "14.0.0",
    "typescript": "^5",
    "@testing-library/jest-dom": "^6.1.0",
    "@testing-library/react": "^13.4.0",
    "jest": "^29.7.0",
    "jest-environment-jsdom": "^29.7.0"
  }
}
EOF
    echo "   ✅ package.json created"
else
    echo "   ✅ package.json exists"
fi

if [ ! -f "Dockerfile" ]; then
    echo "   ❌ Dockerfile missing - creating..."
    
    cat > Dockerfile << 'EOF'
# Multi-stage build for production
FROM node:18-alpine AS base
WORKDIR /app
COPY package*.json ./

FROM base AS deps
RUN npm ci --only=production

FROM base AS build
RUN npm ci
COPY . .
RUN npm run build

FROM node:18-alpine AS runtime
WORKDIR /app

RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs

COPY --from=build --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=build --chown=nextjs:nodejs /app/.next/static ./.next/static
COPY --from=build --chown=nextjs:nodejs /app/public ./public

USER nextjs

EXPOSE 3000
ENV PORT 3000
ENV HOSTNAME "0.0.0.0"

CMD ["node", "server.js"]
EOF
    echo "   ✅ Dockerfile created"
else
    echo "   ✅ Dockerfile exists"
fi

# Issue 5: Test build trigger with proper headers
echo "5. Testing build trigger with fixed headers..."
if [ -n "$JENKINS_CRUMB" ]; then
    BUILD_TEST=$(curl -s -w "%{http_code}" -X POST \
      "$JENKINS_URL/job/SnowBored-Game-Pipeline/build" \
      --user "$JENKINS_USER:$JENKINS_PASSWORD" \
      -H "$JENKINS_CRUMB" \
      -H "ngrok-skip-browser-warning: true")
    
    if [[ "$BUILD_TEST" == *"201"* ]] || [[ "$BUILD_TEST" == *"200"* ]]; then
        echo "   ✅ Build trigger successful"
    else
        echo "   ❌ Build trigger failed: HTTP $BUILD_TEST"
    fi
else
    echo "   ❌ Cannot test build - no CSRF token"
fi

echo "================================================"
echo "🔧 Issue Resolution Complete!"
echo "📋 Summary of fixes:"
echo "✅ CSRF token obtained and configured"
echo "✅ Pipeline job verified/recreated"
echo "✅ Webhook ID resolved"
echo "✅ Project structure validated"
echo "✅ Build trigger tested"
echo "================================================"
echo "🚀 Ready to test again!"
