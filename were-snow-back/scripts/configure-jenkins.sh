#!/bin/bash

echo "⚙️ Configuring Jenkins for SnowBored Game Pipeline"
echo "================================================="

JENKINS_URL="http://localhost:8080"
JENKINS_USER="admin"

echo "Please enter your Jenkins admin password:"
read -s JENKINS_PASSWORD

# Create Jenkins CLI alias
alias jenkins-cli="java -jar jenkins-cli.jar -s $JENKINS_URL -auth $JENKINS_USER:$JENKINS_PASSWORD"

# Download Jenkins CLI
echo "📥 Downloading Jenkins CLI..."
wget -q $JENKINS_URL/jnlpJars/jenkins-cli.jar

# Install required plugins
echo "🔌 Installing required plugins..."
PLUGINS=(
    "git"
    "github"
    "pipeline-stage-view"
    "docker-workflow"
    "nodejs"
    "html-publisher"
    "email-ext"
    "slack"
    "build-timeout"
    "timestamper"
    "ws-cleanup"
    "github-branch-source"
    "pipeline-github-lib"
    "blue-ocean"
)

for plugin in "${PLUGINS[@]}"; do
    echo "Installing $plugin..."
    jenkins-cli install-plugin $plugin
done

echo "🔄 Restarting Jenkins to activate plugins..."
jenkins-cli restart

echo "⏳ Waiting for Jenkins to restart..."
sleep 30

# Configure Global Tools
echo "🛠️ Configuring Global Tools..."

# NodeJS Configuration
cat > nodejs-config.xml << 'EOF'
<?xml version='1.1' encoding='UTF-8'?>
<jenkins.plugins.nodejs.tools.NodeJSInstallation_-DescriptorImpl>
  <installations>
    <jenkins.plugins.nodejs.tools.NodeJSInstallation>
      <name>NodeJS-18</name>
      <home></home>
      <properties>
        <jenkins.plugins.nodejs.tools.NodeJSInstallation_-InstallSourceProperty>
          <installers>
            <jenkins.plugins.nodejs.tools.NodeJSInstaller>
              <id>18.19.0</id>
              <npmPackages>npm@latest</npmPackages>
            </jenkins.plugins.nodejs.tools.NodeJSInstaller>
          </installers>
        </jenkins.plugins.nodejs.tools.NodeJSInstallation_-InstallSourceProperty>
      </properties>
    </jenkins.plugins.nodejs.tools.NodeJSInstallation>
  </installations>
</jenkins.plugins.nodejs.tools.NodeJSInstallation_-DescriptorImpl>
EOF

# Docker Configuration
cat > docker-config.xml << 'EOF'
<?xml version='1.1' encoding='UTF-8'?>
<org.jenkinsci.plugins.docker.commons.tools.DockerTool_-DescriptorImpl>
  <installations>
    <org.jenkinsci.plugins.docker.commons.tools.DockerTool>
      <name>Docker</name>
      <home>/usr/bin/docker</home>
      <properties/>
    </org.jenkinsci.plugins.docker.commons.tools.DockerTool>
  </installations>
</org.jenkinsci.plugins.docker.commons.tools.DockerTool_-DescriptorImpl>
EOF

echo "✅ Jenkins configuration completed!"
echo "🌐 Access Jenkins at: $JENKINS_URL"
echo "📋 Next: Set up GitHub webhook and create pipeline job"
