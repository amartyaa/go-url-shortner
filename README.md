[![Docker](https://github.com/amartyaa/go-url-shortner/actions/workflows/docker-publish.yml/badge.svg)](https://github.com/amartyaa/go-url-shortner/actions/workflows/docker-publish.yml)

# go-url-shortner
Basic URL shortner 
made with go, fiber & redis


### Containers:- 
1. Redis
2. Go build image

#### To run all the containers in detached mode
```bash
docker-compose up -d
```
Application runs on port 8080 {Can be configured in compose file}

#### If you already have redis running over
Change directory to api and build & run docker image
```bash
cd api && docker build . -t aa_url_shotrner_img:v1
docker run -dp 8080:8080 url_shortner_container
```
This will span containers in detached mode
