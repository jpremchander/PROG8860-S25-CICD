pipeline {
    agent any
    
    stages {
        stage('Setup Parameters') {
            steps {
                script {
                    properties([
                        parameters([
                            choice(
                                name: 'Organization_Environment',
                                choices: ['dev', 'prd'],
                                description: 'Select deployment environment (dev/prd)'
                            ),
                            string(
                                name: 'InvoxaAccount',
                                defaultValue: '',
                                description: 'AWS Account (auto-populates to invoxa-dev/invoxa-prd)'
                            ),
                            string(
                                name: 'InvoxaAccountNo',
                                defaultValue: '817998750852',
                                description: 'AWS Account Number'
                            ),
                            string(
                                name: 'RegionName',
                                defaultValue: 'us-east-1',
                                description: 'AWS Region'
                            ),
                            string(
                                name: 'ITCHG',
                                defaultValue: '',
                                description: 'REQUIRED: Enter your change ticket number'
                            )
                        ])
                    ])
                }
            }
        }

        stage('Auto-configure') {
            steps {
                script {
                    // Auto-set values based on environment selection
                    params.InvoxaAccount = "invoxa-${params.Organization_Environment}"
                    echo "✅ Auto-configured for ${params.InvoxaAccount}"
                }
            }
        }

        stage('Validate Inputs') {
            steps {
                script {
                    if (!params.ITCHG?.trim()) {
                        error "❌ Ticket number (ITCHG) is required"
                    }
                }
            }
        }

        stage('Assume AWS Role') {
            steps {
                script {
                    def roleArn = params.Organization_Environment == 'dev' ? 
                        "arn:aws:iam::857736875915:role/RINX_DEVAWS_JENKINS_ADM" :
                        "arn:aws:iam::857736875915:role/RINX_PRDAWS_JENKINS_ADM"
                    
                    def creds = sh(returnStdout: true, script: """
                        aws sts assume-role \
                        --role-arn ${roleArn} \
                        --role-session-name jenkins-${params.Organization_Environment}-${BUILD_NUMBER} \
                        --output json
                    """).trim()
                    
                    def json = readJSON text: creds
                    env.AWS_ACCESS_KEY_ID = json.Credentials.AccessKeyId
                    env.AWS_SECRET_ACCESS_KEY = json.Credentials.SecretAccessKey
                    env.AWS_SESSION_TOKEN = json.Credentials.SessionToken
                }
            }
        }

        stage('Terraform Init') {
            steps {
                script {
                    sh 'terraform init'
                }
            }
        }

        stage('Terraform Plan') {
            steps {
                script {
                    sh "terraform plan -var-file=${params.Organization_Environment}.tfvars -out=tfplan"
                }
            }
        }

        stage('Terraform Apply') {
            steps {
                script {
                    sh 'terraform apply -auto-approve tfplan'
                    echo "✅ Successfully deployed to ${params.InvoxaAccount}"
                }
            }
        }
    }

    post {
        always {
            script {
                // Cleanup credentials only
                env.AWS_ACCESS_KEY_ID = ''
                env.AWS_SECRET_ACCESS_KEY = ''
                env.AWS_SESSION_TOKEN = ''
            }
        }
    }
}