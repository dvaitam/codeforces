package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solveD(a []int) (int, []int) {
	n := len(a)
	tag := true
	flag := false
	for i := 1; i < n; i++ {
		if a[i] != a[i-1] {
			tag = false
		} else {
			flag = true
		}
	}
	if a[0] == a[n-1] {
		flag = true
	}
	k := 0
	colors := make([]int, n)
	if tag {
		k = 1
		for i := 0; i < n; i++ {
			colors[i] = 1
		}
	} else if flag {
		k = 2
		if n%2 == 0 {
			for i := 0; i < n; i++ {
				if i%2 == 0 {
					colors[i] = 2
				} else {
					colors[i] = 1
				}
			}
		} else {
			op := 1
			vis := false
			colors[0] = op + 1
			for i := 1; i < n; i++ {
				if !vis && a[i] == a[i-1] {
					vis = true
				} else {
					op ^= 1
				}
				colors[i] = op + 1
			}
		}
	} else {
		if n%2 == 0 {
			k = 2
			for i := 0; i < n; i++ {
				if i%2 == 0 {
					colors[i] = 2
				} else {
					colors[i] = 1
				}
			}
		} else {
			k = 3
			op := 1
			colors[0] = 2
			for i := 1; i < n-1; i++ {
				op ^= 1
				colors[i] = op + 1
			}
			colors[n-1] = 3
		}
	}
	return k, colors
}

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
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(4)
	for t := 1; t <= 100; t++ {
		n := rand.Intn(20) + 3
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			arr[i] = rand.Intn(5) + 1
		}
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", arr[i]))
		}
		sb.WriteByte('\n')
		k, colors := solveD(arr)
		expect := fmt.Sprintf("%d\n", k)
		for i := 0; i < n; i++ {
			if i > 0 {
				expect += " "
			}
			expect += fmt.Sprintf("%d", colors[i])
		}
		out, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput:\n%s", t, err, sb.String())
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "test %d failed: expected\n%s\ngot\n%s\ninput:\n%s", t, expect, out, sb.String())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
