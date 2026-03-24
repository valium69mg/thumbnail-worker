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
            agent { docker { image 'golang:1.23.7' } }
            steps {
                sh 'go version'
            }
        }

        stage('Build') {
            agent { docker { image 'golang:1.23.7' } }
            steps {
                sh 'go mod download'
                sh 'go build -o app .'
            }
        }

        stage('Test') {
            agent { docker { image 'golang:1.23.7' } }
            steps {
                sh 'go test ./... -coverprofile=coverage.out'
            }
        }

        stage('SonarQube Analysis') {
            agent any
            steps {
                withSonarQubeEnv('SonarQube') {
                    script {
                        def scannerHome = tool 'sonar-scanner'
                        sh "${scannerHome}/bin/sonar-scanner \
                            -Dsonar.projectKey=thumbnail-worker \
                            -Dsonar.sources=. \
                            -Dsonar.go.coverage.reportPaths=coverage.out"
                    }
                }
            }
        }

        stage('Quality Gate') {
            agent any
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