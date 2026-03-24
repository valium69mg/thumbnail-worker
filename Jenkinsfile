pipeline {
    agent any

    environment {
        SONARQUBE = 'SonarQube'
    }

    triggers {
        pollSCM('H/5 * * * *')
    }

    stages {
        stage('Checkout') {
            steps {
                git branch: 'dev',
                    url: 'https://github.com/valium69mg/thumbnail-worker'
            }
        }

        stage('Build') {
            steps {
                sh 'go mod download'
                sh 'go build -o app .'
            }
        }

        stage('Test') {
            steps {
                sh 'go test ./... -coverprofile=coverage.out'
            }
        }

        stage('SonarQube Analysis') {
            steps {
                withSonarQubeEnv('SonarQube') {
                    sh 'sonar-scanner -Dsonar.projectKey=thumbnail-worker -Dsonar.sources=. -Dsonar.go.coverage.reportPaths=coverage.out'
                }
            }
        }

        stage('Quality Gate') {
            steps {
                timeout(time: 2, unit: 'MINUTES') {
                    waitForQualityGate abortPipeline: true
                }
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