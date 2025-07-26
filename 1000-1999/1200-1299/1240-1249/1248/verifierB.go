package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
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

func expectedB(arr []int64) string {
	sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
	var x, y int64
	n := len(arr)
	for i := 0; i < n/2; i++ {
		x += arr[i]
	}
	for i := n / 2; i < n; i++ {
		y += arr[i]
	}
	res := x*x + y*y
	return fmt.Sprintf("%d", res)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}
	rand.Seed(time.Now().UnixNano())
	bin := os.Args[1]
	if bin == "--" && len(os.Args) >= 3 {
		bin = os.Args[2]
	}

	for t := 0; t < 100; t++ {
		n := rand.Intn(20) + 1
		arr := make([]int64, n)
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			arr[i] = int64(rand.Intn(1000))
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(fmt.Sprintf("%d", arr[i]))
		}
		input.WriteByte('\n')
		expected := expectedB(arr)
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
