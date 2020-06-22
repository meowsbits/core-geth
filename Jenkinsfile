pipeline {
    agent {
        label "aws-slave"
    }
    environment {
        GETH_EXPORTS = '/data/ethereum-exports'
        GETH_DATADIR = '/data/geth'
    }
    stages {
        stage('Notify Github of Pending Jobs') {
            steps {
                githubNotify context: 'Kotti Regression', description: 'Assert import of canonical chain data', status: 'PENDING', account: 'meowsbits', repo: 'core-geth', credentialsId: 'meowsbits-github-jenkins', sha: "${GIT_COMMIT}"
                githubNotify context: 'Mordor Regression', description: 'Assert import of canonical chain data', status: 'PENDING', account: 'meowsbits', repo: 'core-geth', credentialsId: 'meowsbits-github-jenkins', sha: "${GIT_COMMIT}"
                githubNotify context: 'Goerli Regression', description: 'Assert import of canonical chain data', status: 'PENDING', account: 'meowsbits', repo: 'core-geth', credentialsId: 'meowsbits-github-jenkins', sha: "${GIT_COMMIT}"
                githubNotify context: 'Classic Regression', description: 'Assert import of canonical chain data', status: 'PENDING', account: 'meowsbits', repo: 'core-geth', credentialsId: 'meowsbits-github-jenkins', sha: "${GIT_COMMIT}"
            }
        }
        stage('Print Context') {
            steps {
                sh 'uname -a'
                sh 'lsb_release -a'
                sh 'go version'
                sh 'go env'
                sh "ls -lshat ${GETH_EXPORTS}"
            }
        }
        stage('Build') {
            steps {
                sh 'make geth'
                sh './build/bin/geth version'
                // This should never happen normally, but in case the instance is halted and the job is unable to clean up after itself, this ensures tabula rasa.
                sh "rm -rf ${GETH_DATADIR}"
            }
        }
        stage('Kotti') {
            steps {
                sh "./build/bin/geth --kotti --cache=2048 --nocompaction --nousb --txlookuplimit=1 --datadir=${GETH_DATADIR} import ${GETH_EXPORTS}/kotti.0-2544960.rlp.gz"
            }
            post {
                always {
                    sh("rm -rf ${GETH_DATADIR}")
                }
                success {
                    githubNotify context: 'Kotti Regression', description: 'Assert import of canonical chain data', status: 'SUCCESS', account: 'meowsbits', repo: 'core-geth', credentialsId: 'meowsbits-github-jenkins', sha: "${GIT_COMMIT}"
                }
                unsuccessful {
                    githubNotify context: 'Kotti Regression', description: 'Assert import of canonical chain data', status: 'FAILURE', account: 'meowsbits', repo: 'core-geth', credentialsId: 'meowsbits-github-jenkins', sha: "${GIT_COMMIT}"
                }
            }
        }
        stage('Mordor') {
            steps {
                sh "./build/bin/geth --mordor --fakepow --cache=2048 --nocompaction --nousb --txlookuplimit=1 --datadir=${GETH_DATADIR} import ${GETH_EXPORTS}/mordor.0-1686858.rlp.gz"
                sh("rm -rf ${GETH_DATADIR}")
            }
            post {
                always {
                    sh("rm -rf ${GETH_DATADIR}")
                }
                success {
                    githubNotify context: 'Mordor Regression', description: 'Assert import of canonical chain data', status: 'SUCCESS', account: 'meowsbits', repo: 'core-geth', credentialsId: 'meowsbits-github-jenkins', sha: "${GIT_COMMIT}"
                }
                unsuccessful {
                    githubNotify context: 'Mordor Regression', description: 'Assert import of canonical chain data', status: 'FAILURE', account: 'meowsbits', repo: 'core-geth', credentialsId: 'meowsbits-github-jenkins', sha: "${GIT_COMMIT}"
                }
            }
        }
        stage('Goerli') {
            steps {
                sh "./build/bin/geth --goerli --cache=1024 --nocompaction --nousb --txlookuplimit=1 --datadir=${GETH_DATADIR} import ${GETH_EXPORTS}/goerli.0-2886512.rlp.gz"
            }
            post {
                always {
                    sh("rm -rf ${GETH_DATADIR}")
                }
                success {
                    githubNotify context: 'Goerli Regression', description: 'Assert import of canonical chain data', status: 'SUCCESS', account: 'meowsbits', repo: 'core-geth', credentialsId: 'meowsbits-github-jenkins', sha: "${GIT_COMMIT}"
                }
                unsuccessful {
                    githubNotify context: 'Goerli Regression', description: 'Assert import of canonical chain data', status: 'FAILURE', account: 'meowsbits', repo: 'core-geth', credentialsId: 'meowsbits-github-jenkins', sha: "${GIT_COMMIT}"
                }
            }
        }
        stage('Classic') {
            steps {
                sh "./build/bin/geth --classic --fakepow --cache=1024 --nocompaction --nousb --txlookuplimit=1 --datadir=${GETH_DATADIR} import ${GETH_EXPORTS}/classic.0-10620587.rlp.gz"
            }
            post {
                always {
                    sh("rm -rf ${GETH_DATADIR}")
                }
                success {
                    githubNotify context: 'Classic Regression', description: 'Assert import of canonical chain data', status: 'SUCCESS', account: 'meowsbits', repo: 'core-geth', credentialsId: 'meowsbits-github-jenkins', sha: "${GIT_COMMIT}"
                }
                unsuccessful {
                    githubNotify context: 'Classic Regression', description: 'Assert import of canonical chain data', status: 'FAILURE', account: 'meowsbits', repo: 'core-geth', credentialsId: 'meowsbits-github-jenkins', sha: "${GIT_COMMIT}"
                }
            }
        }
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