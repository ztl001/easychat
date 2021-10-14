# easychat
1. 广播上线
2. 广播消息
3. 查询当前用户列表
4. 修改用户名
5. 用户主动退出
6. 超时处理功能

# 基本步骤
```shell
1. 主要携程处理用户链接
	1.1 将用户写入map
	1.2 告诉所有在线用户上线信息
	message <- 信息

2. go 发送消息，参数cli
	for msg := range cli.C{
		write(msg)
	}

3. 将消息写入msg
	for{
		msg := <- message
		for _,cli := range onlineMap{
			cli.C <- msg
		}
	}

4. 接收用户的请求，把用户发过来的数据转发
	message <- buf

	对方下线，把当前用户从map中移除
```

