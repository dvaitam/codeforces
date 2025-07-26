package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type Test struct{ input string }

func runExe(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func buildRef() (string, error) {
	ref := "./refE.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1218E.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return ref, nil
}

func genTests() []Test {
	rand.Seed(5)
	tests := make([]Test, 0, 100)
	for i := 0; i < 100; i++ {
		n := rand.Intn(3) + 1
		K := rand.Intn(4) + 1
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, K)
		for j := 0; j < n; j++ {
			fmt.Fprintf(&sb, "%d ", rand.Intn(10))
		}
		sb.WriteByte('\n')
		q := rand.Intn(3) + 1
		fmt.Fprintf(&sb, "%d\n", q)
		for j := 0; j < q; j++ {
			typ := rand.Intn(2) + 1
			if typ == 1 {
				qv := rand.Intn(10)
				idx := rand.Intn(n) + 1
				d := rand.Intn(10)
				fmt.Fprintf(&sb, "1 %d %d %d\n", qv, idx, d)
			} else {
				qv := rand.Intn(10)
				L := rand.Intn(n) + 1
				R := L + rand.Intn(n-L+1)
				d := rand.Intn(10)
				fmt.Fprintf(&sb, "2 %d %d %d %d\n", qv, L, R, d)
			}
		}
		tests = append(tests, Test{sb.String()})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
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
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%sGot:\n%s\n", i+1, tc.input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
