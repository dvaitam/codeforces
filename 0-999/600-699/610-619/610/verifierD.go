package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type cell struct{ x, y int64 }

func solveD(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return ""
	}
	cells := make(map[cell]struct{})
	for i := 0; i < n; i++ {
		var x1, y1, x2, y2 int64
		fmt.Fscan(in, &x1, &y1, &x2, &y2)
		if x1 == x2 {
			if y1 > y2 {
				y1, y2 = y2, y1
			}
			for y := y1; y <= y2; y++ {
				cells[cell{x1, y}] = struct{}{}
			}
		} else {
			if x1 > x2 {
				x1, x2 = x2, x1
			}
			for x := x1; x <= x2; x++ {
				cells[cell{x, y1}] = struct{}{}
			}
		}
	}
	return fmt.Sprint(len(cells))
}

func genTestD(rng *rand.Rand) string {
	n := rng.Intn(5) + 1
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%d\n", n)
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 { // vertical
			x := rng.Intn(11) - 5
			y1 := rng.Intn(11) - 5
			y2 := rng.Intn(11) - 5
			fmt.Fprintf(&buf, "%d %d %d %d\n", x, y1, x, y2)
		} else { // horizontal
			y := rng.Intn(11) - 5
			x1 := rng.Intn(11) - 5
			x2 := rng.Intn(11) - 5
			fmt.Fprintf(&buf, "%d %d %d %d\n", x1, y, x2, y)
		}
	}
	return buf.String()
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(1))
	for i := 1; i <= 100; i++ {
		in := genTestD(rng)
		expect := solveD(in)
		got, err := run(exe, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", i, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Printf("case %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i, in, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
