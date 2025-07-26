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

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return out.String() + stderr.String(), err
	}
	return out.String(), nil
}

type Constraint struct {
	x  int
	op string
	y  int
}

type Test struct {
	n           int
	constraints []Constraint
	input       string
}

func genTest(rng *rand.Rand) Test {
	n := rng.Intn(4) + 1
	k := rng.Intn(6)
	cons := make([]Constraint, k)
	ops := []string{"=", "<", ">", "<=", ">="}
	for i := 0; i < k; i++ {
		x := rng.Intn(n) + 1
		y := rng.Intn(n) + 1
		op := ops[rng.Intn(len(ops))]
		cons[i] = Constraint{x: x, op: op, y: y}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for _, c := range cons {
		sb.WriteString(fmt.Sprintf("%d %s %d\n", c.x, c.op, c.y))
	}
	return Test{n: n, constraints: cons, input: sb.String()}
}

func check(seq []int, cons []Constraint) bool {
	for _, c := range cons {
		vx := seq[c.x-1]
		vy := seq[c.y-1]
		switch c.op {
		case "=":
			if vx != vy {
				return false
			}
		case "<":
			if !(vx < vy) {
				return false
			}
		case ">":
			if !(vx > vy) {
				return false
			}
		case "<=":
			if !(vx <= vy) {
				return false
			}
		case ">=":
			if !(vx >= vy) {
				return false
			}
		}
	}
	return true
}

func dfs(i int, n int, a []int, cons []Constraint, count *int64) {
	if i > n {
		left := make([]int, 0, 2*n)
		for v := 1; v <= n; v++ {
			for j := 0; j < a[v]; j++ {
				left = append(left, v)
			}
		}
		right := make([]int, 0, 2*n)
		for v := n; v >= 1; v-- {
			for j := 0; j < 2-a[v]; j++ {
				right = append(right, v)
			}
		}
		seq := append(left, right...)
		if check(seq, cons) {
			*count++
		}
		return
	}
	for t := 0; t <= 2; t++ {
		a[i] = t
		dfs(i+1, n, a, cons, count)
	}
}

func solve(t Test) string {
	a := make([]int, t.n+1)
	var cnt int64
	dfs(1, t.n, a, t.constraints, &cnt)
	return strconv.FormatInt(cnt, 10)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		t := genTest(rng)
		expected := solve(t)
		out, err := run(bin, t.input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\noutput:\n%s", i+1, err, out)
			os.Exit(1)
		}
		out = strings.TrimSpace(out)
		if out != expected {
			fmt.Printf("test %d failed\ninput:\n%s\nexpected:%s got:%s\n", i+1, t.input, expected, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
