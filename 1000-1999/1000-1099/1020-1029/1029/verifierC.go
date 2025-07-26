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

type TestCase struct {
	Input  string
	Output string
}

type Seg struct {
	l int
	r int
}

func runBinary(bin string, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("timeout")
		}
		return "", fmt.Errorf("%v: %s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveC(segs []Seg) string {
	n := len(segs)
	const INF = int(1e9 + 5)
	L := make([]int, n)
	R := make([]int, n)
	for i := 0; i < n; i++ {
		if i == 0 {
			L[i] = segs[i].l
			R[i] = segs[i].r
		} else {
			if segs[i].l > L[i-1] {
				L[i] = segs[i].l
			} else {
				L[i] = L[i-1]
			}
			if segs[i].r < R[i-1] {
				R[i] = segs[i].r
			} else {
				R[i] = R[i-1]
			}
		}
	}
	Ls := make([]int, n)
	Rs := make([]int, n)
	for i := n - 1; i >= 0; i-- {
		if i == n-1 {
			Ls[i] = segs[i].l
			Rs[i] = segs[i].r
		} else {
			if segs[i].l > Ls[i+1] {
				Ls[i] = segs[i].l
			} else {
				Ls[i] = Ls[i+1]
			}
			if segs[i].r < Rs[i+1] {
				Rs[i] = segs[i].r
			} else {
				Rs[i] = Rs[i+1]
			}
		}
	}
	ans := 0
	for i := 0; i < n; i++ {
		l := -INF
		r := INF
		if i > 0 {
			if L[i-1] > l {
				l = L[i-1]
			}
			if R[i-1] < r {
				r = R[i-1]
			}
		}
		if i+1 < n {
			if Ls[i+1] > l {
				l = Ls[i+1]
			}
			if Rs[i+1] < r {
				r = Rs[i+1]
			}
		}
		if r-l > ans {
			ans = r - l
		}
	}
	if ans < 0 {
		ans = 0
	}
	return fmt.Sprintf("%d", ans)
}

func generateTests() []TestCase {
	rand.Seed(44)
	tests := make([]TestCase, 100)
	for t := 0; t < 100; t++ {
		n := rand.Intn(20) + 2
		segs := make([]Seg, n)
		inputBuilder := strings.Builder{}
		inputBuilder.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			l := rand.Intn(100)
			r := l + rand.Intn(20)
			segs[i] = Seg{l, r}
			inputBuilder.WriteString(fmt.Sprintf("%d %d\n", l, r))
		}
		output := solveC(segs)
		tests[t] = TestCase{Input: inputBuilder.String(), Output: output}
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		got, err := runBinary(bin, tc.Input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != strings.TrimSpace(tc.Output) {
			fmt.Fprintf(os.Stderr, "Test %d failed:\ninput:\n%s\nexpected:%s\n got:%s\n", i+1, tc.Input, tc.Output, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
