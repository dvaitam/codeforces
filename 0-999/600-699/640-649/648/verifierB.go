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

const testCount = 160

func buildOracle() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "648B.go")
	tmp, err := os.CreateTemp("", "oracle648B")
	if err != nil {
		return "", err
	}
	path := tmp.Name()
	tmp.Close()
	os.Remove(path)
	cmd := exec.Command("go", "build", "-o", path, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build oracle: %v\n%s", err, out)
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

func genCase(r *rand.Rand) (string, []int) {
	n := 1 + r.Intn(200)
	sumTarget := 2 + r.Intn(200000-1)
	parts := make([]int, 0, 2*n)
	for i := 0; i < n; i++ {
		minVal := sumTarget - 100000
		if minVal < 1 {
			minVal = 1
		}
		maxVal := sumTarget - 1
		if maxVal > 100000 {
			maxVal = 100000
		}
		if minVal > maxVal {
			minVal, maxVal = 1, 100000
		}
		a := minVal + r.Intn(maxVal-minVal+1)
		b := sumTarget - a
		if b <= 0 || b > 100000 {
			if b <= 0 {
				b = 1
				a = sumTarget - b
			} else {
				b = 100000
				a = sumTarget - b
			}
		}
		parts = append(parts, a, b)
	}
	r.Shuffle(len(parts), func(i, j int) {
		parts[i], parts[j] = parts[j], parts[i]
	})
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < len(parts); i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", parts[i])
	}
	sb.WriteByte('\n')
	return sb.String(), parts
}

func buildFreq(parts []int) map[int]int {
	freq := make(map[int]int)
	for _, v := range parts {
		freq[v]++
	}
	return freq
}

func cloneFreq(src map[int]int) map[int]int {
	dst := make(map[int]int, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func validateOutput(out string, n int, freq map[int]int, target int) error {
	fields := strings.Fields(out)
	if len(fields) != 2*n {
		return fmt.Errorf("expected %d numbers, got %d", 2*n, len(fields))
	}
	freqCopy := cloneFreq(freq)
	for i := 0; i < 2*n; i += 2 {
		x, err := strconv.Atoi(fields[i])
		if err != nil {
			return fmt.Errorf("invalid integer: %v", err)
		}
		y, err := strconv.Atoi(fields[i+1])
		if err != nil {
			return fmt.Errorf("invalid integer: %v", err)
		}
		if freqCopy[x] == 0 {
			return fmt.Errorf("value %d not available", x)
		}
		freqCopy[x]--
		if freqCopy[y] == 0 {
			return fmt.Errorf("value %d not available", y)
		}
		freqCopy[y]--
		if x+y != target {
			return fmt.Errorf("pair %d+%d != %d", x, y, target)
		}
	}
	for _, v := range freqCopy {
		if v != 0 {
			return fmt.Errorf("unused values remain")
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
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
	for t := 0; t < testCount; t++ {
		input, parts := genCase(r)
		n := len(parts) / 2
		sum := 0
		for _, v := range parts {
			sum += v
		}
		target := sum / n
		freq := buildFreq(parts)
		expectStr, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\n", t+1, err)
			os.Exit(1)
		}
		if err := validateOutput(expectStr, n, freq, target); err != nil {
			fmt.Fprintf(os.Stderr, "oracle output invalid on test %d: %v\n", t+1, err)
			os.Exit(1)
		}
		gotStr, err := run(userBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", t+1, err)
			os.Exit(1)
		}
		if err := validateOutput(gotStr, n, freq, target); err != nil {
			fmt.Printf("test %d failed\ninput:\n%s\nerror: %v\n", t+1, input, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", testCount)
}
