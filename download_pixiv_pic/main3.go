package main

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	sql2 "github.com/chen/download_pixiv_pic/common/model/sql"
	"github.com/chen/download_pixiv_pic/database/sql"
	"github.com/lib/pq"
	"log"
)

func main() {
	//var d string
	//fmt.Scanln(&d)
	//if d == "" {
	//	print("ddd")
	//}
	//fmt.Println("123", d, "456")
	G()
}

type Test struct {
	//ID    int                    `gorm:"column:id;type:int;"`
	Age  pq.StringArray `gorm:"column:age;type:text[]"`
	Ages pq.StringArray `gorm:"column:ages;type:text[]"`
	Name string         `gorm:"column:name;type:text;primaryKey"`
	//Local map[string]string `gorm:"column:local;type:json"`
	Local StringAArray `gorm:"column:local;type:text"json:"local"`
	*sql.Aq
}

func (a *A) A1() {
	a.No = "pkpkp"
}

type A struct {
	*sql.Aq
}

func G() {

	t := &Test{
		//ID:    1,
		Age:   []string{"12", "13", "14"},
		Ages:  []string{"aa", "bb", "nn"},
		Name:  "ls",
		Local: StringAArray{"lk": "12121212", "maka": "123"},
		Aq:    &sql.Aq{},
	}
	t.A1()

	t.Create()
	t.Find()
}

func (c *Test) Create() error {
	db := sql2.GetClient()
	db.DB.AutoMigrate(&Test{})
	a := pq.StringArray{}

	b := []pq.StringArray{c.Age, c.Ages}
	for _, _b := range b {
		if _b != nil && len(_b) > 0 {
			for _, v := range _b {
				a = append(a, v)
			}
			_b = a
		}
	}

	if result := db.DB.Where("name=?", c.Name).Save(c); result.Error != nil {
		log.Printf("Error creating company: %s", c.Name)
		return result.Error
	} else {
		fmt.Println(result)
		log.Printf("Successfully created company: %s", c.Name)
		return nil
	}
}

func (c *Test) Find() error {
	db := sql2.GetClient()
	test := []Test{}
	db.DB.Find(&test)
	for i, j := range test {
		log.Printf("%T,%T, %d, %s", i, j, i, j.Local)
	}
	return nil
}

type StringAArray map[string]string

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (j *StringAArray) Scan(value interface{}) (err error) {

	var skills map[string]string
	switch value.(type) {
	case string:
		err = json.Unmarshal([]byte(value.(string)), &skills)
	case []byte:
		err = json.Unmarshal(value.([]byte), &skills)
	default:
		return errors.New("Incompatible type for Skills")
	}
	if err != nil {
		return err
	}
	*j = skills
	return nil
}

// 实现 driver.Valuer 接口，Value 返回 json value
func (j StringAArray) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	if n := len(j); n > 0 {
		// There will be at least two curly brackets, 2*N bytes of quotes,
		// and N-1 bytes of delimiters.
		b := make([]byte, 1, 1+3*n)
		b[0] = '{'
		var _n = 0
		for key, value := range j {
			_n += 1
			b = appendArrayQuotedBytes(b, []byte(key))
			b = append(b, ':')
			b = appendArrayQuotedBytes(b, []byte(value))
			if _n == n {
				continue
			}
			b = append(b, ',')
		}
		return string(append(b, '}')), nil
	}

	return "{}", nil
}

func appendArrayQuotedBytes(b, v []byte) []byte {
	b = append(b, '"')
	for {
		i := bytes.IndexAny(v, `"\`)
		if i < 0 {
			b = append(b, v...)
			break
		}
		if i > 0 {
			b = append(b, v[:i]...)
		}
		b = append(b, '\\', v[i])
		v = v[i+1:]
	}
	return append(b, '"')
}
