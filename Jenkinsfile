pipeline {
  agent {
    docker {
      image 'septemhill/dockergolangit:latest'
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