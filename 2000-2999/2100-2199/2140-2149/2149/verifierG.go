package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

const testCount = 60
const maxN = 400
const maxQ = 400

func buildOracle() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "2149G.go")
	tmp, err := os.CreateTemp("", "oracle2149G")
	if err != nil {
		return "", err
	}
	path := tmp.Name()
	tmp.Close()
	os.Remove(path)
	cmd := exec.Command("go", "build", "-o", path, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build oracle: %v\n%s", err, string(out))
	}
	return path, nil
}

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

func genCase(r *rand.Rand) (string, [][]int, [][][2]int) {
	t := r.Intn(4) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	arrays := make([][]int, t)
	queries := make([][][2]int, t)
	for i := 0; i < t; i++ {
		n := r.Intn(maxN) + 1
		q := r.Intn(maxQ) + 1
		fmt.Fprintf(&sb, "%d %d\n", n, q)
		arr := make([]int, n)
		valRange := r.Intn(50) + 1
		for j := 0; j < n; j++ {
			arr[j] = r.Intn(valRange) + 1
		}
		arrays[i] = arr
		for j, v := range arr {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", v)
		}
		sb.WriteByte('\n')
		qs := make([][2]int, q)
		for j := 0; j < q; j++ {
			l := r.Intn(n) + 1
			rgt := r.Intn(n-l+1) + l
			fmt.Fprintf(&sb, "%d %d\n", l, rgt)
			qs[j] = [2]int{l, rgt}
		}
		queries[i] = qs
	}
	return sb.String(), arrays, queries
}

func parseOutput(out string, totalQ int) ([][]int, error) {
	lines := strings.Split(out, "\n")
	res := make([][]int, 0, totalQ)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if line == "-1" {
			res = append(res, nil)
			continue
		}
		fields := strings.Fields(line)
		vals := make([]int, len(fields))
		for i, f := range fields {
			v, err := strconv.Atoi(f)
			if err != nil {
				return nil, fmt.Errorf("invalid integer in line '%s'", line)
			}
			vals[i] = v
		}
		res = append(res, vals)
	}
	if len(res) != totalQ {
		return nil, fmt.Errorf("expected %d output lines, got %d", totalQ, len(res))
	}
	return res, nil
}

func freqMap(arr []int, l, r int) map[int]int {
	freq := make(map[int]int)
	for i := l - 1; i < r; i++ {
		freq[arr[i]]++
	}
	return freq
}

func checkAnswer(arr []int, q [][2]int, outputs [][]int) error {
	idx := 0
	for _, query := range q {
		l, r := query[0], query[1]
		length := r - l + 1
		threshold := length / 3
		freq := freqMap(arr, l, r)
		line := outputs[idx]
		if line == nil {
			for _, count := range freq {
				if count > threshold {
					return fmt.Errorf("missing value > threshold for query [%d,%d]", l, r)
				}
			}
		} else {
			seen := make(map[int]bool)
			for _, val := range line {
				if seen[val] {
					return fmt.Errorf("duplicate value %d in query [%d,%d]", val, l, r)
				}
				seen[val] = true
				if freq[val] <= threshold {
					return fmt.Errorf("value %d does not exceed threshold in query [%d,%d]", val, l, r)
				}
			}
			for val, count := range freq {
				if count > threshold && !seen[val] {
					return fmt.Errorf("missing value %d in query [%d,%d]", val, l, r)
				}
			}
		}
		idx++
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	userBin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)
	r := rand.New(rand.NewSource(1))
	for tc := 0; tc < testCount; tc++ {
		input, arrays, queries := genCase(r)
		_, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\n", tc+1, err)
			os.Exit(1)
		}
		reader := strings.NewReader(input)
		var t int
		fmt.Fscan(reader, &t)
		totalQ := 0
		for i := 0; i < t; i++ {
			var n, q int
			fmt.Fscan(reader, &n, &q)
			for j := 0; j < n; j++ {
				var tmp int
				fmt.Fscan(reader, &tmp)
			}
			for j := 0; j < q; j++ {
				var l, r int
				fmt.Fscan(reader, &l, &r)
			}
			totalQ += q
		}
		gotStr, err := run(userBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", tc+1, err)
			os.Exit(1)
		}
		gotLines, err := parseOutput(gotStr, totalQ)
		if err != nil {
			fmt.Printf("test %d failed\ninput:\n%s\nerror: %v\n", tc+1, input, err)
			os.Exit(1)
		}
		idx := 0
		for i := 0; i < len(arrays); i++ {
			qCount := len(queries[i])
			if err := checkAnswer(arrays[i], queries[i], gotLines[idx:idx+qCount]); err != nil {
				fmt.Printf("test %d case %d failed\ninput:\n%s\nerror: %v\n", tc+1, i+1, input, err)
				os.Exit(1)
			}
			idx += qCount
		}
	}
	fmt.Printf("All %d tests passed\n", testCount)
}
