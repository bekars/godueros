# goDuerOS - DuerOS Client by Go

golang实现的DuerOS客户端SDK，目前只实现了基本的问答功能，交流和学习使用。

## Pre RUN

安装以下软件包

1. portaudio
2. mpg123

安装以下golang库

1. github.com/gordonklaus/portaudio
2. github.com/bobertlo/go-mpg123/mpg123

## RUN

执行以下命令运行

````bash
go run main/main.go
````

## TODO LIST

* 完整的DCS协议实现
* 完善对event的处理，提高响应灵敏度
* 唤醒机制
