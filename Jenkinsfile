pipeline {
    agent any
    
    stages {
        stage('Test AWS') {
            steps {
                withAWS(credentials: 'INVOXA_AWS_CREDENTIALS', region: 'us-east-1') {
                    sh 'aws sts get-caller-identity'
                }
            }
        }

        stage('Setup Parameters') {
            steps {
                script {
                    properties([
                        parameters([
                            choice(
                                name: 'Organization_Environment',
                                choices: ['dev','prd'],
                                description: "Select Organization Environment"
                            ),
                            string(
                                name: 'InvoxaAccount',
                                defaultValue: '',
                                description: "Enter 'invoxa-dev' or 'invoxa-prd'"
                            ),
                            string(
                                name: 'InvoxaAccountNo',
                                defaultValue: '',
                                description: "Enter AWS account number (817998750852)"
                            ),
                            string(
                                name: 'RegionName',
                                defaultValue: 'us-east-1',
                                description: "Enter deployment region"
                            ),
                            string(
                                name: 'ITCHG',
                                defaultValue: '',
                                description: "Enter Deployment Ticket Number"
                            )
                        ])
                    ])
                }
            }
        }

        stage('Validate Parameters') {
            steps {
                script {
                    // Your original validation logic
                    if (params.InvoxaAccount == 'invoxa-dev') {
                        assert params.InvoxaAccountNo == '817998750852'
                    } else if (params.InvoxaAccount == 'invoxa-prd') {
                        assert params.InvoxaAccountNo == '817998750852'
                    } else {
                        error "Invalid InvoxaAccount: ${params.InvoxaAccount}"
                    }
                }
            }
        }

        stage('Assume AWS Role') {
            steps {
                script {
                    // Your original role assumption code
                    def assumeRole = { roleArn, sessionName ->
                        def assumeRoleOutput = sh(
                            script: """
                                aws sts assume-role --role-arn ${roleArn} --role-session-name ${sessionName} --output json
                            """,
                            returnStdout: true
                        ).trim()
                        def json = readJSON text: assumeRoleOutput
                        return [
                            accessKey: json.Credentials.AccessKeyId,
                            secretKey: json.Credentials.SecretAccessKey,
                            sessionToken: json.Credentials.SessionToken
                        ]
                    }

                    if (params.InvoxaAccount == 'invoxa-dev') {
                        def role1 = assumeRole("arn:aws:iam::857736875915:role/RINX_DEVAWS_JENKINS_ADM", "jenkins-dev-adm-session")
                        env.AWS_ACCESS_KEY_ID = role1.accessKey
                        env.AWS_SECRET_ACCESS_KEY = role1.secretKey
                        env.AWS_SESSION_TOKEN = role1.sessionToken
                    } else if (params.InvoxaAccount == 'invoxa-prd') {
                        def role1 = assumeRole("arn:aws:iam::857736875915:role/RINX_PRDAWS_JENKINS_ADM", "jenkins-prd-adm-session")
                        env.AWS_ACCESS_KEY_ID = role1.accessKey
                        env.AWS_SECRET_ACCESS_KEY = role1.secretKey
                        env.AWS_SESSION_TOKEN = role1.sessionToken
                    }
                }
            }
        }

        /* YOUR ORIGINAL TERRAFORM STAGES - UNCHANGED */
        stage('Terraform Init') {
            steps {
                script {
                    sh 'terraform init'
                    echo 'Terraform initialized successfully.'
                }
            }
        }

        stage('Terraform Plan') {
            steps {
                script {
                    if (params.Organization_Environment == 'dev') {
                        sh 'terraform plan -var-file=dev.tfvars -out=tfplan'
                    } else if (params.Organization_Environment == 'prd') {
                        sh 'terraform plan -var-file=prod.tfvars -out=tfplan'
                    }
                    echo "Terraform plan executed successfully for ${params.Organization_Environment} environment."
                }
            }
        }

        stage('Terraform Apply') {
            steps {
                script {
                    if (params.Organization_Environment == 'dev') {
                        sh 'terraform apply -var-file=dev.tfvars tfplan'
                    } else if (params.Organization_Environment == 'prd') {
                        sh 'terraform apply -var-file=prod.tfvars tfplan'
                    }
                    echo "Terraform apply executed successfully for ${params.Organization_Environment} environment."
                }
            }
        }

        stage('Cleanup') {
            steps {
                script {
                    sh 'rm -f tfplan'
                    echo 'Terraform state files cleaned up.'
                }
            }
        }

        stage('Post-Deployment') {
            steps {
                script {
                    echo "Deployment completed successfully for ${params.Organization_Environment} environment."
                }
            }
        }
    }

    post {
        always {
            script {
                env.AWS_ACCESS_KEY_ID = ''
                env.AWS_SECRET_ACCESS_KEY = ''
                env.AWS_SESSION_TOKEN = ''
            }
        }
        success {
            echo 'Pipeline completed successfully.'
        }
        failure {
            echo 'Pipeline failed.'
        }
    }
}