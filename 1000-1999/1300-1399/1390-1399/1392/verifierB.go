package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solveCase(n int, k int64, arr []int64) []int64 {
	res := make([]int64, n)
	copy(res, arr)
	if k%2 == 1 {
		var mx int64 = res[0]
		for _, v := range res {
			if v > mx {
				mx = v
			}
		}
		for i, v := range res {
			res[i] = mx - v
		}
	} else {
		var mn int64 = res[0]
		for _, v := range res {
			if v < mn {
				mn = v
			}
		}
		for i, v := range res {
			res[i] = v - mn
		}
	}
	return res
}

func generateTest() (string, string) {
	n := rand.Intn(8) + 1
	k := rand.Int63n(20)
	arr := make([]int64, n)
	for i := range arr {
		arr[i] = int64(rand.Intn(21) - 10)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, k)
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')

	res := solveCase(n, k, arr)
	var out strings.Builder
	for i, v := range res {
		if i > 0 {
			out.WriteByte(' ')
		}
		fmt.Fprintf(&out, "%d", v)
	}
	out.WriteByte('\n')
	return sb.String(), out.String()
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
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}
	rand.Seed(2)
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
