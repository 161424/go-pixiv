package util

import (
	"fmt"
	"os"
)

func CreateRootFile(tp string) {
	os.Mkdir("./download"+tp, os.ModePerm)
	fmt.Printf("Creating directory: %s", tp)

}
