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

type query struct {
	typ int
	p   int
	x   int
}

func run(bin string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveB(n int, a []int, queries []query) []int {
	val := make([]int, n+1)
	copy(val[1:], a)
	last := make([]int, n+1)
	payouts := make([]int, len(queries)+1)
	for i, q := range queries {
		idx := i + 1
		if q.typ == 1 {
			val[q.p] = q.x
			last[q.p] = idx
		} else {
			payouts[idx] = q.x
		}
	}
	suf := make([]int, len(queries)+2)
	for i := len(queries); i >= 1; i-- {
		if payouts[i] > suf[i+1] {
			suf[i] = payouts[i]
		} else {
			suf[i] = suf[i+1]
		}
	}
	res := make([]int, n)
	for i := 1; i <= n; i++ {
		mx := suf[last[i]+1]
		if val[i] < mx {
			val[i] = mx
		}
		res[i-1] = val[i]
	}
	return res
}

func generateCase(rng *rand.Rand) ([]byte, []int) {
	n := rng.Intn(5) + 1
	a := make([]int, n)
	for i := range a {
		a[i] = rng.Intn(30)
	}
	qn := rng.Intn(10) + 1
	queries := make([]query, qn)
	var b bytes.Buffer
	fmt.Fprintf(&b, "%d\n", n)
	for i, v := range a {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprint(&b, v)
	}
	b.WriteByte('\n')
	fmt.Fprintf(&b, "%d\n", qn)
	for i := 0; i < qn; i++ {
		if rng.Intn(2) == 0 {
			p := rng.Intn(n) + 1
			x := rng.Intn(30)
			queries[i] = query{typ: 1, p: p, x: x}
			fmt.Fprintf(&b, "1 %d %d\n", p, x)
		} else {
			x := rng.Intn(30)
			queries[i] = query{typ: 2, x: x}
			fmt.Fprintf(&b, "2 %d\n", x)
		}
	}
	expect := solveB(n, a, queries)
	return b.Bytes(), expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		input, expect := generateCase(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, string(input))
			os.Exit(1)
		}
		gotFields := strings.Fields(out)
		if len(gotFields) != len(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %v got %s\ninput:\n%s", i, expect, out, string(input))
			os.Exit(1)
		}
		ok := true
		for j, f := range gotFields {
			if fmt.Sprint(expect[j]) != f {
				ok = false
				break
			}
		}
		if !ok {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %v got %s\ninput:\n%s", i, expect, out, string(input))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
