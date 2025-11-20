package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func buildOracle() (string, error) {
	oracle := filepath.Join(os.TempDir(), "oracleF1851.bin")
	cmd := exec.Command("go", "build", "-o", oracle, "1851F.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("oracle build failed: %v\n%s", err, out)
	}
	return oracle, nil
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

func generate() string {
	const T = 100
	rng := rand.New(rand.NewSource(6))
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", T)
	for i := 0; i < T; i++ {
		n := rng.Intn(5) + 2
		k := rng.Intn(4) + 1
		fmt.Fprintf(&sb, "%d %d\n", n, k)
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", rng.Intn(1<<k))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	cand := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	inputStr := generate()
	expOutput, err := run(oracle, inputStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "oracle failed: %v\n", err)
		os.Exit(1)
	}
	gotOutput, err := run(cand, inputStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	inputScanner := bufio.NewScanner(strings.NewReader(inputStr))
    expScanner := bufio.NewScanner(strings.NewReader(expOutput))
    gotScanner := bufio.NewScanner(strings.NewReader(gotOutput))

    inputScanner.Scan() // Read T
    t, _ := strconv.Atoi(inputScanner.Text())

    for i := 1; i <= t; i++ {
        inputScanner.Scan()
        nk := strings.Fields(inputScanner.Text())
        n, _ := strconv.Atoi(nk[0])
        // k, _ := strconv.Atoi(nk[1])

        inputScanner.Scan()
        aStr := strings.Fields(inputScanner.Text())
        a := make([]int, n)
        for j := 0; j < n; j++ {
            a[j], _ = strconv.Atoi(aStr[j])
        }

        expScanner.Scan()
        expFields := strings.Fields(expScanner.Text())
        exp_i, _ := strconv.Atoi(expFields[0])
        exp_j, _ := strconv.Atoi(expFields[1])
        exp_x, _ := strconv.Atoi(expFields[2])

        gotScanner.Scan()
        gotFields := strings.Fields(gotScanner.Text())
        got_i, _ := strconv.Atoi(gotFields[0])
        got_j, _ := strconv.Atoi(gotFields[1])
        got_x, _ := strconv.Atoi(gotFields[2])

        exp_val := (a[exp_i-1] ^ exp_x) & (a[exp_j-1] ^ exp_x)
        got_val := (a[got_i-1] ^ got_x) & (a[got_j-1] ^ got_x)

        if got_val != exp_val {
            fmt.Printf("wrong answer on testcase %d\n", i)
            fmt.Printf("input:\n%s\n%s\n", strings.Join(nk, " "), strings.Join(aStr, " "))
            fmt.Printf("expected line: %s (value: %d)\n", expScanner.Text(), exp_val)
            fmt.Printf("got line: %s (value: %d)\n", gotScanner.Text(), got_val)
            os.Exit(1)
        }
    }

	fmt.Println("All tests passed")
}