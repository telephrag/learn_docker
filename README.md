# How to run
Works under Linux. Wasn't tested on other platforms.

`sudo docker run --name=learn_docker_container --env-file=.env -v /home/$(id -nu 1000)/volumes/learn_docker_volume:/data -p 1337:8080 learn_docker`

`curl localhost:1337`
