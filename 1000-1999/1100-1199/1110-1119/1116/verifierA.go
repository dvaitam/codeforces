package main

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	expected := []float64{1 / math.Sqrt(3), 1 / math.Sqrt(3), 1 / math.Sqrt(3), 0}
	for i := 0; i < 100; i++ {
		cmd := exec.Command(bin)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "run %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		fields := strings.Fields(out.String())
		if len(fields) != 4 {
			fmt.Fprintf(os.Stderr, "test %d: expected 4 numbers, got %d\n", i+1, len(fields))
			os.Exit(1)
		}
		for j, f := range fields {
			val, err := parseFloat(f)
			if err != nil {
				fmt.Fprintf(os.Stderr, "test %d: invalid float %q\n", i+1, f)
				os.Exit(1)
			}
			if math.Abs(val-expected[j]) > 1e-3 {
				fmt.Fprintf(os.Stderr, "test %d: value %d = %f want %f\n", i+1, j, val, expected[j])
				os.Exit(1)
			}
		}
	}
	fmt.Println("ok")
}

func parseFloat(s string) (float64, error) {
	return strconv.ParseFloat(strings.TrimSpace(s), 64)
}
