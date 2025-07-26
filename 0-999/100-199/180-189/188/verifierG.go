package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func run(binary, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, binary)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func generateArray(r *rand.Rand) ([]int, string) {
	n := r.Intn(10) + 1
	arr := make([]int, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		arr[i] = r.Intn(100) + 1
		sb.WriteString(fmt.Sprintf("%d\n", arr[i]))
	}
	return arr, sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	r := rand.New(rand.NewSource(46))
	for i := 1; i <= 100; i++ {
		arr, input := generateArray(r)
		expectedArr := make([]int, len(arr))
		copy(expectedArr, arr)
		sort.Ints(expectedArr)
		var expBuf bytes.Buffer
		for j, v := range expectedArr {
			if j > 0 {
				expBuf.WriteByte(' ')
			}
			fmt.Fprint(&expBuf, v)
		}
		expected := expBuf.String()
		out, err := run(bin, input)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n", i, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != expected {
			fmt.Printf("wrong answer on test %d: expected %s got %s\n", i, expected, strings.TrimSpace(out))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
