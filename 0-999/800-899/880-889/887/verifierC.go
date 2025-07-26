package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func rotate(c []int, cycles [][]int) []int {
	res := make([]int, len(c))
	copy(res, c)
	for _, cyc := range cycles {
		tmp := c[cyc[0]]
		for i := 0; i < len(cyc)-1; i++ {
			res[cyc[i]] = c[cyc[i+1]]
		}
		res[cyc[len(cyc)-1]] = tmp
	}
	return res
}

func solved(c []int) bool {
	for i := 0; i < 24; i += 4 {
		if !(c[i] == c[i+1] && c[i] == c[i+2] && c[i] == c[i+3]) {
			return false
		}
	}
	return true
}

func solveCase(c []int) string {
	rotations := [][][]int{
		{{0, 1, 3, 2}, {4, 5, 8, 9, 12, 13, 16, 17}},
		{{20, 21, 23, 22}, {6, 7, 18, 19, 14, 15, 10, 11}},
		{{4, 5, 7, 6}, {2, 3, 16, 18, 21, 20, 9, 8}},
		{{12, 13, 15, 14}, {0, 1, 11, 9, 22, 23, 18, 16}},
		{{16, 17, 19, 18}, {0, 2, 4, 6, 20, 22, 15, 13}},
		{{8, 9, 11, 10}, {1, 3, 14, 12, 23, 21, 7, 5}},
	}
	for _, rot := range rotations {
		cw := rotate(c, rot)
		if solved(cw) {
			return "YES"
		}
		cc := rotate(cw, rot)
		cc = rotate(cc, rot)
		if solved(cc) {
			return "YES"
		}
	}
	return "NO"
}

func runCase(bin string, cube []int) error {
	var sb strings.Builder
	for i, v := range cube {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := solveCase(cube)
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(3)
	cases := make([][]int, 100)
	base := []int{1, 1, 1, 1, 2, 2, 2, 2, 3, 3, 3, 3, 4, 4, 4, 4, 5, 5, 5, 5, 6, 6, 6, 6}
	for i := range cases {
		arr := make([]int, 24)
		copy(arr, base)
		rand.Shuffle(24, func(i, j int) { arr[i], arr[j] = arr[j], arr[i] })
		cases[i] = arr
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
