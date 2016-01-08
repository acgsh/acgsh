[![Build Status](https://travis-ci.org/acgshare/acgsh.svg?branch=master)](https://travis-ci.org/acgshare/acgsh)

ACGSH
==================

安装
----

首先请安装好Golang和Twister，并保证Twister已同步完成，而且已建立一个用户账号。

下载编译ACGSH

    go get github.com/acgshare/acgsh
    
建立acgsh工作目录

    mkdir acgsh
    cd acgsh
    
下载acgsh_html目录

    git clone https://github.com/acgshare/acgsh_html.git
    
复制配置文件到acgsh工作目录

    cp $GOPATH/src/github.com/acgshare/acgsh/acgsh.conf ./
    
修改acgsh.conf配置文件

    TwisterServer="http://user:pwd@127.0.0.1:28332" #请根据你的twister core设置修改
    TwisterUsername="your_twister_user_name"        #修改这里为你建立的twister用户账号
    HttpServerPort = "8080"                         #http服务端口号
    
运行ACGSH

    acgsh
    
acgsh会在当前目录下建立acgsh_bolt.db数据库文件并同步你的twister账号follow的账号内容到数据库。    
    
推荐使用[Supervisor](http://supervisord.org/) 来使acgsh在后台工作并监控acgsh工作状态。
    
    
