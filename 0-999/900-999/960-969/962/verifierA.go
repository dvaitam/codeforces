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

type testCase struct {
	n   int
	arr []int
}

func solveCase(tc testCase) int {
	total := 0
	for _, v := range tc.arr {
		total += v
	}
	need := (total + 1) / 2
	sum := 0
	for i, v := range tc.arr {
		sum += v
		if sum >= need {
			return i + 1
		}
	}
	return tc.n
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	arr := make([]int, n)
	input := fmt.Sprintf("%d\n", n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(10) + 1
		if i > 0 {
			input += " "
		}
		input += fmt.Sprintf("%d", arr[i])
	}
	input += "\n"
	tc := testCase{n, arr}
	out := solveCase(tc)
	expected := fmt.Sprintf("%d\n", out)
	return input, expected
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	if strings.TrimSpace(out.String()) != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %s got %s", expected, out.String())
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
