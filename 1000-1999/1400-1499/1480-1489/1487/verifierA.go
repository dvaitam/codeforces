package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func solveA(arr []int) int {
	minVal := 101
	for _, v := range arr {
		if v < minVal {
			minVal = v
		}
	}
	cnt := 0
	for _, v := range arr {
		if v > minVal {
			cnt++
		}
	}
	return cnt
}

func run(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(1))
	for t := 0; t < 100; t++ {
		n := rng.Intn(99) + 2 // 2..100
		arr := make([]int, n)
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(strconv.Itoa(n))
		sb.WriteByte('\n')
		for i := 0; i < n; i++ {
			arr[i] = rng.Intn(100) + 1
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(arr[i]))
		}
		sb.WriteByte('\n')
		input := sb.String()
		expected := strconv.Itoa(solveA(arr))
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n%s", t+1, err, got)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("case %d failed:\ninput:\n%sexpected %s got %s\n", t+1, input, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
