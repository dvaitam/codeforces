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

func solveTime(a []int) int {
	n := len(a)
	d := make([]int, n+3)
	p := make([]int, n+3)
	var even func(int) int
	var odd func(int) int
	max := func(x, y int) int {
		if x > y {
			return x
		}
		return y
	}
	even = func(m int) int {
		if m <= 0 {
			return 0
		}
		if d[m] != 0 {
			return d[m]
		}
		if m == 2 {
			d[m] = max(a[0], a[1])
			p[m] = 0
			return d[m]
		}
		smallest := even(m-2) + max(a[m-2], a[m-1])
		k := m - 2
		sum := 0
		for i := m - 2; i >= 0; i -= 2 {
			cur := even(i) + sum + max(a[i], a[m-1])
			sum += max(a[i], a[i-1])
			if cur < smallest {
				smallest = cur
				k = i
			}
		}
		sum = 0
		for i := m - 3; i >= 0; i -= 2 {
			cur := odd(i) + sum + max(a[i], a[m-1])
			sum += max(a[i], a[i+1])
			if cur < smallest {
				smallest = cur
				k = i
			}
		}
		d[m] = smallest
		p[m] = k
		return smallest
	}
	odd = func(m int) int {
		if m == 1 {
			return max(a[0], a[2])
		}
		if d[m] != 0 {
			return d[m]
		}
		smallest := even(m-1) + max(a[m-1], a[m+1])
		k := m - 1
		sum := 0
		for i := m - 1; i >= 0; i -= 2 {
			cur := even(i) + sum + max(a[i], a[m+1])
			sum += max(a[i], a[i-1])
			if cur < smallest {
				smallest = cur
				k = i
			}
		}
		sum = 0
		for i := m - 2; i >= 0; i -= 2 {
			cur := odd(i) + sum + max(a[i], a[m+1])
			sum += max(a[i], a[i+1])
			if cur < smallest {
				smallest = cur
				k = i
			}
		}
		d[m] = smallest
		p[m] = k
		return smallest
	}
	if n%2 == 0 {
		return even(n)
	}
	return odd(n)
}

func generateCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(6) + 1
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = rng.Intn(20) + 1
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", n)
	for i, v := range a {
		if i > 0 {
			fmt.Fprint(&b, " ")
		}
		fmt.Fprint(&b, v)
	}
	fmt.Fprintln(&b)
	expected := solveTime(append([]int{}, a...))
	return b.String(), expected
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 20; i++ {
		input, expected := generateCase(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		fields := strings.Fields(out)
		if len(fields) == 0 {
			fmt.Printf("case %d: no output\n", i+1)
			os.Exit(1)
		}
		var got int
		fmt.Sscan(fields[0], &got)
		if got != expected {
			fmt.Printf("case %d failed expected %d got %d\ninput:\n%s", i+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
