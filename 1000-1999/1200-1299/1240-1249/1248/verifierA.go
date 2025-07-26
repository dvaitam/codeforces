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

func runProgram(bin string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func expectedA(p, q []int) string {
	evenP, oddP := 0, 0
	for _, v := range p {
		if v%2 == 0 {
			evenP++
		} else {
			oddP++
		}
	}
	evenQ, oddQ := 0, 0
	for _, v := range q {
		if v%2 == 0 {
			evenQ++
		} else {
			oddQ++
		}
	}
	ans := int64(evenP*evenQ + oddP*oddQ)
	return fmt.Sprintf("%d", ans)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		return
	}
	rand.Seed(time.Now().UnixNano())
	bin := os.Args[1]
	if bin == "--" && len(os.Args) >= 3 {
		bin = os.Args[2]
	}

	for t := 0; t < 100; t++ {
		tc := rand.Intn(5) + 1
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d\n", tc))
		expectedOut := make([]string, tc)
		for i := 0; i < tc; i++ {
			n := rand.Intn(10) + 1
			p := make([]int, n)
			input.WriteString(fmt.Sprintf("%d\n", n))
			for j := 0; j < n; j++ {
				p[j] = rand.Intn(1000)
				if j > 0 {
					input.WriteByte(' ')
				}
				input.WriteString(fmt.Sprintf("%d", p[j]))
			}
			input.WriteByte('\n')
			m := rand.Intn(10) + 1
			q := make([]int, m)
			input.WriteString(fmt.Sprintf("%d\n", m))
			for j := 0; j < m; j++ {
				q[j] = rand.Intn(1000)
				if j > 0 {
					input.WriteByte(' ')
				}
				input.WriteString(fmt.Sprintf("%d", q[j]))
			}
			input.WriteByte('\n')
			expectedOut[i] = expectedA(p, q)
		}
		expected := strings.Join(expectedOut, "\n")
		out, err := runProgram(bin, []byte(input.String()))
		if err != nil || strings.TrimSpace(out) != expected {
			fmt.Printf("Test %d failed\n", t+1)
			fmt.Println("Input:\n", input.String())
			fmt.Println("Expected:\n", expected)
			fmt.Println("Output:\n", out)
			if err != nil {
				fmt.Println("Error:", err)
			}
			return
		}
	}
	fmt.Println("All tests passed")
}
