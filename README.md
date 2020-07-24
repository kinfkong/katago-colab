# colab-katago
See how to use: https://colab.research.google.com/drive/1NGGG4t59Atnq2c9uZPOgg28KsCOcQptA?usp=sharing

# Build
```
GOOS=darwin GOARCH=amd64 go build -o ./bin/colab-katago-for-mac 
GOOS=linux GOARCH=amd64 go build -o ./bin/colab-katago-for-linux
GOOS=windows GOARCH=amd64 go build -o ./bin/colab-katago-for-windows
```