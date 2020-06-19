pipeline {
    agent any

    environment {
        GETH_EXPORTS = '/data/ethereum-exports'
        GETH_DATADIR = '/data/geth'
    }
    stages {
        stage('Print Context') {
            steps {
                sh 'pwd'
                sh "ls -lshat ${GETH_EXPORTS}"
                sh 'go version'
                sh 'go env'
            }
        }
        stage('Build') {
            steps {
                sh 'make geth'
                sh './build/bin/geth version'
            }
        }
        stage("Example") {
            steps {
                sh("echo exports dir: ${GETH_EXPORTS}")
                sh("echo geth datadir: ${GETH_DATADIR}")
            }
        }
    //     stage('Kotti') {
    //         steps {
    //             sh "./build/bin/geth --kotti --datadir=${GETH_DATADIR} import ${GETH_EXPORTS}/kotti.0-2544960.rlp.gz"
    //             sh("rm -rf ${GETH_DATADIR}")
    //         }
    //     }
    //     stage('Mordor') {
    //         steps {
    //             sh "./build/bin/geth --mordor --datadir=${GETH_DATADIR} import ${GETH_EXPORTS}/mordor.0-1686858.rlp.gz"
    //             sh("rm -rf ${GETH_DATADIR}")
    //         }
    //     }
    }
    // post {
    //     always {
    //         sh("rm -rf ${GETH_DATADIR}")
    //     }
    // }
}