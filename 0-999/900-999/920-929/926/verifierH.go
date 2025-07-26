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
	k      int
	beauty []int
	colors string
}

func (t Test) Input() string {
	var sb strings.Builder
	n := len(t.beauty)
	sb.WriteString(fmt.Sprintf("%d %d\n", n, t.k))
	for i, v := range t.beauty {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprint(&sb, v)
	}
	sb.WriteString("\n")
	sb.WriteString(t.colors)
	sb.WriteString("\n")
	return sb.String()
}

func buildRef() (string, error) {
	ref := "./refH.bin"
	cmd := exec.Command("go", "build", "-o", ref, "926H.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build ref failed: %v: %s", err, out)
	}
	return ref, nil
}

func runExe(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func genTests() []Test {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]Test, 0, 101)
	colors := []byte{'R', 'O', 'W'}
	for i := 0; i < 100; i++ {
		n := rng.Intn(10) + 1
		k := rng.Intn(n) + 1
		beauty := make([]int, n)
		var sb strings.Builder
		for j := 0; j < n; j++ {
			beauty[j] = rng.Intn(10) + 1
			sb.WriteByte(colors[rng.Intn(3)])
		}
		tests = append(tests, Test{k: k, beauty: beauty, colors: sb.String()})
	}
	tests = append(tests, Test{k: 1, beauty: []int{5}, colors: "R"})
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
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
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runExe(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if exp != got {
			fmt.Printf("Test %d failed\nInput:%sExpected:%s\nGot:%s\n", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
