package systems

//GetMakingInfo ...
func GetMakingInfo(s []int, e int) bool {
	for _, v := range s {
		if e == v {
			return true
		}
	}
	return false
}
