pipeline {
    agent any

    triggers {
        pollSCM('H/5 * * * *')
    }

    tools {
        go 'go-1.22'
    }

    stages {

        stage('Checkout') {
            steps {
                git 'https://github.com/valium69mg/thumbnail-worker'
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
            environment {
                SCANNER_HOME = tool 'sonarqube-scanner'
            }
            steps {
                withSonarQubeEnv('sonarqube-server') {
                    sh '''
                    $SCANNER_HOME/bin/sonar-scanner \
                      -Dsonar.projectKey=thumbnail-worker \
                      -Dsonar.projectName=thumbnail-worker \
                      -Dsonar.sources=. \
                      -Dsonar.go.coverage.reportPaths=coverage.out
                    '''
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