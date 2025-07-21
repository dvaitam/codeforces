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

type testCaseE struct {
	n      int
	digits []int
}

func generateCaseE(rng *rand.Rand) testCaseE {
	n := rng.Intn(6) + 1
	digits := make([]int, 2*n)
	for i := range digits {
		digits[i] = rng.Intn(10)
	}
	return testCaseE{n: n, digits: digits}
}

func maxTotal(digits []int, n int) int64 {
	p10 := make([]int64, n+1)
	p10[0] = 1
	for i := 1; i <= n; i++ {
		p10[i] = p10[i-1] * 10
	}
	f := make([][]int64, n+1)
	for i := range f {
		f[i] = make([]int64, n+1)
	}
	for i := 0; i <= n; i++ {
		for j := 0; j <= n; j++ {
			if i < n {
				idx := i + j
				w := int64(digits[idx]) * p10[n-i-1]
				if v := f[i][j] + w; v > f[i+1][j] {
					f[i+1][j] = v
				}
			}
			if j < n {
				idx := i + j
				w := int64(digits[idx]) * p10[n-j-1]
				if v := f[i][j] + w; v > f[i][j+1] {
					f[i][j+1] = v
				}
			}
		}
	}
	return f[n][n]
}

func totalForOrder(order string, digits []int, n int) (int64, bool) {
	if len(order) != 2*n {
		return 0, false
	}
	hCnt := 0
	mCnt := 0
	var s1, s2 int64
	for idx, c := range order {
		d := int64(digits[idx])
		if c == 'H' {
			if hCnt >= n {
				return 0, false
			}
			s1 = s1*10 + d
			hCnt++
		} else if c == 'M' {
			if mCnt >= n {
				return 0, false
			}
			s2 = s2*10 + d
			mCnt++
		} else {
			return 0, false
		}
	}
	if hCnt != n || mCnt != n {
		return 0, false
	}
	return s1 + s2, true
}

func runCaseE(bin string, tc testCaseE) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", tc.n))
	for _, d := range tc.digits {
		sb.WriteByte(byte('0' + d))
	}
	sb.WriteByte('\n')
	input := sb.String()

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	order := strings.TrimSpace(out.String())
	got, ok := totalForOrder(order, tc.digits, tc.n)
	if !ok {
		return fmt.Errorf("invalid output format")
	}
	if got != maxTotal(tc.digits, tc.n) {
		return fmt.Errorf("output does not achieve maximum")
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCaseE(rng)
		if err := runCaseE(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
