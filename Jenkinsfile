pipeline {
    agent any

    environment {
        // Define any global environment variables here if needed
        IMAGE_NAME = "raijai-backend"
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Build Docker Image') {
            steps {
                script {
                    sh 'docker build -t ${IMAGE_NAME} .'
                }
            }
        }

        stage('Deploy') {
            steps {
                script {
                    withCredentials([
                        string(credentialsId: 'PG_DB_HOST', variable: 'PG_DB_HOST'),
                        string(credentialsId: 'PG_DB_PORT', variable: 'PG_DB_PORT'),
                        string(credentialsId: 'PG_DB_USER', variable: 'PG_DB_USER'),
                        string(credentialsId: 'PG_DB_PASSWORD', variable: 'PG_DB_PASSWORD'),
                        string(credentialsId: 'DB_PG_RAIJAI_DB_NAME', variable: 'DB_PG_RAIJAI_DB_NAME'),
                        string(credentialsId: 'KEYCLOAK_ISSUER', variable: 'KEYCLOAK_ISSUER'),
                        string(credentialsId: 'KEYCLOAK_CLIENT_ID', variable: 'KEYCLOAK_RAIJAI_CLIENT_ID')
                    ]) {
                        // Deploy using docker-compose
                        sh 'docker compose up -d --build'
                    }
                }
            }
        }
    }

    post {
        always {
            cleanWs()
        }
    }
}
