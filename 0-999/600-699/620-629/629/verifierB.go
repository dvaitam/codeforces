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

type Friend struct {
	g    byte
	a, b int
}

func solveB(n int, friends []Friend) int {
	const days = 366
	male := make([]int, days+1)
	female := make([]int, days+1)
	for _, f := range friends {
		for d := f.a; d <= f.b; d++ {
			if f.g == 'M' {
				male[d]++
			} else {
				female[d]++
			}
		}
	}
	ans := 0
	for d := 1; d <= days; d++ {
		m := male[d]
		f := female[d]
		if m < f {
			if ans < m*2 {
				ans = m * 2
			}
		} else {
			if ans < f*2 {
				ans = f * 2
			}
		}
	}
	return ans
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for t := 0; t < 100; t++ {
		n := rand.Intn(50) + 1
		friends := make([]Friend, n)
		var input bytes.Buffer
		fmt.Fprintln(&input, n)
		for i := 0; i < n; i++ {
			g := byte('M')
			if rand.Intn(2) == 0 {
				g = 'F'
			}
			a := rand.Intn(366) + 1
			b := rand.Intn(366-a+1) + a
			friends[i] = Friend{g: g, a: a, b: b}
			fmt.Fprintf(&input, "%c %d %d\n", g, a, b)
		}
		expected := solveB(n, friends)
		cmd := exec.Command(binary)
		cmd.Stdin = &input
		var out bytes.Buffer
		cmd.Stdout = &out
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "test %d: binary error: %v\n", t, err)
			os.Exit(1)
		}
		scanner := bufio.NewScanner(&out)
		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "test %d: no output\n", t)
			os.Exit(1)
		}
		var got int
		if _, err := fmt.Sscan(scanner.Text(), &got); err != nil {
			fmt.Fprintf(os.Stderr, "test %d: invalid output\n", t)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "test %d: expected %d got %d\n", t, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
