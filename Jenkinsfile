def app = 'polymur'
def namespace = 'techops'
def registry = "registry.revinate.net/${namespace}"
def gopath = "/go/src/github.com/revinate/${app}"
def name = "${registry}/${app}"

stage 'Golang build'
node {
  checkout scm

  usingDocker {
    sh "docker run --rm -v `pwd`:${gopath} -w ${gopath} golang:latest make"
  }

  stash name: 'binary', includes: "${app}"
}

stage 'Docker build and push'
node {
  checkout scm
  unstash 'binary'

  version = currentVersion()
  hoister.registry = registry
  hoister.imageName = app
  hoister.buildAndPush version

  stagehandPublish(app, version)
}

stage 'Kubernetes deploy to test'
node {
  deployToTest("${app}", "${namespace}")
}
