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
	"strings"
)

const testCount = 100

// ---- Embedded solver for 2147D ----

func solveAll(input string) string {
	in := bufio.NewReaderSize(strings.NewReader(input), 1<<20)
	var out bytes.Buffer
	w := bufio.NewWriterSize(&out, 1<<20)

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int64, n)
		var sum int64
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
			sum += a[i]
		}

		sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })

		oddCounts := make([]int64, 0)
		for i := 0; i < n; {
			j := i + 1
			for j < n && a[j] == a[i] {
				j++
			}
			if a[i]&1 == 1 {
				oddCounts = append(oddCounts, int64(j-i))
			}
			i = j
		}

		var diff int64
		for i := len(oddCounts) - 1; i >= 0; i-- {
			if oddCounts[i] >= diff {
				diff = oddCounts[i] - diff
			} else {
				diff = diff - oddCounts[i]
			}
		}

		alice := (sum + diff) / 2
		bob := sum - alice
		fmt.Fprintln(w, alice, bob)
	}
	w.Flush()
	return strings.TrimSpace(out.String())
}

// ---- Verifier harness ----

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(r *rand.Rand) string {
	t := r.Intn(5) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	for i := 0; i < t; i++ {
		n := r.Intn(100) + 1
		fmt.Fprintf(&sb, "%d\n", n)
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", r.Intn(200)+1)
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseOutput(out string, t int) ([][2]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != 2*t {
		return nil, fmt.Errorf("expected %d numbers, got %d", 2*t, len(fields))
	}
	res := make([][2]int64, t)
	for i := 0; i < t; i++ {
		a, err := strconv.ParseInt(fields[2*i], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer at case %d", i+1)
		}
		b, err := strconv.ParseInt(fields[2*i+1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer at case %d", i+1)
		}
		res[i] = [2]int64{a, b}
	}
	return res, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	userBin := os.Args[1]

	r := rand.New(rand.NewSource(1))
	for tcase := 0; tcase < testCount; tcase++ {
		input := genCase(r)

		expectStr := solveAll(input)

		gotStr, err := run(userBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", tcase+1, err)
			os.Exit(1)
		}
		reader := strings.NewReader(input)
		var t int
		fmt.Fscan(reader, &t)
		expectRes, err := parseOutput(expectStr, t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid oracle output on test %d: %v\n", tcase+1, err)
			os.Exit(1)
		}
		gotRes, err := parseOutput(gotStr, t)
		if err != nil {
			fmt.Printf("test %d failed\ninput:\n%s\nerror: %v\n", tcase+1, input, err)
			os.Exit(1)
		}
		for i := 0; i < t; i++ {
			if expectRes[i] != gotRes[i] {
				fmt.Printf("test %d case %d failed\ninput:\n%s\nexpected: %d %d\ngot: %d %d\n", tcase+1, i+1, input, expectRes[i][0], expectRes[i][1], gotRes[i][0], gotRes[i][1])
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", testCount)
}
