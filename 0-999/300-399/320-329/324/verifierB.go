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

type queryB struct {
	typ int
	x   int
	y   int
}

type testCaseB struct {
	input    string
	expected []int
}

func solveB(n int, a []int, queries []queryB) []int {
	pos := make([]int, n+1)
	for i := 1; i <= n; i++ {
		pos[a[i]] = i
	}
	res := []int{}
	for _, q := range queries {
		if q.typ == 1 {
			sessions := 1
			for id := q.x; id < q.y; id++ {
				if pos[id+1] < pos[id] {
					sessions++
				}
			}
			res = append(res, sessions)
		} else {
			v1 := a[q.x]
			v2 := a[q.y]
			a[q.x], a[q.y] = v2, v1
			pos[v1], pos[v2] = q.y, q.x
		}
	}
	return res
}

func genCaseB(rng *rand.Rand) testCaseB {
	n := rng.Intn(5) + 1
	perm := rng.Perm(n)
	a := make([]int, n+1)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 1; i <= n; i++ {
		a[i] = perm[i-1] + 1
		if i > 1 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", a[i])
	}
	sb.WriteByte('\n')
	qcnt := rng.Intn(10) + 1
	fmt.Fprintf(&sb, "%d\n", qcnt)
	queries := make([]queryB, qcnt)
	for i := 0; i < qcnt; i++ {
		typ := rng.Intn(2) + 1
		x := rng.Intn(n) + 1
		y := rng.Intn(n) + 1
		if typ == 1 && x > y {
			x, y = y, x
		}
		queries[i] = queryB{typ: typ, x: x, y: y}
		fmt.Fprintf(&sb, "%d %d %d\n", typ, x, y)
	}
	expected := solveB(n, append([]int(nil), a...), queries)
	return testCaseB{input: sb.String(), expected: expected}
}

func runCaseB(bin string, tc testCaseB) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(strings.TrimSpace(out.String()))
	if len(fields) != len(tc.expected) {
		return fmt.Errorf("expected %d numbers got %d", len(tc.expected), len(fields))
	}
	for i, f := range fields {
		var val int
		if _, err := fmt.Sscan(f, &val); err != nil {
			return fmt.Errorf("bad output: %v", err)
		}
		if val != tc.expected[i] {
			return fmt.Errorf("expected %v got %v", tc.expected, fields)
		}
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
		tc := genCaseB(rng)
		if err := runCaseB(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
