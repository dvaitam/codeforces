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

func run(bin, input string) (string, error) {
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

func isPalindrome(arr []int) bool {
	n := len(arr)
	for i := 0; i < n/2; i++ {
		if arr[i] != arr[n-1-i] {
			return false
		}
	}
	return true
}

func permute(arr []int, f func([]int) bool) bool {
	var rec func(int) bool
	rec = func(i int) bool {
		if i == len(arr) {
			tmp := make([]int, len(arr))
			copy(tmp, arr)
			if f(tmp) {
				return true
			}
			return false
		}
		for j := i; j < len(arr); j++ {
			arr[i], arr[j] = arr[j], arr[i]
			if rec(i + 1) {
				return true
			}
			arr[i], arr[j] = arr[j], arr[i]
		}
		return false
	}
	return rec(0)
}

func possible(a []int, l, r int) bool {
	sub := append([]int(nil), a[l:r]...)
	return permute(sub, func(p []int) bool {
		b := append([]int(nil), a...)
		copy(b[l:r], p)
		return isPalindrome(b)
	})
}

func brute(a []int) int {
	n := len(a)
	cnt := 0
	for l := 0; l < n; l++ {
		for r := l; r < n; r++ {
			if possible(a, l, r+1) {
				cnt++
			}
		}
	}
	return cnt
}

func genTest(rng *rand.Rand) (string, int) {
	n := rng.Intn(5) + 1
	a := make([]int, n)
	for i := range a {
		a[i] = rng.Intn(3) + 1
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	exp := brute(a)
	return sb.String(), exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genTest(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		val, err := strconv.Atoi(strings.TrimSpace(out))
		if err != nil || val != exp {
			fmt.Printf("test %d failed: expected %d got %s\ninput:\n%s\n", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
