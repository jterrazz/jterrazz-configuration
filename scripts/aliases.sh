# Docker Remove All Processes
drm() {
    docker rm -vf $(docker ps -aq)
}

# Docker Remove All Images
drmi() {
    docker rmi -f $(docker images -aq)
}
