package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
)

type testD struct {
	arr []uint64
}

func genTestsD() []testD {
	rand.Seed(122004)
	tests := make([]testD, 100)
	for i := range tests {
		n := rand.Intn(10) + 1
		arr := make([]uint64, n)
		for j := range arr {
			arr[j] = uint64(rand.Int63n(1<<30) + 1)
		}
		tests[i] = testD{arr: arr}
	}
	return tests
}

func solveD(tc testD) []uint64 {
	cnt := make([]int, 64)
	for _, v := range tc.arr {
		tz := bits.TrailingZeros64(v)
		if tz < len(cnt) {
			cnt[tz]++
		}
	}
	maxCnt := 0
	tu := 0
	for i, c := range cnt {
		if c > maxCnt {
			maxCnt = c
			tu = i
		}
	}
	res := []uint64{}
	for _, v := range tc.arr {
		if bits.TrailingZeros64(v) != tu {
			res = append(res, v)
		}
	}
	return res
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsD()

	var input bytes.Buffer
	fmt.Fprintln(&input, len(tests))
	for _, tc := range tests {
		fmt.Fprintln(&input, len(tc.arr))
		for i, v := range tc.arr {
			if i > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprint(&input, v)
		}
		input.WriteByte('\n')
	}

	expected := make([][]uint64, len(tests))
	for i, tc := range tests {
		expected[i] = solveD(tc)
	}

	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "runtime error: %v\noutput:\n%s\n", err, out.String())
		os.Exit(1)
	}

	scanner := bufio.NewScanner(bytes.NewReader(out.Bytes()))
	scanner.Split(bufio.ScanWords)
	for i, exp := range expected {
		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "wrong output format on test %d\n", i+1)
			os.Exit(1)
		}
		k, err := strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Fprintf(os.Stderr, "non-integer output on test %d\n", i+1)
			os.Exit(1)
		}
		if k != len(exp) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: expected length %d got %d\n", i+1, len(exp), k)
			os.Exit(1)
		}
		for j := 0; j < k; j++ {
			if !scanner.Scan() {
				fmt.Fprintf(os.Stderr, "wrong output format on test %d\n", i+1)
				os.Exit(1)
			}
			val, err := strconv.ParseUint(scanner.Text(), 10, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "non-integer output on test %d\n", i+1)
				os.Exit(1)
			}
			if val != exp[j] {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d\n", i+1)
				os.Exit(1)
			}
		}
	}
	if scanner.Scan() {
		fmt.Fprintln(os.Stderr, "extra output")
		os.Exit(1)
	}
	fmt.Println("Accepted")
}
