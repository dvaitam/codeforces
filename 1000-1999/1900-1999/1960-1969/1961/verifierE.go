package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func runBinary(bin, input string) (string, error) {
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

func runCase(bin string, arr []int) error {
	input := fmt.Sprintf("1\n%d", len(arr))
	for _, v := range arr {
		input += fmt.Sprintf(" %d", v)
	}
	input += "\n"
	got, err := runBinary(bin, input)
	if err != nil {
		return err
	}
	max := arr[0]
	for _, v := range arr {
		if v > max {
			max = v
		}
	}
	var out int
	if _, err := fmt.Sscan(got, &out); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if out != max {
		return fmt.Errorf("expected %d got %d", max, out)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(62))
	for i := 0; i < 100; i++ {
		n := rng.Intn(10) + 1
		arr := make([]int, n)
		for j := range arr {
			arr[j] = rng.Intn(100) + 1
		}
		if err := runCase(bin, arr); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
