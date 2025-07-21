package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		arr[i] = int64(rng.Intn(21) - 10)
	}
	m := rng.Intn(20) + 1
	var input strings.Builder
	input.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range arr {
		if i > 0 {
			input.WriteByte(' ')
		}
		input.WriteString(fmt.Sprintf("%d", v))
	}
	input.WriteByte('\n')
	input.WriteString(fmt.Sprintf("%d\n", m))
	var expected strings.Builder
	for i := 0; i < m; i++ {
		if rng.Intn(2) == 0 { // query
			lf := rng.Intn(n)
			rg := rng.Intn(n)
			input.WriteString(fmt.Sprintf("%d %d\n", lf, rg))
			minVal := int64(math.MaxInt64)
			if lf <= rg {
				for j := lf; j <= rg; j++ {
					if arr[j] < minVal {
						minVal = arr[j]
					}
				}
			} else {
				for j := lf; j < n; j++ {
					if arr[j] < minVal {
						minVal = arr[j]
					}
				}
				for j := 0; j <= rg; j++ {
					if arr[j] < minVal {
						minVal = arr[j]
					}
				}
			}
			expected.WriteString(fmt.Sprintf("%d\n", minVal))
		} else { // update
			lf := rng.Intn(n)
			rg := rng.Intn(n)
			v := rng.Intn(21) - 10
			input.WriteString(fmt.Sprintf("%d %d %d\n", lf, rg, v))
			if lf <= rg {
				for j := lf; j <= rg; j++ {
					arr[j] += int64(v)
				}
			} else {
				for j := lf; j < n; j++ {
					arr[j] += int64(v)
				}
				for j := 0; j <= rg; j++ {
					arr[j] += int64(v)
				}
			}
		}
	}
	return input.String(), strings.TrimSpace(expected.String())
}

func runCase(bin, input, expected string) error {
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
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outLines := strings.Split(strings.TrimSpace(out.String()), "\n")
	expLines := strings.Split(strings.TrimSpace(expected), "\n")
	if len(outLines) != len(expLines) {
		return fmt.Errorf("expected %d lines got %d", len(expLines), len(outLines))
	}
	for i := range outLines {
		if strings.TrimSpace(outLines[i]) != strings.TrimSpace(expLines[i]) {
			return fmt.Errorf("line %d expected %s got %s", i+1, expLines[i], outLines[i])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
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
