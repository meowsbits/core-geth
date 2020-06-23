pipeline {
    agent any
    environment {
        GETH_EXPORTS = '/data/ethereum-exports'
        GETH_DATADIR = '/data/geth'
    }
    stages {
        stage('Notify Github of Pending Jobs') {
            steps {
                githubNotify context: 'Classic PoW Regression', description: 'Assert import of canonical chain data', status: 'PENDING', account: 'meowsbits', repo: 'core-geth', credentialsId: 'meowsbits-github-jenkins', sha: "${GIT_COMMIT}"
                githubNotify context: 'Kotti Regression', description: 'Assert import of canonical chain data', status: 'PENDING', account: 'meowsbits', repo: 'core-geth', credentialsId: 'meowsbits-github-jenkins', sha: "${GIT_COMMIT}"
                githubNotify context: 'Mordor Regression', description: 'Assert import of canonical chain data', status: 'PENDING', account: 'meowsbits', repo: 'core-geth', credentialsId: 'meowsbits-github-jenkins', sha: "${GIT_COMMIT}"
                githubNotify context: 'Goerli Regression', description: 'Assert import of canonical chain data', status: 'PENDING', account: 'meowsbits', repo: 'core-geth', credentialsId: 'meowsbits-github-jenkins', sha: "${GIT_COMMIT}"
                githubNotify context: 'Classic Regression', description: 'Assert import of canonical chain data', status: 'PENDING', account: 'meowsbits', repo: 'core-geth', credentialsId: 'meowsbits-github-jenkins', sha: "${GIT_COMMIT}"
                githubNotify context: 'Foundation Regression', description: 'Assert import of canonical chain data', status: 'PENDING', account: 'meowsbits', repo: 'core-geth', credentialsId: 'meowsbits-github-jenkins', sha: "${GIT_COMMIT}"
            }
        }
        stage("Run Regression Tests") {
            parallel {
                stage('Classic (Real PoW)') {
                    agent {
                        label "aws-slave-t2-medium"
                    }
                    steps {
                        sh 'make geth'
                        sh './build/bin/geth version'
                        sh "rm -rf ${GETH_DATADIR}-classic-pow"
                        sh "./build/bin/geth --classic --cache=1024 --nocompaction --nousb --txlookuplimit=1 --datadir=${GETH_DATADIR}-classic-pow import ${GETH_EXPORTS}/classic.0-10000.rlp.gz"
                    }
                    post {
                        always {
                            sh "rm -rf ${GETH_DATADIR}-classic-pow"
                        }
                        success {
                            githubNotify context: 'Classic PoW Regression', description: 'Assert import of canonical chain data', status: 'SUCCESS', account: 'meowsbits', repo: 'core-geth', credentialsId: 'meowsbits-github-jenkins', sha: "${GIT_COMMIT}"
                        }
                        unsuccessful {
                            githubNotify context: 'Classic PoW Regression', description: 'Assert import of canonical chain data', status: 'FAILURE', account: 'meowsbits', repo: 'core-geth', credentialsId: 'meowsbits-github-jenkins', sha: "${GIT_COMMIT}"
                        }
                    }
                }
                stage('Kotti') {
                    agent {
                        label "aws-slave-m5-xlarge"
                    }
                    steps {
                        sh 'make geth'
                        sh './build/bin/geth version'
                        sh "rm -rf ${GETH_DATADIR}-kotti"
                        sh "./build/bin/geth --kotti --cache=2048 --nocompaction --nousb --txlookuplimit=1 --datadir=${GETH_DATADIR}-kotti import ${GETH_EXPORTS}/kotti.0-2544960.rlp.gz"
                    }
                    post {
                        always {
                            sh "rm -rf ${GETH_DATADIR}-kotti"
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
                    agent {
                        label "aws-slave-m5-xlarge"
                    }
                    steps {
                        sh 'make geth'
                        sh './build/bin/geth version'
                        sh "rm -rf ${GETH_DATADIR}-mordor"
                        sh "./build/bin/geth --mordor --fakepow --cache=2048 --nocompaction --nousb --txlookuplimit=1 --datadir=${GETH_DATADIR}-mordor import ${GETH_EXPORTS}/mordor.0-1686858.rlp.gz"
                        sh "rm -rf ${GETH_DATADIR}"
                    }
                    post {
                        always {
                            sh "rm -rf ${GETH_DATADIR}-mordor"
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
                    agent {
                        label "aws-slave-m5-xlarge"
                    }
                    steps {
                        sh 'make geth'
                        sh './build/bin/geth version'
                        sh "rm -rf ${GETH_DATADIR}-goerli"
                        sh "./build/bin/geth --goerli --cache=2048 --nocompaction --nousb --txlookuplimit=1 --datadir=${GETH_DATADIR}-goerli import ${GETH_EXPORTS}/goerli.0-2000000.rlp.gz"
                    }
                    post {
                        always {
                            sh "rm -rf ${GETH_DATADIR}-goerli"
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
                    agent {
                        label "aws-slave-m5-xlarge"
                    }
                    steps {
                        sh 'make geth'
                        sh './build/bin/geth version'
                        sh "rm -rf ${GETH_DATADIR}-classic"
                        sh "./build/bin/geth --classic --cache=2048 --nocompaction --nousb --txlookuplimit=1 --datadir=${GETH_DATADIR}-classic import ${GETH_EXPORTS}/classic.0-10620587.rlp.gz"
                    }
                    post {
                        always {
                            sh "rm -rf ${GETH_DATADIR}-classic"
                        }
                        success {
                            githubNotify context: 'Classic Regression', description: 'Assert import of canonical chain data', status: 'SUCCESS', account: 'meowsbits', repo: 'core-geth', credentialsId: 'meowsbits-github-jenkins', sha: "${GIT_COMMIT}"
                        }
                        unsuccessful {
                            sh 'dmesg | grep -i kill'
                            githubNotify context: 'Classic Regression', description: 'Assert import of canonical chain data', status: 'FAILURE', account: 'meowsbits', repo: 'core-geth', credentialsId: 'meowsbits-github-jenkins', sha: "${GIT_COMMIT}"
                        }
                    }
                }
                stage('Foundation') {
                    agent {
                        label "aws-slave-m5-xlarge"
                    }
                    steps {
                        sh 'make geth'
                        sh './build/bin/geth version'
                        sh "rm -rf ${GETH_DATADIR}-foundation"
                        sh "./build/bin/geth --cache=2048 --nocompaction --nousb --txlookuplimit=1 --datadir=${GETH_DATADIR}-foundation import ${GETH_EXPORTS}/ETH.0-10229163.rlp.gz"
                    }
                    post {
                        always {
                            sh "rm -rf ${GETH_DATADIR}-foundation"
                        }
                        success {
                            githubNotify context: 'Foundation Regression', description: 'Assert import of canonical chain data', status: 'SUCCESS', account: 'meowsbits', repo: 'core-geth', credentialsId: 'meowsbits-github-jenkins', sha: "${GIT_COMMIT}"
                        }
                        unsuccessful {
                            sh 'dmesg | grep -i kill'
                            githubNotify context: 'Foundation Regression', description: 'Assert import of canonical chain data', status: 'FAILURE', account: 'meowsbits', repo: 'core-geth', credentialsId: 'meowsbits-github-jenkins', sha: "${GIT_COMMIT}"
                        }
                    }
                }
            }
        }
    }
}


        // stage('Ropsten') {
        //     steps {
        //         sh "./build/bin/geth --ropsten --datadir=${GETH_DATADIR} import ${GETH_EXPORTS}/ropsten.0-8115552.rlp.gz"
        //         sh("rm -rf ${GETH_DATADIR}")
        //     }
        // }

        // stage('Print Context') {
        //     steps {
        //         sh 'hostname'
        //         sh 'uname -a'
        //         sh 'lsb_release -a'
        //         sh 'go version'
        //         sh 'go env'
        //         sh "ls -lshat ${GETH_EXPORTS}"
        //     }
        // }