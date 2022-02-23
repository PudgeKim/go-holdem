package main

import "fmt"

func main() {
	nums := []int{2, 4, 9, 1, 3, 4, 6}
	var answer [][]int
	combinations(nums, &answer, []int{}, 0, 0)
	fmt.Println(answer)

}

func combinations(nums []int, answer *[][]int, tmpList []int, startIdx int, lv int) {
	if lv == 5 {
		*answer = append(*answer, tmpList)
		return
	}

	for i := startIdx; i < len(nums); i++ {
		tmpList = append(tmpList, nums[i])
		combinations(nums, answer, tmpList, i+1, lv+1)
		tmpList = tmpList[:len(tmpList)-1]
	}
}
