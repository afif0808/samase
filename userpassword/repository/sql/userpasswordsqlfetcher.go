package userpasswordsqlrepo

import (
	samasemodels "samase/models"
	"samase/userpassword"
)

func UserPasswordSQLJoin(sf samasemodels.SQLFetcher) *userpassword.UserPassword {
	uspa := &userpassword.UserPassword{}
	dest := []interface{}{&uspa.Hash}
	sf.AddScanDest(dest)
	sf.AddJoins(" INNER JOIN user_password ON user.id = user_password.user_id ")
	sf.AddFields(",user_password.hash")
	return uspa
}
