# Install
- export PATH=$PATH:/home/username/go/bin
- go install


## RESTful
1. /albums
- get, get a list of albums
- post, add a new album 

2. /albums/:id
- get, get an album by id


sudo tee /etc/docker/daemon.json <<EOF
{
  "registry-mirrors": [
    "https://docker.m.daocloud.io"
  ]
}