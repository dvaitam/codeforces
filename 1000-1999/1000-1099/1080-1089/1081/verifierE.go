package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func buildOracle() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "1081E.go")
	bin := filepath.Join(os.TempDir(), "oracle1081E.bin")
	cmd := exec.Command("go", "build", "-o", bin, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return bin, nil
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(r *rand.Rand) string {
	n := r.Intn(10)*2 + 2 // even between 2 and 20
	half := n / 2
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < half; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", r.Intn(100)+1)
	}
	sb.WriteByte('\n')
	return sb.String()
}

func isPerfectSquare(n int64) bool {
	if n < 0 {
		return false
	}
	root := int64(math.Round(math.Sqrt(float64(n))))
	for x := root - 2; x <= root+2; x++ {
		if x >= 0 && x*x == n {
			return true
		}
	}
	return false
}

func verifyOutput(in, out string) error {
	inScanner := bufio.NewScanner(strings.NewReader(in))
	inScanner.Split(bufio.ScanWords)
	if !inScanner.Scan() {
		return fmt.Errorf("bad input")
	}
	nStr := inScanner.Text()
	n, _ := strconv.Atoi(nStr)

	half := n / 2
	expectedEven := make([]int64, half)
	for i := 0; i < half; i++ {
		inScanner.Scan()
		expectedEven[i], _ = strconv.ParseInt(inScanner.Text(), 10, 64)
	}

	outScanner := bufio.NewScanner(strings.NewReader(out))
	outScanner.Split(bufio.ScanWords)
	if !outScanner.Scan() {
		return fmt.Errorf("empty output")
	}
	ans := outScanner.Text()
	if ans == "No" {
		return nil
	}
	if ans != "Yes" {
		return fmt.Errorf("expected Yes/No, got %s", ans)
	}

	var sum int64
	for i := 1; i <= n; i++ {
		if !outScanner.Scan() {
			return fmt.Errorf("not enough numbers")
		}
		val, err := strconv.ParseInt(outScanner.Text(), 10, 64)
		if err != nil {
			return err
		}
		if val <= 0 {
			return fmt.Errorf("number <= 0")
		}
		if val > 10000000000000 {
			return fmt.Errorf("number %d too large", val)
		}
		if i%2 == 0 {
			if val != expectedEven[(i/2)-1] {
				return fmt.Errorf("even position %d mismatch: got %d, expected %d", i, val, expectedEven[(i/2)-1])
			}
		}
		sum += val
		if !isPerfectSquare(sum) {
			return fmt.Errorf("sum at %d is %d, not a perfect square", i, sum)
		}
	}
	return nil
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
	cases := []string{
		"2\n1\n",
		"4\n1 1\n",
		"6\n2 4 6\n",
	}
	for i := 0; i < 97; i++ {
		cases = append(cases, genCase(r))
	}
	for idx, input := range cases {
		want, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		got, err := run(userBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		
		wantIsNo := strings.HasPrefix(want, "No")
		gotIsNo := strings.HasPrefix(got, "No")
		
		if wantIsNo != gotIsNo {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", idx+1, input, want, got)
			os.Exit(1)
		}
		
		if !gotIsNo {
			if err := verifyOutput(input, got); err != nil {
				fmt.Printf("test %d failed validation\ninput:\n%sgot: %s\nerror: %v\n", idx+1, input, got, err)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}