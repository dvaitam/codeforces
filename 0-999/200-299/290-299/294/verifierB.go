package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type book struct{ t, w int }

func solveB(books []book) int {
	n := len(books)
	totalWidth := 0
	for _, b := range books {
		totalWidth += b.w
	}
	maxT := 2 * n
	const negInf = -1000000000
	dp := make([]int, maxT+1)
	for t := 1; t <= maxT; t++ {
		dp[t] = negInf
	}
	dp[0] = 0
	for _, b := range books {
		for t := maxT; t >= b.t; t-- {
			if dp[t-b.t] >= 0 && dp[t-b.t]+b.w > dp[t] {
				dp[t] = dp[t-b.t] + b.w
			}
		}
	}
	ans := maxT
	for t := 0; t <= maxT; t++ {
		if dp[t] >= 0 && dp[t]+t >= totalWidth {
			ans = t
			break
		}
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 1
	books := make([]book, n)
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			books[i].t = 1
		} else {
			books[i].t = 2
		}
		books[i].w = rng.Intn(10) + 1
	}
	ans := solveB(books)
	var in bytes.Buffer
	fmt.Fprintf(&in, "%d\n", n)
	for _, b := range books {
		fmt.Fprintf(&in, "%d %d\n", b.t, b.w)
	}
	out := fmt.Sprintf("%d\n", ans)
	return in.String(), out
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, expect := generateCase(rng)
		got, err := runBinary(bin, in)
		if err != nil {
			fmt.Printf("test %d: execution failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("test %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, in, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
