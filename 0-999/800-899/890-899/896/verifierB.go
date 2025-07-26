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

const numTests = 100

func isSorted(arr []int) bool {
	if len(arr) == 0 {
		return true
	}
	if arr[0] == 0 {
		return false
	}
	for i := 1; i < len(arr); i++ {
		if arr[i] == 0 || arr[i-1] > arr[i] {
			return false
		}
	}
	return true
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]
	rand.Seed(time.Now().UnixNano())

	for t := 1; t <= numTests; t++ {
		n := rand.Intn(5) + 1
		m := rand.Intn(10) + n
		c := rand.Intn(5) + 1
		xs := make([]int, m)
		var input bytes.Buffer
		input.WriteString(fmt.Sprintf("%d %d %d\n", n, m, c))
		for i := 0; i < m; i++ {
			xs[i] = rand.Intn(c) + 1
			input.WriteString(fmt.Sprintf("%d\n", xs[i]))
		}
		inBytes := input.Bytes()

		cmd := exec.Command(target)
		cmd.Stdin = bytes.NewReader(inBytes)
		outBytes, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", t, err)
			os.Exit(1)
		}
		outputs := strings.Fields(string(outBytes))
		arr := make([]int, n)
		idx := 0
		success := false
		for i := 0; i < m; i++ {
			if idx >= len(outputs) {
				fmt.Printf("Test %d: insufficient outputs\n", t)
				os.Exit(1)
			}
			pos, err := strconv.Atoi(outputs[idx])
			idx++
			if err != nil || pos < 1 || pos > n {
				fmt.Printf("Test %d: invalid index %q\n", t, outputs[idx-1])
				os.Exit(1)
			}
			arr[pos-1] = xs[i]
			if isSorted(arr) {
				success = true
				if idx != len(outputs) {
					fmt.Printf("Test %d: extra output after sorted\n", t)
					os.Exit(1)
				}
				break
			}
		}
		if !success && idx != len(outputs) {
			fmt.Printf("Test %d: extra output lines\n", t)
			os.Exit(1)
		}
	}
	fmt.Printf("Passed %d tests\n", numTests)
}
