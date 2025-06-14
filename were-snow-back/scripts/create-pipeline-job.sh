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

# Get Jenkins crumb
JENKINS_CRUMB=$(curl -s "$JENKINS_URL/crumbIssuer/api/xml?xpath=concat(//crumbRequestField,\":\",//crumb)" \
  --user "$JENKINS_USER:$JENKINS_PASSWORD" \
  -H "ngrok-skip-browser-warning: true")

# Create the job
echo "🏗️ Creating pipeline job '$JOB_NAME'..."
CREATE_RESPONSE=$(curl -s -X POST "$JENKINS_URL/createItem?name=$JOB_NAME" \
  --user "$JENKINS_USER:$JENKINS_PASSWORD" \
  -H "$JENKINS_CRUMB" \
  -H "Content-Type: application/xml" \
  -H "ngrok-skip-browser-warning: true" \
  --data-binary @/tmp/snowbored-job.xml)

if [ $? -eq 0 ]; then
    echo "✅ Pipeline job created successfully!"
    echo "🌐 Job URL: $JENKINS_URL/job/$JOB_NAME/"
    
    # Trigger initial build
    echo "🚀 Triggering initial build..."
    sleep 2
    curl -s -X POST "$JENKINS_URL/job/$JOB_NAME/build" \
      --user "$JENKINS_USER:$JENKINS_PASSWORD" \
      -H "$JENKINS_CRUMB" \
      -H "ngrok-skip-browser-warning: true"
    
    echo "✅ Initial build triggered!"
    echo "📊 Monitor build: $JENKINS_URL/job/$JOB_NAME/"
else
    echo "❌ Failed to create job"
    echo "Response: $CREATE_RESPONSE"
fi

# Clean up
rm /tmp/snowbored-job.xml

echo "================================================"
echo "✅ Pipeline Job Setup Complete!"
echo "🌐 Jenkins Dashboard: $JENKINS_URL"
echo "🏗️ Pipeline Job: $JENKINS_URL/job/$JOB_NAME/"
echo "📊 Blue Ocean: $JENKINS_URL/blue/organizations/jenkins/$JOB_NAME/"
echo "================================================"
