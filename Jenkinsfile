pipeline {
    agent any

    environment {
        GETH_EXPORTS = '/data/ethereum-exports'
        GETH_DATADIR = '/data/geth'
    }
    stages {
        stage('Print Context') {
            steps {
                sh 'uname -a'
                sh 'lsb_release -a'
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
        stage('Kotti') {
            steps {
                sh "./build/bin/geth --kotti --datadir=${GETH_DATADIR} import ${GETH_EXPORTS}/kotti.0-2544960.rlp.gz"
                sh("rm -rf ${GETH_DATADIR}")
            }
        }
        stage('Mordor') {
            steps {
                sh "./build/bin/geth --mordor --datadir=${GETH_DATADIR} import ${GETH_EXPORTS}/mordor.0-1686858.rlp.gz"
                sh("rm -rf ${GETH_DATADIR}")
            }
        }
        // stage('Goerli') {
        //     steps {
        //         sh "./build/bin/geth --goerli --datadir=${GETH_DATADIR} import ${GETH_EXPORTS}/goerli.0-2886512.rlp.gz"
        //         sh("rm -rf ${GETH_DATADIR}")
        //     }
        // }
        // stage('Classic') {
        //     steps {
        //         sh "./build/bin/geth --classic --datadir=${GETH_DATADIR} import ${GETH_EXPORTS}/classic.0-10620587.rlp.gz"
        //         sh("rm -rf ${GETH_DATADIR}")
        //     }
        // }
        // stage('Ropsten') {
        //     steps {
        //         sh "./build/bin/geth --ropsten --datadir=${GETH_DATADIR} import ${GETH_EXPORTS}/ropsten.0-8115552.rlp.gz"
        //         sh("rm -rf ${GETH_DATADIR}")
        //     }
        // }
        // stage('Foundation') {
        //     steps {
        //         sh "./build/bin/geth --datadir=${GETH_DATADIR} import ${GETH_EXPORTS}/ETH.0-10229163.rlp.gz"
        //         sh("rm -rf ${GETH_DATADIR}")
        //     }
        // }
    }
    post {
        always {
            sh("rm -rf ${GETH_DATADIR}")
        }
    }
}