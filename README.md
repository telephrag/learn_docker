# How to run
Works under Linux. Wasn't tested on other platforms.

`sudo docker run --name=learn_docker_container --env-file=.env -u $(id -u non-root-username) -p 1337:8080 learn_docker`

`curl localhost:1337`
