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
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(8) + 1
	perm1 := rng.Perm(n)
	perm2 := rng.Perm(n)
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range perm1 {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v + 1))
	}
	sb.WriteByte('\n')
	for i, v := range perm2 {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v + 1))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func lisLength(arr []int) int {
	d := make([]int, 0, len(arr))
	for _, v := range arr {
		lo, hi := 0, len(d)
		for lo < hi {
			mid := (lo + hi) / 2
			if d[mid] >= v {
				hi = mid
			} else {
				lo = mid + 1
			}
		}
		if lo == len(d) {
			d = append(d, v)
		} else {
			d[lo] = v
		}
	}
	return len(d)
}

func expectedOutput(input string) (string, error) {
	in := strings.NewReader(input)
	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return "", err
	}
	out := make([]string, 0, t)
	for ; t > 0; t-- {
		var n int
		if _, err := fmt.Fscan(in, &n); err != nil {
			return "", err
		}
		a := make([]int, n)
		for i := 0; i < n; i++ {
			if _, err := fmt.Fscan(in, &a[i]); err != nil {
				return "", err
			}
		}
		pos := make([]int, n+1)
		for i, v := range a {
			pos[v] = i
		}
		mapped := make([]int, n)
		for i := 0; i < n; i++ {
			var x int
			if _, err := fmt.Fscan(in, &x); err != nil {
				return "", err
			}
			mapped[i] = pos[x]
		}
		ans := n - lisLength(mapped)
		out = append(out, strconv.Itoa(ans))
	}
	return strings.Join(out, "\n"), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	fixedCases := []string{
		"1\n5\n4 5 2 3 1\n5 4 3 1 2\n",
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		fixedCases = append(fixedCases, genCase(rng))
	}

	for i, input := range fixedCases {
		exp, err := expectedOutput(input)
		if err != nil {
			fmt.Fprintln(os.Stderr, "failed to compute expected output:", err)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if exp != got {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:%s\nGot:%s\n", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
