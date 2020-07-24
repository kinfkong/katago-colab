#!/bin/bash
GOOS=darwin GOARCH=amd64 go build -o ./bin/colab-katago-for-mac 
GOOS=linux GOARCH=amd64 go build -o ./bin/colab-katago-for-linux
GOOS=windows GOARCH=amd64 go build -o ./bin/colab-katago-for-windows

cd bin
rm -rf *.zip
cp colab-katago-for-mac colab-katago
zip colab-katago.mac.zip colab-katago
cp colab-katago-for-linux colab-katago
zip colab-katago.linux.zip colab-katago
cp colab-katago-for-windows colab-katago.exe
zip colab-katago.windows.zip colab-katago.exe
rm colab-katago colab-katago.exe
cd -