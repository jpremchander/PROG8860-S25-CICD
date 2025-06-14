// Jenkins Pipeline Job Configuration Script
// Run this in Jenkins Script Console (Manage Jenkins > Script Console)

import jenkins.model.*
import hudson.model.*
import org.jenkinsci.plugins.workflow.job.WorkflowJob
import org.jenkinsci.plugins.workflow.cps.CpsFlowDefinition
import hudson.triggers.SCMTrigger
import org.jenkinsci.plugins.github.GitHubPushTrigger

// Job configuration
def jobName = "SnowBored-Game-Pipeline"
def gitUrl = "https://github.com/YOUR_USERNAME/snowbored-game.git"
def branch = "*/main"

// Create new pipeline job
def job = Jenkins.instance.createProject(WorkflowJob, jobName)

// Configure Git repository
def scm = new hudson.plugins.git.GitSCM([
    new hudson.plugins.git.UserRemoteConfig(gitUrl, null, null, null)
])
scm.branches = [new hudson.plugins.git.BranchSpec(branch)]
job.definition = new CpsFlowDefinition("""
pipeline {
    agent any
    
    environment {
        DOCKER_IMAGE = 'snowbored-game'
        DOCKER_TAG = "\${BUILD_NUMBER}"
        NODE_VERSION = '18'
    }
    
    tools {
        nodejs 'NodeJS-18'
    }
    
    stages {
        stage('Checkout') {
            steps {
                echo 'Checking out source code...'
                checkout scm
                
                script {
                    env.GIT_COMMIT_SHORT = sh(
                        script: "git rev-parse --short HEAD",
                        returnStdout: true
                    ).trim()
                }
            }
        }
        
        stage('Build') {
            steps {
                echo 'Installing dependencies and building application...'
                sh '''
                    echo "Node version: \$(node --version)"
                    echo "NPM version: \$(npm --version)"
                    
                    # Clean install dependencies
                    npm ci
                    
                    # Run linting
                    npm run lint
                    
                    # Build the application
                    npm run build
                '''
            }
            post {
                success {
                    echo 'Build completed successfully!'
                }
                failure {
                    echo 'Build failed!'
                }
            }
        }
        
        stage('Test') {
            parallel {
                stage('Unit Tests') {
                    steps {
                        echo 'Running unit tests...'
                        sh '''
                            # Run tests with coverage
                            npm test -- --coverage --watchAll=false --ci
                        '''
                    }
                    post {
                        always {
                            // Publish test results
                            publishHTML([
                                allowMissing: false,
                                alwaysLinkToLastBuild: true,
                                keepAll: true,
                                reportDir: 'coverage/lcov-report',
                                reportFiles: 'index.html',
                                reportName: 'Coverage Report'
                            ])
                        }
                        success {
                            echo 'All tests passed!'
                        }
                        failure {
                            echo 'Tests failed! Pipeline will be aborted.'
                            error('Test stage failed')
                        }
                    }
                }
                
                stage('Security Scan') {
                    steps {
                        echo 'Running security audit...'
                        sh '''
                            # Run npm audit
                            npm audit --audit-level moderate || true
                        '''
                    }
                }
            }
        }
        
        stage('Docker Build') {
            steps {
                echo 'Building Docker image...'
                script {
                    def dockerImage = docker.build("\${DOCKER_IMAGE}:\${DOCKER_TAG}")
                    
                    // Tag with latest
                    sh "docker tag \${DOCKER_IMAGE}:\${DOCKER_TAG} \${DOCKER_IMAGE}:latest"
                    
                    // Tag with git commit
                    sh "docker tag \${DOCKER_IMAGE}:\${DOCKER_TAG} \${DOCKER_IMAGE}:\${env.GIT_COMMIT_SHORT}"
                }
            }
            post {
                success {
                    echo 'Docker image built successfully!'
                }
                failure {
                    echo 'Docker build failed!'
                }
            }
        }
        
        stage('Deploy to Localhost') {
            steps {
                echo 'Deploying to localhost...'
                sh '''
                    # Stop existing container if running
                    docker stop snowbored-game || true
                    docker rm snowbored-game || true
                    
                    # Run new container
                    docker run -d \\
                        --name snowbored-game \\
                        -p 3000:3000 \\
                        --restart unless-stopped \\
                        \${DOCKER_IMAGE}:\${DOCKER_TAG}
                    
                    # Health check
                    sleep 10
                    curl -f http://localhost:3000 || exit 1
                '''
            }
            post {
                success {
                    echo '🚀 Deployment to localhost successful!'
                    echo '🌐 Application available at: http://localhost:3000'
                }
                failure {
                    echo '❌ Deployment failed!'
                }
            }
        }
        
        stage('Archive Artifacts') {
            steps {
                echo 'Archiving build artifacts...'
                
                sh '''
                    mkdir -p artifacts
                    
                    # Copy build files
                    cp -r .next artifacts/ || true
                    cp -r public artifacts/ || true
                    cp package.json artifacts/
                    cp Dockerfile artifacts/
                    
                    # Create tarball
                    tar -czf snowbored-artifacts-\${BUILD_NUMBER}.tar.gz artifacts/
                '''
                
                // Archive artifacts
                archiveArtifacts artifacts: 'snowbored-artifacts-*.tar.gz', fingerprint: true
                
                // Archive test reports
                archiveArtifacts artifacts: 'coverage/**/*', allowEmptyArchive: true
            }
        }
    }
    
    post {
        always {
            echo 'Pipeline completed!'
            
            // Clean up Docker images to save space
            sh '''
                docker image prune -f
            '''
            
            // Clean workspace
            cleanWs()
        }
        
        success {
            echo '✅ Pipeline executed successfully!'
            echo '🌐 Application deployed at: http://localhost:3000'
        }
        
        failure {
            echo '❌ Pipeline failed!'
        }
        
        unstable {
            echo '⚠️ Pipeline completed with warnings!'
        }
    }
}
""", true)

// Configure triggers
job.addTrigger(new GitHubPushTrigger())

// Save the job
job.save()

println "✅ Pipeline job '${jobName}' created successfully!"
println "🌐 Access at: http://localhost:8080/job/${jobName}/"
