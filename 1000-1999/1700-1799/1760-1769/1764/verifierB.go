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

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	cand := os.Args[1]

	rand.Seed(time.Now().UnixNano())
	const t = 100
	var input bytes.Buffer
	fmt.Fprintln(&input, t)
	expected := make([]int64, t)
	for i := 0; i < t; i++ {
		n := rand.Intn(9) + 2
		fmt.Fprintln(&input, n)
		arr := make([]int64, n)
		base := rand.Int63n(1000) + 1
		arr[0] = base
		fmt.Fprint(&input, arr[0])
		for j := 1; j < n; j++ {
			base += rand.Int63n(10) + 1
			arr[j] = base
			fmt.Fprintf(&input, " %d", arr[j])
		}
		fmt.Fprintln(&input)
		g := arr[0]
		for j := 1; j < n; j++ {
			g = gcd(g, arr[j])
		}
		expected[i] = arr[n-1] / g
	}

	cmd := exec.Command(cand)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	outBytes, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("candidate run error:", err)
		os.Exit(1)
	}
	scanner := bufio.NewScanner(bytes.NewReader(outBytes))
	for i := 0; i < t; i++ {
		if !scanner.Scan() {
			fmt.Println("missing output for test", i+1)
			os.Exit(1)
		}
		var ans int64
		fmt.Sscan(scanner.Text(), &ans)
		if ans != expected[i] {
			fmt.Printf("wrong answer on test %d: expected %d got %d\n", i+1, expected[i], ans)
			os.Exit(1)
		}
	}
	if scanner.Scan() {
		fmt.Println("extra output")
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
