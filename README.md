# qyvx
企业微信hook的后端

```
.
├── conf
│   └── config.yaml	    # 相关配置文件
├── controller
│   ├── code.go	        # 相关返回体状态码设定
│   ├── qyvx.go 				# gin框架的Handler句柄函数
│   ├── response.go			# 返回体设定
│   └── validator.go		# 返回内容的转换
├── dao	# 数据库
│   └── mysql	# 此处代指MOC
│       ├── error.go 		# mysql相关特殊报错设置
│       ├── mysql.go 		# 初始化mysql
│       └── qyvx.go  		# 关于qyvx在mysql当中的数据
├── logger
│   └── logger.go				# log初始化
├── logic
│   ├── errors.go				# 具体业务逻辑处理函数中的特殊error
│   ├── qyvx.go					# 具体业务逻辑处理函数
│   └── request.go			# api接口调用的统一书写文件(密钥存放在此处)
├── models
│   ├── create_tables.sql# 数据库中表结构设计sql语句
│   ├── msgcontent.go		# 企业微信hook的信息结构体
│   └── param.go				# 发送邀请的参数结构体
├── pkg									# 第三方包及自定义工具
│   ├── ihttp						# http官方库的封装
│   │   └── ihttp.go		# 封装http库，满足请求头，token等值的添加
│   ├── snowflask				# 雪花算法
│   │   └── snowflask.go# 雪花算法获取ID
│   ├── utools
│   │   └── github.go		# 关于github相关的处理函数(后续考虑迁移至logic/requests.go中)
│   └── wxbizmsgcrypt
│       └── wxbizmsgcrypt.go	# 关于企业微信信息校验，解码，加密的官方提供库
├── router
│   └── route.go				# 路由设置
├── settings
│   └── settings.go			# config.yaml配置加载设置
├── go.mod
├── go.sum
├── main.go							# 程序入口，进行一系列的初始化和启动
└── README.md
```

