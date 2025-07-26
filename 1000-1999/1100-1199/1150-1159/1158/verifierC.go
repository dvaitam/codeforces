package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func nextFromPerm(p []int) []int {
	n := len(p)
	res := make([]int, n)
	for i := 0; i < n; i++ {
		nxt := n + 1
		for j := i + 1; j < n; j++ {
			if p[j] > p[i] {
				nxt = j + 1
				break
			}
		}
		res[i] = nxt
	}
	return res
}

func check(n int, given []int, out string) bool {
	if strings.TrimSpace(out) == "-1" {
		return false
	}
	parts := strings.Fields(out)
	if len(parts) != n {
		return false
	}
	perm := make([]int, n)
	used := make([]bool, n+1)
	for i, s := range parts {
		v, err := strconv.Atoi(s)
		if err != nil || v < 1 || v > n || used[v] {
			return false
		}
		used[v] = true
		perm[i] = v
	}
	calc := nextFromPerm(perm)
	for i := 0; i < n; i++ {
		if given[i] != -1 && given[i] != calc[i] {
			return false
		}
	}
	return true
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for tcase := 0; tcase < 100; tcase++ {
		n := rng.Intn(6) + 1
		perm := rand.Perm(n)
		for i := range perm {
			perm[i]++
		}
		next := nextFromPerm(perm)
		given := make([]int, n)
		for i := 0; i < n; i++ {
			if rng.Intn(3) == 0 {
				given[i] = -1
			} else {
				given[i] = next[i]
			}
		}
		input := fmt.Sprintf("1\n%d\n", n)
		for i, v := range given {
			if i > 0 {
				input += " "
			}
			input += fmt.Sprint(v)
		}
		input += "\n"
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", tcase+1, err, input)
			os.Exit(1)
		}
		if !check(n, given, got) {
			fmt.Fprintf(os.Stderr, "case %d failed: wrong output\ninput:%s output:%s", tcase+1, input, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
