package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

const mod int64 = 998244353

func factorial(n int) []int64 {
	fact := make([]int64, n+1)
	fact[0] = 1
	for i := 1; i <= n; i++ {
		fact[i] = fact[i-1] * int64(i) % mod
	}
	return fact
}

type testCase struct {
	s string
}

func expected(tc testCase) (int, int64) {
	s := tc.s
	fact := factorial(len(s))
	dp0Len, dp1Len := 0, 0
	var dp0Cnt, dp1Cnt int64
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '0' {
			newLen := 1
			var newCnt int64 = 1
			if dp1Len+1 > newLen {
				newLen = dp1Len + 1
				newCnt = dp1Cnt
			} else if dp1Len+1 == newLen {
				newCnt = (newCnt + dp1Cnt) % mod
			}
			if newLen > dp0Len {
				dp0Len = newLen
				dp0Cnt = newCnt % mod
			} else if newLen == dp0Len {
				dp0Cnt = (dp0Cnt + newCnt) % mod
			}
		} else {
			newLen := 1
			var newCnt int64 = 1
			if dp0Len+1 > newLen {
				newLen = dp0Len + 1
				newCnt = dp0Cnt
			} else if dp0Len+1 == newLen {
				newCnt = (newCnt + dp0Cnt) % mod
			}
			if newLen > dp1Len {
				dp1Len = newLen
				dp1Cnt = newCnt % mod
			} else if newLen == dp1Len {
				dp1Cnt = (dp1Cnt + newCnt) % mod
			}
		}
	}
	L := dp0Len
	cnt := dp0Cnt % mod
	if dp1Len > L {
		L = dp1Len
		cnt = dp1Cnt % mod
	} else if dp1Len == L {
		cnt = (cnt + dp1Cnt) % mod
	}
	deletions := len(s) - L
	ans := cnt * fact[deletions] % mod
	return deletions, ans
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out, errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCases() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]testCase, 0, 102)
	cases = append(cases, testCase{s: "0"})
	cases = append(cases, testCase{s: "1"})
	for len(cases) < 102 {
		n := rng.Intn(20) + 1
		var sb strings.Builder
		for i := 0; i < n; i++ {
			if rng.Intn(2) == 0 {
				sb.WriteByte('0')
			} else {
				sb.WriteByte('1')
			}
		}
		cases = append(cases, testCase{s: sb.String()})
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genCases()
	for i, tc := range tests {
		input := fmt.Sprintf("1\n%s\n", tc.s)
		del, ans := expected(tc)
		want := fmt.Sprintf("%d %d", del, ans)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, want, got, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
