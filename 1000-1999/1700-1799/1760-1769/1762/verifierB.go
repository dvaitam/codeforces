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

type testB struct {
	n   int
	arr []int64
}

func genTestB(rng *rand.Rand) testB {
	n := rng.Intn(10) + 1
	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Int63n(1000) + 1
	}
	return testB{n: n, arr: arr}
}

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
	err := cmd.Run()
	if err != nil {
		return out.String() + errBuf.String(), err
	}
	return strings.TrimSpace(out.String()), nil
}

func isGood(arr []int64) bool {
	n := len(arr)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			a := arr[i]
			b := arr[j]
			if a == 0 || b == 0 {
				return false
			}
			x := a
			y := b
			if x < y {
				x, y = y, x
			}
			if x%y != 0 {
				return false
			}
		}
	}
	return true
}

func verifyOutput(out string, tc testB) error {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return fmt.Errorf("empty output")
	}
	p, err := strconv.Atoi(fields[0])
	if err != nil || p < 0 || p > tc.n {
		return fmt.Errorf("invalid p")
	}
	if len(fields) != 1+2*p {
		return fmt.Errorf("expected %d numbers, got %d", 1+2*p, len(fields))
	}
	arr := make([]int64, len(tc.arr))
	copy(arr, tc.arr)
	idx := 1
	for op := 0; op < p; op++ {
		iVal, err1 := strconv.Atoi(fields[idx])
		xVal, err2 := strconv.ParseInt(fields[idx+1], 10, 64)
		if err1 != nil || err2 != nil || iVal < 1 || iVal > tc.n || xVal < 0 {
			return fmt.Errorf("invalid operation")
		}
		idx += 2
		arr[iVal-1] += xVal
		if arr[iVal-1] > 1e18 {
			return fmt.Errorf("value exceeds limit")
		}
	}
	if !isGood(arr) {
		return fmt.Errorf("final array not good")
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		tc := genTestB(rng)
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		for i, v := range tc.arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		input := sb.String()
		out, err := run(bin, input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n%s", t+1, err, out)
			os.Exit(1)
		}
		if err := verifyOutput(out, tc); err != nil {
			fmt.Printf("test %d failed: %v\ninput:\n%soutput:\n%s\n", t+1, err, input, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", 100)
}
