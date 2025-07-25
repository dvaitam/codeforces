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

func expected(a []int64) int64 {
	n := len(a)
	lmax := make([]int, n)
	rmax := make([]int, n)
	lmin := make([]int, n)
	rmin := make([]int, n)
	stack := make([]int, 0, n)
	for i := 0; i < n; i++ {
		for len(stack) > 0 && a[stack[len(stack)-1]] <= a[i] {
			stack = stack[:len(stack)-1]
		}
		if len(stack) == 0 {
			lmax[i] = -1
		} else {
			lmax[i] = stack[len(stack)-1]
		}
		stack = append(stack, i)
	}
	stack = stack[:0]
	for i := n - 1; i >= 0; i-- {
		for len(stack) > 0 && a[stack[len(stack)-1]] < a[i] {
			stack = stack[:len(stack)-1]
		}
		if len(stack) == 0 {
			rmax[i] = n
		} else {
			rmax[i] = stack[len(stack)-1]
		}
		stack = append(stack, i)
	}
	stack = stack[:0]
	for i := 0; i < n; i++ {
		for len(stack) > 0 && a[stack[len(stack)-1]] >= a[i] {
			stack = stack[:len(stack)-1]
		}
		if len(stack) == 0 {
			lmin[i] = -1
		} else {
			lmin[i] = stack[len(stack)-1]
		}
		stack = append(stack, i)
	}
	stack = stack[:0]
	for i := n - 1; i >= 0; i-- {
		for len(stack) > 0 && a[stack[len(stack)-1]] > a[i] {
			stack = stack[:len(stack)-1]
		}
		if len(stack) == 0 {
			rmin[i] = n
		} else {
			rmin[i] = stack[len(stack)-1]
		}
		stack = append(stack, i)
	}
	var result int64
	for i := 0; i < n; i++ {
		leftMax := i - lmax[i]
		rightMax := rmax[i] - i
		leftMin := i - lmin[i]
		rightMin := rmin[i] - i
		contrib := int64(leftMax)*int64(rightMax) - int64(leftMin)*int64(rightMin)
		result += contrib * a[i]
	}
	return result
}

func runBin(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Int63n(100) - 50
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	sb.WriteByte('\n')
	exp := fmt.Sprintf("%d", expected(arr))
	return sb.String(), exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, exp := genCase(rng)
		out, err := runBin(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\ninput:\n%s\n", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed:\ninput:\n%s\nexpected:%s\ngot:%s\n", i+1, input, exp, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
