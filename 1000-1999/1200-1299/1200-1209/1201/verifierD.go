package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Point struct{ r, c int }

type Test struct {
	n, m      int
	treasures []Point
	safe      []int
}

func (t Test) Input() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d %d\n", t.n, t.m, len(t.treasures), len(t.safe)))
	for _, p := range t.treasures {
		sb.WriteString(fmt.Sprintf("%d %d\n", p.r, p.c))
	}
	for i, v := range t.safe {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runExe(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func buildRef() (string, error) {
	ref := "./refD.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1201D.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return ref, nil
}

func genTests() []Test {
	rand.Seed(3)
	tests := make([]Test, 0, 101)
	for i := 0; i < 100; i++ {
		n := rand.Intn(4) + 2
		m := rand.Intn(4) + 2
		limit := n * m
		k := rand.Intn(min(4, limit)) + 1
		cells := make(map[Point]struct{})
		treasures := make([]Point, 0, k)
		for len(treasures) < k {
			p := Point{rand.Intn(n) + 1, rand.Intn(m) + 1}
			if _, ok := cells[p]; ok {
				continue
			}
			cells[p] = struct{}{}
			treasures = append(treasures, p)
		}
		q := rand.Intn(m) + 1
		cols := make([]int, 0, q)
		colset := map[int]struct{}{}
		for len(cols) < q {
			c := rand.Intn(m) + 1
			if _, ok := colset[c]; ok {
				continue
			}
			colset[c] = struct{}{}
			cols = append(cols, c)
		}
		tests = append(tests, Test{n, m, treasures, cols})
	}
	tests = append(tests, Test{2, 2, []Point{{1, 1}}, []int{1}})
	return tests
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	ref, err := buildRef()
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
			fmt.Printf("reference runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runExe(bin, input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(exp) != strings.TrimSpace(got) {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%sGot:\n%s\n", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
