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
const maxFrequency = 24

func buildOracle() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "1267H.go")
	tmp, err := os.CreateTemp("", "oracle1267H")
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

func genCase(r *rand.Rand) (string, [][]int) {
	t := r.Intn(4) + 1
	inputs := make([][]int, t)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	for i := 0; i < t; i++ {
		n := r.Intn(20) + 1
		if r.Intn(5) == 0 {
			n = r.Intn(60) + 1
		}
		perm := r.Perm(n)
		for i := range perm {
			perm[i]++
		}
		fmt.Fprintf(&sb, "%d\n", n)
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", perm[j])
		}
		sb.WriteByte('\n')
		inputs[i] = make([]int, n)
		copy(inputs[i], perm)
	}
	return sb.String(), inputs
}

func parseFrequencies(output string, t int, ns []int) ([][]int, error) {
	reader := bufio.NewReader(strings.NewReader(output))
	freqs := make([][]int, t)
	for i := 0; i < t; i++ {
		line, err := reader.ReadString('\n')
		if err != nil && len(line) == 0 {
			return nil, fmt.Errorf("missing output for test case %d", i+1)
		}
		fields := strings.Fields(line)
		if len(fields) != ns[i] {
			return nil, fmt.Errorf("test case %d expected %d frequencies, got %d", i+1, ns[i], len(fields))
		}
		freqs[i] = make([]int, ns[i])
		for j := 0; j < ns[i]; j++ {
			val, err := strconv.Atoi(fields[j])
			if err != nil {
				return nil, fmt.Errorf("invalid integer in test case %d: %v", i+1, err)
			}
			if val < 1 || val > maxFrequency {
				return nil, fmt.Errorf("frequency %d out of range at test case %d position %d", val, i+1, j+1)
			}
			freqs[i][j] = val
		}
	}
	if extra, err := reader.ReadString('\n'); err == nil && len(strings.TrimSpace(extra)) > 0 {
		return nil, fmt.Errorf("extra output detected")
	}
	return freqs, nil
}

func validateCase(n int, perm, freq []int) error {
	if len(freq) != n {
		return fmt.Errorf("expected %d frequencies, got %d", n, len(freq))
	}
	order := make([]int, len(perm))
	copy(order, perm)
	// activation order: we switch on according to perm
	active := make([]bool, n)
	for day := 0; day < n; day++ {
		idx := order[day] - 1
		active[idx] = true
		l := 0
		for l < n {
			for l < n && !active[l] {
				l++
			}
			if l == n {
				break
			}
			r := l
			for r+1 < n && active[r+1] {
				r++
			}
			if err := checkSegment(freq[l : r+1]); err != nil {
				return fmt.Errorf("day %d segment [%d,%d]: %v", day+1, l+1, r+1, err)
			}
			l = r + 1
		}
	}
	return nil
}

func checkSegment(segment []int) error {
	if len(segment) == 0 {
		return nil
	}
	counts := make(map[int]int)
	for _, v := range segment {
		counts[v]++
	}
	uniqueFreqs := make([]int, 0, len(counts))
	for freq, cnt := range counts {
		if cnt == 1 {
			uniqueFreqs = append(uniqueFreqs, freq)
		}
	}
	if len(uniqueFreqs) == 0 {
		return fmt.Errorf("segment lacks unique frequency")
	}
	sort.Ints(uniqueFreqs)
	return nil
}

func compareOutputs(expect, got string, t int, ns []int) error {
	expectFreqs, err := parseFrequencies(expect, t, ns)
	if err != nil {
		return fmt.Errorf("oracle output invalid: %v", err)
	}
	gotFreqs, err := parseFrequencies(got, t, ns)
	if err != nil {
		return err
	}
	_ = expectFreqs // just to ensure oracle is well-formed
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierH.go /path/to/binary")
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
	for tcase := 0; tcase < testCount; tcase++ {
		input, perms := genCase(r)
		ns := make([]int, len(perms))
		for i, perm := range perms {
			ns[i] = len(perm)
			sort.Ints(perm)
		}
		expectOut, err := runProgram(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\n", tcase+1, err)
			os.Exit(1)
		}
		if err := compareOutputs(expectOut, expectOut, len(perms), ns); err != nil {
			fmt.Fprintf(os.Stderr, "oracle output invalid on test %d: %v\n", tcase+1, err)
			os.Exit(1)
		}
		gotOut, err := runProgram(userBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", tcase+1, err)
			os.Exit(1)
		}
		freqs, err := parseFrequencies(gotOut, len(perms), ns)
		if err != nil {
			fmt.Printf("test %d failed\ninput:\n%s\nerror: %v\n", tcase+1, input, err)
			os.Exit(1)
		}
		perm := make([]int, len(perms))
		copy(perm, perms[0])
		// re-read perms: we need activation order; revert sorting
		reader := strings.NewReader(input)
		scanner := bufio.NewScanner(reader)
		scanner.Scan() // t
		for i := 0; i < len(perms); i++ {
			scanner.Scan() // n
			scanner.Scan()
			fields := strings.Fields(scanner.Text())
			perms[i] = make([]int, len(fields))
			for j := range fields {
				val, _ := strconv.Atoi(fields[j])
				perms[i][j] = val
			}
		}
		for i := 0; i < len(perms); i++ {
			if err := validateCase(len(perms[i]), perms[i], freqs[i]); err != nil {
				fmt.Printf("test %d case %d failed\ninput:\n%s\nerror: %v\n", tcase+1, i+1, input, err)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", testCount)
}
