package conf

const (
	path          = "./database/sql/"
	authsqlname   = "auth.sql"
	cookiesqlname = "cookie.sql"
)

//func (Pg *PsgrDB) UsRegister() {
//	if err := Pg.DB.CreateTable(path + authsqlname).Error; err != nil {
//		log.Panic("CK table register err")
//	}
//}
//
//func (Pg *PsgrDB) CkRegister() {
//	if err := Pg.DB.CreateTable(path + cookiesqlname).Error; err != nil {
//		log.Panic("CK table register err")
//	}
//}
