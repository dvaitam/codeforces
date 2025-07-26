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

func min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func expected(n int, lengths []int64, s string) string {
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		a[i] = lengths[i] * 2
	}
	var ans, now, nowW, nowG int64
	now = 6
	for i := 0; i < n; i++ {
		switch s[i] {
		case 'W':
			now = 4
			nowW += a[i] / 2
			ans += a[i] * 2
		case 'G':
			tmp := min64(nowW, a[i]/2)
			nowW -= tmp
			a[i] -= tmp * 2
			ans += tmp * 4
			nowG += tmp * 2
			nowG += a[i] / 2
			ans += a[i] * 3
		case 'L':
			tmp := min64(nowW, a[i]/2)
			nowW -= tmp
			a[i] -= tmp * 2
			ans += tmp * 4
			tmp = min64(nowG, a[i]/2)
			nowG -= tmp
			a[i] -= tmp * 2
			ans += tmp * 6
			ans += a[i] * now
		}
	}
	return fmt.Sprintf("%d", ans/2)
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
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for t := 0; t < 100; t++ {
		n := rand.Intn(20) + 1
		lengths := make([]int64, n)
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			lengths[i] = rand.Int63n(100) + 1
			input.WriteString(fmt.Sprintf("%d ", lengths[i]))
		}
		input.WriteByte('\n')
		terrain := make([]byte, n)
		for i := 0; i < n; i++ {
			terrain[i] = "GWL"[rand.Intn(3)]
		}
		s := string(terrain)
		input.WriteString(s)
		input.WriteByte('\n')
		expect := expected(n, lengths, s)
		got, err := run(bin, input.String())
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n", t+1, err)
			fmt.Println("input:\n", input.String())
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("wrong answer on test %d\n", t+1)
			fmt.Println("input:\n", input.String())
			fmt.Printf("expected: %s\n got: %s\n", expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
