#!/bin/bash
GOOS=darwin GOARCH=amd64 go build -o ./bin/colab-katago-for-mac 
GOOS=linux GOARCH=amd64 go build -o ./bin/colab-katago-for-linux
GOOS=windows GOARCH=amd64 go build -o ./bin/colab-katago-for-windows

cd bin
rm -rf *.tar
cp colab-katago-for-mac colab-katago
tar cvf colab-katago.mac.tar colab-katago
cp colab-katago-for-linux colab-katago
tar cvf colab-katago.linux.tar colab-katago
cp colab-katago-for-windows colab-katago.exe
tar cvf colab-katago.windows.tar colab-katago.exe
rm colab-katago colab-katago.exe
cd -