package user

import (
	"go-micro-learning/UserExamples/basis_lib/sql_db"
	proto "go-micro-learning/UserExamples/user_service/proto/user"
	"github.com/micro/go-micro/util/log"
)

func (s *service) QueryUserByName(userName string) (ret *proto.User, err error) {
	queryString := `SELECT id, title, phone FROM repository_userprofile WHERE title = ?`

	// 获取数据库
	o := sql_db.GetDB()

	ret = &proto.User{}

	// 查询
	err = o.QueryRow(queryString, userName).Scan(&ret.Id, &ret.Name, &ret.Pwd)
	if err != nil {
		log.Logf("[QueryUserByName] 查询数据失败，err：%s", err)
		return
	}
	return
}
