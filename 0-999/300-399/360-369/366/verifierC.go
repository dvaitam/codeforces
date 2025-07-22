package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func solve(input string) string {
	r := bufio.NewReader(strings.NewReader(input))
	var n, k int
	if _, err := fmt.Fscan(r, &n, &k); err != nil {
		return ""
	}
	a := make([]int, n)
	b := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(r, &a[i])
	}
	for i := 0; i < n; i++ {
		fmt.Fscan(r, &b[i])
	}
	deltas := make([]int, n)
	negSum, posSum := 0, 0
	for i := 0; i < n; i++ {
		del := a[i] - k*b[i]
		deltas[i] = del
		if del < 0 {
			negSum += del
		} else {
			posSum += del
		}
	}
	base := -negSum
	size := posSum + base + 1
	dp := make([]int, size)
	for i := range dp {
		dp[i] = -1
	}
	dp[base] = 0
	for i := 0; i < n; i++ {
		del := deltas[i]
		ai := a[i]
		if del >= 0 {
			for j := size - 1 - del; j >= 0; j-- {
				if dp[j] >= 0 {
					nj := j + del
					v := dp[j] + ai
					if v > dp[nj] {
						dp[nj] = v
					}
				}
			}
		} else {
			for j := -del; j < size; j++ {
				if dp[j] >= 0 {
					nj := j + del
					v := dp[j] + ai
					if v > dp[nj] {
						dp[nj] = v
					}
				}
			}
		}
	}
	res := dp[base]
	if res > 0 {
		return fmt.Sprintf("%d\n", res)
	}
	return "-1\n"
}

func genTest() (string, string) {
	n := rand.Intn(7) + 1
	k := rand.Intn(5) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, k)
	for i := 0; i < n; i++ {
		fmt.Fprintf(&sb, "%d ", rand.Intn(20)+1)
	}
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		fmt.Fprintf(&sb, "%d ", rand.Intn(20)+1)
	}
	sb.WriteByte('\n')
	inp := sb.String()
	out := solve(inp)
	return inp, out
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		in, exp := genTest()
		got, err := runBinary(bin, in)
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\nInput:\n%sOutput:\n%s\n", i+1, err, in, got)
			return
		}
		if got != strings.TrimSpace(exp) {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%sGot:\n%s\n", i+1, in, exp, got)
			return
		}
	}
	fmt.Println("All tests passed")
}
