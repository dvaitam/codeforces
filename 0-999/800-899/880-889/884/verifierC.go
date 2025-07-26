package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

func solve(n int, p []int) int64 {
	visited := make([]bool, n)
	var cycles []int
	for i := 0; i < n; i++ {
		if !visited[i] {
			cnt := 0
			j := i
			for !visited[j] {
				visited[j] = true
				j = p[j]
				cnt++
			}
			cycles = append(cycles, cnt)
		}
	}
	sort.Ints(cycles)
	if len(cycles) >= 2 {
		a := cycles[len(cycles)-1]
		b := cycles[len(cycles)-2]
		cycles = cycles[:len(cycles)-2]
		cycles = append(cycles, a+b)
	}
	var ans int64
	for _, c := range cycles {
		ans += int64(c) * int64(c)
	}
	return ans
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(3)
	for tc := 0; tc < 100; tc++ {
		n := rand.Intn(50) + 1
		p := rand.Perm(n)
		var input bytes.Buffer
		fmt.Fprintf(&input, "%d\n", n)
		for i := 0; i < n; i++ {
			fmt.Fprintf(&input, "%d ", p[i]+1)
		}
		fmt.Fprintln(&input)
		expected := solve(n, p)
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
		if got != fmt.Sprint(expected) {
			fmt.Println("wrong answer on test", tc+1)
			fmt.Println("input:\n" + input.String())
			fmt.Println("expected:", expected)
			fmt.Println("got:", got)
			os.Exit(1)
		}
	}
	fmt.Println("ok")
}
