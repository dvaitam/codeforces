package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type test struct {
	n int
	k int64
	a []int64
	x []int64
}

func solve(n int, k int64, a, x []int64) string {
	type pair struct {
		t int64
		a int64
	}
	p := make([]pair, n)
	for i := 0; i < n; i++ {
		t := x[i]
		if t < 0 {
			t = -t
		}
		p[i] = pair{t, a[i]}
	}
	sort.Slice(p, func(i, j int) bool { return p[i].t < p[j].t })
	var sum int64
	for _, pr := range p {
		sum += pr.a
		if sum > k*pr.t {
			return "NO"
		}
	}
	return "YES"
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(43))
	tests := make([]test, 0, 100)
	for len(tests) < 100 {
		n := rng.Intn(8) + 1
		k := int64(rng.Intn(10) + 1)
		a := make([]int64, n)
		x := make([]int64, n)
		for i := 0; i < n; i++ {
			a[i] = int64(rng.Intn(5) + 1)
			x[i] = int64(rng.Intn(10) - 5) // -5..4
		}
		tests = append(tests, test{n, k, a, x})
	}
	return tests
}

func run(bin, input string) (string, error) {
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
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d %d\n", t.n, t.k))
		for j, v := range t.a {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteString("\n")
		for j, v := range t.x {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteString("\n")
		exp := solve(t.n, t.k, t.a, t.x)
		out, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, sb.String())
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, sb.String())
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
