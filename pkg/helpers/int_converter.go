package helpers

func GetIntList(list []int32) []int {
	result := make([]int, 0, len(list))

	for _, v := range list {
		result = append(result, int(v))
	}

	return result
}
