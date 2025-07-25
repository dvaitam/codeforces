package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func orientSmallEdges(a []int64) []int64 {
	n := len(a)
	res := make([]int64, n)
	l, r := 0, n-1
	idx := 0
	for l < r {
		res[l] = a[idx]
		res[r] = a[idx+1]
		l++
		r--
		idx += 2
	}
	if l == r {
		res[l] = a[idx]
	}
	return res
}

func orientAlternate(a []int64) []int64 {
	n := len(a)
	res := make([]int64, n)
	l, r := 0, n-1
	i := 0
	for l < r {
		res[l] = a[n-1-i]
		res[r] = a[n-1-i]
		l++
		r--
		i++
		if l > r {
			break
		}
		res[l] = a[i-1]
		res[r] = a[i-1]
		l++
		r--
	}
	if l == r {
		res[l] = a[n/2]
	}
	return res
}

func valid(arr []int64, nums []int64) bool {
	n := len(arr)
	sub := make([]int64, 0, n*(n+1)/2)
	for i := 0; i < n; i++ {
		s := int64(0)
		for j := i; j < n; j++ {
			s += arr[j]
			sub = append(sub, s)
		}
	}
	sort.Slice(sub, func(i, j int) bool { return sub[i] < sub[j] })
	temp := append([]int64(nil), nums...)
	sort.Slice(temp, func(i, j int) bool { return temp[i] < temp[j] })
	i, j, diff := 0, 0, 0
	for i < len(sub) && j < len(temp) {
		if sub[i] == temp[j] {
			i++
			j++
		} else {
			diff++
			i++
		}
	}
	diff += len(sub) - i
	diff += len(temp) - j
	return diff == 1
}

func reconstruct(n int, nums []int64) []int64 {
	sort.Slice(nums, func(i, j int) bool { return nums[i] < nums[j] })
	maxv := nums[len(nums)-1]
	countMax := 0
	for i := len(nums) - 1; i >= 0 && nums[i] == maxv; i-- {
		countMax++
	}
	var arrElems []int64
	if countMax == 1 {
		total := maxv
		rest := nums[:len(nums)-1]
		arrElems = append([]int64(nil), rest[:n-1]...)
		sum := int64(0)
		for _, v := range arrElems {
			sum += v
		}
		arrElems = append(arrElems, total-sum)
	} else {
		arrElems = append([]int64(nil), nums[:n]...)
	}
	sort.Slice(arrElems, func(i, j int) bool { return arrElems[i] < arrElems[j] })
	a1 := orientSmallEdges(arrElems)
	if valid(a1, nums) {
		return a1
	}
	a2 := orientAlternate(arrElems)
	return a2
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		m := n*(n+1)/2 - 1
		nums := make([]int64, m)
		for i := 0; i < m; i++ {
			fmt.Fscan(in, &nums[i])
		}
		ans := reconstruct(n, nums)
		for i, v := range ans {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, v)
		}
		fmt.Fprintln(out)
	}
}
