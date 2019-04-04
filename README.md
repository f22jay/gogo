# entry task
学习golang，新人联系项目

## rpc
和c++的rpc相比，要容易很多，不需要关注事件驱动，并发，通信也特别容易
* 借鉴腾讯taf: server模块长连接管理，client调用方式、svr dispatcher
* 借鉴腾讯svrkit: client工作模式借鉴了svrkit的worker pool，不过svrkit server使用的是worker pool。这里client端使用worker pool。

### 性能测试
在本机mac上做了几组pingpong测试，qps达到18000多，client占用cpu450%,server占用240%,可能是协程调度原因,有待profile追究。
* 进入src/auth/test 目录
* 分别go run authsvr_main.go, go run authcli_main.go


## entry task性能测试
* . run.sh, 开启httpsvr，逻辑服务
* . httptest.sh 进行ab压测。ab只管发请求，模拟随机用户是在httpsvr端模拟的，如果是test用户，真实的username是随机生成的

在mac本机的虚拟机上测试

 qps有波动，在3000-4000之间

| 并发数 | 请求数 |  qps(随机用户)| qps(固定用户)|
| ------ | ------ | ------ |------ |
| 200| 100000 | 3295.14 |3693.98 |
| 2000 | 100000 | 3541.68 |3773.31 |


