package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

const testCount = 120

func buildOracle() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "1431E.go")
	tmp, err := os.CreateTemp("", "oracle1431E")
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

func runProgram(bin, input string) (string, error) {
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

func genCase(r *rand.Rand) (string, [][]int, [][]int) {
	t := r.Intn(4) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	aCases := make([][]int, t)
	bCases := make([][]int, t)
	for i := 0; i < t; i++ {
		n := r.Intn(30) + 1
		if r.Intn(5) == 0 {
			n = r.Intn(80) + 1
		}
		fmt.Fprintf(&sb, "%d\n", n)
		a := make([]int, n)
		b := make([]int, n)
		cur := 0
		for j := 0; j < n; j++ {
			cur += r.Intn(5)
			a[j] = cur
		}
		cur = 0
		for j := 0; j < n; j++ {
			cur += r.Intn(5)
			b[j] = cur
		}
		aCases[i] = a
		bCases[i] = b
		for j, v := range a {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", v)
		}
		sb.WriteByte('\n')
		for j, v := range b {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", v)
		}
		sb.WriteByte('\n')
	}
	return sb.String(), aCases, bCases
}

func parseOutput(out string, t int, ns []int) ([][]int, error) {
	reader := bufio.NewReader(strings.NewReader(out))
	res := make([][]int, t)
	for i := 0; i < t; i++ {
		line, err := reader.ReadString('\n')
		if err != nil && len(line) == 0 {
			return nil, fmt.Errorf("missing output for case %d", i+1)
		}
		fields := strings.Fields(line)
		if len(fields) != ns[i] {
			return nil, fmt.Errorf("expected %d integers at case %d, got %d", ns[i], i+1, len(fields))
		}
		res[i] = make([]int, ns[i])
		used := make([]bool, ns[i])
		for j := 0; j < ns[i]; j++ {
			val, err := strconv.Atoi(fields[j])
			if err != nil {
				return nil, fmt.Errorf("invalid integer at case %d position %d", i+1, j+1)
			}
			if val < 1 || val > ns[i] {
				return nil, fmt.Errorf("index %d out of range at case %d position %d", val, i+1, j+1)
			}
			if used[val-1] {
				return nil, fmt.Errorf("duplicate index %d in case %d", val, i+1)
			}
			used[val-1] = true
			res[i][j] = val - 1
		}
	}
	if extra, err := reader.ReadString('\n'); err == nil && len(strings.TrimSpace(extra)) > 0 {
		return nil, fmt.Errorf("extra output detected")
	}
	return res, nil
}

func fairness(a, b []int, perm []int) int {
	best := 1 << 60
	for i := 0; i < len(a); i++ {
		diff := a[i] - b[perm[i]]
		if diff < 0 {
			diff = -diff
		}
		if diff < best {
			best = diff
		}
	}
	return best
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
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
		input, aCases, bCases := genCase(r)
		ns := make([]int, len(aCases))
		for i := range ns {
			ns[i] = len(aCases[i])
		}
		expectOut, err := runProgram(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\n", tc+1, err)
			os.Exit(1)
		}
		expectPerms, err := parseOutput(expectOut, len(aCases), ns)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle output invalid on test %d: %v\n", tc+1, err)
			os.Exit(1)
		}
		gotOut, err := runProgram(userBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", tc+1, err)
			os.Exit(1)
		}
		gotPerms, err := parseOutput(gotOut, len(aCases), ns)
		if err != nil {
			fmt.Printf("test %d failed\ninput:\n%s\nerror: %v\n", tc+1, input, err)
			os.Exit(1)
		}
		for i := 0; i < len(aCases); i++ {
			best := fairness(aCases[i], bCases[i], expectPerms[i])
			got := fairness(aCases[i], bCases[i], gotPerms[i])
			if got != best {
				fmt.Printf("test %d case %d failed\ninput:\n%s\nexpected fairness: %d\ngot: %d\n", tc+1, i+1, input, best, got)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", testCount)
}
