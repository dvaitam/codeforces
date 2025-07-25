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

type Query struct {
	t, d int64
}

func solveB(n, b int, qs []Query) string {
	q := make([]int64, 0)
	ans := make([]int64, n)
	for i := 0; i < n; i++ {
		t := qs[i].t
		d := qs[i].d
		for len(q) > 0 && q[0] <= t {
			q = q[1:]
		}
		if len(q) > b {
			ans[i] = -1
			continue
		}
		var start int64
		if len(q) == 0 {
			start = t
		} else {
			start = q[len(q)-1]
		}
		finish := start + d
		q = append(q, finish)
		ans[i] = finish
	}
	var sb strings.Builder
	for i, v := range ans {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	b := rng.Intn(5) + 1
	qs := make([]Query, n)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, b)
	cur := int64(rng.Intn(5))
	for i := 0; i < n; i++ {
		cur += int64(rng.Intn(10) + 1)
		d := int64(rng.Intn(10) + 1)
		qs[i] = Query{cur, d}
		fmt.Fprintf(&sb, "%d %d\n", cur, d)
	}
	input := sb.String()
	expected := solveB(n, b, qs)
	return input, expected
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if got != exp {
		return fmt.Errorf("expected:\n%s\ngot:\n%s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
