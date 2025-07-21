package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func generateCase() string {
	n := rand.Intn(99) + 2
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", rand.Intn(1000)+1))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func expected(in string) string {
	scanner := bufio.NewScanner(strings.NewReader(in))
	scanner.Split(bufio.ScanWords)
	var n int
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &n)
	a := make([]int, n)
	for i := 0; i < n; i++ {
		scanner.Scan()
		fmt.Sscan(scanner.Text(), &a[i])
	}
	minDiff := abs(a[n-1] - a[0])
	ans1, ans2 := n, 1
	for i := 1; i < n; i++ {
		d := abs(a[i] - a[i-1])
		if d < minDiff {
			minDiff = d
			ans1 = i
			ans2 = i + 1
		}
	}
	return fmt.Sprintf("%d %d", ans1, ans2)
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	rand.Seed(time.Now().UnixNano())
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		tc := generateCase()
		exp := expected(tc)
		got, err := run(bin, tc)
		if err != nil {
			fmt.Printf("case %d: error executing binary: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("case %d failed\ninput:\n%s\nexpected:%s\ngot:%s\n", i+1, tc, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
