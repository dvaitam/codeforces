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

type testCaseF struct {
	n   int
	q   int
	ops []op
}

type op struct {
	p int
	a int
	b int
}

func generateCaseF(rng *rand.Rand) testCaseF {
	n := rng.Intn(5) + 1
	q := rng.Intn(6) + 1
	ops := make([]op, q)
	for i := 0; i < q; i++ {
		if rng.Intn(2) == 0 {
			c := rng.Intn(n) + 1
			g := rng.Intn(n-c+1) + c
			ops[i] = op{p: 1, a: c, b: g}
		} else {
			l := rng.Intn(n) + 1
			r := rng.Intn(n-l+1) + l
			ops[i] = op{p: 2, a: l, b: r}
		}
	}
	return testCaseF{n: n, q: q, ops: ops}
}

func buildInputF(t testCaseF) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", t.n, t.q))
	for _, op := range t.ops {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", op.p, op.a, op.b))
	}
	return sb.String()
}

func solveF(reader *bufio.Reader) string {
	var n, q int
	fmt.Fscan(reader, &n, &q)
	h := make([]int, n+1)
	w := make([]int, n+1)
	for i := 1; i <= n; i++ {
		h[i] = 0
		w[i] = i
	}
	maxH := 0
	var out strings.Builder
	for ; q > 0; q-- {
		var p int
		fmt.Fscan(reader, &p)
		if p == 1 {
			var c, g int
			fmt.Fscan(reader, &c, &g)
			maxH++
			h[c] = maxH
			w[c] = g
		} else {
			var l, r int
			fmt.Fscan(reader, &l, &r)
			b := make([]int, n+1)
			for j := 1; j <= n; j++ {
				count := 0
				for i := 1; i <= j; i++ {
					if j <= w[i] {
						valid := true
						hi := h[i]
						for k := i + 1; k <= j; k++ {
							if h[k] >= hi {
								valid = false
								break
							}
						}
						if valid {
							count++
						}
					}
				}
				b[j] = count
			}
			sum := 0
			for j := l; j <= r; j++ {
				sum += b[j]
			}
			out.WriteString(fmt.Sprintln(sum))
		}
	}
	return strings.TrimSpace(out.String())
}

func expectedF(t testCaseF) string {
	input := buildInputF(t)
	return solveF(bufio.NewReader(strings.NewReader(input)))
}

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCaseF(rng)
		input := buildInputF(tc)
		expect := expectedF(tc)
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
