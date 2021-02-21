// generated automatically when refresh database models
// don't modify the code manually, changes might be lost

package db

// StoresCommoditySpec is the database model struct
type StoresCommoditySpec struct {
	Id          int64  `json:"id" xorm:"pk autoincr BIGINT(20)"`
	SpecName    string `json:"spec_name" xorm:"not null VARCHAR(64)"`
	Detail      string `json:"detail" xorm:"not null VARCHAR(128)"`
	ParentId    int64  `json:"parent_id" xorm:"not null default 0 index BIGINT(20)"`
	Segmentable int    `json:"segmentable" xorm:"not null default 0 TINYINT(4)"`
	SubId       int64  `json:"sub_id" xorm:"default 0 index BIGINT(20)"`
	SubAmount   int    `json:"sub_amount" xorm:"default 0 INT(11)"`
}
