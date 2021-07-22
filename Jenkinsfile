pipeline {
  agent {
    docker {
      image 'golangci/golangci-lint:latest'
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
        sh 'go get -t'
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
