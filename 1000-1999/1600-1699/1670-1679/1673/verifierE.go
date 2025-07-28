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

func expected(t testCase) string {
	n, k := t.n, t.k
	B := t.b
	if n > 20 {
		return "0"
	}
	var ans uint64
	var dfs func(pos, cnt int, curProd int, curVal uint64)
	dfs = func(pos, cnt int, curProd int, curVal uint64) {
		if pos == n-1 {
			if curProd < LIM {
				curVal ^= 1 << curProd
			}
			if cnt >= k {
				ans ^= curVal
			}
			return
		}
		nextProd := curProd * B[pos+1]
		dfs(pos+1, cnt, nextProd, curVal)
		val := curVal
		if curProd < LIM {
			val ^= 1 << curProd
		}
		dfs(pos+1, cnt+1, B[pos+1], val)
	}
	dfs(0, 0, B[0], 0)
	if ans == 0 {
		return "0"
	}
	out := ""
	for ans > 0 {
		if ans&1 == 1 {
			out = "1" + out
		} else {
			out = "0" + out
		}
		ans >>= 1
	}
	return out
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
