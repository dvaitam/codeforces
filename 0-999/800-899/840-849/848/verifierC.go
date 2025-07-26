package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type op struct {
	typ int
	p   int
	x   int
	l   int
	r   int
}

func runSolution(bin, input string) (string, error) {
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

func parseInputC(in string) (int, int, []int, []op, error) {
	r := bufio.NewReader(strings.NewReader(in))
	var n, m int
	if _, err := fmt.Fscan(r, &n, &m); err != nil {
		return 0, 0, nil, nil, err
	}
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(r, &arr[i])
	}
	ops := make([]op, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(r, &ops[i].typ)
		if ops[i].typ == 1 {
			fmt.Fscan(r, &ops[i].p, &ops[i].x)
		} else {
			fmt.Fscan(r, &ops[i].l, &ops[i].r)
		}
	}
	return n, m, arr, ops, nil
}

func memorySegment(arr []int, l, r int) int64 {
	first := make(map[int]int)
	last := make(map[int]int)
	for i := l; i <= r; i++ {
		v := arr[i]
		if _, ok := first[v]; !ok {
			first[v] = i
		}
		last[v] = i
	}
	var sum int64
	for k := range first {
		sum += int64(last[k] - first[k])
	}
	return sum
}

func solveC(n int, arr []int, ops []op) []int64 {
	a := append([]int(nil), arr...)
	var res []int64
	for _, o := range ops {
		if o.typ == 1 {
			a[o.p-1] = o.x
		} else {
			res = append(res, memorySegment(a, o.l-1, o.r-1))
		}
	}
	return res
}

func verifyC(input, output string) error {
	n, _, arr, ops, err := parseInputC(input)
	if err != nil {
		return fmt.Errorf("input parse: %v", err)
	}
	expected := solveC(n, arr, ops)
	r := bufio.NewReader(strings.NewReader(output))
	for i := 0; i < len(expected); i++ {
		var val int64
		if _, err := fmt.Fscan(r, &val); err != nil {
			return fmt.Errorf("parse output line %d: %v", i+1, err)
		}
		if val != expected[i] {
			return fmt.Errorf("query %d expected %d got %d", i+1, expected[i], val)
		}
	}
	return nil
}

func generateCaseC(rng *rand.Rand) string {
	n := rng.Intn(4) + 1
	m := rng.Intn(6) + 1
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(n) + 1
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", arr[i])
	}
	sb.WriteByte('\n')
	queries := 0
	for i := 0; i < m; i++ {
		if i == m-1 && queries == 0 {
			l := rng.Intn(n) + 1
			r := rng.Intn(n-l+1) + l
			fmt.Fprintf(&sb, "2 %d %d\n", l, r)
			queries++
			continue
		}
		if rng.Intn(2) == 0 {
			p := rng.Intn(n) + 1
			x := rng.Intn(n) + 1
			fmt.Fprintf(&sb, "1 %d %d\n", p, x)
		} else {
			l := rng.Intn(n) + 1
			r := rng.Intn(n-l+1) + l
			fmt.Fprintf(&sb, "2 %d %d\n", l, r)
			queries++
		}
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []string{}
	for len(cases) < 100 {
		cases = append(cases, generateCaseC(rng))
	}
	for i, tc := range cases {
		out, err := runSolution(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
		if err := verifyC(tc, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%soutput:\n%s", i+1, err, tc, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
