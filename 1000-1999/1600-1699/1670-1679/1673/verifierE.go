package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

const LIM = 1 << 20

type testCase struct {
	n int
	k int
	b []int
}

func generateTests() []testCase {
	r := rand.New(rand.NewSource(5))
	tests := []testCase{{n: 1, k: 1, b: []int{1}}, {n: 2, k: 1, b: []int{1, 2}}}
	for len(tests) < 100 {
		n := r.Intn(5) + 1
		k := r.Intn(n) + 1
		b := make([]int, n)
		for i := 0; i < n; i++ {
			b[i] = r.Intn(5) + 1
		}
		tests = append(tests, testCase{n, k, b})
	}
	return tests
}

func f(n, m int) int {
	if m <= 0 {
		if n == 0 {
			return 1
		}
		return 0
	}
	if n < m {
		return 0
	}
	if ((n - 1) & (m - 1)) == (m - 1) {
		return 1
	}
	return 0
}

func expected(t testCase) string {
	n, k := t.n, t.k
	a := make([]int, n+1)
	for i, v := range t.b {
		a[i+1] = v
	}
	ans := make([]byte, LIM)
	maxBit := -1
	for i := 1; i <= n; i++ {
		s := 0
		for j := i; j <= n; {
			if s >= 20 || (LIM-1)>>s < a[i] {
				break
			}
			idx := a[i] << s
			t := 0
			if i > 1 {
				t++
			}
			if j < n {
				t++
			}
			ans[idx] ^= byte(f(n-1-(j-i)-t, k-t))
			if ans[idx] == 1 && idx > maxBit {
				maxBit = idx
			}
			j++
			if j <= n {
				s += a[j]
			}
		}
	}
	for maxBit > 0 && ans[maxBit] == 0 {
		maxBit--
	}
	if maxBit < 0 || (maxBit == 0 && ans[0] == 0) {
		return "0"
	}
	out := make([]byte, maxBit+1)
	for i := maxBit; i >= 0; i-- {
		if ans[i] == 1 {
			out[maxBit-i] = '1'
		} else {
			out[maxBit-i] = '0'
		}
	}
	return string(out)
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("timeout")
		}
		return "", fmt.Errorf("run failed: %v\n%s", err, errb.String())
	}
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("timeout")
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		input := fmt.Sprintf("%d %d\n", t.n, t.k)
		for j, v := range t.b {
			if j > 0 {
				input += " "
			}
			input += fmt.Sprintf("%d", v)
		}
		input += "\n"
		want := expected(t)
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Printf("test %d failed: expected %s got %s\n", i+1, want, strings.TrimSpace(got))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
