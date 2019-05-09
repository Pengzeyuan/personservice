# personservice
go实现的一个grpc调用的简单示例  
  
注:  
客户端代码:  
(1) personclient.go  
服务端代码:  
(1) personserver.go  
  
编译上述代码依赖protobuf和grpc库, 这些可以参考:  
protobuf安装: https://www.cnblogs.com/albizzia/p/10781028.html  
grpc安装: https://www.cnblogs.com/albizzia/p/10803032.html  
  
编译过程如下(在personservice文件夹):  
$ ./compile.sh  
其中生成的personclient和personserver分别是客户端和服务端, 使用./personclient和./personserver启动这两个应用就可以在person-client中看到相关输出.  
  
文档说明可以参考: https://www.cnblogs.com/albizzia/p/10836948.html  
