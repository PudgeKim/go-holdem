package main

import "fmt"

func makeNum(lv int, numList []int, res *[]int) {
	if lv < 5 {
		return
	}
	for i := 0; i < len(numList); i++ {
		*res = append(*res, numList[i])

	}
}
func main() {
	nums := []int{1, 2, 3, 4, 5, 6, 7}
	var res []int
	makeNum(0, nums, &res)
	fmt.Println(res)
}
