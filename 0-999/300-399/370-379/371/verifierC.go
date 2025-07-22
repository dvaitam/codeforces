package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type test struct {
	input    string
	expected string
}

func solve(input string) string {
	reader := strings.NewReader(input)
	var recipe string
	if _, err := fmt.Fscan(reader, &recipe); err != nil {
		return ""
	}
	var nb, ns, nc int64
	fmt.Fscan(reader, &nb, &ns, &nc)
	var pb, ps, pc int64
	fmt.Fscan(reader, &pb, &ps, &pc)
	var r int64
	fmt.Fscan(reader, &r)
	var needB, needS, needC int64
	for _, ch := range recipe {
		switch ch {
		case 'B':
			needB++
		case 'S':
			needS++
		case 'C':
			needC++
		}
	}
	can := func(x int64) bool {
		cost := int64(0)
		if req := needB*x - nb; req > 0 {
			cost += req * pb
		}
		if req := needS*x - ns; req > 0 {
			cost += req * ps
		}
		if req := needC*x - nc; req > 0 {
			cost += req * pc
		}
		return cost <= r
	}
	lo, hi := int64(0), int64(1e13)
	ans := int64(0)
	for lo <= hi {
		mid := (lo + hi) / 2
		if can(mid) {
			ans = mid
			lo = mid + 1
		} else {
			hi = mid - 1
		}
	}
	return fmt.Sprintf("%d", ans)
}

func generateTests() []test {
	rand.Seed(44)
	var tests []test
	fixed := []string{
		"BSC\n1 1 1\n1 1 1\n0\n",
		"BB\n0 0 0\n1 1 1\n10\n",
	}
	for _, in := range fixed {
		tests = append(tests, test{in, solve(in)})
	}
	for len(tests) < 100 {
		l := rand.Intn(5) + 1
		var sb strings.Builder
		for i := 0; i < l; i++ {
			ch := "BSC"[rand.Intn(3)]
			sb.WriteByte(ch)
		}
		recipe := sb.String()
		nb := rand.Int63n(10)
		ns := rand.Int63n(10)
		nc := rand.Int63n(10)
		pb := rand.Int63n(10) + 1
		ps := rand.Int63n(10) + 1
		pc := rand.Int63n(10) + 1
		r := rand.Int63n(1000)
		inp := fmt.Sprintf("%s\n%d %d %d\n%d %d %d\n%d\n", recipe, nb, ns, nc, pb, ps, pc, r)
		tests = append(tests, test{inp, solve(inp)})
	}
	return tests
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
	out, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("time limit")
	}
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, out)
	}
	return strings.TrimSpace(string(out)), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		out, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(t.expected) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\ninput:\n%sexpected:%s\n got:%s\n", i+1, t.input, t.expected, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
