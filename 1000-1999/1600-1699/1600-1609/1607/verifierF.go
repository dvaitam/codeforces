package main

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func runCmd(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = io.Discard
	err := cmd.Run()
	return out.String(), err
}

func randBoard() (int, int, []string) {
	n := rand.Intn(5) + 1
	m := rand.Intn(5) + 1
	board := make([]string, n)
	dirs := "LRUD"
	for i := 0; i < n; i++ {
		b := make([]byte, m)
		for j := 0; j < m; j++ {
			b[j] = dirs[rand.Intn(4)]
		}
		board[i] = string(b)
	}
	return n, m, board
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		return
	}
	candidate := os.Args[1]

	refBin := "./refF.bin"
	if err := exec.Command("go", "build", "-o", refBin, "1607F.go").Run(); err != nil {
		fmt.Println("failed to build reference solution:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	rand.Seed(6)
	type tc struct {
		n, m  int
		board []string
	}
	cases := []tc{
		{1, 1, []string{"U"}},
	}
	for len(cases) < 100 {
		n, m, b := randBoard()
		cases = append(cases, tc{n, m, b})
	}

	for i, c := range cases {
		input := fmt.Sprintf("1\n%d %d\n", c.n, c.m)
		for _, row := range c.board {
			input += row + "\n"
		}
		exp, err := runCmd(refBin, input)
		if err != nil {
			fmt.Println("reference solution failed:", err)
			os.Exit(1)
		}
		got, err := runCmd(candidate, input)
		if err != nil {
			fmt.Printf("test %d: candidate runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		exp = strings.TrimSpace(exp)
		got = strings.TrimSpace(got)
		if exp != got {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
