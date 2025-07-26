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

func expected(n int, coords [][2]int64) string {
	var sx, sy int64
	for _, p := range coords {
		sx += p[0]
		sy += p[1]
	}
	return fmt.Sprintf("%d %d", sx/int64(n), sy/int64(n))
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
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for t := 0; t < 100; t++ {
		n := rand.Intn(10) + 1
		coords := make([][2]int64, 2*n)
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < 2*n; i++ {
			x := rand.Int63n(2000001) - 1000000
			y := rand.Int63n(2000001) - 1000000
			coords[i] = [2]int64{x, y}
			input.WriteString(fmt.Sprintf("%d %d\n", x, y))
		}
		expect := expected(n, coords)
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
