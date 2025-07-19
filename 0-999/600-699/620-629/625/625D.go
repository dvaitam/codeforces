package main

import (
	"bufio"
	"fmt"
	"os"
)

func check(sum []int, n int) (bool, string) {
	ans := make([]byte, n)
	i := 0
	for i < n/2 {
		if sum[i] == sum[n-1-i] {
			i++
		} else if sum[i] == sum[n-1-i]+1 || sum[i] == sum[n-1-i]+11 {
			sum[i]--
			sum[i+1] += 10
		} else if sum[i] == sum[n-1-i]+10 {
			sum[n-2-i]--
			sum[n-1-i] += 10
		} else {
			return false, ""
		}
	}
	if n%2 == 1 {
		mid := n / 2
		if sum[mid]%2 != 0 || sum[mid] > 18 || sum[mid] < 0 {
			return false, ""
		}
		ans[mid] = byte(sum[mid]/2) + '0'
	}
	for j := 0; j < n/2; j++ {
		if sum[j] > 18 || sum[j] < 0 {
			return false, ""
		}
		ans[j] = byte((sum[j]+1)/2) + '0'
		ans[n-1-j] = byte(sum[j]/2) + '0'
	}
	if ans[0] <= '0' {
		return false, ""
	}
	return true, string(ans)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var s string
	if _, err := fmt.Fscan(reader, &s); err != nil {
		return
	}
	n := len(s)
	sum := make([]int, n)
	for i := 0; i < n; i++ {
		sum[i] = int(s[i] - '0')
	}
	if ok, res := check(append([]int{}, sum...), n); ok {
		fmt.Println(res)
		return
	}
	if n > 1 && s[0] == '1' {
		n1 := n - 1
		sum1 := make([]int, n1)
		for i := 0; i < n1; i++ {
			sum1[i] = int(s[i+1] - '0')
		}
		sum1[0] += 10
		if ok, res := check(sum1, n1); ok {
			fmt.Println(res)
			return
		}
	}
	fmt.Println("0")
}
