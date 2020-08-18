package useremailsqlrepo

import (
	samasemodels "samase/models"
	"samase/useremail"
)

func UserEmailSQLJoin(sf samasemodels.SQLFetcher) *useremail.UserEmail {
	usem := &useremail.UserEmail{}
	dest := []interface{}{&usem.Value, usem.Verified}
	sf.AddScanDest(dest)
	sf.AddJoins(" INNER JOIN user_email ON user.id = user_email.user_id ")
	sf.AddFields(",user_email.value,user_email.verified")
	return usem
}
