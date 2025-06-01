pipeline {
  agent any

  environment {
    REGISTRY               = 'harbor.suphasan.site'
    REGISTRY_PROJECT_NAME = 'ok'

    IMAGE      = 'my-app'
    TAG        = "${env.BUILD_NUMBER}"
    FULL_IMAGE = "${REGISTRY}/${REGISTRY_PROJECT_NAME}/${IMAGE}:${TAG}"
    IMAGE_PORT = '5000'

    TARGET_USER = 'ubuntu'
    TARGET_IP = credentials('production-host-ip') // ต้องเป็นแค่ IP

    // Jenkins credential ID แบบ username/password
    DOCKER_CREDS = credentials('harbor-credentials')
  }

  stages {

    stage('Load Secrets from Vault') {
      steps {
        script {
          withVault(
            vaultSecrets: [
              [
                path: 'secret/ok',
                secretValues: [
                  [envVar: 'GOGO', vaultKey: 'gogo'],
                  [envVar: 'YOYO', vaultKey: 'yoyo']
                ]
              ]
            ]
          ) {
            // ทดสอบ echo ค่า Vault
            sh 'echo Loaded from Vault: GOGO=$GOGO, YOYO=$YOYO'
          }
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
        sh 'echo $DOCKER_CREDS_PSW | docker login $REGISTRY -u $DOCKER_CREDS_USR --password-stdin'
      }
    }

    stage('Push Docker Image') {
      steps {
        sh 'docker push $FULL_IMAGE'
      }
    }

    stage('Add Jenkins SSH Public Key to Target') {
      steps {
        sh """
          echo "[INFO] Copying Jenkins public key to target machine..."
          PUB_KEY_PATH=/var/jenkins_home/.ssh/id_rsa.pub

          ssh -o StrictHostKeyChecking=no ${TARGET_USER}@${TARGET_IP} "mkdir -p ~/.ssh && touch ~/.ssh/authorized_keys && chmod 600 ~/.ssh/authorized_keys"
          cat \$PUB_KEY_PATH | ssh -o StrictHostKeyChecking=no ${TARGET_USER}@${TARGET_IP} "grep -qxF '\$(cat)' ~/.ssh/authorized_keys || echo '\$(cat)' >> ~/.ssh/authorized_keys"
        """
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
            --extra-vars "registry=${REGISTRY} registry_project_name=${REGISTRY_PROJECT_NAME} image_name=${IMAGE} tag=${TAG} image_port=${IMAGE_PORT} host_ip=${TARGET_IP} ansible_user=${TARGET_USER}"
        """
      }
    }
  }
}