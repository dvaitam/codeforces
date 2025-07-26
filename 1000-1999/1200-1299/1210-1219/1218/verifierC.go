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
	ref := "./refC.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1218C.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return ref, nil
}

func genTests() []Test {
	rand.Seed(3)
	tests := make([]Test, 0, 100)
	for i := 0; i < 100; i++ {
		N := rand.Intn(4) + 1
		M := rand.Intn(4) + 1
		K := rand.Intn(5)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d %d\n", N, M, K)
		for j := 0; j < K; j++ {
			x := rand.Intn(N)
			y := rand.Intn(M)
			maxD := N - 1 - x
			if y < maxD {
				if maxD > y {
					maxD = y
				}
			}
			if M-1-y < maxD {
				maxD = M - 1 - y
			}
			if maxD < 0 {
				maxD = 0
			}
			d := rand.Intn(maxD + 1)
			t := rand.Intn(10)
			e := rand.Intn(100)
			fmt.Fprintf(&sb, "%d %d %d %d %d\n", x, y, d, t, e)
		}
		tests = append(tests, Test{sb.String()})
	}
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
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%sGot:\n%s\n", i+1, tc.input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
