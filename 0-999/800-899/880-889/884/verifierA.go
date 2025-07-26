package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solve(n int, t int, a []int) int {
	for i := 0; i < n; i++ {
		t -= 86400 - a[i]
		if t <= 0 {
			return i + 1
		}
	}
	return n
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(1)
	for tc := 0; tc < 100; tc++ {
		n := rand.Intn(100) + 1
		t := rand.Intn(1000000) + 1
		a := make([]int, n)
		var input bytes.Buffer
		fmt.Fprintf(&input, "%d %d\n", n, t)
		for i := 0; i < n; i++ {
			a[i] = rand.Intn(86401)
			fmt.Fprintf(&input, "%d ", a[i])
		}
		fmt.Fprintln(&input)
		expected := solve(n, t, a)
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewReader(input.Bytes())
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Println("error running binary:", err)
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		exp := fmt.Sprint(expected)
		if got != exp {
			fmt.Println("wrong answer on test", tc+1)
			fmt.Println("input:\n" + input.String())
			fmt.Println("expected:", exp)
			fmt.Println("got:", got)
			os.Exit(1)
		}
	}
	fmt.Println("ok")
}
