package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func runBinary(bin, input string) (string, error) {
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
	return out.String(), err
}

func digitsNeeded(x int) int {
	d := 1
	for x >= 7 {
		x /= 7
		d++
	}
	return d
}

func toBase7(x, d int) []int {
	res := make([]int, d)
	for i := d - 1; i >= 0; i-- {
		res[i] = x % 7
		x /= 7
	}
	return res
}

func distinct(a, b []int) bool {
	used := make([]bool, 7)
	for _, v := range a {
		if used[v] {
			return false
		}
		used[v] = true
	}
	for _, v := range b {
		if used[v] {
			return false
		}
		used[v] = true
	}
	return true
}

func expected(n, m int) int {
	d1 := digitsNeeded(n - 1)
	d2 := digitsNeeded(m - 1)
	if d1+d2 > 7 {
		return 0
	}
	ans := 0
	for h := 0; h < n; h++ {
		hd := toBase7(h, d1)
		for mm := 0; mm < m; mm++ {
			md := toBase7(mm, d2)
			if distinct(hd, md) {
				ans++
			}
		}
	}
	return ans
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go <binary>")
		return
	}
	bin := os.Args[1]
	rand.Seed(42)
	for i := 0; i < 100; i++ {
		n := rand.Intn(10) + 1
		m := rand.Intn(10) + 1
		input := fmt.Sprintf("%d %d\n", n, m)
		want := expected(n, m)
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\noutput:%s", i+1, err, out)
			return
		}
		out = strings.TrimSpace(out)
		got, err := strconv.Atoi(out)
		if err != nil {
			fmt.Printf("invalid output on test %d: %q\n", i+1, out)
			return
		}
		if got != want {
			fmt.Printf("wrong answer on test %d\ninput:%sexpected:%d got:%d\n", i+1, input, want, got)
			return
		}
	}
	fmt.Println("All tests passed")
}
