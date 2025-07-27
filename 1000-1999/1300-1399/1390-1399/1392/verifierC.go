package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solveCase(arr []int64) int64 {
	var ans int64
	for i := 1; i < len(arr); i++ {
		if arr[i] < arr[i-1] {
			ans += arr[i-1] - arr[i]
		}
	}
	return ans
}

func generateTest() (string, string) {
	n := rand.Intn(10) + 1
	arr := make([]int64, n)
	for i := range arr {
		arr[i] = int64(rand.Intn(21))
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')

	ans := solveCase(arr)
	exp := fmt.Sprintf("%d\n", ans)
	return sb.String(), exp
}

func referenceIO(t int) (string, string) {
	var in strings.Builder
	var out strings.Builder
	fmt.Fprintf(&in, "%d\n", t)
	for i := 0; i < t; i++ {
		ti, to := generateTest()
		in.WriteString(ti)
		out.WriteString(to)
	}
	return in.String(), out.String()
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return buf.String(), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		return
	}
	rand.Seed(3)
	in, exp := referenceIO(100)
	out, err := runBinary(os.Args[1], in)
	if err != nil {
		fmt.Println("Runtime error:", err)
		os.Exit(1)
	}
	if strings.TrimSpace(out) != strings.TrimSpace(exp) {
		fmt.Println("Wrong Answer")
		fmt.Println("Expected:\n" + exp)
		fmt.Println("Got:\n" + out)
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
