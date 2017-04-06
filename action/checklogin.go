package action

import (
	"tribe/sqldb"
)

// 判断玩家的帐号和密码
//	@parames
//		uid	string 	玩家UID
//		pwd string	玩家密码
//	@return
//		int					(0成功 1帐号不存在 2密码不正确)
//		map[string]string	返回玩家对象信息
func CheckLogin(uid, pwd string) (int, map[string]string) {
	rss, err := sqldb.Querys("u_user", "*", "uid='"+uid+"'")
	if err != nil {
		return 3, nil
	}
	if len(rss) < 1 {
		return 1, nil
	}
	rs := rss[0]
	if rs["pass"] != pwd {
		return 2, nil
	}
	return 0, rs
}
