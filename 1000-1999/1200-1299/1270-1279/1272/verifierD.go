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

type testD struct {
	a []int
}

func genTestsD() []testD {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testD, 100)
	for i := range tests {
		n := r.Intn(100) + 1
		arr := make([]int, n)
		for j := range arr {
			arr[j] = r.Intn(1_000_000_000) + 1
		}
		tests[i] = testD{a: arr}
	}
	return tests
}

func solveD(tc testD) int {
	n := len(tc.a)
	if n == 0 {
		return 0
	}
	left := make([]int, n)
	left[0] = 1
	for i := 1; i < n; i++ {
		if tc.a[i] > tc.a[i-1] {
			left[i] = left[i-1] + 1
		} else {
			left[i] = 1
		}
	}
	right := make([]int, n)
	right[n-1] = 1
	for i := n - 2; i >= 0; i-- {
		if tc.a[i] < tc.a[i+1] {
			right[i] = right[i+1] + 1
		} else {
			right[i] = 1
		}
	}
	ans := 1
	for i := 0; i < n; i++ {
		if left[i] > ans {
			ans = left[i]
		}
	}
	if n > 1 {
		if right[1] > ans {
			ans = right[1]
		}
		if left[n-2] > ans {
			ans = left[n-2]
		}
	}
	for i := 1; i+1 < n; i++ {
		if tc.a[i+1] > tc.a[i-1] {
			if left[i-1]+right[i+1] > ans {
				ans = left[i-1] + right[i+1]
			}
		}
	}
	if ans > n {
		ans = n
	}
	return ans
}

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", bin)
	} else {
		cmd = exec.CommandContext(ctx, bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("timeout")
	}
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsD()
	for i, tc := range tests {
		var input strings.Builder
		fmt.Fprintf(&input, "%d\n", len(tc.a))
		for j, v := range tc.a {
			if j > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprintf(&input, "%d", v)
		}
		input.WriteByte('\n')
		expected := fmt.Sprintf("%d", solveD(tc))
		out, err := runBinary(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != expected {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected:%s\ngot:%s\n", i+1, input.String(), expected, out)
			os.Exit(1)
		}
	}
	fmt.Println("Accepted")
}
