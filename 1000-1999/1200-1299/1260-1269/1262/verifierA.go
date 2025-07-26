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

func solve(n int, seg [][2]int) int {
	lmax := seg[0][0]
	rmin := seg[0][1]
	for i := 1; i < n; i++ {
		if seg[i][0] > lmax {
			lmax = seg[i][0]
		}
		if seg[i][1] < rmin {
			rmin = seg[i][1]
		}
	}
	if lmax <= rmin {
		return 0
	}
	return lmax - rmin
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
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func runCase(bin string, seg [][2]int) error {
	n := len(seg)
	input := fmt.Sprintf("1\n%d\n", n)
	for _, p := range seg {
		input += fmt.Sprintf("%d %d\n", p[0], p[1])
	}
	expect := solve(n, seg)
	out, err := run(bin, input)
	if err != nil {
		return err
	}
	var got int
	if _, err := fmt.Sscan(out, &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got != expect {
		return fmt.Errorf("expected %d got %d", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := [][][2]int{
		{{4, 5}, {5, 9}, {7, 7}},
		{{11, 19}, {4, 17}, {16, 16}, {3, 12}, {14, 17}},
		{{1, 10}},
		{{1, 1}},
	}
	for idx, seg := range cases {
		if err := runCase(bin, seg); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
	}
	for i := len(cases); i < 100; i++ {
		n := rng.Intn(10) + 1
		seg := make([][2]int, n)
		for j := 0; j < n; j++ {
			a := rng.Intn(20)
			b := rng.Intn(20)
			if a > b {
				a, b = b, a
			}
			seg[j] = [2]int{a, b}
		}
		if err := runCase(bin, seg); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
