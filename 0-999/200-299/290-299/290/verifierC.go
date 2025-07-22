package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func solve(arr []float64) float64 {
	sum := 0.0
	best := arr[0]
	for i, v := range arr {
		sum += v
		avg := sum / float64(i+1)
		if avg > best {
			best = avg
		}
	}
	return best
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		n := rand.Intn(10) + 1
		vals := make([]float64, n)
		for j := 0; j < n; j++ {
			vals[j] = float64(rand.Intn(100))
		}
		expected := solve(vals)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", n)
		for j, v := range vals {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", int(v))
		}
		sb.WriteByte('\n')

		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewBufferString(sb.String())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\noutput:\n%s\n", i+1, err, string(out))
			os.Exit(1)
		}
		outStr := strings.TrimSpace(string(out))
		got, err2 := strconv.ParseFloat(outStr, 64)
		if err2 != nil || mathAbs(got-expected) > 1e-4 {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d expected %.5f got %s\n", i+1, expected, outStr)
			os.Exit(1)
		}
	}
	fmt.Println("Accepted")
}

func mathAbs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
