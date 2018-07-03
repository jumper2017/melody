package network

//定义peer接口，　表示的是各种协议的connector/acceptor，　考虑其由多个部分来组成：
//1. 协议相关属性，　2. 协议相关通信过程（connect / accept），　3.　错误处理策略（发送/接受错误）