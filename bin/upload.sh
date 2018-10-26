#! /bin/bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o weituan ../main.go
scp -i "/Users/zhangshengchao/Documents/aws.pem" weituan ubuntu@ec2-13-250-117-241.ap-southeast-1.compute.amazonaws.com:/home/ubuntu/weituan
