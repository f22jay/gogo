#!/usr/bin/env bash
################################################################################
#
# All Rights Reserved
#
################################################################################
#Breif:
#Authors: zhangfangjie (f22jay@163.com)
#Date:    2019/04/01 19:52:58

# update
curl -F "username=jack" -F "nickname=哈哈哈" --cookie  "token=test" -F "image=@httpsvr.go" http://127.0.0.1:10086/user?username=jack

# get
curl  --cookie  "token=test" http://127.0.0.1:10086/user?username=test

# bench get usser
ab -n100000 -c200 -k -C token=test   http://127.0.0.1:10086/user?username=test
ab -n100000 -c2000 -k -C token=test   http://127.0.0.1:10086/user?username=test
