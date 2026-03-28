pipeline {
    agent none

    environment {
        SONARQUBE = 'SonarQube'
        DOCKER_IMAGE = 'carlostranquilinocr98/single-vendor-ecommerce-thumbnail-worker'
    }

    triggers {
        pollSCM('H/5 * * * *')
    }

    stages {
        stage('Verify Go') {
            agent { docker { image 'golang:1.25.7' } }
            steps {
                sh 'go version'
            }
        }

        stage('Build') {
            agent { docker { image 'golang:1.25.7' } }
            steps {
                sh 'go mod download'
                sh 'go build -o app .'
            }
        }

        stage('Test') {
            agent { docker { image 'golang:1.25.7' } }
            steps {
                sh 'go mod download'
                sh 'go test ./... -v -coverprofile=coverage.out -covermode=atomic'
                sh 'go tool cover -func=coverage.out'
            }
        }

        stage('SonarQube Analysis') {
            agent {
                docker {
                    image 'sonarsource/sonar-scanner-cli:latest'
                    args '--network=jenkins_default -u 0:0'
                }
            }
            environment {
                SONAR_HOST_URL = 'http://sonarqube:9000'
                SONAR_LOGIN    = credentials('sonar-token')
            }
            steps {
                withSonarQubeEnv('SonarQube') {
                    sh '''
                        sonar-scanner \
                        -Dsonar.projectKey=thumbnail-worker \
                        -Dsonar.sources=. \
                        -Dsonar.go.coverage.reportPaths=coverage.out \
                        -Dsonar.host.url=$SONAR_HOST_URL \
                        -Dsonar.login=$SONAR_LOGIN
                    '''
                }
            }
        }

        stage('Wait for SonarQube Quality Gate') {
            steps {
                echo 'Waiting for SonarQube Quality Gate...'
                timeout(time: 10, unit: 'MINUTES') {
                    waitForQualityGate abortPipeline: true
                }
            }
        }

        stage('Build & Push Docker Image') {
            agent any
            steps {
                echo 'Building Docker image...'
                sh "docker build -t carlostranquilinocr98/single-vendor-ecommerce-thumbnail-worker:latest ."

                echo 'Pushing Docker image to Docker Hub...'
                withCredentials([usernamePassword(credentialsId: 'docker-hub-credentials',
                                                  usernameVariable: 'DOCKER_USER',
                                                  passwordVariable: 'DOCKER_PASS')]) {
                    sh "echo $DOCKER_PASS | docker login -u $DOCKER_USER --password-stdin"
                    sh "docker push carlostranquilinocr98/single-vendor-ecommerce-thumbnail-worker:latest"
                }
            }
        }

        stage('Deploy') {
            agent any
            steps {
                sshagent(['ubuntu-server-ssh']) {
                    sh '''
                        ssh -o StrictHostKeyChecking=no carlostr@192.168.100.50 \
                        "cd /home/carlostr/single-vendor-ecommerce && \
                        docker compose pull && \
                        docker compose up -d --remove-orphans"
                    '''
                }
            }
        }
    }

    post {
        success {
            echo '✅ Build, tests, SonarQube analysis, quality gate, and Docker push succeeded!'
        }
        failure {
            echo '❌ Build, tests, SonarQube analysis, quality gate, or Docker push failed!'
        }
    }
}
