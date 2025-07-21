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

func expected(nums []int) string {
	cntEven := 0
	for _, v := range nums {
		if v%2 == 0 {
			cntEven++
		}
	}
	wantEven := cntEven == 1
	for i, v := range nums {
		if (v%2 == 0) == wantEven {
			return fmt.Sprintf("%d", i+1)
		}
	}
	return "1"
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	
	for i := 0; i < 100; i++ {
		n := rng.Intn(98) + 3
		nums := make([]int, n)
		
		var evenCount, oddCount int
		for j := 0; j < n; j++ {
			nums[j] = rng.Intn(100) + 1
			if nums[j]%2 == 0 {
				evenCount++
			} else {
				oddCount++
			}
		}
		
		if evenCount == 1 || oddCount == 1 {
		} else {
			diffPos := rng.Intn(n)
			if evenCount > 1 {
				nums[diffPos] = 2*rng.Intn(50) + 1
			} else {
				nums[diffPos] = 2 * (rng.Intn(50) + 1)
			}
		}
		
		input := fmt.Sprintf("%d\n", n)
		for j, v := range nums {
			if j > 0 {
				input += " "
			}
			input += fmt.Sprintf("%d", v)
		}
		input += "\n"
		
		expectedOut := expected(nums)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if got != expectedOut {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expectedOut, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}