package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCaseC struct {
	n   int
	arr []int
}

func generateCaseC(rng *rand.Rand) testCaseC {
	n := rng.Intn(10) + 1
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(5)
	}
	return testCaseC{n: n, arr: arr}
}

func buildInputC(t testCaseC) string {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprint(t.n))
	sb.WriteByte('\n')
	for i, v := range t.arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func solveC(reader *bufio.Reader) string {
	var T int
	fmt.Fscan(reader, &T)
	out := strings.Builder{}
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(reader, &n)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
		}
		b := make([]int, 0, n)
		for _, v := range arr {
			if len(b) == 0 || b[len(b)-1] != v {
				b = append(b, v)
			}
		}
		m := len(b)
		if m == 1 {
			out.WriteString("0\n")
			continue
		}
		pos := make(map[int][]int)
		for i, v := range b {
			pos[v] = append(pos[v], i)
		}
		dp := make([][]int, m)
		for i := range dp {
			dp[i] = make([]int, m)
		}
		for length := 2; length <= m; length++ {
			for l := 0; l+length-1 < m; l++ {
				r := l + length - 1
				best := dp[l+1][r] + 1
				val := b[l]
				for _, k := range pos[val] {
					if k <= l || k > r {
						continue
					}
					cand := dp[l][k-1] + dp[k][r]
					if cand < best {
						best = cand
					}
				}
				dp[l][r] = best
			}
		}
		out.WriteString(fmt.Sprintf("%d\n", dp[0][m-1]))
	}
	return strings.TrimSpace(out.String())
}

func expectedC(t testCaseC) string {
	input := buildInputC(t)
	return solveC(bufio.NewReader(strings.NewReader(input)))
}

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", bin)
	} else {
		cmd = exec.CommandContext(ctx, bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return out.String(), fmt.Errorf("timeout")
	}
	if err != nil {
		return out.String(), err
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCaseC(rng)
		input := buildInputC(tc)
		expect := expectedC(tc)
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\nOutput:%s", i+1, err, out)
			os.Exit(1)
		}
		got := strings.TrimSpace(out)
		exp := strings.TrimSpace(expect)
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
