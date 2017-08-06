package modules

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

type CustomerG7sList []CustomerG7s

func (cg CustomerG7s) TableName() string {
	return "t_customer_g7s"
}

func (cg CustomerG7s) PrintDuration(printParam map[string]interface{}) func() {
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

func (cgl *CustomerG7sList) BatchFetchByCreateTimeList(db *gorm.DB, createTimeList []time.Time) error {
	defer CustomerG7s{}.PrintDuration(map[string]interface{}{
		"request": "[DB]CustomerG7s.BatchFetchByCreateTimeList",
	})()

	if len(createTimeList) == 0 {
		return nil
	}

	err := db.Table(CustomerG7s{}.TableName()).Where("F_create_time in (?) and F_enabled = ?", createTimeList, 1).Find(cgl).Error
	return err
}

func (cgl *CustomerG7sList) BatchFetchByCustomerIDList(db *gorm.DB, customerIDList []uint64) error {
	defer CustomerG7s{}.PrintDuration(map[string]interface{}{
		"request": "[DB]CustomerG7s.BatchFetchByCustomerIDList",
	})()

	if len(customerIDList) == 0 {
		return nil
	}

	err := db.Table(CustomerG7s{}.TableName()).Where("F_customer_id in (?) and F_enabled = ?", customerIDList, 1).Find(cgl).Error
	return err
}

func (cgl *CustomerG7sList) BatchFetchByG7sUserIDList(db *gorm.DB, g7sUserIDList []string) error {
	defer CustomerG7s{}.PrintDuration(map[string]interface{}{
		"request": "[DB]CustomerG7s.BatchFetchByG7sUserIDList",
	})()

	if len(g7sUserIDList) == 0 {
		return nil
	}

	err := db.Table(CustomerG7s{}.TableName()).Where("F_g7s_user_id in (?) and F_enabled = ?", g7sUserIDList, 1).Find(cgl).Error
	return err
}

func (cgl *CustomerG7sList) BatchFetchByUpdateTimeList(db *gorm.DB, updateTimeList []time.Time) error {
	defer CustomerG7s{}.PrintDuration(map[string]interface{}{
		"request": "[DB]CustomerG7s.BatchFetchByUpdateTimeList",
	})()

	if len(updateTimeList) == 0 {
		return nil
	}

	err := db.Table(CustomerG7s{}.TableName()).Where("F_update_time in (?) and F_enabled = ?", updateTimeList, 1).Find(cgl).Error
	return err
}

func (cg *CustomerG7s) Create(db *gorm.DB) error {
	defer cg.PrintDuration(map[string]interface{}{
		"request": "[DB]CustomerG7s.Create",
	})()

	if cg.CreateTime.IsZero() {
		cg.CreateTime = time.Now()
	}

	if cg.UpdateTime.IsZero() {
		cg.UpdateTime = time.Now()
	}

	cg.Enabled = uint8(1)
	return db.Table(cg.TableName()).Create(cg).Error
}

func (cg *CustomerG7s) DeleteByCustomerIDAndG7sOrgCodeAndG7sUserID(db *gorm.DB) error {
	defer cg.PrintDuration(map[string]interface{}{
		"request": "[DB]CustomerG7s.DeleteByCustomerIDAndG7sOrgCodeAndG7sUserID",
	})()

	err := db.Table(cg.TableName()).Where("F_customer_id = ? and F_g7s_org_code = ? and F_g7s_user_id = ? and F_enabled = ?", cg.CustomerID, cg.G7sOrgCode, cg.G7sUserID, 1).Delete(cg).Error
	return err
}

func (cg *CustomerG7s) DeleteByG7sUserIDAndCustomerID(db *gorm.DB) error {
	defer cg.PrintDuration(map[string]interface{}{
		"request": "[DB]CustomerG7s.DeleteByG7sUserIDAndCustomerID",
	})()

	err := db.Table(cg.TableName()).Where("F_g7s_user_id = ? and F_customer_id = ? and F_enabled = ?", cg.G7sUserID, cg.CustomerID, 1).Delete(cg).Error
	return err
}

func (cgl *CustomerG7sList) FetchByCreateTime(db *gorm.DB, createTime time.Time) error {
	defer CustomerG7s{}.PrintDuration(map[string]interface{}{
		"request": "[DB]CustomerG7s.FetchByCreateTime",
	})()

	err := db.Table(CustomerG7s{}.TableName()).Where("F_create_time = ? and F_enabled = ?", createTime, 1).Find(cgl).Error
	return err
}

func (cgl *CustomerG7sList) FetchByCustomerID(db *gorm.DB, customerID uint64) error {
	defer CustomerG7s{}.PrintDuration(map[string]interface{}{
		"request": "[DB]CustomerG7s.FetchByCustomerID",
	})()

	err := db.Table(CustomerG7s{}.TableName()).Where("F_customer_id = ? and F_enabled = ?", customerID, 1).Find(cgl).Error
	return err
}

func (cgl *CustomerG7sList) FetchByCustomerIDAndG7sOrgCode(db *gorm.DB, customerID uint64, g7sOrgCode string) error {
	defer CustomerG7s{}.PrintDuration(map[string]interface{}{
		"request": "[DB]CustomerG7s.FetchByCustomerIDAndG7sOrgCode",
	})()

	err := db.Table(CustomerG7s{}.TableName()).Where("F_customer_id = ? and F_g7s_org_code = ? and F_enabled = ?", customerID, g7sOrgCode, 1).Find(cgl).Error
	return err
}

func (cg *CustomerG7s) FetchByCustomerIDAndG7sOrgCodeAndG7sUserID(db *gorm.DB) error {
	defer cg.PrintDuration(map[string]interface{}{
		"request": "[DB]CustomerG7s.FetchByCustomerIDAndG7sOrgCodeAndG7sUserID",
	})()

	err := db.Table(cg.TableName()).Where("F_customer_id = ? and F_g7s_org_code = ? and F_g7s_user_id = ? and F_enabled = ?", cg.CustomerID, cg.G7sOrgCode, cg.G7sUserID, 1).Find(cg).Error
	return err
}

func (cg *CustomerG7s) FetchByCustomerIDAndG7sOrgCodeAndG7sUserIDForUpdate(db *gorm.DB) error {
	defer cg.PrintDuration(map[string]interface{}{
		"request": "[DB]CustomerG7s.FetchByCustomerIDAndG7sOrgCodeAndG7sUserIDForUpdate",
	})()

	err := db.Table(cg.TableName()).Where("F_customer_id = ? and F_g7s_org_code = ? and F_g7s_user_id = ? and F_enabled = ?", cg.CustomerID, cg.G7sOrgCode, cg.G7sUserID, 1).Set("gorm:query_option", "FOR UPDATE").Find(cg).Error
	return err
}

func (cgl *CustomerG7sList) FetchByG7sUserID(db *gorm.DB, g7sUserID string) error {
	defer CustomerG7s{}.PrintDuration(map[string]interface{}{
		"request": "[DB]CustomerG7s.FetchByG7sUserID",
	})()

	err := db.Table(CustomerG7s{}.TableName()).Where("F_g7s_user_id = ? and F_enabled = ?", g7sUserID, 1).Find(cgl).Error
	return err
}

func (cg *CustomerG7s) FetchByG7sUserIDAndCustomerID(db *gorm.DB) error {
	defer cg.PrintDuration(map[string]interface{}{
		"request": "[DB]CustomerG7s.FetchByG7sUserIDAndCustomerID",
	})()

	err := db.Table(cg.TableName()).Where("F_g7s_user_id = ? and F_customer_id = ? and F_enabled = ?", cg.G7sUserID, cg.CustomerID, 1).Find(cg).Error
	return err
}

func (cg *CustomerG7s) FetchByG7sUserIDAndCustomerIDForUpdate(db *gorm.DB) error {
	defer cg.PrintDuration(map[string]interface{}{
		"request": "[DB]CustomerG7s.FetchByG7sUserIDAndCustomerIDForUpdate",
	})()

	err := db.Table(cg.TableName()).Where("F_g7s_user_id = ? and F_customer_id = ? and F_enabled = ?", cg.G7sUserID, cg.CustomerID, 1).Set("gorm:query_option", "FOR UPDATE").Find(cg).Error
	return err
}

func (cgl *CustomerG7sList) FetchByUpdateTime(db *gorm.DB, updateTime time.Time) error {
	defer CustomerG7s{}.PrintDuration(map[string]interface{}{
		"request": "[DB]CustomerG7s.FetchByUpdateTime",
	})()

	err := db.Table(CustomerG7s{}.TableName()).Where("F_update_time = ? and F_enabled = ?", updateTime, 1).Find(cgl).Error
	return err
}

func (cgl *CustomerG7sList) FetchList(db *gorm.DB, size, offset int32, query ...map[string]interface{}) (int32, error) {
	defer CustomerG7s{}.PrintDuration(map[string]interface{}{
		"request": "[DB]CustomerG7s.FetchList",
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

	err = db.Table(CustomerG7s{}.TableName()).Where(query[0]).Count(&count).Limit(size).Offset(offset).Order("F_create_time desc").Find(cgl).Error

	return int32(count), err
}

func (cg *CustomerG7s) SoftDeleteByCustomerIDAndG7sOrgCodeAndG7sUserID(db *gorm.DB) error {
	defer cg.PrintDuration(map[string]interface{}{
		"request": "[DB]CustomerG7s.SoftDeleteByCustomerIDAndG7sOrgCodeAndG7sUserID",
	})()

	var updateMap = map[string]interface{}{}
	updateMap["F_enabled"] = 2

	if cg.UpdateTime.IsZero() {
		cg.UpdateTime = time.Now()
	}

	err := db.Table(cg.TableName()).Where("F_customer_id = ? and F_g7s_org_code = ? and F_g7s_user_id = ? and F_enabled = ?", cg.CustomerID, cg.G7sOrgCode, cg.G7sUserID, 1).Updates(updateMap).Error
	return err
}

func (cg *CustomerG7s) SoftDeleteByG7sUserIDAndCustomerID(db *gorm.DB) error {
	defer cg.PrintDuration(map[string]interface{}{
		"request": "[DB]CustomerG7s.SoftDeleteByG7sUserIDAndCustomerID",
	})()

	var updateMap = map[string]interface{}{}
	updateMap["F_enabled"] = 2

	if cg.UpdateTime.IsZero() {
		cg.UpdateTime = time.Now()
	}

	err := db.Table(cg.TableName()).Where("F_g7s_user_id = ? and F_customer_id = ? and F_enabled = ?", cg.G7sUserID, cg.CustomerID, 1).Updates(updateMap).Error
	return err
}

func (cg *CustomerG7s) UpdateByCustomerIDAndG7sOrgCodeAndG7sUserIDWithMap(db *gorm.DB, updateMap map[string]interface{}) error {
	defer cg.PrintDuration(map[string]interface{}{
		"request": "[DB]CustomerG7s.UpdateByCustomerIDAndG7sOrgCodeAndG7sUserIDWithMap",
	})()

	if _, ok := updateMap["F_update_time"]; !ok {
		cg.UpdateTime = time.Now()

	}
	dbRet := db.Table(cg.TableName()).Where("F_customer_id = ? and F_g7s_org_code = ? and F_g7s_user_id = ? and F_enabled = ?", cg.CustomerID, cg.G7sOrgCode, cg.G7sUserID, 1).Updates(updateMap)
	err := dbRet.Error
	if err != nil {
		return err
	} else {
		if dbRet.RowsAffected == 0 {
			return db.Table(cg.TableName()).Where("F_customer_id = ? and F_g7s_org_code = ? and F_g7s_user_id = ? and F_enabled = ?", cg.CustomerID, cg.G7sOrgCode, cg.G7sUserID, 1).Find(&CustomerG7s{}).Error
		} else {
			return nil
		}
	}
}

func (cg *CustomerG7s) UpdateByCustomerIDAndG7sOrgCodeAndG7sUserIDWithStruct(db *gorm.DB) error {
	defer cg.PrintDuration(map[string]interface{}{
		"request": "[DB]CustomerG7s.UpdateByCustomerIDAndG7sOrgCodeAndG7sUserIDWithStruct",
	})()

	if cg.UpdateTime.IsZero() {
		cg.UpdateTime = time.Now()
	}

	dbRet := db.Table(cg.TableName()).Where("F_customer_id = ? and F_g7s_org_code = ? and F_g7s_user_id = ? and F_enabled = ?", cg.CustomerID, cg.G7sOrgCode, cg.G7sUserID, 1).Updates(cg)
	err := dbRet.Error
	if err != nil {
		return err
	} else {
		if dbRet.RowsAffected == 0 {
			return db.Table(cg.TableName()).Where("F_customer_id = ? and F_g7s_org_code = ? and F_g7s_user_id = ? and F_enabled = ?", cg.CustomerID, cg.G7sOrgCode, cg.G7sUserID, 1).Find(&CustomerG7s{}).Error
		} else {
			return nil
		}
	}
}

func (cg *CustomerG7s) UpdateByG7sUserIDAndCustomerIDWithMap(db *gorm.DB, updateMap map[string]interface{}) error {
	defer cg.PrintDuration(map[string]interface{}{
		"request": "[DB]CustomerG7s.UpdateByG7sUserIDAndCustomerIDWithMap",
	})()

	if _, ok := updateMap["F_update_time"]; !ok {
		cg.UpdateTime = time.Now()

	}
	dbRet := db.Table(cg.TableName()).Where("F_g7s_user_id = ? and F_customer_id = ? and F_enabled = ?", cg.G7sUserID, cg.CustomerID, 1).Updates(updateMap)
	err := dbRet.Error
	if err != nil {
		return err
	} else {
		if dbRet.RowsAffected == 0 {
			return db.Table(cg.TableName()).Where("F_g7s_user_id = ? and F_customer_id = ? and F_enabled = ?", cg.G7sUserID, cg.CustomerID, 1).Find(&CustomerG7s{}).Error
		} else {
			return nil
		}
	}
}

func (cg *CustomerG7s) UpdateByG7sUserIDAndCustomerIDWithStruct(db *gorm.DB) error {
	defer cg.PrintDuration(map[string]interface{}{
		"request": "[DB]CustomerG7s.UpdateByG7sUserIDAndCustomerIDWithStruct",
	})()

	if cg.UpdateTime.IsZero() {
		cg.UpdateTime = time.Now()
	}

	dbRet := db.Table(cg.TableName()).Where("F_g7s_user_id = ? and F_customer_id = ? and F_enabled = ?", cg.G7sUserID, cg.CustomerID, 1).Updates(cg)
	err := dbRet.Error
	if err != nil {
		return err
	} else {
		if dbRet.RowsAffected == 0 {
			return db.Table(cg.TableName()).Where("F_g7s_user_id = ? and F_customer_id = ? and F_enabled = ?", cg.G7sUserID, cg.CustomerID, 1).Find(&CustomerG7s{}).Error
		} else {
			return nil
		}
	}
}
