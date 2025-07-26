package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
)

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func solveC(arr []int) string {
	expected := make([]int, len(arr))
	copy(expected, arr)
	sort.Slice(expected, func(i, j int) bool { return expected[i] > expected[j] })

	a := make([]int, len(arr))
	copy(a, arr)
	output := make([]int, 0, len(arr))

	for len(a) > 0 {
		changed := true
		for changed {
			changed = false
			for i := 0; i < len(a)-1; i++ {
				if a[i]-a[i+1] >= 2 {
					a[i]--
					a[i+1]++
					changed = true
				} else if a[i+1]-a[i] >= 2 {
					a[i+1]--
					a[i]++
					changed = true
				}
			}
		}
		idx := 0
		for i := 1; i < len(a); i++ {
			if a[i] > a[idx] {
				idx = i
			}
		}
		output = append(output, a[idx])
		a = append(a[:idx], a[idx+1:]...)
	}

	for i := 0; i < len(arr); i++ {
		if output[i] != expected[i] {
			return "NO"
		}
	}
	return "YES"
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := [][]int{
		{1},
		{2},
		{1, 1},
		{2, 1},
		{1, 2, 3},
		{10, 10, 10},
		{3, 1, 3, 1},
		{4, 5, 6, 7},
		{9, 8, 7, 6, 5},
		{1, 2, 1, 2, 1},
		{3, 3, 3, 3, 3, 3},
		{1, 3, 2, 4, 3, 5, 6},
		{2, 4, 6, 8, 10, 12, 14, 16},
		{16, 15, 14, 13, 12, 11, 10, 9},
		{1, 1, 2, 2, 3, 3, 4, 4, 5},
		{5, 4, 3, 2, 1, 1, 2, 3, 4},
		{10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
		{5, 5, 5, 5, 5, 5, 5, 5, 5, 5},
		{1, 100, 1, 100, 1, 100, 1, 100, 1, 100},
		{100, 99, 98, 97, 96, 95, 94, 93, 92, 91},
	}

	for i, arr := range tests {
		input := fmt.Sprintf("%d\n", len(arr))
		for j, v := range arr {
			if j > 0 {
				input += " "
			}
			input += fmt.Sprintf("%d", v)
		}
		input += "\n"
		want := solveC(arr)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Printf("Test %d failed: expected %q, got %q\n", i+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
