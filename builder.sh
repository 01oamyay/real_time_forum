#remove unused object (free cash)
docker system prune -a
#build img with the Dockerfile 
docker image build -f Dockerfile -t forum-img .
#run a container with the previous image
docker container run  -dp 8080:8080 --name forum forum-img
#execute the bash in the running container
docker exec -it forum /bin/bash