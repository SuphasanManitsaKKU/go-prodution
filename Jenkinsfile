pipeline {
    agent any

    environment {
        DOCKER_BUILDKIT = '1'
    }

    stages {
        stage('Check Build Info') {
            steps {
                sh '''
                    echo "🔧 BUILD_NUMBER = $BUILD_NUMBER"
                    echo "🔧 JOB_NAME = $JOB_NAME"
                    echo "🔧 BUILD_ID = $BUILD_ID"
                '''
            }
        }

        stage('Load Configs from Vault') {
            steps {
                script {
                    withVault(vaultSecrets: [[
                        path: 'pipeline/env',
                        secretValues: [
                            [envVar: 'DOCKER_PASSWORD', vaultKey: 'DOCKER_PASSWORD'],
                            [envVar: 'DOCKER_USERNAME', vaultKey: 'DOCKER_USERNAME'],
                            [envVar: 'REGISTRY', vaultKey: 'REGISTRY'],
                            [envVar: 'REGISTRY_PROJECT_NAME', vaultKey: 'REGISTRY_PROJECT_NAME'],
                            [envVar: 'IMAGE', vaultKey: 'IMAGE'],
                            [envVar: 'IMAGE_OUTPUT_PORT', vaultKey: 'IMAGE_OUTPUT_PORT'],
                            [envVar: 'TARGET_USER', vaultKey: 'TARGET_USER'],
                            [envVar: 'TARGET_IP', vaultKey: 'TARGET_IP']
                        ]
                    ]]) {
                        // assign เข้า env pipeline
                        env.DOCKER_PASSWORD = env.DOCKER_PASSWORD
                        env.DOCKER_USERNAME = env.DOCKER_USERNAME
                        env.REGISTRY = env.REGISTRY
                        env.REGISTRY_PROJECT_NAME = env.REGISTRY_PROJECT_NAME
                        env.IMAGE = env.IMAGE
                        env.IMAGE_OUTPUT_PORT = env.IMAGE_OUTPUT_PORT
                        env.TARGET_USER = env.TARGET_USER
                        env.TARGET_IP = env.TARGET_IP
                        env.FULL_IMAGE = "${env.REGISTRY}/${env.REGISTRY_PROJECT_NAME}/${env.IMAGE}:"
                    }
                }
            }
        }

        stage('Setup Image Tag') {
            steps {
                script {
                    env.TAG = env.BUILD_NUMBER?.toString() ?: 'latest'
                    env.FULL_IMAGE = "${env.REGISTRY}/${env.REGISTRY_PROJECT_NAME}/${env.IMAGE}:${env.TAG}"
                    echo "✅ FULL_IMAGE: ${env.FULL_IMAGE}"
                }
            }
        }

        stage('Build Docker Image') {
            steps {
                sh 'docker build -t $FULL_IMAGE .'
            }
        }

        stage('Login to Registry') {
            steps {
                sh '''
                    echo "🧪 Logging in to Docker Registry"
                    echo "Username: $DOCKER_USERNAME"
                    echo "Registry: $REGISTRY"
                    echo "Image: $FULL_IMAGE"

                    echo $DOCKER_PASSWORD | docker login $REGISTRY -u $DOCKER_USERNAME --password-stdin
                '''
            }
        }

        stage('Push Docker Image') {
            steps {
                sh 'docker push $FULL_IMAGE'
            }
        }

        stage('Clean Up Local Image') {
            steps {
                sh '''
                    echo "🧹 Cleaning up local image: $FULL_IMAGE"
                    docker rmi $FULL_IMAGE || echo "⚠️ Failed to remove image"
                '''
            }
        }

        stage('Trigger Ansible CD') {
            steps {
                sh """
                    echo "[INFO] Running Ansible deployment..."

                    mkdir -p ~/.ssh
                    cp /var/jenkins_home/.ssh/id_rsa ~/.ssh/id_rsa
                    chmod 600 ~/.ssh/id_rsa
                    ssh-keyscan -H $TARGET_IP >> ~/.ssh/known_hosts

                    ansible-playbook -i /var/jenkins_home/ansible/inventory.ini \
                        /var/jenkins_home/ansible/playbooks/deploy_app.yml \
                        --extra-vars \\
                        "registry=${REGISTRY} \\
                         registry_project_name=${REGISTRY_PROJECT_NAME} \\
                         image_name=${IMAGE} \\
                         tag=${TAG} \\
                         image_port=${IMAGE_OUTPUT_PORT} \\
                         host_ip=${TARGET_IP} \\
                         ansible_user=${TARGET_USER}"
                """
            }
        }
    }
}
