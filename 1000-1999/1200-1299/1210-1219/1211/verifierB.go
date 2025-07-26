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
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func expectedB(a []int64) int64 {
	n := int64(len(a))
	var ans int64
	for i, v := range a {
		if v == 0 {
			continue
		}
		candidate := int64(i+1) + n*(v-1)
		if candidate > ans {
			ans = candidate
		}
	}
	return ans
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}
	rand.Seed(time.Now().UnixNano())
	bin := os.Args[1]

	for t := 0; t < 100; t++ {
		n := rand.Intn(20) + 1
		aVals := make([]int64, n)
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d\n", n))
		allZero := true
		for i := 0; i < n; i++ {
			v := int64(rand.Intn(10))
			aVals[i] = v
			if v > 0 {
				allZero = false
			}
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(fmt.Sprintf("%d", v))
		}
		if allZero {
			aVals[0] = 1
			input.Reset()
			input.WriteString(fmt.Sprintf("%d\n", n))
			for i := 0; i < n; i++ {
				if i > 0 {
					input.WriteByte(' ')
				}
				input.WriteString(fmt.Sprintf("%d", aVals[i]))
			}
		}
		input.WriteByte('\n')
		inBytes := []byte(input.String())
		expected := fmt.Sprintf("%d", expectedB(aVals))
		out, err := runProgram(bin, inBytes)
		if err != nil || strings.TrimSpace(out) != expected {
			fmt.Println("Test", t+1, "failed")
			fmt.Println("Input:\n", input.String())
			fmt.Println("Expected:", expected)
			fmt.Println("Output:", out)
			if err != nil {
				fmt.Println("Error:", err)
			}
			return
		}
	}
	fmt.Println("All tests passed")
}
