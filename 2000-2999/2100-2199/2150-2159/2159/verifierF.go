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

func buildOracle() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "2159F.go")
	tmp, err := os.CreateTemp("", "oracle2159F")
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

func genPath(n int, r *rand.Rand) string {
	down := n - 1
	right := n - 1
	path := make([]byte, 0, 2*n-2)
	path = append(path, 'D')
	down--
	for down > 0 || right > 0 {
		if down == 0 {
			path = append(path, 'R')
			right--
		} else if right == 0 {
			path = append(path, 'D')
			down--
		} else if r.Intn(2) == 0 {
			path = append(path, 'D')
			down--
		} else {
			path = append(path, 'R')
			right--
		}
	}
	return string(path)
}

func genInput(r *rand.Rand) string {
	t := r.Intn(3) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	totalN := 0
	for i := 0; i < t; i++ {
		remaining := 50 - totalN
		maxN := 20
		if remaining < maxN {
			maxN = remaining
		}
		if maxN < 2 {
			maxN = 2
		}
		n := r.Intn(maxN-1) + 2
		totalN += n
		m := r.Intn(n*(2*n-1)) + 1
		fmt.Fprintf(&sb, "%d %d\n", n, m)
		perm := make([]int, n*n)
		for j := 0; j < n*n; j++ {
			perm[j] = j + 1
		}
		r.Shuffle(len(perm), func(a, b int) {
			perm[a], perm[b] = perm[b], perm[a]
		})
		for row := 0; row < n; row++ {
			for col := 0; col < n; col++ {
				fmt.Fprintf(&sb, "%d ", perm[row*n+col])
			}
			sb.WriteByte('\n')
		}
		for l := 1; l <= n; l++ {
			path := genPath(n, r)
			fmt.Fprintln(&sb, path)
		}
	}
	return sb.String()
}

func parseOutputs(out string, mList []int) ([][]int, error) {
	tokens := strings.Fields(out)
	total := 0
	for _, m := range mList {
		total += m
	}
	if len(tokens) != total {
		return nil, fmt.Errorf("expected %d numbers, got %d", total, len(tokens))
	}
	res := make([][]int, len(mList))
	idx := 0
	for i, m := range mList {
		res[i] = make([]int, m)
		for j := 0; j < m; j++ {
			val, err := strconv.Atoi(tokens[idx])
			if err != nil {
				return nil, fmt.Errorf("invalid integer at test %d position %d", i+1, j+1)
			}
			res[i][j] = val
			idx++
		}
		for j := 1; j < m; j++ {
			if res[i][j] < res[i][j-1] {
				return nil, fmt.Errorf("output for test %d not non-decreasing", i+1)
			}
		}
	}
	return res, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
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
		input := genInput(r)
		expectStr, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\n", tc+1, err)
			os.Exit(1)
		}
		gotStr, err := run(userBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", tc+1, err)
			os.Exit(1)
		}
		reader := strings.NewReader(input)
		var t int
		fmt.Fscan(reader, &t)
		mList := make([]int, t)
		for i := 0; i < t; i++ {
			var n, m int
			fmt.Fscan(reader, &n, &m)
			mList[i] = m
			for row := 0; row < n; row++ {
				for col := 0; col < n; col++ {
					var tmp int
					fmt.Fscan(reader, &tmp)
				}
			}
			for l := 0; l < n; l++ {
				var path string
				fmt.Fscan(reader, &path)
			}
		}
		expectVals, err := parseOutputs(expectStr, mList)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid oracle output on test %d: %v\n", tc+1, err)
			os.Exit(1)
		}
		gotVals, err := parseOutputs(gotStr, mList)
		if err != nil {
			fmt.Printf("test %d failed\ninput:\n%s\nerror: %v\n", tc+1, input, err)
			os.Exit(1)
		}
		for i := 0; i < len(mList); i++ {
			if len(expectVals[i]) != len(gotVals[i]) {
				fmt.Printf("test %d case %d failed\ninput:\n%s\nexpected length %d got %d\n", tc+1, i+1, input, len(expectVals[i]), len(gotVals[i]))
				os.Exit(1)
			}
			for j := 0; j < len(expectVals[i]); j++ {
				if expectVals[i][j] != gotVals[i][j] {
					fmt.Printf("test %d case %d failed\ninput:\n%s\nexpected: %v\ngot: %v\n", tc+1, i+1, input, expectVals[i], gotVals[i])
					os.Exit(1)
				}
			}
		}
	}
	fmt.Printf("All %d tests passed\n", testCount)
}
