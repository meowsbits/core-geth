pipeline {
    agent any

    stages {
        stage('Print Context') {
            steps {
                sh 'pwd'
                sh 'git status'
                sh 'ls -lshat .'
                sh 'ls -lshat /data/ethereum-exports'
                sh 'go version'
                sh 'go env'
            }
        }
        stage('Build') {
            steps {
                echo "Building 1234567..."
                sh 'make geth'
                sh './build/bin/geth version'
            }
        }
        stage('Kotti') {
            steps {
                sh './build/bin/geth --help'
            }
        }
    }
}