#!/bin/bash

set -e

# 协议代码和文档生成工具
#
# 使用方法:
#
# 1. 安装 npm, 然后使用 npm 安装 github-markdown
#
#	$ sudo npm install -g github-markdown
#
# 2. protogen 的使用
#
#	Options:
#	
#	  -h, --help                        显示帮助
#	  -i[=proto.go]                     输入源文件,多个文件以逗号分隔
#	  -v[=i]                            日志等级: t,d,i,w,e,f
#	      --readme-file[=./README.md]   生成的 README.md 文档文件目录
#
#	输入源文件为go语言源代码文件,而且只允许有 常量(const),结构体(struct)，不能有其他包导入依赖
#
#	1) 常量使用以 `enum: XXX` 开头的注释来定义一个枚举类型 XXX,如下定义了一个用户状态的枚举 UserStatus
#
#		// enum: UserStatus
#		// 用户的状态
#		const (
#			OrderStatusDisable          = 0 // 禁用
#			OrderStatusEnable           = 1 // 启用
#		)
#
#	2) 结构体使用如下的方式声明：
#   使用多行注释：
#       以结构名开头的行概述结构使用目的；
#       以“path:”开头的行表示该数据结构所适用的接口路径；
#       以“method:”开头的行描述访问方法，‘*’表示任一方法；
#       以“redirect:”开头的行表示该接口访问可能会产生重定向动作或响应
#	一般访问模式分别用 `Req` 和 `Res` 结尾，分别表示请求和响应；Ajax模式访问则分别用 `ReqAjax` 和 `ResAjax` 结尾； 同一接口请求和响应按照相同前缀尽量成对组织
#	请求协议都必须继承 `ReqCommon`
#	比如如下的协议 LoginReq 和 LoginReq:
#
#		// LoginReq 获取登录页面
#		// path: '/public/login'
#		// redirect: 如果是已登录状态将重定向到'/public/index'
#		// method: *
#		type LoginReq struct {
#			ReqCommon
#		}
#		
#		// LoginRes is the response for login request
#		type LoginRes struct {
#		}
#

if [ ! -f "$0" ]; then
    echo 'protocol generation must be run within its container folder' 1>&2
    exit 1
fi

echo "refresh protocol tools"
CURDIR=`pwd`
go build -o ../../../../../../bin/tools/protogen github.com/lucifinil-long/stores/tools/protogen

echo "gofmt files"
gofmt -w ../../proto/proto.go

echo "protogen -i ../../proto/proto.go"
../../../../../../bin/tools/protogen -i ../../proto/proto.go

echo "github-markdown ./README.md"
github-markdown ./README.md