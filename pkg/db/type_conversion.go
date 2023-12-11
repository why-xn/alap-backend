package db

import "github.com/why-xn/alap-backend/pkg/db/model"

type DbFindObj interface{}

func ConvertToUserArray(dbFindObjList []DbFindObj) []model.User {
	var userList []model.User
	userList = make([]model.User, len(dbFindObjList))
	for i := range dbFindObjList {
		userList[i] = *dbFindObjList[i].(*model.User)
	}
	return userList
}
