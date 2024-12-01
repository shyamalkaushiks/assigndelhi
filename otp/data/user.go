package data

import "otppro/model"

func GetUSerbyno(phoneno string) (model.UserRegisterwe, error) {
	db := model.DBConn
	var user model.UserRegisterwe
	err := db.Where("phone_number=?", phoneno).Error
	if err != nil {
		return user, err
	}
	return user, nil
}
