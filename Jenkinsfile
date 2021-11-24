pipeline {
    agent {
        docker {
            image 'golang:1.17.3'
        }
    }
    stages {
        stage('build') {
            steps {
                sh 'go build'
            }
        }
    }
}
