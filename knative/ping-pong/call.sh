#!/bin/bash

NODE_PORT=$(kubectl --namespace kong get service kong-proxy -o go-template='{{(index .spec.ports 0).nodePort}}')

curl -H "Host: f1.default.example.com" http://$(curl ifconfig.me):${NODE_PORT}
