// generated automatically when refresh database models
// don't modify the code manually, changes might be lost

package db

// StoresCommoditySpec is the database model struct
type StoresCommoditySpec struct {
	Id            int64  `json:"id" xorm:"pk autoincr BIGINT(20)"`
	SpecName      string `json:"spec_name" xorm:"not null VARCHAR(64)"`
	Detail        string `json:"detail" xorm:"not null VARCHAR(128)"`
	Segmentable   int    `json:"segmentable" xorm:"not null default 0 TINYINT(4)"`
	SegmentId     int64  `json:"segment_id" xorm:"BIGINT(20)"`
	SegmentAmount int    `json:"segment_amount" xorm:"INT(11)"`
}
