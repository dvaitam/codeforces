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

type TestCase struct {
	Input  string
	Output string
}

func runBinary(bin string, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("timeout")
		}
		return "", fmt.Errorf("%v: %s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func digits(x int) int {
	d := 0
	for x > 0 {
		d++
		x /= 10
	}
	return d
}

func solveD(a []int, k int) string {
	n := len(a)
	aMod := make([]int, n)
	lens := make([]int, n)
	v := make([][]int, 11)
	for i := 0; i < n; i++ {
		d := digits(a[i])
		lens[i] = d
		m := a[i] % k
		aMod[i] = m
		v[d] = append(v[d], m)
	}
	for d := 1; d <= 10; d++ {
		if len(v[d]) > 1 {
			sort.Ints(v[d])
		}
	}
	var res int64
	for i := 0; i < n; i++ {
		x := aMod[i]
		for j := 1; j <= 10; j++ {
			x = (x * 10) % k
			arr := v[j]
			if len(arr) == 0 {
				continue
			}
			y := (k - x) % k
			l := sort.Search(len(arr), func(idx int) bool { return arr[idx] >= y })
			r := sort.Search(len(arr), func(idx int) bool { return arr[idx] > y })
			res += int64(r - l)
			if lens[i] == j && aMod[i] == y {
				res--
			}
		}
	}
	return fmt.Sprintf("%d", res)
}

func generateTests() []TestCase {
	rand.Seed(45)
	tests := make([]TestCase, 100)
	for t := 0; t < 100; t++ {
		n := rand.Intn(20) + 1
		k := rand.Intn(98) + 2
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			arr[i] = rand.Intn(1000) + 1
		}
		inputBuilder := strings.Builder{}
		inputBuilder.WriteString(fmt.Sprintf("%d %d\n", n, k))
		for i := 0; i < n; i++ {
			if i > 0 {
				inputBuilder.WriteByte(' ')
			}
			inputBuilder.WriteString(fmt.Sprintf("%d", arr[i]))
		}
		inputBuilder.WriteByte('\n')
		output := solveD(arr, k)
		tests[t] = TestCase{Input: inputBuilder.String(), Output: output}
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		got, err := runBinary(bin, tc.Input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != strings.TrimSpace(tc.Output) {
			fmt.Fprintf(os.Stderr, "Test %d failed:\ninput:\n%s\nexpected:%s\n got:%s\n", i+1, tc.Input, tc.Output, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
