package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

type caseB struct {
	n      int
	colors []int
}

func solveB(n int, arr []int) []int {
	cnt := make([]int, n+1)
	last := make([]int, n+1)
	for i := range last {
		last[i] = -1
	}
	for i, c := range arr {
		p := (i + 1) % 2
		if last[c] == -1 || last[c] != p {
			cnt[c]++
			last[c] = p
		}
	}
	return cnt[1:]
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	const T = 100
	tests := make([]caseB, T)
	var input bytes.Buffer
	fmt.Fprintln(&input, T)
	expected := make([][]int, T)
	for i := 0; i < T; i++ {
		n := rng.Intn(20) + 1
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = rng.Intn(n) + 1
		}
		tests[i] = caseB{n: n, colors: arr}
		fmt.Fprintln(&input, n)
		for j, v := range arr {
			if j+1 == len(arr) {
				fmt.Fprintf(&input, "%d\n", v)
			} else {
				fmt.Fprintf(&input, "%d ", v)
			}
		}
		expected[i] = solveB(n, arr)
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
	scanner.Split(bufio.ScanWords)
	for i := 0; i < T; i++ {
		for j := 0; j < tests[i].n; j++ {
			if !scanner.Scan() {
				fmt.Fprintln(os.Stderr, "insufficient output")
				os.Exit(1)
			}
			var got int
			fmt.Sscan(scanner.Text(), &got)
			if got != expected[i][j] {
				fmt.Fprintf(os.Stderr, "mismatch on test %d color %d: expected %d got %d\n", i+1, j+1, expected[i][j], got)
				os.Exit(1)
			}
		}
	}
	if scanner.Scan() {
		fmt.Fprintln(os.Stderr, "extra output after", T, "tests")
		os.Exit(1)
	}
	fmt.Println("all tests passed")
}
