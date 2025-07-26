package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func maxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func solveD(x, y []int64, h int64) int64 {
	n := len(x)
	i, j := 0, 0
	sum := y[0] - x[0]
	budget := h
	for j+1 < n {
		gap := x[j+1] - y[j]
		if budget-gap > 0 {
			budget -= gap
		} else {
			break
		}
		j++
		sum += gap + (y[j] - x[j])
	}
	ans := sum + budget
	if j == n-1 {
		return ans
	}
	for i < n {
		i++
		if i >= n {
			break
		}
		sum -= x[i] - x[i-1]
		budget += x[i] - y[i-1]
		for j+1 < n {
			gap := x[j+1] - y[j]
			if budget-gap > 0 {
				budget -= gap
			} else {
				break
			}
			j++
			sum += gap + (y[j] - x[j])
		}
		ans = maxInt64(ans, sum+budget)
		if j == n-1 {
			break
		}
	}
	return ans
}

func run(bin, input string) (string, error) {
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
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(4)
	for t := 1; t <= 100; t++ {
		n := rand.Intn(20) + 1
		h := rand.Int63n(100) + 1
		x := make([]int64, n)
		y := make([]int64, n)
		cur := int64(rand.Intn(10))
		for i := 0; i < n; i++ {
			length := int64(rand.Intn(10) + 1)
			x[i] = cur
			y[i] = cur + length
			cur = y[i] + int64(rand.Intn(5)) + 1
		}
		input := fmt.Sprintf("%d %d\n", n, h)
		for i := 0; i < n; i++ {
			input += fmt.Sprintf("%d %d\n", x[i], y[i])
		}
		expect := fmt.Sprintf("%d", solveD(append([]int64(nil), x...), append([]int64(nil), y...), h))
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
