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

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type testCaseB struct {
	input    string
	expected []int64
}

func computeB(w, h []int64) []int64 {
	n := len(w)
	sumW := int64(0)
	for _, v := range w {
		sumW += v
	}
	var maxH, secondH int64
	idxMax := 0
	for i, hi := range h {
		if hi > maxH {
			secondH = maxH
			maxH = hi
			idxMax = i
		} else if hi > secondH {
			secondH = hi
		}
	}
	res := make([]int64, n)
	for i := 0; i < n; i++ {
		curH := maxH
		if i == idxMax {
			curH = secondH
		}
		res[i] = (sumW - w[i]) * curH
	}
	return res
}

func generateCaseB() testCaseB {
	n := rand.Intn(9) + 2 // 2..10
	w := make([]int64, n)
	h := make([]int64, n)
	for i := 0; i < n; i++ {
		w[i] = int64(rand.Intn(1000) + 1)
		h[i] = int64(rand.Intn(1000) + 1)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", w[i], h[i]))
	}
	return testCaseB{input: sb.String(), expected: computeB(w, h)}
}

func parseInt64s(s string) ([]int64, error) {
	parts := strings.Fields(s)
	vals := make([]int64, len(parts))
	for i, p := range parts {
		v, err := strconv.ParseInt(p, 10, 64)
		if err != nil {
			return nil, err
		}
		vals[i] = v
	}
	return vals, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for i := 1; i <= 100; i++ {
		tc := generateCaseB()
		out, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, tc.input)
			os.Exit(1)
		}
		vals, err := parseInt64s(out)
		if err != nil || len(vals) != len(tc.expected) {
			fmt.Fprintf(os.Stderr, "case %d: invalid output\ninput:\n%s", i, tc.input)
			os.Exit(1)
		}
		for j, v := range vals {
			if v != tc.expected[j] {
				fmt.Fprintf(os.Stderr, "case %d: expected %v got %v\ninput:\n%s", i, tc.expected, vals, tc.input)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
