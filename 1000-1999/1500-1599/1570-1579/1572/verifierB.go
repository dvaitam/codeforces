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

type testCaseB struct {
	n   int
	arr []int
}

func generateCaseB(rng *rand.Rand) testCaseB {
	n := rng.Intn(8) + 3 // 3..10
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(2)
	}
	return testCaseB{n: n, arr: arr}
}

func buildInputB(t testCaseB) string {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprint(t.n))
	sb.WriteString("\n")
	for i, v := range t.arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteString("\n")
	return sb.String()
}

func solveB(reader *bufio.Reader) string {
	var T int
	fmt.Fscan(reader, &T)
	out := strings.Builder{}
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		totalXor := 0
		hasZero := false
		for _, v := range a {
			totalXor ^= v
			if v == 0 {
				hasZero = true
			}
		}
		if totalXor != 0 || !hasZero {
			out.WriteString("NO\n")
			continue
		}
		cnt := make([]int, n-2)
		for i := 0; i < n-2; i++ {
			cnt[i] = a[i] + a[i+1] + a[i+2]
		}
		queue := make([]int, 0, n)
		inq := make([]bool, n-2)
		for i := 0; i < n-2; i++ {
			if cnt[i] == 1 || cnt[i] == 2 {
				queue = append(queue, i)
				inq[i] = true
			}
		}
		ops := make([]int, 0, n)
		head := 0
		for head < len(queue) && len(ops) <= n {
			i := queue[head]
			head++
			if i < 0 || i >= n-2 {
				continue
			}
			c := cnt[i]
			if c != 1 && c != 2 {
				continue
			}
			p := c & 1
			ops = append(ops, i+1)
			for j := i; j < i+3; j++ {
				a[j] = p
			}
			for k := i - 2; k <= i+2; k++ {
				if k >= 0 && k < n-2 {
					cnt[k] = a[k] + a[k+1] + a[k+2]
					if !inq[k] && (cnt[k] == 1 || cnt[k] == 2) {
						queue = append(queue, k)
						inq[k] = true
					}
				}
			}
		}
		ok := true
		for _, v := range a {
			if v != 0 {
				ok = false
				break
			}
		}
		if !ok || len(ops) > n {
			out.WriteString("NO\n")
		} else {
			out.WriteString("YES\n")
			out.WriteString(fmt.Sprintln(len(ops)))
			for i, v := range ops {
				if i > 0 {
					out.WriteByte(' ')
				}
				out.WriteString(fmt.Sprint(v))
			}
			out.WriteByte('\n')
		}
	}
	return strings.TrimSpace(out.String())
}

func expectedB(t testCaseB) string {
	input := buildInputB(t)
	return solveB(bufio.NewReader(strings.NewReader(input)))
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCaseB(rng)
		input := buildInputB(tc)
		expect := expectedB(tc)
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\nOutput:%s", i+1, err, out)
			os.Exit(1)
		}
		got := strings.TrimSpace(out)
		exp := strings.TrimSpace(expect)
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed:\nexpected:\n%s\ngot:\n%s\ninput:\n%s", i+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
