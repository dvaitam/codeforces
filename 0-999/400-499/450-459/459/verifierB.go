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

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func expectedAnswer(arr []int64) (int64, int64) {
	minv := arr[0]
	maxv := arr[0]
	for _, v := range arr {
		if v < minv {
			minv = v
		}
		if v > maxv {
			maxv = v
		}
	}
	diff := maxv - minv
	var cntMin, cntMax int64
	for _, v := range arr {
		if v == minv {
			cntMin++
		}
		if v == maxv {
			cntMax++
		}
	}
	var ways int64
	if diff == 0 {
		n := int64(len(arr))
		ways = n * (n - 1) / 2
	} else {
		ways = cntMin * cntMax
	}
	return diff, ways
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go <binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for t := 0; t < 100; t++ {
		n := rand.Intn(20) + 2
		arr := make([]int64, n)
		for i := range arr {
			arr[i] = int64(rand.Intn(1000))
		}
		input := fmt.Sprintf("%d\n", n)
		for i, v := range arr {
			if i > 0 {
				input += " "
			}
			input += fmt.Sprintf("%d", v)
		}
		input += "\n"

		expectedDiff, expectedWays := expectedAnswer(arr)
		outStr, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n", t+1, err)
			os.Exit(1)
		}
		var gotDiff, gotWays int64
		nRead, err := fmt.Sscanf(outStr, "%d %d", &gotDiff, &gotWays)
		if err != nil || nRead != 2 {
			fmt.Printf("invalid output on test %d\ninput:\n%soutput:\n%s\n", t+1, input, outStr)
			os.Exit(1)
		}
		if gotDiff != expectedDiff || gotWays != expectedWays {
			fmt.Printf("wrong answer on test %d\ninput:\n%sexpected: %d %d\n got: %d %d\n", t+1, input, expectedDiff, expectedWays, gotDiff, gotWays)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
