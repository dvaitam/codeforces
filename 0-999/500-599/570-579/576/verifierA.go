package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

func buildRef() (string, error) {
	ref := "./refA_bin"
	cmd := exec.Command("go", "build", "-o", ref, "576A.go")
	cmd.Stdout = new(bytes.Buffer)
	cmd.Stderr = new(bytes.Buffer)
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return ref, nil
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = new(bytes.Buffer)
	err := cmd.Run()
	return out.String(), err
}

func genTests() []string {
	rand.Seed(1)
	tests := make([]string, 100)
	for i := 0; i < 100; i++ {
		n := rand.Intn(1000) + 1
		tests[i] = fmt.Sprintf("%d\n", n)
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	cand := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Println("failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := genTests()
	for idx, input := range tests {
		exp, err1 := runBinary(ref, input)
		out, err2 := runBinary(cand, input)
		if err1 != nil || err2 != nil {
			fmt.Printf("runtime error on test %d\n", idx+1)
			os.Exit(1)
		}
		expToks := strings.Fields(exp)
		outToks := strings.Fields(out)
		if len(expToks) != len(outToks) {
			fmt.Printf("wrong answer on test %d\ninput:\n%sexpected:\n%s\ngot:\n%s", idx+1, input, exp, out)
			os.Exit(1)
		}

		// Compare outputs ignoring the order of numbers. The first token
		// is the count of numbers in the sequence; the remaining tokens
		// are the actual question values, which may appear in any order.
		if len(expToks) == 0 {
			fmt.Printf("wrong answer on test %d: empty output\n", idx+1)
			os.Exit(1)
		}
		if expToks[0] != outToks[0] {
			fmt.Printf("wrong answer on test %d\ninput:\n%sexpected:\n%s\ngot:\n%s", idx+1, input, exp, out)
			os.Exit(1)
		}
		expNums := expToks[1:]
		outNums := outToks[1:]
		if len(expNums) != len(outNums) {
			fmt.Printf("wrong answer on test %d\ninput:\n%sexpected:\n%s\ngot:\n%s", idx+1, input, exp, out)
			os.Exit(1)
		}
		sort.Strings(expNums)
		sort.Strings(outNums)
		for i := range expNums {
			if expNums[i] != outNums[i] {
				fmt.Printf("wrong answer on test %d\ninput:\n%sexpected:\n%s\ngot:\n%s", idx+1, input, exp, out)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
