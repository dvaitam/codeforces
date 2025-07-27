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
	cmd := exec.Command("go", "build", "-o", ref, "1458C.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return ref, nil
}

func randOps(n int) string {
	letters := []byte{'R', 'L', 'D', 'U', 'I', 'C'}
	var b strings.Builder
	for i := 0; i < n; i++ {
		b.WriteByte(letters[rand.Intn(len(letters))])
	}
	return b.String()
}

func genTests() []Test {
	rand.Seed(2)
	tests := make([]Test, 0, 110)
	for i := 0; i < 100; i++ {
		t := rand.Intn(3) + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", t))
		for c := 0; c < t; c++ {
			n := rand.Intn(3) + 1
			m := rand.Intn(6) + 1
			sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
			for r := 0; r < n; r++ {
				for c2 := 0; c2 < n; c2++ {
					if c2 > 0 {
						sb.WriteByte(' ')
					}
					val := rand.Intn(n) + 1
					sb.WriteString(fmt.Sprintf("%d", val))
				}
				sb.WriteByte('\n')
			}
			sb.WriteString(randOps(m) + "\n")
		}
		tests = append(tests, Test{sb.String()})
	}
	// edge simple
	var sb strings.Builder
	sb.WriteString("1\n1 1\n1\nR\n")
	tests = append(tests, Test{sb.String()})
	return tests
}

func normalize(s string) string {
	return strings.TrimSpace(s)
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
		if normalize(exp) != normalize(got) {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%sGot:\n%s\n", i+1, tc.input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
