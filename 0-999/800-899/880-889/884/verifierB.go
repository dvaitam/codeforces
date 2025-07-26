package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solve(n int, x int, a []int) string {
	sum := 0
	for _, v := range a {
		sum += v
	}
	if sum+n-1 == x {
		return "YES"
	}
	return "NO"
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(2)
	for tc := 0; tc < 100; tc++ {
		n := rand.Intn(50) + 1
		a := make([]int, n)
		sum := 0
		for i := 0; i < n; i++ {
			a[i] = rand.Intn(10000) + 1
			sum += a[i]
		}
		var x int
		if rand.Intn(2) == 0 {
			x = sum + n - 1
		} else {
			delta := rand.Intn(10) + 1
			if rand.Intn(2) == 0 {
				x = sum + n - 1 + delta
			} else {
				x = sum + n - 1 - delta
				if x < 1 {
					x = 1
				}
			}
		}
		var input bytes.Buffer
		fmt.Fprintf(&input, "%d %d\n", n, x)
		for i := 0; i < n; i++ {
			fmt.Fprintf(&input, "%d ", a[i])
		}
		fmt.Fprintln(&input)
		expected := solve(n, x, a)
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
		if got != expected {
			fmt.Println("wrong answer on test", tc+1)
			fmt.Println("input:\n" + input.String())
			fmt.Println("expected:", expected)
			fmt.Println("got:", got)
			os.Exit(1)
		}
	}
	fmt.Println("ok")
}
