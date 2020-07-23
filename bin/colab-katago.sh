#!/bin/bash
USER_NAME=$1
if [ "$USER_NAME" == "" ];
then
    echo "USER_NAME is missing"
    exit -1
fi
curl --silent https://kata-config.oss-cn-beijing.aliyuncs.com/$USER_NAME.ssh.txt -o /tmp/ssh.txt
SSH_CMD=`cat /tmp/ssh.txt`
ssh $SSH_CMD