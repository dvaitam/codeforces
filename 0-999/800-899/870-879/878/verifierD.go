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

type Op struct {
	t int
	x int
	y int
}

type Test struct {
	n     int
	k     int
	q     int
	a     [][]int
	ops   []Op
	input string
}

func genTest(rng *rand.Rand) Test {
	n := rng.Intn(3) + 1
	k := rng.Intn(3) + 1
	q := rng.Intn(6) + 1
	a := make([][]int, k)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, k, q))
	for i := 0; i < k; i++ {
		a[i] = make([]int, n)
		for j := 0; j < n; j++ {
			a[i][j] = rng.Intn(10)
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(a[i][j]))
		}
		sb.WriteByte('\n')
	}
	ops := make([]Op, q)
	cur := k
	for i := 0; i < q; i++ {
		ttype := rng.Intn(3) + 1
		var x, y int
		if ttype == 1 || ttype == 2 {
			x = rng.Intn(cur) + 1
			y = rng.Intn(cur) + 1
			ops[i] = Op{t: ttype, x: x, y: y}
			sb.WriteString(fmt.Sprintf("%d %d %d\n", ttype, x, y))
			cur++
		} else {
			x = rng.Intn(cur) + 1
			y = rng.Intn(n) + 1
			ops[i] = Op{t: 3, x: x, y: y}
			sb.WriteString(fmt.Sprintf("3 %d %d\n", x, y))
		}
	}
	return Test{n: n, k: k, q: q, a: a, ops: ops, input: sb.String()}
}

func solve(t Test) string {
	creatures := make([][]int, t.k+t.q+1)
	for i := 1; i <= t.k; i++ {
		creatures[i] = append([]int(nil), t.a[i-1]...)
	}
	cur := t.k + 1
	var sb strings.Builder
	for _, op := range t.ops {
		switch op.t {
		case 1:
			vals := make([]int, t.n)
			for j := 0; j < t.n; j++ {
				x := creatures[op.x][j]
				y := creatures[op.y][j]
				if x >= y {
					vals[j] = x
				} else {
					vals[j] = y
				}
			}
			creatures[cur] = vals
			cur++
		case 2:
			vals := make([]int, t.n)
			for j := 0; j < t.n; j++ {
				x := creatures[op.x][j]
				y := creatures[op.y][j]
				if x <= y {
					vals[j] = x
				} else {
					vals[j] = y
				}
			}
			creatures[cur] = vals
			cur++
		case 3:
			sb.WriteString(fmt.Sprintf("%d\n", creatures[op.x][op.y-1]))
		}
	}
	return strings.TrimSpace(sb.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
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
			fmt.Printf("test %d failed\ninput:\n%s\nexpected:\n%s\nactual:\n%s\n", i+1, t.input, expected, out)
			os.Exit(1)
		}
	}
	fmt.Println("ok 100 tests")
}
