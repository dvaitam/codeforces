package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCase struct {
	arr []int64
}

func compute(arr []int64) (int64, int64) {
	even, odd := int64(1), int64(0)
	pos, neg := int64(0), int64(0)
	parity := 0
	for _, x := range arr {
		if x < 0 {
			parity ^= 1
		}
		if parity == 0 {
			pos += even
			neg += odd
			even++
		} else {
			pos += odd
			neg += even
			odd++
		}
	}
	return neg, pos
}

func run(binary, input string) (string, error) {
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
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]

	rand.Seed(43)
	const T = 100
	for i := 0; i < T; i++ {
		n := rand.Intn(20) + 1
		arr := make([]int64, n)
		for j := range arr {
			v := rand.Intn(9) + 1
			if rand.Intn(2) == 0 {
				v = -v
			}
			arr[j] = int64(v)
		}
		input := fmt.Sprintf("%d\n", n)
		for j, v := range arr {
			if j > 0 {
				input += " "
			}
			input += fmt.Sprintf("%d", v)
		}
		input += "\n"

		expNeg, expPos := compute(arr)
		out, err := run(binary, input)
		if err != nil {
			fmt.Printf("test %d: execution failed: %v\ninput: %s\n", i+1, err, input)
			os.Exit(1)
		}
		var gotNeg, gotPos int64
		if _, err := fmt.Sscanf(out, "%d %d", &gotNeg, &gotPos); err != nil {
			fmt.Printf("test %d: failed to parse output '%s'\n", i+1, out)
			os.Exit(1)
		}
		if gotNeg != expNeg || gotPos != expPos {
			fmt.Printf("test %d failed: expected %d %d got %d %d\ninput:%s\n", i+1, expNeg, expPos, gotNeg, gotPos, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
