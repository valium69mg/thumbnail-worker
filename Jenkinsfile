pipeline {
    agent none

    environment {
        SONARQUBE = 'SonarQube'
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
                    args '-u 0:0'
                }
            }
            environment {
                SONAR_HOST_URL = 'http://sonarqube:9000'
                SONAR_LOGIN    = credentials('sonar-token')
            }
            steps {
                echo 'Running SonarQube analysis...'
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

    post {
        success {
            echo '✅ Build, tests, and SonarQube analysis succeeded!'
        }
        failure {
            echo '❌ Build, tests, or SonarQube analysis failed!'
        }
    }
}