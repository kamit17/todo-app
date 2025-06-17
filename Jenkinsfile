pipeline {
    agent none

    stages{
        stage('Build backend'){
            agent {
                docker {
                    image 'golang:1.21-alpine'
                    args '-u root'
                }
            }
            steps {
                echo 'Building GO backend...'
                dir('backend') {
                    sh 'go mod tidy'
                    sh 'go build -o todo-app'
                    sh './todo-app & sleep 1 && curl http://localhost:8080 || echo "Server not reachable "'
                }
            }
        }

        stage('Validate Frontend') {
            agent {
                docker {
                    image 'node:16-alpine'
                    args '-u root'
                }
            }
            steps {
                echo 'Validating frontend'
                dir('frontend') {
                    sh 'ls -lh index.html script.js style.css'
                }
            }
        }       
    }

    post {
        success {
            echo 'Todo App pipeline ran succesfully.'
        }
        failure {
            echo 'Pipeline failed. Check logs'
        }
    }
}