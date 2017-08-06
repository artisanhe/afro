package modules

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

type PhysicsDeleteByUniquustomerG7sList []PhysicsDeleteByUniquustomerG7s

func (pdbug PhysicsDeleteByUniquustomerG7s) TableName() string {
	return "t_physics_delete_by_uniquustomer_g7s"
}

func (pdbug PhysicsDeleteByUniquustomerG7s) PrintDuration(printParam map[string]interface{}) func() {
	now := time.Now()
	return func() {
		duration := time.Now().Sub(now)
		var printStr string
		for key, value := range printParam {
			printStr += fmt.Sprintf("%s:%s, ", key, value)
		}
		fmt.Printf("%sCost:%sms\n", printStr, fmt.Sprintf("%0.3f", float64(duration/time.Millisecond)))
	}
}

func (pdbugl *PhysicsDeleteByUniquustomerG7sList) BatchFetchByCreateTimeList(db *gorm.DB, createTimeList []time.Time) error {
	defer PhysicsDeleteByUniquustomerG7s{}.PrintDuration(map[string]interface{}{
		"request": "[DB]PhysicsDeleteByUniquustomerG7s.BatchFetchByCreateTimeList",
	})()

	if len(createTimeList) == 0 {
		return nil
	}

	err := db.Table(PhysicsDeleteByUniquustomerG7s{}.TableName()).Where("F_create_time in (?) and F_enabled = ?", createTimeList, 1).Find(pdbugl).Error
	return err
}

func (pdbugl *PhysicsDeleteByUniquustomerG7sList) BatchFetchByCustomerIDList(db *gorm.DB, customerIDList []uint64) error {
	defer PhysicsDeleteByUniquustomerG7s{}.PrintDuration(map[string]interface{}{
		"request": "[DB]PhysicsDeleteByUniquustomerG7s.BatchFetchByCustomerIDList",
	})()

	if len(customerIDList) == 0 {
		return nil
	}

	err := db.Table(PhysicsDeleteByUniquustomerG7s{}.TableName()).Where("F_customer_id in (?) and F_enabled = ?", customerIDList, 1).Find(pdbugl).Error
	return err
}

func (pdbugl *PhysicsDeleteByUniquustomerG7sList) BatchFetchByG7sOrgCodeList(db *gorm.DB, g7sOrgCodeList []string) error {
	defer PhysicsDeleteByUniquustomerG7s{}.PrintDuration(map[string]interface{}{
		"request": "[DB]PhysicsDeleteByUniquustomerG7s.BatchFetchByG7sOrgCodeList",
	})()

	if len(g7sOrgCodeList) == 0 {
		return nil
	}

	err := db.Table(PhysicsDeleteByUniquustomerG7s{}.TableName()).Where("F_g7s_org_code in (?) and F_enabled = ?", g7sOrgCodeList, 1).Find(pdbugl).Error
	return err
}

func (pdbugl *PhysicsDeleteByUniquustomerG7sList) BatchFetchByG7sUserIDList(db *gorm.DB, g7sUserIDList []string) error {
	defer PhysicsDeleteByUniquustomerG7s{}.PrintDuration(map[string]interface{}{
		"request": "[DB]PhysicsDeleteByUniquustomerG7s.BatchFetchByG7sUserIDList",
	})()

	if len(g7sUserIDList) == 0 {
		return nil
	}

	err := db.Table(PhysicsDeleteByUniquustomerG7s{}.TableName()).Where("F_g7s_user_id in (?) and F_enabled = ?", g7sUserIDList, 1).Find(pdbugl).Error
	return err
}

func (pdbugl *PhysicsDeleteByUniquustomerG7sList) BatchFetchByUpdateTimeList(db *gorm.DB, updateTimeList []time.Time) error {
	defer PhysicsDeleteByUniquustomerG7s{}.PrintDuration(map[string]interface{}{
		"request": "[DB]PhysicsDeleteByUniquustomerG7s.BatchFetchByUpdateTimeList",
	})()

	if len(updateTimeList) == 0 {
		return nil
	}

	err := db.Table(PhysicsDeleteByUniquustomerG7s{}.TableName()).Where("F_update_time in (?) and F_enabled = ?", updateTimeList, 1).Find(pdbugl).Error
	return err
}

func (pdbug *PhysicsDeleteByUniquustomerG7s) Create(db *gorm.DB) error {
	defer pdbug.PrintDuration(map[string]interface{}{
		"request": "[DB]PhysicsDeleteByUniquustomerG7s.Create",
	})()

	if pdbug.CreateTime.IsZero() {
		pdbug.CreateTime = time.Now()
	}

	if pdbug.UpdateTime.IsZero() {
		pdbug.UpdateTime = time.Now()
	}

	pdbug.Enabled = uint8(1)
	return db.Table(pdbug.TableName()).Create(pdbug).Error
}

func (pdbug *PhysicsDeleteByUniquustomerG7s) DeleteByG7sUserIDAndCustomerID(db *gorm.DB) error {
	defer pdbug.PrintDuration(map[string]interface{}{
		"request": "[DB]PhysicsDeleteByUniquustomerG7s.DeleteByG7sUserIDAndCustomerID",
	})()

	err := db.Table(pdbug.TableName()).Where("F_g7s_user_id = ? and F_customer_id = ? and F_enabled = ?", pdbug.G7sUserID, pdbug.CustomerID, 1).Delete(pdbug).Error
	return err
}

func (pdbugl *PhysicsDeleteByUniquustomerG7sList) FetchByCreateTime(db *gorm.DB, createTime time.Time) error {
	defer PhysicsDeleteByUniquustomerG7s{}.PrintDuration(map[string]interface{}{
		"request": "[DB]PhysicsDeleteByUniquustomerG7s.FetchByCreateTime",
	})()

	err := db.Table(PhysicsDeleteByUniquustomerG7s{}.TableName()).Where("F_create_time = ? and F_enabled = ?", createTime, 1).Find(pdbugl).Error
	return err
}

func (pdbugl *PhysicsDeleteByUniquustomerG7sList) FetchByCustomerID(db *gorm.DB, customerID uint64) error {
	defer PhysicsDeleteByUniquustomerG7s{}.PrintDuration(map[string]interface{}{
		"request": "[DB]PhysicsDeleteByUniquustomerG7s.FetchByCustomerID",
	})()

	err := db.Table(PhysicsDeleteByUniquustomerG7s{}.TableName()).Where("F_customer_id = ? and F_enabled = ?", customerID, 1).Find(pdbugl).Error
	return err
}

func (pdbugl *PhysicsDeleteByUniquustomerG7sList) FetchByG7sOrgCode(db *gorm.DB, g7sOrgCode string) error {
	defer PhysicsDeleteByUniquustomerG7s{}.PrintDuration(map[string]interface{}{
		"request": "[DB]PhysicsDeleteByUniquustomerG7s.FetchByG7sOrgCode",
	})()

	err := db.Table(PhysicsDeleteByUniquustomerG7s{}.TableName()).Where("F_g7s_org_code = ? and F_enabled = ?", g7sOrgCode, 1).Find(pdbugl).Error
	return err
}

func (pdbugl *PhysicsDeleteByUniquustomerG7sList) FetchByG7sUserID(db *gorm.DB, g7sUserID string) error {
	defer PhysicsDeleteByUniquustomerG7s{}.PrintDuration(map[string]interface{}{
		"request": "[DB]PhysicsDeleteByUniquustomerG7s.FetchByG7sUserID",
	})()

	err := db.Table(PhysicsDeleteByUniquustomerG7s{}.TableName()).Where("F_g7s_user_id = ? and F_enabled = ?", g7sUserID, 1).Find(pdbugl).Error
	return err
}

func (pdbug *PhysicsDeleteByUniquustomerG7s) FetchByG7sUserIDAndCustomerID(db *gorm.DB) error {
	defer pdbug.PrintDuration(map[string]interface{}{
		"request": "[DB]PhysicsDeleteByUniquustomerG7s.FetchByG7sUserIDAndCustomerID",
	})()

	err := db.Table(pdbug.TableName()).Where("F_g7s_user_id = ? and F_customer_id = ? and F_enabled = ?", pdbug.G7sUserID, pdbug.CustomerID, 1).Find(pdbug).Error
	return err
}

func (pdbug *PhysicsDeleteByUniquustomerG7s) FetchByG7sUserIDAndCustomerIDForUpdate(db *gorm.DB) error {
	defer pdbug.PrintDuration(map[string]interface{}{
		"request": "[DB]PhysicsDeleteByUniquustomerG7s.FetchByG7sUserIDAndCustomerIDForUpdate",
	})()

	err := db.Table(pdbug.TableName()).Where("F_g7s_user_id = ? and F_customer_id = ? and F_enabled = ?", pdbug.G7sUserID, pdbug.CustomerID, 1).Set("gorm:query_option", "FOR UPDATE").Find(pdbug).Error
	return err
}

func (pdbugl *PhysicsDeleteByUniquustomerG7sList) FetchByUpdateTime(db *gorm.DB, updateTime time.Time) error {
	defer PhysicsDeleteByUniquustomerG7s{}.PrintDuration(map[string]interface{}{
		"request": "[DB]PhysicsDeleteByUniquustomerG7s.FetchByUpdateTime",
	})()

	err := db.Table(PhysicsDeleteByUniquustomerG7s{}.TableName()).Where("F_update_time = ? and F_enabled = ?", updateTime, 1).Find(pdbugl).Error
	return err
}

func (pdbugl *PhysicsDeleteByUniquustomerG7sList) FetchList(db *gorm.DB, size, offset int32, query ...map[string]interface{}) (int32, error) {
	defer PhysicsDeleteByUniquustomerG7s{}.PrintDuration(map[string]interface{}{
		"request": "[DB]PhysicsDeleteByUniquustomerG7s.FetchList",
	})()

	var count int32
	if len(query) == 0 {
		query = append(query, map[string]interface{}{"F_enabled": 1})
	} else {
		if _, ok := query[0]["F_enabled"]; !ok {
			query[0]["F_enabled"] = 1
		}
	}

	if size <= 0 {
		size = -1
		offset = -1
	}
	var err error

	err = db.Table(PhysicsDeleteByUniquustomerG7s{}.TableName()).Where(query[0]).Count(&count).Limit(size).Offset(offset).Order("F_create_time desc").Find(pdbugl).Error

	return int32(count), err
}

func (pdbug *PhysicsDeleteByUniquustomerG7s) SoftDeleteByG7sUserIDAndCustomerID(db *gorm.DB) error {
	defer pdbug.PrintDuration(map[string]interface{}{
		"request": "[DB]PhysicsDeleteByUniquustomerG7s.SoftDeleteByG7sUserIDAndCustomerID",
	})()

	var updateMap = map[string]interface{}{}
	updateMap["F_enabled"] = 2

	if pdbug.UpdateTime.IsZero() {
		pdbug.UpdateTime = time.Now()
	}

	err := db.Table(pdbug.TableName()).Where("F_g7s_user_id = ? and F_customer_id = ? and F_enabled = ?", pdbug.G7sUserID, pdbug.CustomerID, 1).Updates(updateMap).Error
	return err
}

func (pdbug *PhysicsDeleteByUniquustomerG7s) UpdateByG7sUserIDAndCustomerIDWithMap(db *gorm.DB, updateMap map[string]interface{}) error {
	defer pdbug.PrintDuration(map[string]interface{}{
		"request": "[DB]PhysicsDeleteByUniquustomerG7s.UpdateByG7sUserIDAndCustomerIDWithMap",
	})()

	if _, ok := updateMap["F_update_time"]; !ok {
		pdbug.UpdateTime = time.Now()

	}
	dbRet := db.Table(pdbug.TableName()).Where("F_g7s_user_id = ? and F_customer_id = ? and F_enabled = ?", pdbug.G7sUserID, pdbug.CustomerID, 1).Updates(updateMap)
	err := dbRet.Error
	if err != nil {
		return err
	} else {
		if dbRet.RowsAffected == 0 {
			return db.Table(pdbug.TableName()).Where("F_g7s_user_id = ? and F_customer_id = ? and F_enabled = ?", pdbug.G7sUserID, pdbug.CustomerID, 1).Find(&PhysicsDeleteByUniquustomerG7s{}).Error
		} else {
			return nil
		}
	}
}

func (pdbug *PhysicsDeleteByUniquustomerG7s) UpdateByG7sUserIDAndCustomerIDWithStruct(db *gorm.DB) error {
	defer pdbug.PrintDuration(map[string]interface{}{
		"request": "[DB]PhysicsDeleteByUniquustomerG7s.UpdateByG7sUserIDAndCustomerIDWithStruct",
	})()

	if pdbug.UpdateTime.IsZero() {
		pdbug.UpdateTime = time.Now()
	}

	dbRet := db.Table(pdbug.TableName()).Where("F_g7s_user_id = ? and F_customer_id = ? and F_enabled = ?", pdbug.G7sUserID, pdbug.CustomerID, 1).Updates(pdbug)
	err := dbRet.Error
	if err != nil {
		return err
	} else {
		if dbRet.RowsAffected == 0 {
			return db.Table(pdbug.TableName()).Where("F_g7s_user_id = ? and F_customer_id = ? and F_enabled = ?", pdbug.G7sUserID, pdbug.CustomerID, 1).Find(&PhysicsDeleteByUniquustomerG7s{}).Error
		} else {
			return nil
		}
	}
}
