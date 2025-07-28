package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type caseA struct {
	n, m int
	arr  []int
}

func solveA(n, m int, arr []int) string {
	res := make([]byte, m)
	for i := range res {
		res[i] = 'B'
	}
	for _, v := range arr {
		left := v
		right := m + 1 - v
		if left > right {
			left, right = right, left
		}
		if res[left-1] == 'B' {
			res[left-1] = 'A'
		} else {
			res[right-1] = 'A'
		}
	}
	return string(res)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	const T = 100
	tests := make([]caseA, T)
	expected := make([]string, T)
	var input bytes.Buffer
	fmt.Fprintln(&input, T)
	for i := 0; i < T; i++ {
		n := rng.Intn(50) + 1
		m := rng.Intn(50) + 1
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = rng.Intn(m) + 1
		}
		tests[i] = caseA{n: n, m: m, arr: arr}
		fmt.Fprintf(&input, "%d %d\n", n, m)
		for j, v := range arr {
			if j+1 == len(arr) {
				fmt.Fprintf(&input, "%d\n", v)
			} else {
				fmt.Fprintf(&input, "%d ", v)
			}
		}
		expected[i] = solveA(n, m, arr)
	}
	cmd := exec.Command(binary)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "error running binary:", err)
		os.Exit(1)
	}
	scanner := bufio.NewScanner(bytes.NewReader(out.Bytes()))
	for i := 0; i < T; i++ {
		if !scanner.Scan() {
			fmt.Fprintln(os.Stderr, "insufficient output")
			os.Exit(1)
		}
		got := strings.TrimSpace(scanner.Text())
		if got != expected[i] {
			fmt.Printf("mismatch on test %d: expected %s got %s\n", i+1, expected[i], got)
			os.Exit(1)
		}
	}
	if scanner.Scan() {
		fmt.Fprintln(os.Stderr, "extra output after", T, "tests")
		os.Exit(1)
	}
	fmt.Println("all tests passed")
}
