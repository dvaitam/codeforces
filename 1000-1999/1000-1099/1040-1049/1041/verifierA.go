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

func solveA(a []int64) int64 {
	sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
	var c int64
	for i := 0; i+1 < len(a); i++ {
		diff := a[i+1] - a[i] - 1
		if diff > 0 {
			c += diff
		}
	}
	return c
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(1)
	for t := 1; t <= 100; t++ {
		n := rand.Intn(100) + 1
		set := make(map[int64]bool)
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			v := rand.Int63n(1000)
			for set[v] {
				v = rand.Int63n(1000)
			}
			arr[i] = v
			set[v] = true
		}
		input := fmt.Sprintf("%d\n", n)
		for i := 0; i < n; i++ {
			if i > 0 {
				input += " "
			}
			input += fmt.Sprintf("%d", arr[i])
		}
		input += "\n"
		expect := fmt.Sprintf("%d", solveA(append([]int64(nil), arr...)))
		out, err := run(bin, input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n", t, err)
			os.Exit(1)
		}
		if out != expect {
			fmt.Printf("test %d failed: expected %s got %s\n", t, expect, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
