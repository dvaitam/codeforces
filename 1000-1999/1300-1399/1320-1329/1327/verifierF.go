package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const mod int64 = 998244353

func solve(n, k, m int, L, R, X []int) int64 {
	ans := int64(1)
	for bit := 0; bit < k; bit++ {
		diff := make([]int, n+2)
		left := make([]int, n+1)
		zeroSegL := make([]int, 0)
		zeroSegR := make([]int, 0)
		for i := 0; i < m; i++ {
			if (X[i]>>bit)&1 == 1 {
				diff[L[i]]++
				diff[R[i]+1]--
			} else {
				if left[R[i]] < L[i] {
					left[R[i]] = L[i]
				}
				zeroSegL = append(zeroSegL, L[i])
				zeroSegR = append(zeroSegR, R[i])
			}
		}
		forced := make([]bool, n+1)
		prefOnes := make([]int, n+1)
		cur := 0
		for i := 1; i <= n; i++ {
			cur += diff[i]
			if cur > 0 {
				forced[i] = true
				prefOnes[i] = prefOnes[i-1] + 1
			} else {
				prefOnes[i] = prefOnes[i-1]
			}
		}
		valid := true
		for idx := range zeroSegL {
			l := zeroSegL[idx]
			r := zeroSegR[idx]
			if prefOnes[r]-prefOnes[l-1] == r-l+1 {
				valid = false
				break
			}
			if left[r] < l {
				left[r] = l
			}
		}
		if !valid {
			ans = 0
			break
		}
		dp := make([]int64, n+1)
		dp[0] = 1
		pointer := 0
		sumValid := int64(1)
		requirement := 0
		for i := 1; i <= n; i++ {
			validPrev := sumValid
			if !forced[i] {
				dp[i] = validPrev % mod
			}
			sumValid = (sumValid + dp[i]) % mod
			if left[i] > requirement {
				requirement = left[i]
			}
			for pointer < requirement {
				sumValid -= dp[pointer]
				sumValid %= mod
				if sumValid < 0 {
					sumValid += mod
				}
				pointer++
			}
		}
		ans = ans * (sumValid % mod) % mod
	}
	return ans % mod
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	rand.Seed(6)
	t := 100
	for idx := 0; idx < t; idx++ {
		n := rand.Intn(6) + 1
		k := rand.Intn(4) + 1
		m := rand.Intn(6)
		L := make([]int, m)
		R := make([]int, m)
		X := make([]int, m)
		for i := 0; i < m; i++ {
			l := rand.Intn(n) + 1
			r := rand.Intn(n-l+1) + l
			L[i] = l
			R[i] = r
			X[i] = rand.Intn(1 << k)
		}
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d %d %d\n", n, k, m))
		for i := 0; i < m; i++ {
			input.WriteString(fmt.Sprintf("%d %d %d\n", L[i], R[i], X[i]))
		}
		want := fmt.Sprintf("%d\n", solve(n, k, m, L, R, X))

		cmd := exec.Command(binary)
		cmd.Stdin = strings.NewReader(input.String())
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			fmt.Printf("Runtime error on test %d: %v\n%s", idx+1, err, out.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != strings.TrimSpace(want) {
			fmt.Printf("Wrong answer on test %d\nExpected: %s\nGot: %s\n", idx+1, strings.TrimSpace(want), got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
