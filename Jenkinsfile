pipeline {
    agent any
    stages {
        stage ('Lint') {
            steps {
                echo 'Linting application'
                sh 'make golang.lint'
            }
        }
        stage ('Build') {
            steps {
                echo 'Building'
                sh 'make local.build'
            }
        }
    }
}