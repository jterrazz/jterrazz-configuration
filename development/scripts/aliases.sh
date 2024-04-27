# Docker Processes

drm() {
    docker rm -vf $(docker ps -aq)
}

# Docker Images

drmi() {
    docker rmi -f $(docker images -aq)
}
