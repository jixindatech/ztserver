## ztserver
这是零信任网关的主要组件，包含一个chrome插件，用来实现打开网络通信的功能，同时提供了web接口进行网关管理告警配置等。欢迎issue和star!

## 基本框架

## 安装
- 环境依赖：需要vue和Go环境。
- make 即可编译所有文件, make run 即可运行。
- 支持Docker, 编译镜像: docker build -t ztserver . (. 代表当前目录)。
- ztserver 数据存储在es中，需要配置es index(mapping.json 在server/doc/ 目录，需要自己优化)。
- ztserver会联动网关防火墙，需要提前在网关启用firewall 和ipset， 建立的ipset名称在配置文件中配置。
- ztserver会开启websocket，前面必须配置一个web网关(可配置websocket的保护策略)，在 x-forwarded-for 头部字段中添加客户端的IP。

## 配置说明
```
# server mode
mode: release   #debug 
# Length of secret Must be 32, encrypt data for token
secret: 'ch@ngfengpol@nghuiyoushi!!!!!!!!'
# the web api server
web:
addr : 0.0.0.0:8000
user : admin
password : 111111
jwtkey: 123456sdfq345sdfaf
identitykey: 'id'

#log_path : ./logs
# log level[debug|info|warn|error],default error
log_level : debug
# db_path for sqlite3
db_path: ztserver.db
# ssh for gateway hosts
ssh:
-
fw: firewall  # firewall type, if use iptables, you would add fw type for iptables
ipset: test  # ipset
host: 192.168.91.100:22 # ip:port for ssh
user: root # ssh user
password: 123456 # ssh password
server: # servers proxied by this host
- hr.jixindatech.com
- b.jixindatech.com
-
fw: firewall
ipset: test
host: 127.0.0.1:22
user: root
password: 123456
server:
    - c.jixindatech.com
# redis for data cache
redis:
addr: 192.168.91.100:6379
password: 123456
db: 0

es:
host:
- http://192.168.91.100:9200
username: admin
password: guessme
gwindex: gwindex
wsindex: wsindex
```
- web/user 和web/password 是web管理用户， web/jwtkey 是jwt的key ,建议修改这三项
- ztserver 需要联动防火墙开启和关闭端口，只支持firewall, 目前经过centos8测试， 
  其中需要修改 ssh/ipset,ztserver联动防火墙的ip集合用的是ipset,所以要提前在网关上建立ipset，
  ipset的名称在此处指定(firewall建立ipset，需要指定--permanent)，比如(port是开放的资源端口): 
  rule family="ipv4" source ipset="test" port port="8080" protocol="tcp" accept
  rule family="ipv4" source ipset="test" port port="8443" protocol="tcp" accept
  ssh/host 和ssh/user是
  ztserver登陆网关的ssh用户凭证，ssh/server是该网关可以允许访问的web资源，也就是域名，目前不支持
  通配符，需要跟管理后台配置的资源一致(后台管理的资源管理菜单配置域名)。
- redis 需要配置，所有的数据都会存储在redis里面。
- es 需要配置，chrome插件的日志存储在es里面，同时 web端会通过es查询网关的日志记录。gwindex 和wsindex 
分别是网关日志和插件日志的index。

## chrome-extension 配置
- 配置manifest.json 中的permissions，需要将资源的域名加入，比如 *://xxxx.com/* ，第一个*http https均可，
xxx.com 为实际的域名，第二个*为uri，具体可参考chrom插件资料。
- 配置 js/const.js 中的 WS_URI，连接ztserver的websocket。
- 插件需要配置用户的email和token。

## 其他
参考 doc 下面的说明文档。

## Contributing
PRs accepted.

## Discussion Group
QQ群: 254210748

##License
Unlicense


