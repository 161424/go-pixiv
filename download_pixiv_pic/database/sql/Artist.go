package sql

type Aq struct {
	No string `gorm:"column:no;type:text;primaryKey"`
}

func (a *Aq) A1() {
	a.No = "pkpkp"
}
