package modules

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

type UserTestList []UserTest

func (ut UserTest) PrintDuration(printParam map[string]interface{}) func() {
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

func (utl *UserTestList) BatchFetchByCityIDList(db *gorm.DB, cityIDList []uint64) error {
	defer UserTest{}.PrintDuration(map[string]interface{}{
		"request": "[DB]UserTest.BatchFetchByCityIDList",
	})()

	if len(cityIDList) == 0 {
		return nil
	}

	err := db.Table(UserTest{}.TableName()).Where("F_city_id in (?) and F_enabled = ?", cityIDList, 1).Find(utl).Error
	return err
}

func (utl *UserTestList) BatchFetchByIdList(db *gorm.DB, idList []uint64) error {
	defer UserTest{}.PrintDuration(map[string]interface{}{
		"request": "[DB]UserTest.BatchFetchByIdList",
	})()

	if len(idList) == 0 {
		return nil
	}

	err := db.Table(UserTest{}.TableName()).Where("F_id in (?) and F_enabled = ?", idList, 1).Find(utl).Error
	return err
}

func (utl *UserTestList) BatchFetchByPhoneList(db *gorm.DB, phoneList []string) error {
	defer UserTest{}.PrintDuration(map[string]interface{}{
		"request": "[DB]UserTest.BatchFetchByPhoneList",
	})()

	if len(phoneList) == 0 {
		return nil
	}

	err := db.Table(UserTest{}.TableName()).Where("F_phone in (?) and F_enabled = ?", phoneList, 1).Find(utl).Error
	return err
}

func (utl *UserTestList) BatchFetchByUserIDList(db *gorm.DB, userIDList []uint64) error {
	defer UserTest{}.PrintDuration(map[string]interface{}{
		"request": "[DB]UserTest.BatchFetchByUserIDList",
	})()

	if len(userIDList) == 0 {
		return nil
	}

	err := db.Table(UserTest{}.TableName()).Where("F_user_id in (?) and F_enabled = ?", userIDList, 1).Find(utl).Error
	return err
}

func (ut *UserTest) Create(db *gorm.DB) error {
	defer ut.PrintDuration(map[string]interface{}{
		"request": "[DB]UserTest.Create",
	})()

	if ut.CreateTime.IsZero() {
		ut.CreateTime = time.Now()
	}

	if ut.UpdateTime.IsZero() {
		ut.UpdateTime = time.Now()
	}

	ut.Enabled = uint8(1)
	return db.Table(ut.TableName()).Create(ut).Error
}

func (ut *UserTest) DeleteById(db *gorm.DB) error {
	defer ut.PrintDuration(map[string]interface{}{
		"request": "[DB]UserTest.DeleteById",
	})()

	err := db.Table(ut.TableName()).Where("F_id = ? and F_enabled = ?", ut.Id, 1).Delete(ut).Error
	return err
}

func (ut *UserTest) DeleteByPhone(db *gorm.DB) error {
	defer ut.PrintDuration(map[string]interface{}{
		"request": "[DB]UserTest.DeleteByPhone",
	})()

	err := db.Table(ut.TableName()).Where("F_phone = ? and F_enabled = ?", ut.Phone, 1).Delete(ut).Error
	return err
}

func (ut *UserTest) DeleteByUserIDAndName(db *gorm.DB) error {
	defer ut.PrintDuration(map[string]interface{}{
		"request": "[DB]UserTest.DeleteByUserIDAndName",
	})()

	err := db.Table(ut.TableName()).Where("F_user_id = ? and F_name = ? and F_enabled = ?", ut.UserID, ut.Name, 1).Delete(ut).Error
	return err
}

func (utl *UserTestList) FetchByCityID(db *gorm.DB, cityID uint64) error {
	defer UserTest{}.PrintDuration(map[string]interface{}{
		"request": "[DB]UserTest.FetchByCityID",
	})()

	err := db.Table(UserTest{}.TableName()).Where("F_city_id = ? and F_enabled = ?", cityID, 1).Find(utl).Error
	return err
}

func (utl *UserTestList) FetchByCityIDAndAreaID(db *gorm.DB, cityID uint64, areaID uint64) error {
	defer UserTest{}.PrintDuration(map[string]interface{}{
		"request": "[DB]UserTest.FetchByCityIDAndAreaID",
	})()

	err := db.Table(UserTest{}.TableName()).Where("F_city_id = ? and F_area_id = ? and F_enabled = ?", cityID, areaID, 1).Find(utl).Error
	return err
}

func (ut *UserTest) FetchById(db *gorm.DB) error {
	defer ut.PrintDuration(map[string]interface{}{
		"request": "[DB]UserTest.FetchById",
	})()

	err := db.Table(ut.TableName()).Where("F_id = ? and F_enabled = ?", ut.Id, 1).Find(ut).Error
	return err
}

func (ut *UserTest) FetchByIdForUpdate(db *gorm.DB) error {
	defer ut.PrintDuration(map[string]interface{}{
		"request": "[DB]UserTest.FetchByIdForUpdate",
	})()

	err := db.Table(ut.TableName()).Where("F_id = ? and F_enabled = ?", ut.Id, 1).Set("gorm:query_option", "FOR UPDATE").Find(ut).Error
	return err
}

func (ut *UserTest) FetchByPhone(db *gorm.DB) error {
	defer ut.PrintDuration(map[string]interface{}{
		"request": "[DB]UserTest.FetchByPhone",
	})()

	err := db.Table(ut.TableName()).Where("F_phone = ? and F_enabled = ?", ut.Phone, 1).Find(ut).Error
	return err
}

func (ut *UserTest) FetchByPhoneForUpdate(db *gorm.DB) error {
	defer ut.PrintDuration(map[string]interface{}{
		"request": "[DB]UserTest.FetchByPhoneForUpdate",
	})()

	err := db.Table(ut.TableName()).Where("F_phone = ? and F_enabled = ?", ut.Phone, 1).Set("gorm:query_option", "FOR UPDATE").Find(ut).Error
	return err
}

func (utl *UserTestList) FetchByUserID(db *gorm.DB, userID uint64) error {
	defer UserTest{}.PrintDuration(map[string]interface{}{
		"request": "[DB]UserTest.FetchByUserID",
	})()

	err := db.Table(UserTest{}.TableName()).Where("F_user_id = ? and F_enabled = ?", userID, 1).Find(utl).Error
	return err
}

func (ut *UserTest) FetchByUserIDAndName(db *gorm.DB) error {
	defer ut.PrintDuration(map[string]interface{}{
		"request": "[DB]UserTest.FetchByUserIDAndName",
	})()

	err := db.Table(ut.TableName()).Where("F_user_id = ? and F_name = ? and F_enabled = ?", ut.UserID, ut.Name, 1).Find(ut).Error
	return err
}

func (ut *UserTest) FetchByUserIDAndNameForUpdate(db *gorm.DB) error {
	defer ut.PrintDuration(map[string]interface{}{
		"request": "[DB]UserTest.FetchByUserIDAndNameForUpdate",
	})()

	err := db.Table(ut.TableName()).Where("F_user_id = ? and F_name = ? and F_enabled = ?", ut.UserID, ut.Name, 1).Set("gorm:query_option", "FOR UPDATE").Find(ut).Error
	return err
}

func (utl *UserTestList) FetchList(db *gorm.DB, size, offset int32, query ...map[string]interface{}) (int32, error) {
	defer UserTest{}.PrintDuration(map[string]interface{}{
		"request": "[DB]UserTest.FetchList",
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

	err = db.Table(UserTest{}.TableName()).Where(query[0]).Count(&count).Limit(size).Offset(offset).Order("F_create_time desc").Find(utl).Error

	return int32(count), err
}

func (ut *UserTest) SoftDeleteById(db *gorm.DB) error {
	defer ut.PrintDuration(map[string]interface{}{
		"request": "[DB]UserTest.SoftDeleteById",
	})()

	var updateMap = map[string]interface{}{}
	updateMap["F_enabled"] = 2

	if ut.UpdateTime.IsZero() {
		ut.UpdateTime = time.Now()
	}

	err := db.Table(ut.TableName()).Where("F_id = ? and F_enabled = ?", ut.Id, 1).Updates(updateMap).Error
	return err
}

func (ut *UserTest) SoftDeleteByPhone(db *gorm.DB) error {
	defer ut.PrintDuration(map[string]interface{}{
		"request": "[DB]UserTest.SoftDeleteByPhone",
	})()

	var updateMap = map[string]interface{}{}
	updateMap["F_enabled"] = 2

	if ut.UpdateTime.IsZero() {
		ut.UpdateTime = time.Now()
	}

	err := db.Table(ut.TableName()).Where("F_phone = ? and F_enabled = ?", ut.Phone, 1).Updates(updateMap).Error
	return err
}

func (ut *UserTest) SoftDeleteByUserIDAndName(db *gorm.DB) error {
	defer ut.PrintDuration(map[string]interface{}{
		"request": "[DB]UserTest.SoftDeleteByUserIDAndName",
	})()

	var updateMap = map[string]interface{}{}
	updateMap["F_enabled"] = 2

	if ut.UpdateTime.IsZero() {
		ut.UpdateTime = time.Now()
	}

	err := db.Table(ut.TableName()).Where("F_user_id = ? and F_name = ? and F_enabled = ?", ut.UserID, ut.Name, 1).Updates(updateMap).Error
	return err
}

func (ut *UserTest) UpdateByIdWithMap(db *gorm.DB, updateMap map[string]interface{}) error {
	defer ut.PrintDuration(map[string]interface{}{
		"request": "[DB]UserTest.UpdateByIdWithMap",
	})()

	if _, ok := updateMap["F_update_time"]; !ok {
		ut.UpdateTime = time.Now()

	}
	dbRet := db.Table(ut.TableName()).Where("F_id = ? and F_enabled = ?", ut.Id, 1).Updates(updateMap)
	err := dbRet.Error
	if err != nil {
		return err
	} else {
		if dbRet.RowsAffected == 0 {
			return db.Table(ut.TableName()).Where("F_id = ? and F_enabled = ?", ut.Id, 1).Find(&UserTest{}).Error
		} else {
			return nil
		}
	}
}

func (ut *UserTest) UpdateByIdWithStruct(db *gorm.DB) error {
	defer ut.PrintDuration(map[string]interface{}{
		"request": "[DB]UserTest.UpdateByIdWithStruct",
	})()

	if ut.UpdateTime.IsZero() {
		ut.UpdateTime = time.Now()
	}

	dbRet := db.Table(ut.TableName()).Where("F_id = ? and F_enabled = ?", ut.Id, 1).Updates(ut)
	err := dbRet.Error
	if err != nil {
		return err
	} else {
		if dbRet.RowsAffected == 0 {
			return db.Table(ut.TableName()).Where("F_id = ? and F_enabled = ?", ut.Id, 1).Find(&UserTest{}).Error
		} else {
			return nil
		}
	}
}

func (ut *UserTest) UpdateByPhoneWithMap(db *gorm.DB, updateMap map[string]interface{}) error {
	defer ut.PrintDuration(map[string]interface{}{
		"request": "[DB]UserTest.UpdateByPhoneWithMap",
	})()

	if _, ok := updateMap["F_update_time"]; !ok {
		ut.UpdateTime = time.Now()

	}
	dbRet := db.Table(ut.TableName()).Where("F_phone = ? and F_enabled = ?", ut.Phone, 1).Updates(updateMap)
	err := dbRet.Error
	if err != nil {
		return err
	} else {
		if dbRet.RowsAffected == 0 {
			return db.Table(ut.TableName()).Where("F_phone = ? and F_enabled = ?", ut.Phone, 1).Find(&UserTest{}).Error
		} else {
			return nil
		}
	}
}

func (ut *UserTest) UpdateByPhoneWithStruct(db *gorm.DB) error {
	defer ut.PrintDuration(map[string]interface{}{
		"request": "[DB]UserTest.UpdateByPhoneWithStruct",
	})()

	if ut.UpdateTime.IsZero() {
		ut.UpdateTime = time.Now()
	}

	dbRet := db.Table(ut.TableName()).Where("F_phone = ? and F_enabled = ?", ut.Phone, 1).Updates(ut)
	err := dbRet.Error
	if err != nil {
		return err
	} else {
		if dbRet.RowsAffected == 0 {
			return db.Table(ut.TableName()).Where("F_phone = ? and F_enabled = ?", ut.Phone, 1).Find(&UserTest{}).Error
		} else {
			return nil
		}
	}
}

func (ut *UserTest) UpdateByUserIDAndNameWithMap(db *gorm.DB, updateMap map[string]interface{}) error {
	defer ut.PrintDuration(map[string]interface{}{
		"request": "[DB]UserTest.UpdateByUserIDAndNameWithMap",
	})()

	if _, ok := updateMap["F_update_time"]; !ok {
		ut.UpdateTime = time.Now()

	}
	dbRet := db.Table(ut.TableName()).Where("F_user_id = ? and F_name = ? and F_enabled = ?", ut.UserID, ut.Name, 1).Updates(updateMap)
	err := dbRet.Error
	if err != nil {
		return err
	} else {
		if dbRet.RowsAffected == 0 {
			return db.Table(ut.TableName()).Where("F_user_id = ? and F_name = ? and F_enabled = ?", ut.UserID, ut.Name, 1).Find(&UserTest{}).Error
		} else {
			return nil
		}
	}
}

func (ut *UserTest) UpdateByUserIDAndNameWithStruct(db *gorm.DB) error {
	defer ut.PrintDuration(map[string]interface{}{
		"request": "[DB]UserTest.UpdateByUserIDAndNameWithStruct",
	})()

	if ut.UpdateTime.IsZero() {
		ut.UpdateTime = time.Now()
	}

	dbRet := db.Table(ut.TableName()).Where("F_user_id = ? and F_name = ? and F_enabled = ?", ut.UserID, ut.Name, 1).Updates(ut)
	err := dbRet.Error
	if err != nil {
		return err
	} else {
		if dbRet.RowsAffected == 0 {
			return db.Table(ut.TableName()).Where("F_user_id = ? and F_name = ? and F_enabled = ?", ut.UserID, ut.Name, 1).Find(&UserTest{}).Error
		} else {
			return nil
		}
	}
}
