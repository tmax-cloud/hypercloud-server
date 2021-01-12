#!/bin/sh
img=tmaxcloudck/hypercloud-server:b5.0.0.5
docker rmi $img 
docker build -t $img  . 
docker push $img 
