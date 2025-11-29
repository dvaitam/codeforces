package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
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
	tests := make([]testC, 20)
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
		sort.Ints(queries)
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

	for i, tc := range tests {
		var input bytes.Buffer
		fmt.Fprintf(&input, "%d %d\n", tc.n, tc.m)
		for j, v := range tc.rooms {
			if j > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprint(&input, v)
		}
		input.WriteByte('\n')
		for j, v := range tc.queries {
			if j > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprint(&input, v)
		}
		input.WriteByte('\n')

		expected := solveC(tc)

		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewReader(input.Bytes())
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\noutput:\n%s\n", i+1, err, out.String())
			os.Exit(1)
		}

		scanner := bufio.NewScanner(bytes.NewReader(out.Bytes()))
		scanner.Split(bufio.ScanWords)
		for j, exp := range expected {
			if !scanner.Scan() {
				fmt.Fprintf(os.Stderr, "test %d: wrong output format (missing output)\n", i+1)
				os.Exit(1)
			}
			dorm, err := strconv.Atoi(scanner.Text())
			if err != nil {
				fmt.Fprintf(os.Stderr, "test %d: non-integer output\n", i+1)
				os.Exit(1)
			}
			if !scanner.Scan() {
				fmt.Fprintf(os.Stderr, "test %d: wrong output format (missing room)\n", i+1)
				os.Exit(1)
			}
			room, err := strconv.Atoi(scanner.Text())
			if err != nil {
				fmt.Fprintf(os.Stderr, "test %d: non-integer output\n", i+1)
				os.Exit(1)
			}
			if dorm != exp[0] || room != exp[1] {
				fmt.Fprintf(os.Stderr, "test %d: wrong answer on query %d. Expected %d %d, got %d %d\nInput:\n%s\n", i+1, j+1, exp[0], exp[1], dorm, room, input.String())
				os.Exit(1)
			}
		}
		if scanner.Scan() {
			fmt.Fprintf(os.Stderr, "test %d: extra output\n", i+1)
			os.Exit(1)
		}
	}
	fmt.Println("Accepted")
}
