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

var answer = [][]int{
	{}, {},
	{2, 1}, {3, 2}, {4, 3}, {5, 3}, {6, 5}, {7, 6}, {8, 6, 5, 4}, {9, 5},
	{10, 7}, {11, 9}, {12, 11, 10, 4}, {13, 12, 11, 8}, {14, 13, 12, 2},
	{15, 14}, {16, 14, 13, 11}, {17, 14}, {18, 11}, {19, 6, 2, 1},
	{20, 17}, {21, 19}, {22, 21}, {23, 18}, {24, 23, 22, 17},
	{25, 22}, {26, 6, 2, 1}, {27, 5, 2, 1}, {28, 25}, {29, 27},
	{30, 6, 4, 1}, {31, 28}, {32, 22, 2, 1}, {33, 20}, {34, 27, 2, 1},
	{35, 33}, {36, 25}, {37, 5, 4, 3, 2, 1}, {38, 6, 5, 1}, {39, 35},
	{40, 38, 21, 19}, {41, 38}, {42, 41, 20, 19}, {43, 42, 38, 37},
	{44, 43, 18, 17}, {45, 44, 42, 41}, {46, 45, 26, 25}, {47, 42},
	{48, 47, 21, 20}, {49, 40}, {50, 49, 24, 23}, {},
}

func expected(k int) string {
	arr := make([]int, k)
	for _, x := range answer[k] {
		if x > 0 && x <= k {
			arr[x-1] = 1
		}
	}
	parts := make([]string, k)
	for i, v := range arr {
		parts[i] = fmt.Sprint(v)
	}
	line := strings.Join(parts, " ")
	return line + "\n" + line + "\n"
}

func generateCase(rng *rand.Rand) (string, string) {
	k := rng.Intn(49) + 2
	return fmt.Sprintf("%d\n", k), expected(k)
}

func runCase(bin, input, exp string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	if out.String() != exp {
		return fmt.Errorf("expected %q got %q", exp, out.String())
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
