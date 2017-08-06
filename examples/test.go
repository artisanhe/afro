package modules

import (
	"time"
)

//@IgnoreTableName
type UserTest struct {
	Id uint64 `gorm:"primary_key;column:F_id" sql:"type:bigint(64) unsigned auto_increment;not null" json:"-"`
	////Id1        uint64                 `gorm:"primary_key;column:F_id1" sql:"type:bigint(64) unsigned auto_increment;not null" json:"-"`
	Phone      string                 `gorm:"column:F_phone" sql:"type:varchar(16);not null;unique_index:I_phone" json:"phone"`
	UserID     uint64                 `gorm:"column:F_user_id" sql:"type:bigint(64) unsigned;not null;unique_index:I_user_id" json:"user_id"`
	Name       string                 `gorm:"column:F_name" sql:"type:varchar(32) ;not null;unique_index:I_user_id" json:"name"`
	CityID     uint64                 `gorm:"column:F_city_id" sql:"type:bigint(64) ;not null;index:I_city_id;default:0" json:"city_id"`
	AreaID     uint64 `gorm:"column:F_area_id" sql:"type:bigint(64) ;not null;index:I_city_id;default:0" json:"area_id"`
	Enabled    uint8                  `gorm:"column:F_enabled" sql:"type:tinyint(8) unsigned;not null;default:1" json:"-"`
	CreateTime time.Time              `gorm:"column:F_create_time" sql:"type:bigint(64) unsigned;not null;default:0" json:"-"`
	UpdateTime time.Time `gorm:"column:F_update_time" sql:"type:bigint(64) unsigned;not null;default:0" json:"-"`
}

func (ut UserTest) TableName() string {
    return "db.user_test"
}

type CustomerG7s struct {
	// 客户id
	CustomerID uint64 `json:"customerID,string" gorm:"primary_key;column:F_customer_id" sql:"type:bigint(64) unsigned;not null;unique_index:I_user_customer[1]"`
	// g7s机构
	G7sOrgCode string `json:"orgCode" gorm:"primary_key;column:F_g7s_org_code" sql:"type:varchar(32);not null`
	// g7s用户id
	G7sUserID string `json:"userID,string" gorm:"primary_key;column:F_g7s_user_id" sql:"type:varchar(32);not null;unique_index:I_user_customer[0]"`
	// 记录创建时间
	CreateTime time.Time `json:"createTime" gorm:"column:F_create_time" sql:"type:bigint(64);not null;index:I_create_time"`
	// 记录更新时间
	UpdateTime time.Time `json:"updateTime" gorm:"column:F_update_time" sql:"type:bigint(64);not null;index:I_update_time"`
	Enabled    uint8 `json:"-" gorm:"primary_key;column:F_enabled" sql:"type:tinyint(8) unsigned;not null;default:1;unique_index:I_user_customer[2]"`
}

type PhysicsDeleteByUniquustomerG7s struct {
	// g7s用户id
	G7sUserID string `json:"userID,string" gorm:"primary_key;column:F_g7s_user_id" sql:"type:varchar(32);not null"`
	//// 客户id
	CustomerID uint64 `json:"customerID,string" gorm:"primary_key;column:F_customer_id" sql:"type:bigint(64) unsigned;not null;index:I_customer"`
	// g7s机构
	G7sOrgCode string `json:"orgCode" gorm:"column:F_g7s_org_code" sql:"type:varchar(32);not null;index:I_org_code"`
	// 记录创建时间
	CreateTime time.Time `json:"createTime" gorm:"column:F_create_time" sql:"type:bigint(64);not null;index:I_create_time"`
	// 记录更新时间
	UpdateTime time.Time `json:"updateTime" gorm:"column:F_update_time" sql:"type:bigint(64);not null;index:I_update_time"`
	Enabled    uint8 `json:"-" gorm:"primary_key;column:F_enabled" sql:"type:tinyint(8) unsigned;not null;default:1"`
}
