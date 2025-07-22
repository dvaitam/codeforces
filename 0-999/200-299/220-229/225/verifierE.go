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

func runBinary(bin string, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("timeout")
		}
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func solve(n int) int64 {
	const mod = 1000000007
	count := 0
	for z := 1; ; z++ {
		found := false
		for x := 1; x <= 2*z; x++ {
			f := x / 2
			num := z - f
			if num <= 0 {
				break
			}
			denom := x + 1
			if num%denom == 0 {
				y := num / denom
				if y > 0 {
					found = true
					break
				}
			}
		}
		if !found {
			count++
			if count == n {
				return int64(z) % mod
			}
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	r := rand.New(rand.NewSource(5))
	for t := 1; t <= 100; t++ {
		n := r.Intn(20) + 1
		input := fmt.Sprintf("%d\n", n)
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("Test %d runtime error: %v\nInput:%s", t, err, input)
			return
		}
		expected := solve(n)
		got, errp := strconv.ParseInt(strings.TrimSpace(out), 10, 64)
		if errp != nil || got != expected {
			fmt.Printf("Test %d FAILED\nInput:%sExpected:%d Got:%s\n", t, input, expected, out)
			return
		}
	}
	fmt.Println("All tests passed")
}
