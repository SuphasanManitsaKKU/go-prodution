pipeline {
  agent any

  stages {
    stage('Load All Configs from Vault') {
      steps {
        script {
          withVault(
            vaultSecrets: [
              [
                path: 'pipeline/env', // <-- ตรงกับที่คุณ `vault kv put`
                secretValues: [
                  [envVar: 'REGISTRY', vaultKey: 'REGISTRY'],
                  [envVar: 'REGISTRY_PROJECT_NAME', vaultKey: 'REGISTRY_PROJECT_NAME'],
                  [envVar: 'IMAGE', vaultKey: 'IMAGE'],
                  [envVar: 'IMAGE_OUTPUT_PORT', vaultKey: 'IMAGE_OUTPUT_PORT'],
                  [envVar: 'TARGET_USER', vaultKey: 'TARGET_USER'],
                  [envVar: 'TARGET_IP', vaultKey: 'TARGET_IP'],
                  [envVar: 'DOCKER_USERNAME', vaultKey: 'DOCKER_USERNAME'],
                  [envVar: 'DOCKER_PASSWORD', vaultKey: 'DOCKER_PASSWORD']
                ]
              ]
            ]
          ) {
            env.TAG = env.BUILD_NUMBER
            env.FULL_IMAGE = "${env.REGISTRY}/${env.REGISTRY_PROJECT_NAME}/${env.IMAGE}:${env.TAG}"

            sh 'echo ✅ Loaded all Vault secrets.'
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
        sh 'echo $DOCKER_PASSWORD | docker login $REGISTRY -u $DOCKER_USERNAME --password-stdin'
      }
    }

    stage('Push Docker Image') {
      steps {
        sh 'docker push $FULL_IMAGE'
      }
    }

    // stage('Add Jenkins SSH Public Key to Target') {
    //   steps {
    //     sh """
    //       echo "[INFO] Copying Jenkins public key to target machine..."
    //       PUB_KEY_PATH=/var/jenkins_home/.ssh/id_rsa.pub

    //       ssh -o StrictHostKeyChecking=no ${TARGET_USER}@${TARGET_IP} "mkdir -p ~/.ssh && touch ~/.ssh/authorized_keys && chmod 600 ~/.ssh/authorized_keys"
    //       cat \$PUB_KEY_PATH | ssh -o StrictHostKeyChecking=no ${TARGET_USER}@${TARGET_IP} "grep -qxF '\$(cat)' ~/.ssh/authorized_keys || echo '\$(cat)' >> ~/.ssh/authorized_keys"
    //     """
    //   }
    // }

    // stage('Trigger Ansible CD') {
    //   steps {
    //     sh """
    //       echo "[INFO] Running Ansible deployment..."

    //       mkdir -p ~/.ssh
    //       cp /var/jenkins_home/.ssh/id_rsa ~/.ssh/id_rsa
    //       chmod 600 ~/.ssh/id_rsa
    //       ssh-keyscan -H $TARGET_IP >> ~/.ssh/known_hosts

    //       ansible-playbook -i /var/jenkins_home/ansible/inventory.ini \
    //         /var/jenkins_home/ansible/playbooks/deploy_app.yml \
    //         --extra-vars "registry=${REGISTRY} registry_project_name=${REGISTRY_PROJECT_NAME} image_name=${IMAGE} tag=${TAG} image_port=${IMAGE_OUTPUT_PORT} host_ip=${TARGET_IP} ansible_user=${TARGET_USER}"
    //     """
    //   }
    // }
  }
}