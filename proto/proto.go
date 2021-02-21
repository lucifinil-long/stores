package proto

/*
 *	go语言源代码文件,而且只允许有 常量(const),结构体(struct)，不能有其他包导入依赖
 *
 *	1) 常量使用以 `enum: XXX` 开头的注释来定义一个枚举类型 XXX,如下定义了一个用户状态的枚举 UserStatus
 *
 *		// enum: UserStatus
 *		// 用户的状态
 *		const (
 *			OrderStatusDisable          = 0 // 禁用
 *			OrderStatusEnable           = 1 // 启用
 *		)
 *
 *	2) 页面请求相关结构体使用如下的方式声明：
 *   使用多行注释：
 *       以结构名开头的行概述结构使用目的；
 *       以“path:”开头的行表示该数据结构所适用的接口路径；
 *       以“method:”开头的行描述访问方法，‘*’表示任一方法；
 *	结构名以`Page`结尾
 *
 *	3) API请求相关结构体使用如下的方式声明：
 *   使用多行注释：
 *       以结构名开头的行概述结构使用目的；
 *       以“path:”开头的行表示该数据结构所适用的接口路径；
 *       以“method:”开头的行描述访问方法，‘*’表示任一方法；
 *	API访问模式相关的结构名称应该分别用 `Req` 和 `Res` 结尾，分别表示请求和响应；同一接口请求和响应按照相同前缀成对组织
 *	请求协议都必须继承 `ReqCommon`
 *	比如如下的协议 LoginReq 和 LoginReq:
 *
 *		// LoginReq 提交登录请求
 *		// path: '/public/login'
 *		// method: *
 *		type LoginReq struct {
 *			ReqCommon
 *			UserName string `json:"username"` // 必填项，用户名
 *			Password string `json:"password"` // 必填项，用户密码（需要对用户原始明文密码进行MD5计算后填充）
 *		}
 *
 *		// LoginRes is the response for login request
 *		type LoginRes struct {
 *		}
 *
 */

// enum: ReturnStatus
// API接口请求返回状态码
const (
	ReturnStatusFailed       = 0   // 操作失败
	ReturnStatusSuccess      = 1   // 操作成功
	ReturnStatusTempMoved    = 302 // 请求位置被临时移动，需要访问protocol字段指定的新位置
	ReturnStatusTempRedirect = 307 // 需要重定向到protocol字段指定的新位置
	ReturnStatusNotAuthorize = 401 // 请求要求身份验证
	ReturnStatusNeedLogin    = 418 // 该请求需要登录访问
)

// Response 是非页面接口api 统一返回数据结构
type Response struct {
	Status      int         `json:"status"`   // 返回码(参见枚举 ReturnStatus)
	Description string      `json:"info"`     // 返回码描述
	Protocol    interface{} `json:"protocol"` // 协议数据
}

// ReqCommon is common base request struct
type ReqCommon struct {
}

// ResCommon is common base return struct
type ResCommon struct {
}

// enum: SortDirection
// 排序方向枚举
const (
	SortDirectAsc  = "asc"  // 升序排列
	SortDirectDesc = "desc" // 降序排列
)

// PageReqCommon store common page request parameters
type PageReqCommon struct {
	ReqCommon
	Page  int    `json:"page"`  // 请求分页数据的第几页数据，必填项
	Rows  int    `json:"rows"`  // 每页数据的分页条数，必填项
	Sort  string `json:"sort"`  // 排序字段名，取值为下发的邮件服务配置数据的字段，可选项
	Order string `json:"order"` // 排序方向，参见SortDirection定义，可选项
}

//
// 获取页面接口
//

// IndexPage 获取首页页面
// path: '/', '/pages/index'
// method: *
type IndexPage struct {
}

// LoginPage 获取登录页面
// path: '/pages/login'
// method: *
type LoginPage struct {
}

// AdminUsersPage 获取管理后台用户列表页面
// path: '/pages/admin/user'
// method: *
type AdminUsersPage struct {
}

// AdminOperationLogsPage 获取管理后台操作日志列表页面
// path: '/pages/admin/operations'
// method: *
type AdminOperationLogsPage struct {
}

//
// API 请求
//

// TreeMenuReq 获取菜单目录树数据
// path: '/public/treemenu'
// method: *
type TreeMenuReq struct {
	ReqCommon
}

// TreeMenuNode is struct for tree menu node
type TreeMenuNode struct {
	ID         int64           `json:"id"`         // 节点ID
	PID        int64           `json:"pid"`        // 父节点ID
	Text       string          `json:"text"`       // 显示文本
	IconCls    string          `json:"iconCls"`    // 图标
	Attributes Attributes      `json:"attributes"` // 属性
	Children   []*TreeMenuNode `json:"children"`   // 子节点列表，与父节点结构相同
}

// Attributes is attributes struct
type Attributes struct {
	URL string `json:"url"` // 节点对应的页面url
}

// TreeMenuRes is the response for TreeMenuReq request
type TreeMenuRes []TreeMenuNode

// IsLoggedInReq 查询当前浏览器的登录状态
// path: '/public/isloggedin'
// method: *
type IsLoggedInReq struct {
	ReqCommon
}

// IsLoggedInRes is the response for login request
type IsLoggedInRes struct {
	ResCommon
	Status   bool   `json:"status"`   // 登录状态， true为已登录；false为未登录
	Redirect string `json:"redirect"` // 登录成功状态应该转向的页面
}

// LoginReq 提交登录数据
// path: '/public/login'
// method: *
type LoginReq struct {
	ReqCommon
	UserName string `json:"username"` // 必填项，用户名
	Password string `json:"password"` // 必填项，用户密码（需要对用户原始明文密码进行MD5计算后填充）
}

// LoginRes is the response for login request
type LoginRes struct {
	ResCommon
	User     User   `json:"user"`     // 登录成功的用户信息
	Redirect string `json:"redirect"` // 登录成功状态应该转向的页面
}

// ModifyPasswordReq 处理修改当前用户密码的请求
// path: '/public/changepwd'
// method: *
// redirect: 未登录状态将重定向到'/public/login'
type ModifyPasswordReq struct {
	ReqCommon
	Old    string `json:"old"`    // 必填项，旧密码
	New    string `json:"new"`    // 必填项，新密码
	Repeat string `json:"repeat"` // 必填项，重复新密码
}

// ModifyPasswordRes is the response for ModifyPasswordReq
type ModifyPasswordRes struct {
	ResCommon
}

// enum: AuthStatus
const (
	AuthStatusNone = 0 // 无需授权验证
	AuthStatusNeed = 1 // 需要授权验证
)

// AccessTreeNode is struct for tree node of access
type AccessTreeNode struct {
	ID       int64             `json:"id"`       // 权限节点ID
	Title    string            `json:"title"`    // 权限节点显示标题
	Path     string            `json:"path"`     // 权限节点关联的路径
	Level    int               `json:"level"`    // 节点在权限树中的层次
	Pid      int64             `json:"pid"`      // 父节点ID
	Auth     int               `json:"auth"`     // 授权需要，参见AuthStatus定义
	Icon     string            `json:"icon"`     // 图标
	Remark   string            `json:"remark"`   // 备注
	Children []*AccessTreeNode `json:"children"` // 子节点列表，与父节点结构相同
}

// AccessListReq 获取当前仓储管理系统定义的权限列表
// path: '/public/accesslist'
// method: *
type AccessListReq struct {
	ReqCommon
}

// AccessListRes is the response for AccessListReq
type AccessListRes []AccessTreeNode

// User is the admin user entry
type User struct {
	ID            int64  `json:"id"`              // 用户ID
	Nickname      string `json:"nickname"`        // 用户昵称
	Mobile        int64  `json:"mobile"`          // 用户手机
	Remark        string `json:"remark"`          // 备注
	Role          string `json:"role"`            // 用户角色
	LastLoginTime string `json:"last_login_time"` // 最后登录时间，格式为2017-12-05T15:48:31+08:00
	CreatedTime   string `json:"created_time"`    // 创建时间，格式为2017-12-05T15:48:31+08:00
	Password      string // 用户密码
}

// UserListReq 获取仓储管理系统用户列表数据
// path: '/admin/user/list'
// method: *
type UserListReq struct {
	PageReqCommon
}

// UserListRes is the response for UserListReq
type UserListRes struct {
	ResCommon
	Total int64  `json:"total"` // 总的管理用户数
	Rows  []User `json:"rows"`  // 管理用户列表指定页数据
}

// NewUser is new user entry
type NewUser struct {
	Mobile   int64  `json:"mobile"`   // 用户手机，必填项
	Nickname string `json:"nickname"` // 用户昵称
	Password string `json:"password"` // 用户密码
	Remark   string `json:"remark"`   // 备注
	ID       int64  // 用户ID
}

// AddUserReq 提交添加一个新用户的请求
// path: '/admin/user/add'
// method: *
type AddUserReq struct {
	ReqCommon
	Insert NewUser `json:"insert"` // 添加的用户信息
}

// AddUserRes is the response for AddUserReq
type AddUserRes struct {
	ResCommon
}

// UserUpdate is user update entry
type UserUpdate struct {
	ID       int    `json:"id"`       // 用户ID，必填项
	Nickname string `json:"nickname"` // 用户昵称
	Password string `json:"password"` // 用户密码
	Mobile   string `json:"mobile"`   // 用户手机
	Remark   string `json:"remark"`   // 备注
	Status   int    `json:"status"`   // 用户状态，0为禁用，1为启用
	Access   []int  `json:"access"`   // 用户权限ID列表
}

// UpdateUserReq 提交更新一个用户的请求
// path: '/admin/user/update'
// method: *
type UpdateUserReq struct {
	ReqCommon
	Update UserUpdate `json:"update"` //  需要更新的用户信息
}

// UpdateUserRes is the response for UpdateUserReq
type UpdateUserRes struct {
	ResCommon
}

// DeleteUserReq 提交删除一个用户的请求
// path: '/admin/user/delete'
// method: *
type DeleteUserReq struct {
	ReqCommon
	UID int `json:"uid"` //  需要删除的用户ID
}

// DeleteUserRes is the response for DeleteUserReq
type DeleteUserRes struct {
	ResCommon
}

// UserAccessesReq 获取指定用户的权限ID列表
// path: '/admin/user/access'
// method: *
type UserAccessesReq struct {
	ReqCommon
	UID int `json:"uid"` //  需要查询用户权限的用户ID
}

// UserAccessesRes is the response for UserAccessesReq
type UserAccessesRes struct {
	ResCommon
	List []int `json:"list"` // 指定用户的权限ID列表
}

// OperationLog is the operation log entry
type OperationLog struct {
	ID          int64  `json:"id"`           // 操作记录ID
	UserID      int    `json:"user_id"`      // 操作用户ID
	Username    string `json:"username"`     // 操作用户名
	From        string `json:"from"`         // 操作用户的来源IP
	Action      string `json:"action"`       // 执行的动作简介
	Detail      string `json:"detail"`       // 操作详情
	CreatedTime string `json:"created_time"` // 操作发生的时间，格式为2017-12-05T15:48:31+08:00
}

// OperationLogListReq 获取管理系统操作日志列表数据
// path: '/admin/operations/list'
// method: *
type OperationLogListReq struct {
	PageReqCommon
}

// OperationLogListRes is the response for OperationLogListReq
type OperationLogListRes struct {
	ResCommon
	Total int64          `json:"total"` // 总的操作日志记录数
	Rows  []OperationLog `json:"rows"`  // 操作日志列表指定页数据
}

// Shelf is the shelf in depot
type Shelf struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Layers int    `json:"layers"`
	Detail string `json:"detail"`
}

// Depot is the depot
type Depot struct {
	ID     int64   `json:"id"`
	Name   string  `json:"name"`
	Detail string  `json:"detail"`
	Shelfs []Shelf `json:"shelfs"`
}

// DepotsListReq 获取仓库数据
// path: '/location/depots/list'
// method: *
type DepotsListReq struct {
	PageReqCommon
}

// DepotsListRes is the response for DepotsListReq
type DepotsListRes struct {
	ResCommon
	Total int64   `json:"total"` // 总的操作日志记录数
	Rows  []Depot `json:"rows"`  // 操作日志列表指定页数据
}

// AddDepotReq 提交添加一个新仓库的请求
// path: '/location/depots/add'
// method: *
type AddDepotReq struct {
	ReqCommon
	Insert Depot `json:"insert"` // 添加的仓库信息
}

// AddDepotRes is the response for AddDepotReq
type AddDepotRes struct {
	ResCommon
}

// UpdateDepotReq 提交更新一个仓库的请求, 仅更新仓库属性，不更新货架信息
// path: '/location/depots/update'
// method: *
type UpdateDepotReq struct {
	ReqCommon
	Update Depot `json:"update"` //  需要更新的仓库信息
}

// UpdateDepotRes is the response for UpdateDepotReq
type UpdateDepotRes struct {
	ResCommon
}

// DeleteDepotReq 提交删除仓库的请求
// path: '/location/depots/update'
// method: *
type DeleteDepotReq struct {
	ReqCommon
	DepotIds []int64 `json:"depot_ids"` //  需要删除的仓库信息
}

// DeleteDepotRes is the response for DeleteDepotReq
type DeleteDepotRes struct {
	ResCommon
}

// ShelfsListReq 获取仓库货架数据
// path: '/location/shelfs/list'
// method: *
type ShelfsListReq struct {
	PageReqCommon
	DepotID int64 `json:"depot_id"` // 货架所在的仓库ID
}

// ShelfsListRes is the response for ShelfsListReq
type ShelfsListRes struct {
	ResCommon
	Total int64   `json:"total"` // 总的操作日志记录数
	Rows  []Shelf `json:"rows"`  // 操作日志列表指定页数据
}

// AddShelfsReq 提交添加一个新仓库货架的请求
// path: '/location/shelfs/add'
// method: *
type AddShelfsReq struct {
	ReqCommon
	DepotID int64   `json:"depot_id"`
	Shelfs  []Shelf `json:"shelfs"` // 添加的仓库信息
}

// AddShelfsRes is the response for AddShelfsReq
type AddShelfsRes struct {
	ResCommon
}

// UpdateShelfReq 提交更新一个仓库货架的请求
// path: '/location/shelfs/update'
// method: *
type UpdateShelfReq struct {
	ReqCommon
	Update Shelf `json:"update"` //  需要更新的仓库货架信息
}

// UpdateShelfRes is the response for UpdateShelfReq
type UpdateShelfRes struct {
	ResCommon
}

// DeleteShelfReq 提交删除仓库货架的请求
// path: '/location/shelfs/delete'
// method: *
type DeleteShelfReq struct {
	ReqCommon
	ShelfIds []int64 `json:"shelf_ids"` //  需要删除的仓库货架信息
}

// DeleteShelfRes is the response for DeleteShelfReq
type DeleteShelfRes struct {
	ResCommon
}

// SpecEntry is the commodity specification
type SpecEntry struct {
	ID          int64  `json:"id"`          // 规格ID
	Name        string `json:"name"`        // 规格名称
	Detail      string `json:"detail"`      // 规格详情
	ParentID    int64  `json:"parent_id"`   // 父级规格ID
	Amount      int    `json:"amount"`      // 父级规格ID可以拆分为本级规格的数量
	Segmentable int    `json:"segmentable"` // 是否可拆分为其他规格
	SubID       int64  `json:"sub_id"`      // 可拆分的子级规格ID
	SubAmount   int    `json:"sub_amount"`  // 拆分为子级规格的数量
}

// Specification is the commodity specification
type Specification struct {
	ID          int64          `json:"id"`          // 规格ID
	Name        string         `json:"name"`        // 规格名称
	Detail      string         `json:"detail"`      // 规格详情
	Parent      *Specification `json:"parent"`      // 父级规格
	Segmentable int            `json:"segmentable"` // 是否可拆分为其他规格
	Sub         *Specification `json:"sub"`         // 子级规格
	SubAmount   int            `json:"sub_amount"`  // 可拆分为子级规格数量
}

// SpecificationListReq 获取规格数据
// path: '/admin/specifications/list'
// method: *
type SpecificationListReq struct {
	PageReqCommon
}

// SpecificationListRes is the response for SpecificationListReq
type SpecificationListRes struct {
	ResCommon
	Total int64           `json:"total"` // 总的操作日志记录数
	Rows  []Specification `json:"rows"`  // 操作日志列表指定页数据
}

// AddSpecificationReq 提交添加一个新规格的请求
// path: '/admin/specifications/add'
// method: *
type AddSpecificationReq struct {
	ReqCommon
	Insert SpecEntry `json:"insert"` // 添加的仓库信息
}

// AddSpecificationRes is the response for AddSpecificationReq
type AddSpecificationRes struct {
	ResCommon
}

// UpdateSpecificationReq 提交更新一个规格的请求
// path: '/admin/specifications/update'
// method: *
type UpdateSpecificationReq struct {
	ReqCommon
	Update SpecEntry `json:"update"` //  需要更新的规格信息
}

// UpdateSpecificationRes is the response for UpdateSpecificationReq
type UpdateSpecificationRes struct {
	ResCommon
}

// DeleteSpecificationReq 提交删除规格的请求
// path: '/admin/specifications/delete'
// method: *
type DeleteSpecificationReq struct {
	ReqCommon
	SpecDetailID []int64 `json:"spec_ids"` //  需要删除的规格ID列表
}

// DeleteSpecificationRes is the response for DeleteSpecificationReq
type DeleteSpecificationRes struct {
	ResCommon
}
