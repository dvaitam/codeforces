package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func kAmazing(a []int) []int {
	n := len(a)
	lastPos := make([]int, n+1)
	maxGap := make([]int, n+1)
	seen := make([]bool, n+1)
	for i := 1; i <= n; i++ {
		x := a[i-1]
		gap := i - lastPos[x]
		if gap > maxGap[x] {
			maxGap[x] = gap
		}
		lastPos[x] = i
		seen[x] = true
	}
	for x := 1; x <= n; x++ {
		if !seen[x] {
			continue
		}
		gap := (n + 1) - lastPos[x]
		if gap > maxGap[x] {
			maxGap[x] = gap
		}
	}
	const inf = int(1e9)
	best := make([]int, n+2)
	for i := range best {
		best[i] = inf
	}
	for x := 1; x <= n; x++ {
		if !seen[x] {
			continue
		}
		k := maxGap[x]
		if x < best[k] {
			best[k] = x
		}
	}
	res := make([]int, n)
	curr := inf
	for k := 1; k <= n; k++ {
		if best[k] < curr {
			curr = best[k]
		}
		if curr == inf {
			res[k-1] = -1
		} else {
			res[k-1] = curr
		}
	}
	return res
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	rand.Seed(1)
	for t := 1; t <= 100; t++ {
		n := rand.Intn(20) + 1
		a := make([]int, n)
		for i := range a {
			a[i] = rand.Intn(n) + 1
		}
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i, v := range a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		input := sb.String()
		got, err := runBinary(binary, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\n", t, err)
			os.Exit(1)
		}
		expectArr := kAmazing(a)
		expectFields := make([]string, len(expectArr))
		for i, v := range expectArr {
			expectFields[i] = fmt.Sprintf("%d", v)
		}
		expectStr := strings.Join(expectFields, " ")
		if got != expectStr {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\nexpected: %s\ngot: %s\n", t, expectStr, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
