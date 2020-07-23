#!/bin/bash
USER_NAME=$1
if [ "$USER_NAME" == "" ];
then
    echo "USER_NAME is missing"
    exit -1
fi
curl --silent https://kata-config.oss-cn-beijing.aliyuncs.com/$USER_NAME.ssh.txt -o /tmp/ssh.txt
curl --silent https://kata-config.oss-cn-beijing.aliyuncs.com/$USER_NAME.pem -o /tmp/ssh.pem

chmod 600 /tmp/ssh.pem

SSH_CMD=`cat /tmp/ssh.txt`

ssh -i /tmp/ssh.pem $SSH_CMD