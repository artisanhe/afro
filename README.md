* 为了在新建module时不写过多的重复性代码，直接根据module里面每个字段的索引生成指定的代码，使用方法如下:  
	
	```
	cd /xxx/afro/
	make install
	afro -d /examples/
	afro -f test.go
	``` 
* 默认情况下会生成TableName()函数，如果不想生成默认的TableName()函数，可以在Model添加注释(@IgnoreTableName):

	```
	//@IgnoreTableName
  	type UserTest struct {
  	  Id uint64 `gorm:"primary_key;column:F_id" sql:"type:bigint(64) unsigned auto_increment;not null" json:"-"`
      Phone      string                 `gorm:"column:F_phone" sql:"type:varchar(16);not null;unique_index:I_phone" json:"phone"`
      UserID     uint64                 `gorm:"column:F_user_id" sql:"type:bigint(64) unsigned;not null;unique_index:I_user_id" json:"user_id"`
      Name       string                 `gorm:"column:F_name" sql:"type:varchar(32) ;not null;unique_index:I_user_id" json:"name"`
      CityID     uint64                 `gorm:"column:F_city_id" sql:"type:bigint(64) ;not null;index:I_city_id;default:0" json:"city_id"`
      AreaID     uint64                 `gorm:"column:F_area_id" sql:"type:bigint(64) ;not null;index:I_city_id;default:0" json:"area_id"`
      Enabled    uint8                  `gorm:"column:F_enabled" sql:"type:tinyint(8) unsigned;not null;default:1" json:"-"`
      CreateTime time.Time              `gorm:"column:F_create_time" sql:"type:bigint(64) unsigned;not null;default:0" json:"-"`
      UpdateTime timelib.MySQLTimestamp `gorm:"column:F_update_time" sql:"type:bigint(64) unsigned;not null;default:0" json:"-"`
  	}
	```
