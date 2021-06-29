pipeline {
    agent any
    stages {
        stage ('Lint') {
            steps {
                echo 'Linting application'
                sh 'make coverage.lint'
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