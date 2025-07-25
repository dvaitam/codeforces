package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type Query struct {
	t int
	l int
	r int
	x int
	y int
}

type Test struct {
	n   int
	q   int
	arr []int
	qs  []Query
}

func (tc Test) Input() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.q))
	for i, v := range tc.arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	for _, qu := range tc.qs {
		if qu.t == 1 {
			sb.WriteString(fmt.Sprintf("1 %d %d %d %d\n", qu.l, qu.r, qu.x, qu.y))
		} else {
			sb.WriteString(fmt.Sprintf("2 %d %d\n", qu.l, qu.r))
		}
	}
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
	ref := "./refF.bin"
	cmd := exec.Command("go", "build", "-o", ref, "794F.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return ref, nil
}

func genTests() []Test {
	rand.Seed(5)
	tests := make([]Test, 0, 100)
	for i := 0; i < 100; i++ {
		n := rand.Intn(5) + 1
		qnum := rand.Intn(5) + 1
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = rand.Intn(100)
		}
		qs := make([]Query, qnum)
		for j := 0; j < qnum; j++ {
			if rand.Intn(2) == 0 {
				l := rand.Intn(n) + 1
				r := rand.Intn(n-l+1) + l
				x := rand.Intn(10)
				y := rand.Intn(10)
				qs[j] = Query{t: 1, l: l, r: r, x: x, y: y}
			} else {
				l := rand.Intn(n) + 1
				r := rand.Intn(n-l+1) + l
				qs[j] = Query{t: 2, l: l, r: r}
			}
		}
		tests = append(tests, Test{n: n, q: qnum, arr: arr, qs: qs})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierF.go /path/to/binary")
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
			fmt.Printf("Test %d failed\nInput:%sExpected:%sGot:%s\n", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
