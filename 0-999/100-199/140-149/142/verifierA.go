package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strings"
)

func divisors(n int64) []int64 {
	var ds []int64
	lim := int64(math.Sqrt(float64(n)))
	for i := int64(1); i <= lim; i++ {
		if n%i == 0 {
			ds = append(ds, i)
			if j := n / i; j != i {
				ds = append(ds, j)
			}
		}
	}
	return ds
}

func solve(n int64) (int64, int64) {
	const inf64 int64 = 1<<63 - 1
	minStolen := inf64
	maxStolen := int64(0)
	for _, x := range divisors(n) {
		m := n / x
		for _, y := range divisors(m) {
			z := m / y
			s := 2*x*y + 2*x*z + 4*x + y*z + 2*y + 2*z + 4
			if s < minStolen {
				minStolen = s
			}
			if s > maxStolen {
				maxStolen = s
			}
		}
	}
	return minStolen, maxStolen
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesA.txt")
	if err != nil {
		fmt.Println("could not open testcasesA.txt:", err)
		os.Exit(1)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		var n int64
		if _, err := fmt.Sscan(line, &n); err != nil {
			fmt.Printf("bad test case on line %d\n", idx)
			os.Exit(1)
		}
		expMin, expMax := solve(n)
		input := fmt.Sprintf("%d\n", n)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		cmd.Stdout = &out
		var errBuf bytes.Buffer
		cmd.Stderr = &errBuf
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, errBuf.String())
			os.Exit(1)
		}
		var gotMin, gotMax int64
		if _, err := fmt.Fscan(bytes.NewReader(out.Bytes()), &gotMin, &gotMax); err != nil {
			fmt.Printf("test %d: failed to parse output %q\n", idx, out.String())
			os.Exit(1)
		}
		if gotMin != expMin || gotMax != expMax {
			fmt.Printf("test %d failed: expected %d %d got %d %d\n", idx, expMin, expMax, gotMin, gotMax)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
