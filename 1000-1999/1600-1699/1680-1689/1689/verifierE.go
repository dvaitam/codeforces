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

type Test struct {
	arr []int
}

func (t Test) Input() string {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", len(t.arr)))
	for i, v := range t.arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func buildOracle() (string, error) {
	ref := "oracleE"
	cmd := exec.Command("go", "build", "-o", ref, "1689E.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, string(out))
	}
	return ref, nil
}

func runExe(path, input string) (string, error) {
	if !strings.Contains(path, "/") {
		path = "./" + path
	}
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

type dsu struct {
	p []int
	s []int
}

func newDsu(n int) *dsu {
	p := make([]int, n)
	s := make([]int, n)
	for i := 0; i < n; i++ {
		p[i] = i
		s[i] = 1
	}
	return &dsu{p: p, s: s}
}

func (d *dsu) find(x int) int {
	for x != d.p[x] {
		d.p[x] = d.p[d.p[x]]
		x = d.p[x]
	}
	return x
}

func (d *dsu) union(a, b int) {
	ra := d.find(a)
	rb := d.find(b)
	if ra == rb {
		return
	}
	if d.s[ra] < d.s[rb] {
		ra, rb = rb, ra
	}
	d.p[rb] = ra
	d.s[ra] += d.s[rb]
}

func connected(a []int) bool {
	n := len(a)
	start := -1
	for i := 0; i < n; i++ {
		if a[i] != 0 {
			start = i
			break
		}
	}
	if start == -1 {
		return false
	}
	d := newDsu(n)
	for b := 0; b < 30; b++ {
		prev := -1
		for i := 0; i < n; i++ {
			if ((a[i] >> b) & 1) != 0 {
				if prev == -1 {
					prev = i
				} else {
					d.union(prev, i)
				}
			}
		}
	}
	root := d.find(start)
	for i := 0; i < n; i++ {
		if a[i] != 0 && d.find(i) != root {
			return false
		}
		if a[i] == 0 {
			return false
		}
	}
	return true
}

func genTests() []Test {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]Test, 100)
	for i := 0; i < 100; i++ {
		n := rng.Intn(8) + 2 // 2..9
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = rng.Intn(64) // 0..63
		}
		tests[i] = Test{arr: arr}
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	ref, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := genTests()
	for i, tc := range tests {
		input := tc.Input()
		exp, err := runExe(ref, input)
		if err != nil {
			fmt.Printf("oracle runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runExe(bin, input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if exp != got {
			expLines := strings.Split(exp, "\n")
			gotLines := strings.Split(got, "\n")

			if len(gotLines) < 2 {
				fmt.Printf("test %d failed (output too short)\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, input, exp, got)
				os.Exit(1)
			}

			if expLines[0] != gotLines[0] {
				fmt.Printf("test %d failed (wrong ops)\ninput:\n%sexpected ops: %s\ngot ops: %s\n", i+1, input, expLines[0], gotLines[0])
				os.Exit(1)
			}

			var gotArr []int
			for _, s := range strings.Fields(gotLines[1]) {
				var v int
				fmt.Sscanf(s, "%d", &v)
				gotArr = append(gotArr, v)
			}

			ops := 0
			for j, v := range gotArr {
				diff := v - tc.arr[j]
				if diff < 0 {
					diff = -diff
				}
				ops += diff
			}
			if fmt.Sprintf("%d", ops) != expLines[0] {
				fmt.Printf("test %d failed (ops sum mismatch)\ninput:\n%sexpected ops: %s\nactual ops sum: %d\n", i+1, input, expLines[0], ops)
				os.Exit(1)
			}

			if !connected(gotArr) {
				fmt.Printf("test %d failed (not connected)\ninput:\n%sgot array: %v\n", i+1, input, gotArr)
				os.Exit(1)
			}
		}
	}
	fmt.Println("all tests passed")
}
