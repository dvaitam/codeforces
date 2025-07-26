package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
)

type testC struct {
	n       int
	m       int
	rooms   []int
	queries []int
}

func genTestsC() []testC {
	rand.Seed(44)
	tests := make([]testC, 100)
	for i := range tests {
		n := rand.Intn(10) + 1
		m := rand.Intn(10) + 1
		rooms := make([]int, n)
		for j := range rooms {
			rooms[j] = rand.Intn(10) + 1
		}
		sum := 0
		for _, v := range rooms {
			sum += v
		}
		queries := make([]int, m)
		for j := range queries {
			queries[j] = rand.Intn(sum) + 1
		}
		tests[i] = testC{n: n, m: m, rooms: rooms, queries: queries}
	}
	return tests
}

func solveC(tc testC) [][2]int {
	pref := make([]int, tc.n)
	total := 0
	for i, v := range tc.rooms {
		total += v
		pref[i] = total
	}
	res := make([][2]int, tc.m)
	for i, q := range tc.queries {
		idx := 0
		for idx < tc.n && q > pref[idx] {
			idx++
		}
		prev := 0
		if idx > 0 {
			prev = pref[idx-1]
		}
		res[i] = [2]int{idx + 1, q - prev}
	}
	return res
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsC()

	var input bytes.Buffer
	fmt.Fprintln(&input, len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&input, "%d %d\n", tc.n, tc.m)
		for i, v := range tc.rooms {
			if i > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprint(&input, v)
		}
		input.WriteByte('\n')
		for i, v := range tc.queries {
			if i > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprint(&input, v)
		}
		input.WriteByte('\n')
	}

	expected := make([][][2]int, len(tests))
	for i, tc := range tests {
		expected[i] = solveC(tc)
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
		for j := range exp {
			if !scanner.Scan() {
				fmt.Fprintf(os.Stderr, "wrong output format on test %d\n", i+1)
				os.Exit(1)
			}
			dorm, err := strconv.Atoi(scanner.Text())
			if err != nil {
				fmt.Fprintf(os.Stderr, "non-integer output on test %d\n", i+1)
				os.Exit(1)
			}
			if !scanner.Scan() {
				fmt.Fprintf(os.Stderr, "wrong output format on test %d\n", i+1)
				os.Exit(1)
			}
			room, err := strconv.Atoi(scanner.Text())
			if err != nil {
				fmt.Fprintf(os.Stderr, "non-integer output on test %d\n", i+1)
				os.Exit(1)
			}
			if dorm != exp[j][0] || room != exp[j][1] {
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
