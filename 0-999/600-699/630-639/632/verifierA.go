package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func runBinary(path string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func expected(n int, p int64, ops []string) int64 {
	var money, cur int64
	for i := n - 1; i >= 0; i-- {
		if ops[i] == "half" {
			money += cur * p
			cur *= 2
		} else {
			money += cur*p + p/2
			cur = cur*2 + 1
		}
	}
	return money
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for t := 0; t < 100; t++ {
		n := rand.Intn(40) + 1
		p := int64(2 * (rand.Intn(500) + 1))
		ops := make([]string, n)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, p))
		for i := 0; i < n; i++ {
			if rand.Intn(2) == 0 {
				ops[i] = "half"
			} else {
				ops[i] = "halfplus"
			}
			sb.WriteString(ops[i])
			if i+1 < n {
				sb.WriteByte('\n')
			}
		}
		input := sb.String()
		exp := expected(n, p, ops)
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\nInput:\n%s\nOutput:\n%s\n", t+1, err, input, out)
			os.Exit(1)
		}
		val, err := strconv.ParseInt(strings.Fields(out)[0], 10, 64)
		if err != nil || val != exp {
			fmt.Printf("Test %d failed\nInput:\n%s\nExpected: %d\nGot: %s\n", t+1, input, exp, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
