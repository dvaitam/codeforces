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

func expectedA(a, b, x, y int) int {
	left := x * b
	right := (a - x - 1) * b
	top := y * a
	bottom := (b - y - 1) * a
	ans := left
	if right > ans {
		ans = right
	}
	if top > ans {
		ans = top
	}
	if bottom > ans {
		ans = bottom
	}
	return ans
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
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
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		a := rand.Intn(10000) + 1
		b := rand.Intn(10000) + 1
		if a+b <= 2 {
			a = 2
			b = 1
		}
		x := rand.Intn(a)
		y := rand.Intn(b)
		input := fmt.Sprintf("1\n%d %d %d %d\n", a, b, x, y)
		expect := expectedA(a, b, x, y)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		if got != fmt.Sprint(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\ninput:%s", i+1, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
