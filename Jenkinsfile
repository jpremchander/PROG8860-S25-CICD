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
                   String sectionHeaderStyle = '''
                   color: white;
                   background: dimgrey;
                   font-family: Roboto, sans-serif !important;
                   padding: 5px;
                   text-align: center;
                   '''
                   String separatorStyle = '''
                       border: 0;
                       border-bottom: 3px;
                       background: #999;
                   '''
                properties([
                               $class: 'ParameterSeparator',
                               name: 'Tag_HEADER',
                               sectionHeader: 'INVOXA',
                               separatorStyle: separatorStyle,
                               sectionHeaderStyle: sectionHeaderStyle
                           ],
                           choice(
                               choices: ['dev','prd'],
                               name: 'Organization_Environment',
                               description: "Select Organization Environment"
                           ),
                           [
                               $class: 'CascadeChoiceParameter',
                               choiceType: 'PT_SINGLE_SELECT',
                               description: 'Invoxa Account where application will be deployed',
                               filterLength: 1, filterable: false,
                               name: 'InvoxaAccount',
                               randomName: 'choice-parameter-01',
                               referencedParameters: 'Organization_Environment',
                               script:
                               [
                                   $class: 'GroovyScript',
                                   fallbackScript:
                                   [
                                       classpath: [],
                                       sandbox: true,
                                       script: 'return[""]'],
                                   script: [
                                       classpath: [],
                                       sandbox: true,
                                       script: '''
                                           if(Organization_Environment.equals('dev'))
                                           {
                                           return['invoxa-dev']
                                           }
                                           else if(Organization_Environment.equals('prd'))
                                           {
                                           return['invoxa-prd']
                                           }
                                           else
                                           {
                                           return['NA']
                                           }    
                                        '''
                                           ]
                               ]
                           ],
                           [
                               $class: 'CascadeChoiceParameter',
                               choiceType: 'PT_SINGLE_SELECT',
                               description: 'Invoxa Account Number where application will be deployed',
                               filterLength: 1, filterable: false,
                               name: 'InvoxaAccountNo',
                               randomName: 'choice-parameter-02',
                               referencedParameters: 'InvoxaAccount',
                               script: [
                                   $class: 'GroovyScript',
                                   fallbackScript: [
                                       classpath: [],
                                       sandbox: true,
                                       script: 'return[""]'
                                   ],
                                   script: [
                                       classpath: [],
                                       sandbox: true,
                                       script: '''
                                           if(InvoxaAccount.equals('invoxa-dev')) {
                                               return['817998750852']
                                           }
                                           else if(InvoxaAccount.equals('invoxa-prd')) {
                                               return['817998750852']
                                           }
                                           else {
                                               return['NA']
                                           }
                                       '''
                                   ]
                               ]
                           ],
                           [
                               $class: 'CascadeChoiceParameter',
                               choiceType: 'PT_SINGLE_SELECT',
                               description: 'Select Region where IAM Role will be Created',
                               filterLength: 1, filterable: false,
                               name: 'RegionName',
                               randomName: 'choice-parameter-03',
                               referencedParameters: 'InvoxaAccount',
                               script: [
                                   $class: 'GroovyScript',
                                   fallbackScript: [
                                       classpath: [],
                                       sandbox: true,
                                       script: 'return[""]'
                                   ],
                                   script: [
                                       classpath: [],
                                       sandbox: true,
                                       script: '''
                                           if(InvoxaAccount.equals('invoxa-dev') || InvoxaAccount.equals('invoxa-prd')) {
                                               return['us-east-1']
                                           }
                                           else {
                                               return['NA']
                                           }    
                                       '''
                                   ]
                               ]
                           ],
                           [
                               $class: 'ParameterSeparator',
                               name: 'Tag_HEADER',
                               sectionHeader: 'Tag Details',
                               separatorStyle: separatorStyle,
                               sectionHeaderStyle: sectionHeaderStyle
                           ],
                            string(
                               name: 'JIRA_Ticket_Number',
                               defaultValue: '',
                               description: 'Enter Deployment Ticket Number'
                           )
                   ])
               }
           }
       }

        stage('Assume AWS Role') {
           steps {
               script {
                   // role assume
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

                   // Assuming IAM Roles based On Organization Environment
                   if (params.InvoxaAccount == 'invoxa-dev') {
                       def devrole = assumeRole("arn:aws:iam::857736875915:role/RINX_DEVAWS_JENKINS_ADM", "jenkins-dev-adm-session")
                       env.AWS_ACCESS_KEY_ID = devrole.accessKey
                       env.AWS_SECRET_ACCESS_KEY = devrole.secretKey
                       env.AWS_SESSION_TOKEN = devrole.sessionToken
                       //CTRGB ASSUME ROLE
                     } else if (params.InvoxaAccount == 'invoxa-prd') {
                       def prdrole = assumeRole("arn:aws:iam::857736875915:role/RINX_PRDAWS_JENKINS_ADM", "jenkins-prd-adm-session")
                       env.AWS_ACCESS_KEY_ID = prdrole.accessKey
                       env.AWS_SECRET_ACCESS_KEY = prdrole.secretKey
                       env.AWS_SESSION_TOKEN = prdrole.sessionToken
                   } else {
                       error "Unsupported Organization Environment: ${params.Organization_Environment}"
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