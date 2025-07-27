package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solveCase(n int, arr []int) int {
	same := true
	for i := 1; i < n; i++ {
		if arr[i] != arr[0] {
			same = false
			break
		}
	}
	if same {
		return n
	}
	return 1
}

func generateTest() (string, string) {
	n := rand.Intn(10) + 1
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rand.Intn(5)
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

	ans := solveCase(n, arr)
	exp := fmt.Sprintf("%d\n", ans)
	return sb.String(), exp
}

func referenceIO(t int) (string, string) {
	var in strings.Builder
	var out strings.Builder
	fmt.Fprintf(&in, "%d\n", t)
	for i := 0; i < t; i++ {
		tcIn, tcOut := generateTest()
		in.WriteString(tcIn)
		out.WriteString(tcOut)
	}
	return in.String(), out.String()
}

func runBinary(path string, input string) (string, error) {
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
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		return
	}
	rand.Seed(1)
	in, expected := referenceIO(100)
	out, err := runBinary(os.Args[1], in)
	if err != nil {
		fmt.Println("Runtime error:", err)
		os.Exit(1)
	}
	if strings.TrimSpace(out) != strings.TrimSpace(expected) {
		fmt.Println("Wrong Answer")
		fmt.Println("Expected:\n" + expected)
		fmt.Println("Got:\n" + out)
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
