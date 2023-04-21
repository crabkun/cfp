# cfp

## 说明
一键转发指定CF节点并且限制指定host防止被扫，注意本工具只支持http或者websocket协议。

## 防止被扫的意义是？
我发现很多人在自己的vps/nat vps上面开放端口转发Cloudflare，并且没做任何限制。这样做其实有非常大的安全隐患，直接转发Cloudflare相当于你的VPS变成了Cloudflare节点。
人人都可以拿你的端口套到他们的host上，跑的是你的流量，俗称白嫖。

## 使用方法
```
cfp
  -addr
      本地监听的地址，默认0.0.0.0代表全网卡
  -port
    	本地监听的端口，后续客户端连接这个端口
  -cfaddr
    	cloudflare节点地址，可以不带端口（默认80）
  -hosts
    	host白名单，也就是v2客户端上面的ws host啦。不在白名单内的host会被断开连接。多个host用逗号分隔
  -debug
    	打开调试模式，会显示一大堆log
```
## 例子

在你的vps上面开放12345转发CF节点104.17.219.1，并且限制只能ntt.mydomain.com和hkt.mydomain.com这两个域名(ws host)访问
```bash
./cfp -addr 0.0.0.0 -port 12345 -cfaddr 104.17.219.1 -hosts ntt.mydomain.com,hkt.mydomain.com
```

想要后台运行，可以使用nohup命令
```bash
nohup ./cfp -addr 0.0.0.0 -port 12345 -cfaddr 104.17.219.1 -hosts ntt.mydomain.com,hkt.mydomain.com &
```
PS:想要弄开机自启可以自己查资料动手~

## 遇到了问题？
欢迎提issue或Pull Request

## 开源协议
BSD 3-Clause License
