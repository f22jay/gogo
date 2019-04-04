#!/usr/bin/env bash
################################################################################
#
# All Rights Reserved
#
################################################################################
#Breif:
#Authors: zhangfangjie (f22jay@163.com)
#Date:    2019/04/01 19:21:31
# kill process
ps -ef | grep usersvr | awk '{print "kill -9 " $2}' | sh
ps -ef | grep authsvr | awk '{print "kill -9 " $2}' | sh
ps -ef | grep httpsvr | awk '{print "kill -9 " $2}' | sh
nohup go run user/test/usersvr_main.go > user.log &
nohup go run auth/test/authsvr_main.go > auth.log &
nohup go run httpsvr.go > http.log &
