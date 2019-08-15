pipeline {
  agent {
        docker {
            image 'golang:1.12.8-alpine3.10'
        }
  }
  stages {
    stage('Build') {
      steps {
        echo 'Building..'
      }
    }
    stage('Test') {
      steps {
        sh 'go test -v'
      }
    }
    stage('Deploy') {
      steps {
        echo 'Deploying....'
      }
    }
  }
}
