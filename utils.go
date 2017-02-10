package lazycache

import "strings"
import "fmt"
import "sort"

func MungeHostname(hostname string) []string {
	splitHN := strings.Split(hostname, ".")
	fmt.Println(splitHN)

	for i, j := 0, len(splitHN)-1; i < j; i, j = i+1, j-1 {
		sort.StringSlice(splitHN).Swap(i, j)
	}
	return splitHN
}

func stripBlankElementsRight(slice []string) []string {
	if len(slice) > 0 && len(slice[len(slice)-1]) == 0 {
		return stripBlankElementsRight(slice[:len(slice)-1])
	}
	return slice
}
