package handle

func Map2String(mp []string, sp string) string {
	vm := ""
	for _, i := range mp {
		vm += i
		vm += sp
	}
	return vm
}
