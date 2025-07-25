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

func runBinary(path, input string) (string, error) {
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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func isMagic(n int, a [][]int) bool {
	for i := 0; i < n; i++ {
		if a[i][i] != 0 {
			return false
		}
		for j := 0; j < n; j++ {
			if a[i][j] != a[j][i] {
				return false
			}
		}
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			for k := 0; k < n; k++ {
				if a[i][j] > max(a[i][k], a[k][j]) {
					return false
				}
			}
		}
	}
	return true
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for t := 0; t < 100; t++ {
		n := rand.Intn(5) + 1
		a := make([][]int, n)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			a[i] = make([]int, n)
			for j := 0; j < n; j++ {
				if i == j {
					a[i][j] = 0
				} else {
					a[i][j] = rand.Intn(10)
				}
				sb.WriteString(fmt.Sprintf("%d ", a[i][j]))
			}
			sb.WriteByte('\n')
		}
		input := sb.String()
		exp := "NOT MAGIC"
		if isMagic(n, a) {
			exp = "MAGIC"
		}
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("Test %d runtime error: %v\nInput:\n%s\nOutput:\n%s\n", t+1, err, input, out)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Printf("Test %d failed\nInput:\n%s\nExpected: %s\nGot: %s\n", t+1, input, exp, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
