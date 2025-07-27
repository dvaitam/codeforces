package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = nil
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	rand.Seed(2)
	for tcase := 1; tcase <= 100; tcase++ {
		n := rand.Intn(4) + 2 // 2..5
		a := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			a[i] = int64(rand.Intn(20) + 1)
		}
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 1; i <= n; i++ {
			if i > 1 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", a[i]))
		}
		sb.WriteByte('\n')
		input := sb.String()
		out, err := runBinary(binary, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\n", tcase, err)
			os.Exit(1)
		}
		sum := int64(0)
		for i := 1; i <= n; i++ {
			sum += a[i]
		}
		scanner := bufio.NewScanner(strings.NewReader(out))
		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "invalid output on test %d\n", tcase)
			os.Exit(1)
		}
		firstFields := strings.Fields(scanner.Text())
		if sum%int64(n) != 0 {
			if len(firstFields) != 1 || firstFields[0] != "-1" {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d: expected -1\n", tcase)
				os.Exit(1)
			}
			continue
		}
		if len(firstFields) != 1 {
			fmt.Fprintf(os.Stderr, "invalid output on test %d\n", tcase)
			os.Exit(1)
		}
		k, err := strconv.Atoi(firstFields[0])
		if err != nil || k < 0 || k > 3*n {
			fmt.Fprintf(os.Stderr, "invalid k on test %d\n", tcase)
			os.Exit(1)
		}
		for op := 0; op < k; op++ {
			if !scanner.Scan() {
				fmt.Fprintf(os.Stderr, "not enough operations on test %d\n", tcase)
				os.Exit(1)
			}
			f := strings.Fields(scanner.Text())
			if len(f) != 3 {
				fmt.Fprintf(os.Stderr, "invalid operation format on test %d\n", tcase)
				os.Exit(1)
			}
			i, _ := strconv.Atoi(f[0])
			j, _ := strconv.Atoi(f[1])
			x, _ := strconv.Atoi(f[2])
			if i < 1 || i > n || j < 1 || j > n || x < 0 {
				fmt.Fprintf(os.Stderr, "operation out of range on test %d\n", tcase)
				os.Exit(1)
			}
			val := int64(x) * int64(i)
			if a[i] < val {
				fmt.Fprintf(os.Stderr, "negative value on test %d\n", tcase)
				os.Exit(1)
			}
			a[i] -= val
			a[j] += val
		}
		if scanner.Scan() {
			extra := strings.TrimSpace(scanner.Text())
			if extra != "" {
				fmt.Fprintf(os.Stderr, "extra output on test %d\n", tcase)
				os.Exit(1)
			}
		}
		avg := a[1]
		for i := 2; i <= n; i++ {
			if a[i] != avg {
				fmt.Fprintf(os.Stderr, "array not equal on test %d\n", tcase)
				os.Exit(1)
			}
		}
	}
	fmt.Println("all tests passed")
}
