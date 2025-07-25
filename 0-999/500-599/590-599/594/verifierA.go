package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

func solveA(a []int) int {
	sort.Ints(a)
	half := len(a) / 2
	ans := a[half] - a[0]
	for i := 1; i < half; i++ {
		if d := a[i+half] - a[i]; d < ans {
			ans = d
		}
	}
	return ans
}

func run(binary string, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, binary)
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
	for t := 0; t < 100; t++ {
		n := rand.Intn(50)*2 + 2 // even between 2 and 100
		seen := map[int]bool{}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			v := rand.Intn(1000000)
			for seen[v] {
				v = rand.Intn(1000000)
			}
			seen[v] = true
			arr[i] = v
		}
		input := fmt.Sprintf("%d\n", n)
		for i, v := range arr {
			if i > 0 {
				input += " "
			}
			input += strconv.Itoa(v)
		}
		input += "\n"
		expected := strconv.Itoa(solveA(append([]int{}, arr...)))
		got, err := run(bin, input)
		if err != nil {
			fmt.Println("test", t, "runtime error:", err)
			fmt.Println("output:", got)
			os.Exit(1)
		}
		if got != expected {
			fmt.Println("test", t, "failed")
			fmt.Println("input:\n" + input)
			fmt.Println("expected:", expected, "got:", got)
			os.Exit(1)
		}
	}
	fmt.Println("OK")
}
