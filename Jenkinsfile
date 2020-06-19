pipeline {
    agent any

    stages {
        stage('Print Context') {
            steps {
                pwd
                git status
                ls -lshat .
                ls -lshat /data/ethereum-exports
                go version
                go env
            }
        }
        stage('Install') {
            steps {
                go get -v -t -d ./...
            }
        }
        stage('Build') {
            steps {
                make geth
                ./build/bin/geth version
            }
        }
        stage('Kotti') {
            steps {
                ./build/bin/geth --help
            }
        }
    }
}