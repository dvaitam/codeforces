package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type op struct {
	t int
	l int
	r int
}

func mex(set map[int]bool) int {
	m := 1
	for {
		if !set[m] {
			return m
		}
		m++
	}
}

func simulate(ops []op) []int {
	set := map[int]bool{}
	res := []int{}
	for _, op := range ops {
		switch op.t {
		case 1:
			for i := op.l; i <= op.r; i++ {
				set[i] = true
			}
		case 2:
			for i := op.l; i <= op.r; i++ {
				delete(set, i)
			}
		case 3:
			for i := op.l; i <= op.r; i++ {
				set[i] = !set[i]
			}
		}
		res = append(res, mex(set))
	}
	return res
}

func runBin(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func genCase(rng *rand.Rand) (string, string) {
	q := rng.Intn(20) + 1
	ops := make([]op, 0, q)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", q))
	for i := 0; i < q; i++ {
		t := rng.Intn(3) + 1
		l := rng.Intn(50) + 1
		r := rng.Intn(50) + 1
		if l > r {
			l, r = r, l
		}
		ops = append(ops, op{t: t, l: l, r: r})
		sb.WriteString(fmt.Sprintf("%d %d %d\n", t, l, r))
	}
	results := simulate(ops)
	var exp strings.Builder
	for i, v := range results {
		if i > 0 {
			exp.WriteByte('\n')
		}
		exp.WriteString(strconv.Itoa(v))
	}
	return sb.String(), exp.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, exp := genCase(rng)
		out, err := runBin(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\ninput:\n%s\n", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed:\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", i+1, input, exp, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
