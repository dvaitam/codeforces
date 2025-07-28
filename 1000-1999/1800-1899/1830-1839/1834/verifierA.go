package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCaseA struct {
	arr []int
}

func expectedA(arr []int) string {
	n := len(arr)
	neg := 0
	for _, x := range arr {
		if x == -1 {
			neg++
		}
	}
	sum := n - 2*neg
	ops := 0
	if neg%2 == 1 {
		ops++
		neg--
		sum += 2
	}
	if sum < 0 {
		deficit := -sum
		pairs := (deficit + 3) / 4
		ops += pairs * 2
	}
	return fmt.Sprint(ops)
}

func genTestsA() []testCaseA {
	rand.Seed(1)
	tests := make([]testCaseA, 0, 100)
	for len(tests) < 100 {
		n := rand.Intn(10) + 1
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			if rand.Intn(2) == 0 {
				arr[i] = -1
			} else {
				arr[i] = 1
			}
		}
		tests = append(tests, testCaseA{arr: arr})
	}
	return tests
}

func runCase(bin string, tc testCaseA) error {
	var input strings.Builder
	input.WriteString("1\n")
	input.WriteString(fmt.Sprintf("%d\n", len(tc.arr)))
	for i, v := range tc.arr {
		if i > 0 {
			input.WriteByte(' ')
		}
		input.WriteString(fmt.Sprint(v))
	}
	input.WriteByte('\n')

	cmd := exec.Command(bin)
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	}
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	expect := expectedA(tc.arr)
	if got != expect {
		return fmt.Errorf("expected %s got %s for %v", expect, got, tc.arr)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsA()
	for i, tc := range tests {
		if err := runCase(bin, tc); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
