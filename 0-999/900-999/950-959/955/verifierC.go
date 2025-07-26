package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type Query struct{ L, R int64 }

type Test struct {
	q  int
	qs []Query
}

func (t Test) Input() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t.q))
	for _, qq := range t.qs {
		sb.WriteString(fmt.Sprintf("%d %d\n", qq.L, qq.R))
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
	ref := "./refC.bin"
	cmd := exec.Command("go", "build", "-o", ref, "955C.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	return ref, nil
}

func genTests() []Test {
	rng := rand.New(rand.NewSource(2))
	tests := make([]Test, 0, 105)
	for i := 0; i < 100; i++ {
		q := rng.Intn(5) + 1
		qs := make([]Query, q)
		for j := 0; j < q; j++ {
			l := rng.Int63n(1_000_000) + 1
			r := l + rng.Int63n(1_000_000)
			qs[j] = Query{l, r}
		}
		tests = append(tests, Test{q: q, qs: qs})
	}
	tests = append(tests,
		Test{q: 1, qs: []Query{{1, 1}}},
		Test{q: 2, qs: []Query{{1, 10}, {100, 1000}}},
		Test{q: 1, qs: []Query{{999999, 1000000}}},
		Test{q: 3, qs: []Query{{5, 5}, {10, 20}, {30, 40}}},
		Test{q: 1, qs: []Query{{500, 1500}}},
	)
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		return
	}
	candidate := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := genTests()
	for i, tc := range tests {
		input := tc.Input()
		want, err := runExe(ref, input)
		if err != nil {
			fmt.Printf("reference failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runExe(candidate, input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			fmt.Println("input:\n" + input)
			os.Exit(1)
		}
		if strings.TrimSpace(want) != strings.TrimSpace(got) {
			fmt.Printf("test %d failed\ninput:%sexpected:%sgot:%s\n", i+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
