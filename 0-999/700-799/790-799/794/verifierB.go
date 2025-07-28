package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type Test struct {
	n int
	h float64
}

func (t Test) Input() string {
	return fmt.Sprintf("%d %.2f\n", t.n, t.h)
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
	ref := "./refB.bin"
	cmd := exec.Command("go", "build", "-o", ref, "794B.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return ref, nil
}

func genTests() []Test {
	rand.Seed(1)
	tests := make([]Test, 0, 100)
	for i := 0; i < 100; i++ {
		n := rand.Intn(4) + 2 // 2..5
		h := rand.Float64()*20 + 1
		tests = append(tests, Test{n: n, h: h})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
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
		expFields := strings.Fields(exp)
		gotFields := strings.Fields(got)
		if len(expFields) != len(gotFields) {
			fmt.Printf("Test %d failed\nInput:%sExpected %d numbers got %d\n%s", i+1, input, len(expFields), len(gotFields), got)
			os.Exit(1)
		}
		for j := range expFields {
			var ef, gf float64
			if _, err := fmt.Sscan(expFields[j], &ef); err != nil {
				fmt.Printf("bad reference output on test %d: %v\n", i+1, err)
				os.Exit(1)
			}
			if _, err := fmt.Sscan(gotFields[j], &gf); err != nil {
				fmt.Printf("bad candidate output on test %d: %v\n", i+1, err)
				os.Exit(1)
			}
			diff := math.Abs(ef - gf)
			tol := 1e-6 * math.Max(1, math.Abs(ef))
			if diff > tol {
				fmt.Printf("Test %d failed\nInput:%sExpected:%sGot:%s\n", i+1, input, exp, got)
				os.Exit(1)
			}
		}
	}
	fmt.Println("all tests passed")
}
