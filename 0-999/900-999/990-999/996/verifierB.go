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

func expectedB(a []int64) int {
	n := len(a)
	nn := int64(n)
	minTime := int64(1<<63 - 1)
	ans := 0
	for i := 0; i < n; i++ {
		t := int64(i)
		if a[i] > int64(i) {
			diff := a[i] - int64(i)
			cycles := (diff + nn - 1) / nn
			t += cycles * nn
		}
		if t < minTime {
			minTime = t
			ans = i
		}
	}
	return ans + 1
}

func runCase(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func buildInput(arr []int64) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(arr)))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	fixed := [][]int64{
		{2, 3, 2, 0},
		{10, 10},
		{5, 2, 6, 5, 7, 4},
		make([]int64, 1000),
	}

	caseNum := 0
	for _, arr := range fixed {
		input := buildInput(arr)
		expect := fmt.Sprintf("%d", expectedB(arr))
		got, err := runCase(bin, input)
		caseNum++
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", caseNum, err, input)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", caseNum, expect, got, input)
			os.Exit(1)
		}
	}

	for ; caseNum < 100; caseNum++ {
		n := rng.Intn(20) + 1
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			arr[i] = rng.Int63n(1000000000)
		}
		input := buildInput(arr)
		expect := fmt.Sprintf("%d", expectedB(arr))
		got, err := runCase(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", caseNum+1, err, input)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", caseNum+1, expect, got, input)
			os.Exit(1)
		}
	}

	fmt.Println("All tests passed")
}
