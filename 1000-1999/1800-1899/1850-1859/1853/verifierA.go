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

func expectedA(arr []int) int {
	n := len(arr)
	for i := 0; i < n-1; i++ {
		if arr[i] > arr[i+1] {
			return 0
		}
	}
	minDiff := arr[1] - arr[0]
	for i := 1; i < n-1; i++ {
		d := arr[i+1] - arr[i]
		if d < minDiff {
			minDiff = d
		}
	}
	return minDiff/2 + 1
}

func genCase(rng *rand.Rand) []int {
	n := rng.Intn(8) + 2 // 2..9
	arr := make([]int, n)
	if rng.Intn(2) == 0 {
		// generate array that is not sorted
		v := rng.Intn(1000)
		for i := 0; i < n; i++ {
			v += rng.Intn(20)
			arr[i] = v
		}
		idx := rng.Intn(n - 1)
		arr[idx] = arr[idx+1] + rng.Intn(5) + 1
	} else {
		// generate non-decreasing array
		v := rng.Intn(1000)
		for i := 0; i < n; i++ {
			v += rng.Intn(20)
			arr[i] = v
		}
	}
	return arr
}

func runBinary(bin, input string) (string, string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var outBuf bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	err := cmd.Run()
	return outBuf.String(), errBuf.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	const tests = 100
	cases := make([][]int, tests)
	expected := make([]int, tests)
	for i := 0; i < tests; i++ {
		arr := genCase(rng)
		cases[i] = arr
		expected[i] = expectedA(arr)
	}

	var input strings.Builder
	fmt.Fprintf(&input, "%d\n", tests)
	for _, arr := range cases {
		fmt.Fprintf(&input, "%d\n", len(arr))
		for j, v := range arr {
			if j > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(v))
		}
		input.WriteByte('\n')
	}

	out, errOut, err := runBinary(bin, input.String())
	if err != nil {
		fmt.Fprintf(os.Stderr, "runtime error: %v\n%s", err, errOut)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(strings.NewReader(out))
	scanner.Split(bufio.ScanWords)
	for i := 0; i < tests; i++ {
		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "missing output for test %d\n", i+1)
			os.Exit(1)
		}
		val, err := strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid integer on test %d\n", i+1)
			os.Exit(1)
		}
		if val != expected[i] {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: expected %d got %d\n", i+1, expected[i], val)
			os.Exit(1)
		}
	}
	if scanner.Scan() {
		fmt.Fprintln(os.Stderr, "extra output")
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
