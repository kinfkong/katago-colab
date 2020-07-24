#!/bin/bash
#CONFIG_PATH=./config
#TARGET_PATH=.
CONFIG_PATH=/content/katago-colab/config
TARGET_PATH=/content

SUFFIX=$1
CHANGE_VALUE=${SUFFIX%?}
if [[ "$SUFFIX" == *s ]]
then
    CMD1="s|^.*maxTime = .*$|maxTime = $CHANGE_VALUE|g"
    CMD2="s|^.*maxVisits = .*$|# maxVisits = <PENDING>|g"
else
    CMD1="s|^.*maxVisits = .*$|maxVisits = $CHANGE_VALUE|g"
    CMD2="s|^.*maxTime = .*$|# maxTime = <PENDING>|g"
fi

cp $CONFIG_PATH/gtp_colab.cfg $TARGET_PATH/gtp_colab_$SUFFIX.cfg

if [[ "$OSTYPE" == "darwin"* ]]
then
    sed -i '' -e "$CMD1" $TARGET_PATH/gtp_colab_$SUFFIX.cfg
    sed -i '' -e "$CMD2" $TARGET_PATH/gtp_colab_$SUFFIX.cfg
else
    sed -i -e "$CMD1" $TARGET_PATH/gtp_colab_$SUFFIX.cfg
    sed -i -e "$CMD2" $TARGET_PATH/gtp_colab_$SUFFIX.cfg
fi