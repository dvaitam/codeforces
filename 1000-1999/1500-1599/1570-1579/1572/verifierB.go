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
		a := make([]int, n+1)
		pos := 0
		for i := 1; i <= n; i++ {
			fmt.Fscan(reader, &a[i])
			a[i] ^= a[i-1]
			if i%2 == 1 && a[i] == 0 {
				pos = i
			}
		}
		if pos != 0 && a[n] == 0 {
			out.WriteString("YES\n")
			m := n - 1
			if n != pos {
				m--
			}
			out.WriteString(fmt.Sprintln(m))
			for i := pos + 1; i+2 <= n; i += 2 {
				out.WriteString(fmt.Sprintf("%d ", i))
			}
			for i := pos; ; {
				i -= 2
				if i < 1 {
					break
				}
				out.WriteString(fmt.Sprintf("%d ", i))
			}
			for i := 1; i+2 <= n; i += 2 {
				out.WriteString(fmt.Sprintf("%d ", i))
			}
			out.WriteByte('\n')
		} else {
			out.WriteString("NO\n")
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
