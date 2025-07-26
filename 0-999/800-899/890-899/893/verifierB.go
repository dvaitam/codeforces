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

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func solveCase(n int) int {
	beautiful := []int{}
	for k := 0; ; k++ {
		val := (1<<(k+1) - 1) * (1 << k)
		if val > 100000 {
			break
		}
		beautiful = append(beautiful, val)
	}
	ans := 1
	for i := len(beautiful) - 1; i >= 0; i-- {
		if n%beautiful[i] == 0 {
			ans = beautiful[i]
			break
		}
	}
	return ans
}

func generateCase(r *rand.Rand) (string, string) {
	n := r.Intn(100000) + 1
	input := fmt.Sprintf("%d\n", n)
	expect := fmt.Sprint(solveCase(n))
	return input, expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(r)
		out, err := runBinary(bin, in)
		if err != nil {
			fmt.Printf("Test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Printf("Test %d failed.\nInput:\n%sExpected: %s\nGot: %s\n", i+1, in, exp, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
