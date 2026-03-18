package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type testCase struct {
	n       int
	arr     []int
	queries [][2]int
}

func genCase(rng *rand.Rand) testCase {
	n := rng.Intn(8) + 2 // 2..9
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = i + 1
	}
	rng.Shuffle(n, func(i, j int) { arr[i], arr[j] = arr[j], arr[i] })
	q := rng.Intn(5) + 1
	queries := make([][2]int, q)
	for i := 0; i < q; i++ {
		l := rng.Intn(n-1) + 1 // 1..n-1
		r := l + rng.Intn(n-l) + 1 // l+1..n
		queries[i] = [2]int{l, r}
	}
	return testCase{n: n, arr: arr, queries: queries}
}

func solve(tc testCase) string {
	var sb strings.Builder
	for idx, qr := range tc.queries {
		l := qr[0] - 1
		r := qr[1] - 1
		good := false
		for i := l; i < r && !good; i++ {
			maxLeft := tc.arr[l]
			for j := l; j <= i; j++ {
				if tc.arr[j] > maxLeft {
					maxLeft = tc.arr[j]
				}
			}
			minRight := tc.arr[i+1]
			for j := i + 1; j <= r; j++ {
				if tc.arr[j] < minRight {
					minRight = tc.arr[j]
				}
			}
			if maxLeft < minRight {
				good = true
			}
		}
		if good {
			sb.WriteString("Yes")
		} else {
			sb.WriteString("No")
		}
		if idx+1 < len(tc.queries) {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(tc.n))
	sb.WriteByte('\n')
	for i, v := range tc.arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	sb.WriteString(strconv.Itoa(len(tc.queries)))
	sb.WriteByte('\n')
	for _, q := range tc.queries {
		sb.WriteString(strconv.Itoa(q[0]))
		sb.WriteByte(' ')
		sb.WriteString(strconv.Itoa(q[1]))
		sb.WriteByte('\n')
	}
	return sb.String()
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
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierD /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	numCases := 200

	for idx := 0; idx < numCases; idx++ {
		tc := genCase(rng)
		input := buildInput(tc)
		expected := solve(tc)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d: %v\n", idx+1, err)
			os.Exit(1)
		}

		// Normalize comparison: compare line by line, case-insensitive
		gotScanner := bufio.NewScanner(strings.NewReader(got))
		expScanner := bufio.NewScanner(strings.NewReader(expected))
		line := 0
		for expScanner.Scan() {
			line++
			if !gotScanner.Scan() {
				fmt.Printf("case %d failed: expected more output at line %d\n", idx+1, line)
				os.Exit(1)
			}
			if !strings.EqualFold(strings.TrimSpace(gotScanner.Text()), strings.TrimSpace(expScanner.Text())) {
				fmt.Printf("case %d failed: line %d expected %s got %s\n", idx+1, line, expScanner.Text(), gotScanner.Text())
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", numCases)
}
