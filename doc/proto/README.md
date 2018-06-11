
# 仓储管理服务 接口文档

## 所有用到的枚举值

```c
// API接口请求返回状态码
enum ReturnStatus {
	ReturnStatusFailed       =   0, // 操作失败
	ReturnStatusSuccess      =   1, // 操作成功
	ReturnStatusTempMoved    = 302, // 请求位置被临时移动，需要访问protocol字段指定的新位置
	ReturnStatusTempRedirect = 307, // 需要重定向到protocol字段指定的新位置
	ReturnStatusNotAuthorize = 401, // 请求要求身份验证
	ReturnStatusNeedLogin    = 418, // 该请求需要登录访问
};
```

```c
// 排序方向枚举
enum SortDirection {
	SortDirectAsc  =  "asc", // 升序排列
	SortDirectDesc = "desc", // 降序排列
};
```

```c
enum AuthStatus {
	AuthStatusNone = 0, // 无需授权验证
	AuthStatusNeed = 1, // 需要授权验证
};
```



## HTTP API接口 请求/响应参数
访问方式为请求路径加请求表单方式传递参数
每个请求参数结构的第一层每个字段对应一个表单参数
请求或响应参数为空的，表示无需请求数据或没有响应数据

API请求接口使用统一的返回JSON结构：

```js
{
	"status": 0, // 返回码(参见枚举 ReturnStatus)
	"info": "value", // 返回码描述
	"protocol": {} // 协议数据
}
```

### API接口'/public/treemenu' 获取菜单目录树数据
 访问方法: *

```js
// 请求表单参数示例
无需参数

// 返回JSON示例(仅Response.Protocol部分)
[
	{
		"id": 0, // 节点ID
		"pid": 0, // 父节点ID
		"text": "value", // 显示文本
		"iconCls": "value", // 图标
		"attributes": {
			"url": "value" // 节点对应的页面url
		},
		"children": [] // 子节点列表，与父节点结构相同
	}
]
```

### API接口'/public/isloggedin' 查询当前浏览器的登录状态
 访问方法: *

```js
// 请求表单参数示例
无需参数

// 返回JSON示例(仅Response.Protocol部分)
{
	"status": false, // 登录状态， true为已登录；false为未登录
	"redirect": "value" // 登录成功状态应该转向的页面
}
```

### API接口'/public/login' 提交登录数据
 访问方法: *

```js
// 请求表单参数示例
username=xxx&password=xxx
// 请求表单参数内容格式说明(参数值为json对象的，值对应为json对象的字符串内容)
{
	"username": "value", // 必填项，用户名
	"password": "value" // 必填项，用户密码（需要对用户原始明文密码进行MD5计算后填充）
}

// 返回JSON示例(仅Response.Protocol部分)
{
	"user": {
		"id": 0, // 用户ID
		"nickname": "value", // 用户昵称
		"mobile": 0, // 用户手机
		"remark": "value", // 备注
		"role": "value", // 用户角色
		"last_login_time": "value", // 最后登录时间，格式为2017-12-05T15:48:31+08:00
		"created_time": "value" // 创建时间，格式为2017-12-05T15:48:31+08:00
	},
	"redirect": "value" // 登录成功状态应该转向的页面
}
```

### API接口'/public/changepwd' 处理修改当前用户密码的请求
 访问方法: *

```js
// 请求表单参数示例
old=xxx&new=xxx&repeat=xxx
// 请求表单参数内容格式说明(参数值为json对象的，值对应为json对象的字符串内容)
{
	"old": "value", // 必填项，旧密码
	"new": "value", // 必填项，新密码
	"repeat": "value" // 必填项，重复新密码
}

// 返回JSON示例(仅Response.Protocol部分)
本接口不返回json格式数据
```

### API接口'/public/accesslist' 获取当前仓储管理系统定义的权限列表
 访问方法: *

```js
// 请求表单参数示例
无需参数

// 返回JSON示例(仅Response.Protocol部分)
[
	{
		"id": 0, // 权限节点ID
		"title": "value", // 权限节点显示标题
		"path": "value", // 权限节点关联的路径
		"level": 0, // 节点在权限树中的层次
		"pid": 0, // 父节点ID
		"auth": 0, // 授权需要，参见AuthStatus定义
		"icon": "value", // 图标
		"remark": "value", // 备注
		"children": [] // 子节点列表，与父节点结构相同
	}
]
```

### API接口'/admin/user/list' 获取仓储管理系统用户列表数据
 访问方法: *

```js
// 请求表单参数示例
page=xxx&rows=xxx&sort=xxx&order=xxx
// 请求表单参数内容格式说明(参数值为json对象的，值对应为json对象的字符串内容)
{
	"page": 0, // 请求分页数据的第几页数据，必填项
	"rows": 0, // 每页数据的分页条数，必填项
	"sort": "value", // 排序字段名，取值为下发的邮件服务配置数据的字段，可选项
	"order": "value", // 排序方向，参见SortDirection定义，可选项
}

// 返回JSON示例(仅Response.Protocol部分)
{
	"total": 0, // 总的管理用户数
	"rows": [
	{
			"id": 0, // 用户ID
			"nickname": "value", // 用户昵称
			"mobile": 0, // 用户手机
			"remark": "value", // 备注
			"role": "value", // 用户角色
			"last_login_time": "value", // 最后登录时间，格式为2017-12-05T15:48:31+08:00
			"created_time": "value" // 创建时间，格式为2017-12-05T15:48:31+08:00
		}
	] // 管理用户列表指定页数据
}
```

### API接口'/admin/user/add' 提交添加一个新用户的请求
 访问方法: *

```js
// 请求表单参数示例
insert=xxx
// 请求表单参数内容格式说明(参数值为json对象的，值对应为json对象的字符串内容)
{
	"insert": {
		"mobile": 0, // 用户手机，必填项
		"nickname": "value", // 用户昵称
		"password": "value", // 用户密码
		"remark": "value" // 备注
	}
}

// 返回JSON示例(仅Response.Protocol部分)
本接口不返回json格式数据
```

### API接口'/admin/user/update' 提交更新一个用户的请求
 访问方法: *

```js
// 请求表单参数示例
update=xxx
// 请求表单参数内容格式说明(参数值为json对象的，值对应为json对象的字符串内容)
{
	"update": {
		"id": 0, // 用户ID，必填项
		"nickname": "value", // 用户昵称
		"password": "value", // 用户密码
		"mobile": "value", // 用户手机
		"remark": "value", // 备注
		"status": 0, // 用户状态，0为禁用，1为启用
		"access": [0] // 用户权限ID列表
	}
}

// 返回JSON示例(仅Response.Protocol部分)
本接口不返回json格式数据
```

### API接口'/admin/user/delete' 提交删除一个用户的请求
 访问方法: *

```js
// 请求表单参数示例
uid=xxx
// 请求表单参数内容格式说明(参数值为json对象的，值对应为json对象的字符串内容)
{
	"uid": 0 //  需要删除的用户ID
}

// 返回JSON示例(仅Response.Protocol部分)
本接口不返回json格式数据
```

### API接口'/admin/user/access' 获取指定用户的权限ID列表
 访问方法: *

```js
// 请求表单参数示例
uid=xxx
// 请求表单参数内容格式说明(参数值为json对象的，值对应为json对象的字符串内容)
{
	"uid": 0 //  需要查询用户权限的用户ID
}

// 返回JSON示例(仅Response.Protocol部分)
{
	"list": [0] // 指定用户的权限ID列表
}
```

### API接口'/admin/operations/list' 获取管理系统操作日志列表数据
 访问方法: *

```js
// 请求表单参数示例
page=xxx&rows=xxx&sort=xxx&order=xxx
// 请求表单参数内容格式说明(参数值为json对象的，值对应为json对象的字符串内容)
{
	"page": 0, // 请求分页数据的第几页数据，必填项
	"rows": 0, // 每页数据的分页条数，必填项
	"sort": "value", // 排序字段名，取值为下发的邮件服务配置数据的字段，可选项
	"order": "value", // 排序方向，参见SortDirection定义，可选项
}

// 返回JSON示例(仅Response.Protocol部分)
{
	"total": 0, // 总的操作日志记录数
	"rows": [
	{
			"id": 0, // 操作记录ID
			"user_id": 0, // 操作用户ID
			"username": "value", // 操作用户名
			"from": "value", // 操作用户的来源IP
			"action": "value", // 执行的动作简介
			"detail": "value", // 操作详情
			"created_time": "value" // 操作发生的时间，格式为2017-12-05T15:48:31+08:00
		}
	] // 操作日志列表指定页数据
}
```

### API接口'/location/depots/list' 获取仓库数据
 访问方法: *

```js
// 请求表单参数示例
page=xxx&rows=xxx&sort=xxx&order=xxx
// 请求表单参数内容格式说明(参数值为json对象的，值对应为json对象的字符串内容)
{
	"page": 0, // 请求分页数据的第几页数据，必填项
	"rows": 0, // 每页数据的分页条数，必填项
	"sort": "value", // 排序字段名，取值为下发的邮件服务配置数据的字段，可选项
	"order": "value", // 排序方向，参见SortDirection定义，可选项
}

// 返回JSON示例(仅Response.Protocol部分)
{
	"total": 0, // 总的操作日志记录数
	"rows": [
	{
			"id": 0,
			"name": "value",
			"detail": "value",
			"shelfs": [
	{
					"id": 0,
					"name": "value",
					"layers": 0,
					"detail": "value"
				}
			]
		}
	] // 操作日志列表指定页数据
}
```

### API接口'/location/depots/add' 提交添加一个新仓库的请求
 访问方法: *

```js
// 请求表单参数示例
insert=xxx
// 请求表单参数内容格式说明(参数值为json对象的，值对应为json对象的字符串内容)
{
	"insert": {
		"id": 0,
		"name": "value",
		"detail": "value",
		"shelfs": [
	{
				"id": 0,
				"name": "value",
				"layers": 0,
				"detail": "value"
			}
		]
	}
}

// 返回JSON示例(仅Response.Protocol部分)
本接口不返回json格式数据
```

### API接口'/location/depots/update' 提交更新一个仓库的请求, 仅更新仓库属性，不更新货架信息
 访问方法: *

```js
// 请求表单参数示例
update=xxx
// 请求表单参数内容格式说明(参数值为json对象的，值对应为json对象的字符串内容)
{
	"update": {
		"id": 0,
		"name": "value",
		"detail": "value",
		"shelfs": [
	{
				"id": 0,
				"name": "value",
				"layers": 0,
				"detail": "value"
			}
		]
	}
}

// 返回JSON示例(仅Response.Protocol部分)
本接口不返回json格式数据
```

### API接口'/location/depots/update' 提交删除仓库的请求
 访问方法: *

```js
// 请求表单参数示例
depot_ids=xxx
// 请求表单参数内容格式说明(参数值为json对象的，值对应为json对象的字符串内容)
{
	"depot_ids": [0] //  需要删除的仓库信息
}

// 返回JSON示例(仅Response.Protocol部分)
本接口不返回json格式数据
```

### API接口'/location/shelfs/list' 获取仓库货架数据
 访问方法: *

```js
// 请求表单参数示例
page=xxx&rows=xxx&sort=xxx&order=xxx&depot_id=xxx
// 请求表单参数内容格式说明(参数值为json对象的，值对应为json对象的字符串内容)
{
	"page": 0, // 请求分页数据的第几页数据，必填项
	"rows": 0, // 每页数据的分页条数，必填项
	"sort": "value", // 排序字段名，取值为下发的邮件服务配置数据的字段，可选项
	"order": "value", // 排序方向，参见SortDirection定义，可选项
	"depot_id": 0 // 货架所在的仓库ID
}

// 返回JSON示例(仅Response.Protocol部分)
{
	"total": 0, // 总的操作日志记录数
	"rows": [
	{
			"id": 0,
			"name": "value",
			"layers": 0,
			"detail": "value"
		}
	] // 操作日志列表指定页数据
}
```

### API接口'/location/shelfs/add' 提交添加一个新仓库货架的请求
 访问方法: *

```js
// 请求表单参数示例
depot_id=xxx&shelfs=xxx
// 请求表单参数内容格式说明(参数值为json对象的，值对应为json对象的字符串内容)
{
	"depot_id": 0,
	"shelfs": [
	{
			"id": 0,
			"name": "value",
			"layers": 0,
			"detail": "value"
		}
	] // 添加的仓库信息
}

// 返回JSON示例(仅Response.Protocol部分)
本接口不返回json格式数据
```

### API接口'/location/shelfs/update' 提交更新一个仓库货架的请求
 访问方法: *

```js
// 请求表单参数示例
update=xxx
// 请求表单参数内容格式说明(参数值为json对象的，值对应为json对象的字符串内容)
{
	"update": {
		"id": 0,
		"name": "value",
		"layers": 0,
		"detail": "value"
	}
}

// 返回JSON示例(仅Response.Protocol部分)
本接口不返回json格式数据
```

### API接口'/location/shelfs/delete' 提交删除仓库货架的请求
 访问方法: *

```js
// 请求表单参数示例
shelf_ids=xxx
// 请求表单参数内容格式说明(参数值为json对象的，值对应为json对象的字符串内容)
{
	"shelf_ids": [0] //  需要删除的仓库货架信息
}

// 返回JSON示例(仅Response.Protocol部分)
本接口不返回json格式数据
```



## 页面请求

### 访问'/', '/pages/index'可以获取首页页面
 访问方法: *

```js
访问指定路径可以获取指定html页面
```

### 访问'/pages/login'可以获取登录页面
 访问方法: *

```js
访问指定路径可以获取指定html页面
```

### 访问'/pages/admin/user'可以获取管理后台用户列表页面
 访问方法: *

```js
访问指定路径可以获取指定html页面
```

### 访问'/pages/admin/operations'可以获取管理后台操作日志列表页面
 访问方法: *

```js
访问指定路径可以获取指定html页面
```


