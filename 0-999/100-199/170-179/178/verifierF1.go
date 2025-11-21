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

var (
	k    int
	strs []string
)

// dfs mirrors the reference solution logic from 178F1.go.
func dfs(l, r, d int) ([]int64, int) {
	total := 0
	dp := make([]int64, 1)
	dp[0] = 0

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
		for picked := 0; picked <= leafCount && picked <= k; picked++ {
			newDp[picked] = 0
		}
		dp = newDp
		total = newTotal
	}

	for j := i; j < r; {
		next := j + 1
		for next < r && len(strs[next]) > d && strs[next][d] == strs[j][d] {
			next++
		}
		childDp, childCnt := dfs(j, next, d+1)
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
		j = next
	}

	if d > 0 {
		for x := 0; x <= total; x++ {
			dp[x] += int64(x * (x - 1) / 2)
		}
	}

	return dp, total
}

func solveRef(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	var n int
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return ""
	}
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

type testCase struct {
	name   string
	input  string
	expect string
}

func handcraftedTests() []testCase {
	var tests []testCase

	tests = append(tests, makeCase("single_string", []string{"a"}, 1))
	tests = append(tests, makeCase("no_common_prefix", []string{"abc", "xyz", "def"}, 2))
	tests = append(tests, makeCase("all_equal", []string{"aaa", "aaa", "aaa"}, 3))
	tests = append(tests, makeCase("mix_lengths", []string{"a", "ab", "abc", "abcd"}, 3))
	tests = append(tests, makeCase("high_k_small_n", []string{"abc", "abd", "acd"}, 3))

	return tests
}

func makeCase(name string, arr []string, kVal int) testCase {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", len(arr), kVal))
	for i, s := range arr {
		sb.WriteString(s)
		if i+1 < len(arr) {
			sb.WriteByte('\n')
		}
	}
	sb.WriteByte('\n')
	input := sb.String()
	expect := solveRef(input)
	return testCase{name: name, input: input, expect: expect}
}

func genRandomTests() []testCase {
	r := rand.New(rand.NewSource(178))
	letters := []byte("abcdxyz")
	var tests []testCase

	appendRandom := func(name string, cnt, nLo, nHi, lenLo, lenHi int) {
		for i := 0; i < cnt; i++ {
			n := nLo + r.Intn(nHi-nLo+1)
			if n <= 0 {
				n = 1
			}
			kVal := 1 + r.Intn(n)
			var sb strings.Builder
			sb.WriteString(fmt.Sprintf("%d %d\n", n, kVal))
			for j := 0; j < n; j++ {
				curLen := lenLo + r.Intn(lenHi-lenLo+1)
				if curLen <= 0 {
					curLen = 1
				}
				var s strings.Builder
				for t := 0; t < curLen; t++ {
					s.WriteByte(letters[r.Intn(len(letters))])
				}
				sb.WriteString(s.String())
				sb.WriteByte('\n')
			}
			input := sb.String()
			expect := solveRef(input)
			tests = append(tests, testCase{
				name:   fmt.Sprintf("%s_%d", name, i+1),
				input:  input,
				expect: expect,
			})
		}
	}

	appendRandom("small", 120, 1, 7, 1, 4)
	appendRandom("medium", 80, 5, 20, 1, 8)
	appendRandom("large", 30, 50, 200, 1, 20)

	return tests
}

func runCandidate(bin, input string) (string, error) {
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
		fmt.Println("usage: go run verifierF1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := append(handcraftedTests(), genRandomTests()...)
	for idx, tc := range tests {
		got, err := runCandidate(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d (%s) runtime error: %v\ninput:\n%s\n", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		if got != tc.expect {
			fmt.Printf("test %d (%s) failed\ninput:\n%s\nexpect:%s\nactual:%s\n", idx+1, tc.name, tc.input, tc.expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
