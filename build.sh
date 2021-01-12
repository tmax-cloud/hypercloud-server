#!/bin/sh
docker rmi hypercloud-server:5.0 
docker build -t hypercloud-server:5.0  . 
