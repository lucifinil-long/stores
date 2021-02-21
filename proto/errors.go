package proto

import "errors"

var (
	// ErrReloginNeed is the error when user session is expired
	ErrReloginNeed = errors.New("用户会话已过期，需要重新登陆才能执行该操作")
	// ErrNotEnoughPasswordStrength indicates the password strength is not enough
	ErrNotEnoughPasswordStrength = errors.New("密码强度不够，请使用至少5位长度以上，并且包含特殊字符和数字及大小写混合的密码")

	// ErrUserNotExist is the error for user is not exists
	ErrUserNotExist = errors.New("用户不存在")
	// ErrUserDisabled is the error for user is disable
	ErrUserDisabled = errors.New("用户处于被禁用状态")
	// ErrUserWrongPwd is the error for wrong password
	ErrUserWrongPwd = errors.New("密码错误")
	// ErrCanNotDeleteSelf is the error for in case user wants to delete self
	ErrCanNotDeleteSelf = errors.New("用户不能自己删除自己")
	// ErrCanNotUpdateHighLevelUser is the error for low level user wants to update high level user's information
	ErrCanNotUpdateHighLevelUser = errors.New("低级用户无权更新高级用户信息")
)

var (
	// ErrCanNotDeleteDepotsForStock is the error for delete depots failed for stock
	ErrCanNotDeleteDepotsForStock = errors.New("待删除的仓库尚有库存，请清空库存后再删除")
	// ErrDepotIsNotExist is the error for the depot is not exist
	ErrDepotIsNotExist = errors.New("指定仓库不存在")
	// ErrEmptyShelfName is the error for empty shelf name
	ErrEmptyShelfName = errors.New("货架名不能为空")
	// ErrCanNotDeleteShelfsForStock is the error for delete depots failed for stock
	ErrCanNotDeleteShelfsForStock = errors.New("待删除的仓库货架尚有库存，请清空库存后再删除")
)

var (
	// ErrParentSpecIsNotExist is the error for parent spect is not exist
	ErrParentSpecIsNotExist = errors.New("指定父级规格不存在")
	// ErrParentHasChild is the error for parent spect has child yet
	ErrParentHasChild = errors.New("指定父级规格已经存在子级规格")
	// ErrChildSpecIsNotExist is the error for child spect is not exist
	ErrChildSpecIsNotExist = errors.New("指定子级规格不存在")
	// ErrChildHasParent is the error for child spect has parent yet
	ErrChildHasParent = errors.New("指定父级规格已经存在子级规格")
	// ErrDeleteSkuBeforeSpec is the error that indicates must delete related sku before spec
	ErrDeleteSkuBeforeSpec = errors.New("指定规格已经关联到SKU，请先删除关联SKU再删除规格")
)

var (
	// ErrCommonInvalidParam is common invalid parameter error
	ErrCommonInvalidParam = errors.New("错误参数")
	// ErrCommonInternalError is common internal error
	ErrCommonInternalError = errors.New("内部错误")
)

var (
	// ErrDatabase is common database error
	ErrDatabase = errors.New("数据库操作错误。")
	// ErrDupKey is duplicate key record error
	ErrDupKey = errors.New("已经存在相同关键字的记录。")
	// ErrDataTooLong is data too long error
	ErrDataTooLong = errors.New("内容长度超过字段长度限制。")
	// ErrNotStructType is error for object is not struct type
	ErrNotStructType = errors.New("非结构类型。")
	// ErrInvalidValueForNutNullField is the error for there is nil for not null field
	ErrInvalidValueForNutNullField = errors.New("非空字段必须填写有效值。")
	// ErrRecordIsNotExist is error for record is not exist
	ErrRecordIsNotExist = errors.New("记录不存在")
)
