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

type Query struct {
	a, b int
}

type Test struct {
	arr     []int
	queries []Query
}

func (t Test) Input() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(t.arr)))
	for i, v := range t.arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	sb.WriteString(fmt.Sprintf("%d\n", len(t.queries)))
	for _, q := range t.queries {
		sb.WriteString(fmt.Sprintf("%d %d\n", q.a, q.b))
	}
	return sb.String()
}

func buildOracle() (string, error) {
	ref := "oracleF"
	cmd := exec.Command("go", "build", "-o", ref, "1771F.go")
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

func genTests(ref string) []Test {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]Test, 100)
	for i := 0; i < 100; i++ {
		n := rng.Intn(8) + 1
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = rng.Intn(50) + 1
		}
		q := rng.Intn(8) + 1
		queries := make([]Query, q)
		prev := 0
		base := fmt.Sprintf("%d\n", n)
		for j, v := range arr {
			if j > 0 {
				base += " "
			}
			base += fmt.Sprintf("%d", v)
		}
		base += "\n"
		for j := 0; j < q; j++ {
			l := rng.Intn(n) + 1
			r := rng.Intn(n-l+1) + l
			a := l ^ prev
			b := r ^ prev
			queries[j] = Query{a: a, b: b}
			// build partial input to get new answer from oracle
			var sb strings.Builder
			sb.WriteString(base)
			sb.WriteString(fmt.Sprintf("%d\n", j+1))
			for k := 0; k <= j; k++ {
				sb.WriteString(fmt.Sprintf("%d %d\n", queries[k].a, queries[k].b))
			}
			out, err := runExe(ref, sb.String())
			if err != nil {
				// if oracle fails, just keep prev as 0
				prev = 0
			} else {
				lines := strings.Split(strings.TrimSpace(out), "\n")
				val, err := strconv.Atoi(lines[len(lines)-1])
				if err != nil {
					prev = 0
				} else {
					prev = val
				}
			}
		}
		tests[i] = Test{arr: arr, queries: queries}
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	ref, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := genTests(ref)
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
			fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
