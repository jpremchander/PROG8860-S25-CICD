#!/bin/bash

echo "🏗️ Creating SnowBored Game Pipeline Job"
echo "======================================"

# Load configuration
if [ -f ".jenkins-config" ]; then
    source .jenkins-config
    echo "✅ Loaded Jenkins configuration"
else
    echo "❌ No configuration found. Run configure-with-ngrok.sh first"
    exit 1
fi

JOB_NAME="SnowBored-Game-Pipeline"

# Create job configuration XML
cat > /tmp/snowbored-job.xml << EOF
<?xml version='1.1' encoding='UTF-8'?>
<flow-definition plugin="workflow-job">
  <actions/>
  <description>SnowBored Game CI/CD Pipeline - Automated build, test, and deploy</description>
  <keepDependencies>false</keepDependencies>
  <properties>
    <org.jenkinsci.plugins.workflow.job.properties.PipelineTriggersJobProperty>
      <triggers>
        <com.cloudbees.jenkins.GitHubPushTrigger plugin="github">
          <spec></spec>
        </com.cloudbees.jenkins.GitHubPushTrigger>
      </triggers>
    </org.jenkinsci.plugins.workflow.job.properties.PipelineTriggersJobProperty>
    <hudson.model.ParametersDefinitionProperty>
      <parameterDefinitions>
        <hudson.model.StringParameterDefinition>
          <name>BRANCH</name>
          <description>Git branch to build</description>
          <defaultValue>main</defaultValue>
          <trim>false</trim>
        </hudson.model.StringParameterDefinition>
        <hudson.model.BooleanParameterDefinition>
          <name>DEPLOY</name>
          <description>Deploy to localhost after successful build</description>
          <defaultValue>true</defaultValue>
        </hudson.model.BooleanParameterDefinition>
      </parameterDefinitions>
    </hudson.model.ParametersDefinitionProperty>
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
          <name>*/\${BRANCH}</name>
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

# Get Jenkins crumb with proper format
echo "🔐 Getting CSRF token..."
CRUMB_RESPONSE=$(curl -s "$JENKINS_URL/crumbIssuer/api/json" \
  --user "$JENKINS_USER:$JENKINS_PASSWORD" \
  -H "ngrok-skip-browser-warning: true")

if [ -z "$CRUMB_RESPONSE" ]; then
    echo "❌ Failed to get crumb response"
    exit 1
fi

# Extract crumb field and value properly
CRUMB_FIELD=$(echo "$CRUMB_RESPONSE" | grep -o '"crumbRequestField":"[^"]*"' | cut -d'"' -f4)
CRUMB_VALUE=$(echo "$CRUMB_RESPONSE" | grep -o '"crumb":"[^"]*"' | cut -d'"' -f4)

if [ -z "$CRUMB_FIELD" ] || [ -z "$CRUMB_VALUE" ]; then
    echo "❌ Failed to extract CSRF token"
    exit 1
fi

echo "✅ CSRF token obtained: $CRUMB_FIELD"

# Check if job already exists
echo "🔍 Checking if job already exists..."
JOB_CHECK=$(curl -s -w "%{http_code}" -o /dev/null \
  --user "$JENKINS_USER:$JENKINS_PASSWORD" \
  -H "ngrok-skip-browser-warning: true" \
  "$JENKINS_URL/job/$JOB_NAME/api/json")

if [ "$JOB_CHECK" = "200" ]; then
    echo "⚠️  Job already exists. Updating configuration..."
    
    # Update existing job
    UPDATE_RESPONSE=$(curl -s -w "%{http_code}" -X POST \
      "$JENKINS_URL/job/$JOB_NAME/config.xml" \
      --user "$JENKINS_USER:$JENKINS_PASSWORD" \
      -H "$JENKINS_CRUMB" \
      -H "Content-Type: application/xml" \
      -H "ngrok-skip-browser-warning: true" \
      --data-binary @/tmp/snowbored-job.xml)
    
    if [[ "$UPDATE_RESPONSE" == *"200"* ]]; then
        echo "✅ Pipeline job updated successfully!"
    else
        echo "⚠️  Update response: $UPDATE_RESPONSE"
    fi
else
    # Create new job with proper CSRF header
    echo "🏗️ Creating new pipeline job '$JOB_NAME'..."
    CREATE_RESPONSE=$(curl -s -w "%{http_code}" -X POST \
      "$JENKINS_URL/createItem?name=$JOB_NAME" \
      --user "$JENKINS_USER:$JENKINS_PASSWORD" \
      -H "$CRUMB_FIELD: $CRUMB_VALUE" \
      -H "Content-Type: application/xml" \
      -H "ngrok-skip-browser-warning: true" \
      --data-binary @/tmp/snowbored-job.xml)

    if [[ "$CREATE_RESPONSE" == *"200"* ]]; then
        echo "✅ Pipeline job created successfully!"
    else
        echo "⚠️  Create response: $CREATE_RESPONSE"
    fi
fi

# Clean up
rm /tmp/snowbored-job.xml

echo "================================================"
echo "✅ Pipeline Job Setup Complete!"
echo "🌐 Jenkins Dashboard: $JENKINS_URL"
echo "🏗️ Pipeline Job: $JENKINS_URL/job/$JOB_NAME/"
echo "📊 Blue Ocean: $JENKINS_URL/blue/organizations/jenkins/$JOB_NAME/"
echo "================================================"
echo "🚀 Ready to trigger builds!"
