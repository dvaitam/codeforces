package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type Test struct {
	input string
}

func buildRef() (string, error) {
	ref := "refC.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1200C.go")
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
	n := int64(rng.Intn(50) + 1)
	m := int64(rng.Intn(50) + 1)
	q := rng.Intn(5) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, q))
	for i := 0; i < q; i++ {
		sx := rng.Intn(2) + 1
		var sy int64
		if sx == 1 {
			sy = int64(rng.Intn(int(n)) + 1)
		} else {
			sy = int64(rng.Intn(int(m)) + 1)
		}
		ex := rng.Intn(2) + 1
		var ey int64
		if ex == 1 {
			ey = int64(rng.Intn(int(n)) + 1)
		} else {
			ey = int64(rng.Intn(int(m)) + 1)
		}
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", sx, sy, ex, ey))
	}
	return Test{sb.String()}
}

func genTests() []Test {
	rng := rand.New(rand.NewSource(2))
	tests := make([]Test, 0, 105)
	for i := 0; i < 100; i++ {
		tests = append(tests, randCase(rng))
	}
	tests = append(tests, Test{"1 1 1\n1 1 1 1\n"})
	tests = append(tests, Test{"2 3 2\n1 1 2 2\n2 3 2 1\n"})
	tests = append(tests, Test{"5 5 1\n1 5 2 1\n"})
	tests = append(tests, Test{"10 1 1\n2 1 2 1\n"})
	tests = append(tests, Test{"7 9 3\n1 1 1 7\n2 9 1 1\n2 3 2 8\n"})
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
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
