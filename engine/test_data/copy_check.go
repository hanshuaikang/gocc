package test_data

// test file
// nolint
func bubbleSort(array []int) []int {
	length := len(array)

	for i := 0; i < length; i++ {
		for j := 0; j < length-i-1; j++ {
			if array[j] > array[j+1] {
				array[j+1], array[j] = array[j], array[j+1]
			}
		}
	}
	return array
}

// nolint
func bubbleSort2(array []int) []int {
	length := len(array)

	for i := 0; i < length; i++ {
		for j := 0; j < length-i-1; j++ {
			if array[j] > array[j+1] {
				array[j+1], array[j] = array[j], array[j+1]
			}
		}
	}
	return array
}
