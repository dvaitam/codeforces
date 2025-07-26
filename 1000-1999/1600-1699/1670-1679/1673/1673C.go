package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int = 1e9 + 7
const MAXN int = 40000

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	pals := generatePalindromes(MAXN)
	dp := make([]int, MAXN+1)
	dp[0] = 1
	for _, v := range pals {
		for i := v; i <= MAXN; i++ {
			dp[i] += dp[i-v]
			if dp[i] >= MOD {
				dp[i] -= MOD
			}
		}
	}

	var t, n int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		fmt.Fscan(reader, &n)
		fmt.Fprintln(writer, dp[n])
	}
}

func generatePalindromes(limit int) []int {
	res := make([]int, 0)
	for i := 1; i <= limit; i++ {
		if isPalindrome(i) {
			res = append(res, i)
		}
	}
	return res
}

func isPalindrome(x int) bool {
	orig := x
	rev := 0
	for x > 0 {
		rev = rev*10 + x%10
		x /= 10
	}
	return orig == rev
}
