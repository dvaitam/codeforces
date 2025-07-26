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

type Test struct {
	input string
}

func buildRef() (string, error) {
	ref := "refF.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1200F.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
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

func randCase(rng *rand.Rand) Test {
	n := rng.Intn(5) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(rng.Intn(11) - 5))
	}
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		m := rng.Intn(3) + 1
		sb.WriteString(fmt.Sprintf("%d\n", m))
		for j := 0; j < m; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(rng.Intn(n) + 1))
		}
		sb.WriteByte('\n')
	}
	q := rng.Intn(5) + 1
	sb.WriteString(fmt.Sprintf("%d\n", q))
	for i := 0; i < q; i++ {
		x := rng.Intn(n) + 1
		y := rng.Intn(11) - 5
		sb.WriteString(fmt.Sprintf("%d %d\n", x, y))
	}
	return Test{sb.String()}
}

func genTests() []Test {
	rng := rand.New(rand.NewSource(5))
	tests := make([]Test, 0, 105)
	for i := 0; i < 100; i++ {
		tests = append(tests, randCase(rng))
	}
	tests = append(tests, Test{"1\n0\n1\n1\n1\n0\n"})
	tests = append(tests, Test{"2\n1 1\n1\n1\n1\n2\n1\n1 0\n"})
	tests = append(tests, Test{"3\n0 0 0\n1\n1\n1\n2\n1 2\n3\n1 0\n2 -1\n3 2\n"})
	tests = append(tests, Test{"1\n5\n1\n1\n2\n1\n1 1\n"})
	tests = append(tests, Test{"2\n-1 2\n1\n2\n1 2\n1\n1 0\n"})
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
		exp, err := runExe(ref, tc.input)
		if err != nil {
			fmt.Printf("reference runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runExe(bin, tc.input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if exp != got {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%s\nGot:\n%s\n", i+1, tc.input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
