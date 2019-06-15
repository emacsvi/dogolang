# 利用docker运行
[镜像rabbitmq里面很详细](https://hub.docker.com/_/rabbitmq) 有不同的版本，这里主要使用提供管理界面的版本`rabbitmq:management`:
```bash
docker run -d --hostname rabbit-svr --name rabbit -p 5672:5672 -p 15672:15672 -p 25672:25672 -v /Users/liwei/g/coding/rabbit/data:/var/lib/rabbitmq rabbitmq:management
```
然后访问页面：[http://127.0.0.1:15672](http://127.0.0.1:15672) 用户名和密码都是`guest`即可登录。

