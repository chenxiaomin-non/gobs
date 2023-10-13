package kmp

// search for all occurrences of substrings in string
// return a list of indexes
func SearchAllOccurrences(substring string, stringToSearch string) []int {
	// create the failure table
	failureTable := createFailureTable(substring)
	// search for all occurrences
	var result []int
	var i int = 0
	var j int = 0
	for i < len(stringToSearch) {
		if stringToSearch[i] == substring[j] {
			i++
			j++
			if j == len(substring) {
				result = append(result, i-j)
				j = failureTable[j-1]
			}
		} else {
			if j == 0 {
				i++
			} else {
				j = failureTable[j-1]
			}
		}
	}
	return result
}

// create the failure table
func createFailureTable(substring string) []int {
	var failureTable []int = make([]int, len(substring))
	var i int = 1
	var j int = 0
	for i < len(substring) {
		if substring[i] == substring[j] {
			failureTable[i] = j + 1
			i++
			j++
		} else {
			if j == 0 {
				failureTable[i] = 0
				i++
			} else {
				j = failureTable[j-1]
			}
		}
	}
	return failureTable
}
