package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

var a, b, ta, tb []int

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func solve(as, ae, bs, be int) int {
	cnt0, cnt1 := 0, 0
	for i := as; i <= ae; i++ {
		if a[i]&1 == 0 {
			cnt0++
		} else {
			cnt1++
		}
	}
	for i := bs; i <= be; i++ {
		if (b[i]&1)^1 == 0 {
			cnt0++
		} else {
			cnt1++
		}
	}
	ret := maxInt(cnt0, cnt1)
	if as > ae || bs > be {
		return ret
	}
	if as == ae && bs == be {
		return maxInt(ret, 2)
	}
	j := as
	for i := as; i <= ae; i++ {
		if a[i]&1 == 1 {
			ta[j] = a[i] >> 1
			j++
		}
	}
	A := j
	for i := as; i <= ae; i++ {
		if a[i]&1 == 0 {
			ta[j] = a[i] >> 1
			j++
		}
	}
	for i := as; i <= ae; i++ {
		a[i] = ta[i]
	}
	k := bs
	for i := bs; i <= be; i++ {
		if b[i]&1 == 1 {
			tb[k] = b[i] >> 1
			k++
		}
	}
	B := k
	for i := bs; i <= be; i++ {
		if b[i]&1 == 0 {
			tb[k] = b[i] >> 1
			k++
		}
	}
	for i := bs; i <= be; i++ {
		b[i] = tb[i]
	}
	ret = maxInt(ret, solve(as, A-1, bs, B-1))
	ret = maxInt(ret, solve(A, ae, B, be))
	return ret
}

func expected(n int, arrA []int, m int, arrB []int) int {
	a = arrA
	b = arrB
	ta = make([]int, n)
	tb = make([]int, m)
	return solve(0, n-1, 0, m-1)
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
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(6)
	for t := 1; t <= 100; t++ {
		n := rand.Intn(10) + 1
		m := rand.Intn(10) + 1
		y1 := rand.Intn(100)
		y2 := y1 + rand.Intn(100) + 1
		arrA := make([]int, n)
		for i := 0; i < n; i++ {
			arrA[i] = rand.Intn(1000)
		}
		arrB := make([]int, m)
		for i := 0; i < m; i++ {
			arrB[i] = rand.Intn(1000)
		}
		input := fmt.Sprintf("%d %d\n", n, y1)
		for i := 0; i < n; i++ {
			if i > 0 {
				input += " "
			}
			input += fmt.Sprintf("%d", arrA[i])
		}
		input += fmt.Sprintf("\n%d %d\n", m, y2)
		for i := 0; i < m; i++ {
			if i > 0 {
				input += " "
			}
			input += fmt.Sprintf("%d", arrB[i])
		}
		input += "\n"
		expect := fmt.Sprintf("%d", expected(n, append([]int(nil), arrA...), m, append([]int(nil), arrB...)))
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
