group "default" {
  targets = [
    "orchestrator",
    "runner"
  ]
}

variable "REGISTRY" {
    default = "harbor.petite-patate.ovh/library"  
}

variable "TAG" {
  default = "latest"
}

variable "ENV_FLAG" {
  default = "-dev"
}

target "orchestrator" {
  context = "."
  dockerfile = "Dockerfile.orchestrator"
  tags = [
    "${REGISTRY}/orchestrator:${TAG}${ENV_FLAG}"
  ]

  cache-from = [
    "type=registry,ref=${REGISTRY}/orchestrator:cache"
  ]

  cache-to = [
    "type=registry,ref=${REGISTRY}/orchestrator:cache,mode=max"
  ]

  platforms = [ "linux/amd64", "linux/arm64" ]
}

target "runner" {
  context = "."
  dockerfile = "Dockerfile.runner"
  tags = [
    "${REGISTRY}/runner:${TAG}${ENV_FLAG}"
  ]

  cache-from = [
    "type=registry,ref=${REGISTRY}/runner:cache"
  ]

  cache-to = [
    "type=registry,ref=${REGISTRY}/runner:cache,mode=max"
  ]

  platforms = [ "linux/amd64", "linux/arm64" ]
}