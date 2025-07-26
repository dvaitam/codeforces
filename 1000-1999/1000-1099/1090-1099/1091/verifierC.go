package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func expected(n int64) string {
	var divisors []int64
	for i := int64(1); i*i <= n; i++ {
		if n%i == 0 {
			divisors = append(divisors, i)
			if i*i != n {
				divisors = append(divisors, n/i)
			}
		}
	}
	var ans []int64
	for _, d := range divisors {
		l := n + 1 - d
		x := (l-1)/d + 1
		y := (1 + l) * x / 2
		ans = append(ans, y)
	}
	sort.Slice(ans, func(i, j int) bool { return ans[i] < ans[j] })
	var sb strings.Builder
	for i, v := range ans {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	return sb.String()
}

func run(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for t := 0; t < 100; t++ {
		n := rand.Int63n(1e9) + 2
		input := fmt.Sprintf("%d\n", n)
		expect := expected(n)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n", t+1, err)
			fmt.Println("input:", input)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("wrong answer on test %d\n", t+1)
			fmt.Println("input:", input)
			fmt.Printf("expected: %s\n got: %s\n", expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
