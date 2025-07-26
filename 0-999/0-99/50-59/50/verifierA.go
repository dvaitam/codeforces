package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	type test struct{ m, n int }
	var tests []test
	for m := 1; m <= 16; m++ {
		for n := m; n <= 16; n++ {
			tests = append(tests, test{m, n})
		}
	}
	for i, t := range tests {
		input := fmt.Sprintf("%d %d\n", t.m, t.n)
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewBufferString(input)
		var out bytes.Buffer
		cmd.Stdout = &out
		if err := cmd.Run(); err != nil {
			fmt.Printf("Test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		var got int
		if _, err := fmt.Fscan(&out, &got); err != nil {
			fmt.Printf("Test %d failed to parse output: %v\n", i+1, err)
			os.Exit(1)
		}
		want := (t.m * t.n) / 2
		if got != want {
			fmt.Printf("Test %d failed: expected %d got %d\n", i+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
