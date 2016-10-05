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

  stash name: 'binary-polymur', includes: "polymur"
  stash name: 'binary-polymur-gateway', includes: "polymur-gateway"
  stash name: 'binary-polymur-proxy', includes: "polymur-proxy"
}

stage 'Docker build and push :: polymur'
node {
  checkout scm
  unstash 'binary-polymur'

  version = currentVersion()
  hoister.registry = registry
  hoister.imageName = polymur
  hoister.buildAndPush version

  stagehandPublish(polymur, version)
}

stage 'Docker build and push :: polymur-gateway'
node {
  checkout scm
  unstash 'binary-polymur-gateway'

  version = currentVersion()
  hoister.registry = registry
  hoister.imageName = polymur-gateway
  hoister.buildAndPush version

  stagehandPublish(polymur-gateway, version)
}

stage 'Docker build and push :: polymur-proxy'
node {
  checkout scm
  unstash 'binary-polymur-proxy'

  version = currentVersion()
  hoister.registry = registry
  hoister.imageName = polymur-proxy
  hoister.buildAndPush version

  stagehandPublish(polymur-proxy, version)
}


stage 'Kubernetes deploy to test'
node {
  deployToTest("${app}", "${namespace}")
}
