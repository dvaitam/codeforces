package main

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strings"
)

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
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func runCase(bin string, h, m int) error {
	input := fmt.Sprintf("%02d:%02d\n", h, m)
	out, err := run(bin, input)
	if err != nil {
		return err
	}
	var ah, am float64
	if _, err := fmt.Sscan(out, &ah, &am); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	expH := float64(h%12)*30.0 + float64(m)/2.0
	expM := float64(m) * 6.0
	if math.Abs(ah-expH) > 1e-6 || math.Abs(am-expM) > 1e-6 {
		return fmt.Errorf("expected %.9f %.9f got %.9f %.9f", expH, expM, ah, am)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	total := 0
	for h := 0; h < 24; h++ {
		for m := 0; m < 60; m++ {
			if err := runCase(bin, h, m); err != nil {
				fmt.Fprintf(os.Stderr, "case %d failed: %v (time %02d:%02d)\n", total+1, err, h, m)
				os.Exit(1)
			}
			total++
		}
	}
	fmt.Printf("All %d tests passed\n", total)
}
