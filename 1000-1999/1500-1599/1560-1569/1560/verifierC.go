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

type testCase struct {
	k int64
}

func position(k int64) (int64, int64) {
	root := int64(math.Sqrt(float64(k)))
	if root*root < k {
		root++
	}
	prev := (root - 1) * (root - 1)
	diff := k - prev
	if diff <= root {
		return diff, root
	}
	return root, root*2 - diff
}

func genTests() []testCase {
	r := rand.New(rand.NewSource(3))
	tests := make([]testCase, 100)
	for i := range tests {
		tests[i].k = int64(r.Intn(1000000) + 1)
	}
	return tests
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return out.String() + errBuf.String(), fmt.Errorf("runtime error: %v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	tests := genTests()
	for i, tc := range tests {
		input := fmt.Sprintf("1\n%d\n", tc.k)
		r, c := position(tc.k)
		want := fmt.Sprintf("%d %d", r, c)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			return
		}
		if strings.TrimSpace(got) != want {
			fmt.Printf("case %d failed: expected %s got %s\n", i+1, want, got)
			fmt.Printf("input:\n%s", input)
			return
		}
	}
	fmt.Println("All tests passed")
}
