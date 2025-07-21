package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func generateCase() string {
	n := rand.Intn(100) + 1
	m := rand.Intn(n) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", rand.Intn(2001)-1000))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func expected(in string) string {
	scanner := bufio.NewScanner(strings.NewReader(in))
	scanner.Split(bufio.ScanWords)
	var n, m int
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &n)
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &m)
	a := make([]int, n)
	for i := 0; i < n; i++ {
		scanner.Scan()
		fmt.Sscan(scanner.Text(), &a[i])
	}
	sort.Ints(a)
	res := 0
	for i := 0; i < n && i < m; i++ {
		if a[i] < 0 {
			res += -a[i]
		} else {
			break
		}
	}
	return fmt.Sprintf("%d", res)
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
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
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
