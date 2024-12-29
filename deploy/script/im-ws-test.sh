#!/bin/bash
reso_addr='crpi-dhjl0efo9cuwus90.cn-hangzhou.personal.cr.aliyuncs.com/gozero-chat/im-ws-dev'
tag='latest'

container_name="easy-chat-im-ws-test"

docker stop ${container_name}

docker rm ${container_name}

docker rmi ${reso_addr}:${tag}

docker pull ${reso_addr}:${tag}


# 如果需要指定配置文件的
# docker run -p 10001:8080 --network imooc_easy-im -v /easy-im/config/user-rpc:/user/conf/ --name=${container_name} -d ${reso_addr}:${tag}
docker run  --name=${container_name} -d ${reso_addr}:${tag}