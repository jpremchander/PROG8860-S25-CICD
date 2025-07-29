pipeline {
    agent any

    environment {
        // Auto-set based on branch
        TF_ENV = "${env.BRANCH_NAME == 'release' ? 'prd' : 'dev'}"
        TF_VARS_FILE = "${TF_ENV}.tfvars"
        TF_DIR = "terraform"
    }
    
    stages {
        stage('Setup Parameters') {
            steps {
                script {
                    // Validate environment selection before setting parameters
                    if (!params.Organization_Environment in ['dev', 'prd']) {
                        error "Invalid environment selected: ${params.Organization_Environment}"
                    }

                    properties([
                        parameters([
                            choice(
                                name: 'Organization_Environment',
                                choices: ['dev', 'prd'],
                                description: 'Select deployment environment (dev/prd)'
                            ),
                            string(
                                name: 'InvoxaAccount',
                                defaultValue: params.Organization_Environment == 'dev' ? 'invoxa-dev' : 'invoxa-prd',
                                description: 'AWS Account'
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
                                name: 'JIRA_TICKET',
                                defaultValue: '',
                                description: 'REQUIRED: Enter your change ticket number'
                            )
                        ])
                    ])
                }
            }
        }

        stage('Validate Inputs') {
            steps {
                script {                   
                    // Cross-validate parameters
                    if (params.Organization_Environment == 'dev' && params.InvoxaAccount != 'invoxa-dev') {
                        error "Invalid account for dev environment"
                    }
                    if (params.Organization_Environment == 'prd' && params.InvoxaAccount != 'invoxa-prd') {
                        error "Invalid account for prod environment"
                    }
                }
            }
        }

        stage('Assume AWS Role') {
            steps {
                script {
                    if (params.Organization_Environment == 'dev') {
                        env.AWS_ROLE_ARN = "arn:aws:iam::${params.InvoxaAccountNo}:role/RINX_DEVAWS_JENKINS_ADM"
                    sh '''
                        aws sts assume-role \
                        --role-arn ${env.AWS_ROLE_ARN} \
                        --role-session-name jenkins-${params.Organization_Environment}-${BUILD_NUMBER} \
                        --output json
                    '''
                        def creds = readJSON text: sh(script: 'aws sts assume-role --role-arn ${env.AWS_ROLE_ARN} --role-session-name jenkins-${params.Organization_Environment}-${BUILD_NUMBER} --output json', returnStdout: true)
                        env.AWS_ACCESS_KEY_ID = creds.Credentials.AccessKeyId
                        env.AWS_SECRET_ACCESS_KEY = creds.Credentials.SecretAccessKey
                        env.AWS_SESSION_TOKEN = creds.Credentials.SessionToken
                    } else {
                        env.AWS_ROLE_ARN = "arn:aws:iam::${params.InvoxaAccountNo}:role/RINX_PRDAWS_JENKINS_ADM"
                    sh '''
                        aws sts assume-role \ 
                        --role-arn ${env.AWS_ROLE_ARN} \
                        --role-session-name jenkins-${params.Organization_Environment}-${BUILD_NUMBER} \    
                        --output json
                    '''
                        def creds = readJSON text: sh(script: 'aws sts assume-role --role-arn ${env.AWS_ROLE_ARN} --role-session-name jenkins-${params.Organization_Environment}-${BUILD_NUMBER} --output json', returnStdout: true)
                        env.AWS_ACCESS_KEY_ID = creds.Credentials.AccessKeyId
                        env.AWS_SECRET_ACCESS_KEY = creds.Credentials.SecretAccessKey
                        env.AWS_SESSION_TOKEN = creds.Credentials.SessionToken
                    }
                }
            }
        }

        stage('Check Terraform Changes') {
            steps {
                script {
                    // Compare against origin/main to catch all unreviewed changes
                    env.TF_CHANGES = sh(
                        script: """
                        git fetch origin main
                        changed=\$(git diff --name-only HEAD origin/main -- ${TF_DIR}/)
                        echo "\$changed" | grep -q ".tf\$" && echo "true" || echo "false"
                        """,
                        returnStdout: true
                    ).trim()

                    echo "Terraform changes detected: ${env.TF_CHANGES}"
                    if (env.TF_CHANGES == 'true') {
                        env.TF_CHANGED_FILES = sh(
                            script: "git diff --name-only HEAD origin/main -- ${TF_DIR}/ | grep '.tf\$'",
                            returnStdout: true
                        ).trim()
                        echo "Changed files:\n${env.TF_CHANGED_FILES}"
                    }
                }
            }
        }

        stage('Terraform Init/Plan') {
            when {
                expression { env.TF_CHANGES == 'true' }
            }
            steps {
                dir(TF_DIR) {
                    sh """
                    terraform init -backend-config=backend-${TF_ENV}.conf
                    terraform plan -var-file=${TF_VARS_FILE} -out=tfplan
                    """
                    archiveArtifacts artifacts: 'tfplan'
                }
            }
        }

        stage('Terraform Apply (Conditional)') {
            when {
                expression { env.TF_CHANGES == 'true' }
            }
            steps {
                dir(TF_DIR) {
                    script {
                        if (TF_ENV == 'prd') {
                            timeout(time: 15, unit: 'MINUTES') {
                                input(message: "Approve PROD deployment?", ok: "Deploy")
                            }
                        }
                        sh "terraform apply ${TF_ENV == 'dev' ? '-auto-approve' : ''} tfplan"
                    }
                }
            }
        }

        stage('Build & Deploy App') {
            steps {
                script {
                    echo "Deploying app to ${params.Organization_Environment} environment"
                    
                    // Get outputs from Terraform if available
                    def clusterName = ""
                    def serviceName = ""
                    if (env.TF_CHANGES == 'true') {
                        dir(TF_DIR) {
                            clusterName = sh(
                                script: 'terraform output -raw ecs_cluster_name',
                                returnStdout: true
                            ).trim()
                            serviceName = sh(
                                script: 'terraform output -raw ecs_service_name',
                                returnStdout: true
                            ).trim()
                        }
                    }
                    
                    // Fallback to naming convention if no Terraform outputs
                    clusterName = clusterName ?: "${params.Organization_Environment}-cluster"
                    serviceName = serviceName ?: "${params.Organization_Environment}-service"
                    
                    sh """
                    aws ecs update-service \
                        --cluster ${clusterName} \
                        --service ${serviceName} \
                        --region ${params.RegionName} \
                        --force-new-deployment
                    """
                }
            }
        }
    }

    post {
        always {
            script {
                // Cleanup credentials
                env.AWS_ACCESS_KEY_ID = ''
                env.AWS_SECRET_ACCESS_KEY = ''
                env.AWS_SESSION_TOKEN = ''
                
                // Notify deployment result
                echo "Deployment to ${params.Organization_environment} completed with status: ${currentBuild.currentResult}"
                cleanWs()
            }
        }
    }
}