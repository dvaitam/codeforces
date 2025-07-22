package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

var k int
var strs []string

func dfs(l, r, d int) ([]int64, int) {
	total := 0
	dp := make([]int64, 1)
	for i := l; i < r && len(strs[i]) == d; i++ {
	}
	i := l
	for i < r && len(strs[i]) == d {
		i++
	}
	leafCount := i - l
	if leafCount > 0 {
		newTotal := leafCount
		if newTotal > k {
			newTotal = k
		}
		newDp := make([]int64, newTotal+1)
		for t1 := 0; t1 <= leafCount && t1 <= k; t1++ {
			newDp[t1] = 0
		}
		dp = newDp
		total = newTotal
	}
	for j := i; j < r; {
		ch := strs[j][d]
		m := j + 1
		for m < r && len(strs[m]) > d && strs[m][d] == ch {
			m++
		}
		childDp, childCnt := dfs(j, m, d+1)
		newTotal := total + childCnt
		if newTotal > k {
			newTotal = k
		}
		newDp := make([]int64, newTotal+1)
		for t0 := 0; t0 <= total; t0++ {
			for t1 := 0; t1 <= childCnt && t0+t1 <= k; t1++ {
				val := dp[t0] + childDp[t1]
				if val > newDp[t0+t1] {
					newDp[t0+t1] = val
				}
			}
		}
		dp = newDp
		total = newTotal
		j = m
	}
	if d > 0 {
		for x := 0; x <= total; x++ {
			dp[x] += int64(x * (x - 1) / 2)
		}
	}
	return dp, total
}

func solveF(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	var n int
	fmt.Fscan(reader, &n, &k)
	strs = make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &strs[i])
	}
	sort.Strings(strs)
	dp, _ := dfs(0, n, 0)
	if k < len(dp) {
		return fmt.Sprintf("%d", dp[k])
	}
	return "0"
}

type testF struct{ input, expect string }

func genTests() []testF {
	rand.Seed(42)
	var tests []testF
	letters := "abc"
	for i := 0; i < 100; i++ {
		n := rand.Intn(5) + 1
		kVal := rand.Intn(n) + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, kVal))
		for j := 0; j < n; j++ {
			l := rand.Intn(3) + 1
			var s strings.Builder
			for t := 0; t < l; t++ {
				s.WriteByte(letters[rand.Intn(len(letters))])
			}
			sb.WriteString(s.String() + "\n")
		}
		input := sb.String()
		k = kVal
		expect := solveF(input)
		tests = append(tests, testF{input, expect})
	}
	return tests
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()
	for i, t := range tests {
		got, err := run(bin, t.input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != t.expect {
			fmt.Printf("test %d failed\ninput:\n%sexpect:%s\nactual:%s\n", i+1, t.input, t.expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
