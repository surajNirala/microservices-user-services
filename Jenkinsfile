pipeline {
    agent any 

    environment {
        DOCKER_HUB_REPO = 'snirala1995/user_services'
        DOCKER_IMAGE_TAG = "${DOCKER_HUB_REPO}:${env.BUILD_NUMBER}"
        CONTAINER_NAME = 'user_services' // Name for your Docker container
        CONTAINER_PORT = '9091' // Port inside the Docker container
        USER_SERVICE_ENV = credentials('USER_SERVICE_ENV') // For secret file
        SNIRALA_DOCKERHUB_CREDENTIAL = 'snirala-dockerhub-credentials'
        SERVER_1 = '34.131.139.0' 
        CREDENTIALS_SERVER_2 = 'credentials-server-2'
        HOST_PORT = '9091'
    }

    stages {

        stage('Check Existing Container') {
            steps {
                script {
                    echo "Checking if the container already exists"
                    def existingContainer = sh(script: "docker ps -aqf name=${CONTAINER_NAME}-${HOST_PORT}", returnStdout: true).trim()
                    if (existingContainer) {
                        echo "Stopping and removing the existing container: ${CONTAINER_NAME}-${HOST_PORT}"
                        sh "docker rm -f ${CONTAINER_NAME}-${HOST_PORT}"
                    }
                }
            }
        }

        // stage('Prepare .env File') {
        //     steps {
        //         echo "Removing the existing .env file if it exists"
        //         sh 'rm -f .env'
        //         echo "Copying the new .env file"
        //         sh "cp ${USER_SERVICE_ENV} .env"
        //         echo "=================env file copied successfully================"
        //     }
        // }

        stage('Build Docker Image') {
            steps {
                script {
                    echo "Building the Docker image"
                    docker.build(DOCKER_IMAGE_TAG)
                    echo "Docker image built successfully."
                }
            }
        }

        stage('Push Docker Image To Docker Hub') {
            steps {
                script {
                    try {
                        echo "Pushing Docker image to DockerHub."
                        docker.withRegistry('https://registry.hub.docker.com', SNIRALA_DOCKERHUB_CREDENTIAL) {
                            docker.image(DOCKER_IMAGE_TAG).push()
                        }
                        echo "Docker image pushed to DockerHub successfully."
                    } catch (Exception e) {
                        echo "Failed to push Docker image: ${e.message}"
                    }
                }
            }
        }

        stage('Deploy') {
            steps {
                script {
                    echo "Deploying to ================= SRJ-SERVER ============== (${SERVER_1})"
                    sshagent([CREDENTIALS_SERVER_2]) {
                        echo "Deploying to ${SERVER_1} on port ${HOST_PORT} with image ${DOCKER_IMAGE_TAG}"
                        withCredentials([usernamePassword(credentialsId: SNIRALA_DOCKERHUB_CREDENTIAL, usernameVariable: 'DOCKER_USERNAME', passwordVariable: 'DOCKER_PASSWORD')]){
                        sh """
                            echo "Connecting to ${SERVER_1}..."
                            ssh -o StrictHostKeyChecking=no srj@${SERVER_1} <<EOF
                            echo "Remote server connected successfully!"

                            echo "Logging into DockerHub"
                            echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin

                            echo "Pulling Docker image from DockerHub: ${DOCKER_IMAGE_TAG}"
                            docker pull ${DOCKER_IMAGE_TAG}

                            echo "Stopping and removing any existing container"
                            docker rm -f ${CONTAINER_NAME}-${HOST_PORT} || true

                            echo "Running the Docker container"
                            
                            docker run -d --init -p ${HOST_PORT}:${CONTAINER_PORT} --name ${CONTAINER_NAME}-${HOST_PORT} --env-file .env ${DOCKER_IMAGE_TAG}

                            echo "Docker image ${DOCKER_IMAGE_TAG} run successfully."
                            exit
                        """
                        }
                        //docker run -d --init -p ${HOST_PORT}:${CONTAINER_PORT} --name ${CONTAINER_NAME}-${HOST_PORT} ${DOCKER_IMAGE_TAG}
                    }
                       
                }
            }
        }
    }

    post {
        success {
            script {
                echo "Docker image ${DOCKER_IMAGE_TAG} successfully pushed to Docker Hub."
                echo "Container running on port: ${HOST_PORT}"
                echo "Pipeline completed successfully."
                echo "Click the following link to check the website live: ${SERVER_1}:"
            }
        }
        failure {
            script {
                echo "Pipeline failed. Check logs for details."
            }
        }
    }
}
