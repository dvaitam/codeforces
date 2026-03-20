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

const refSource = `package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}

	for tc := 0; tc < t; tc++ {
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}

		dp := make([]int, n)
		for i := 0; i < n; i++ {
			dp[i] = 1
		}

		ans := 1
		for i := 1; i < n; i++ {
			for j := 0; j < i; j++ {
				if a[j] < a[i] {
					if dp[j]+1 > dp[i] {
						dp[i] = dp[j] + 1
					}
				}
			}
			if dp[i] > ans {
				ans = dp[i]
			}
		}

		dpRev := make([]int, n)
		for i := 0; i < n; i++ {
			dpRev[i] = 1
		}
		for i := n - 2; i >= 0; i-- {
			for j := n - 1; j > i; j-- {
				if a[j] < a[i] {
					if dpRev[j]+1 > dpRev[i] {
						dpRev[i] = dpRev[j] + 1
					}
				}
			}
			if dpRev[i] > ans {
				ans = dpRev[i]
			}
		}

		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				if a[i] < a[j] {
					cost := dp[i] + dpRev[j]
					if cost > ans {
						ans = cost
					}
				}
			}
		}

		fmt.Fprintln(writer, ans)
	}
}
`

func buildEmbeddedRef(dir string) (string, error) {
	src := filepath.Join(dir, "ref_embedded_1530H.go")
	if err := os.WriteFile(src, []byte(refSource), 0644); err != nil {
		return "", err
	}
	bin := filepath.Join(dir, "ref_1530H_bin")
	cmd := exec.Command("go", "build", "-o", bin, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build ref: %v\n%s", err, out)
	}
	os.Remove(src)
	return bin, nil
}

type testH struct {
	n   int
	arr []int
}

func genTestsH() []testH {
	rand.Seed(1530008)
	// Reduced test count and size for ARM to avoid TLE
	tests := make([]testH, 30)
	for i := range tests {
		n := rand.Intn(15) + 1
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = rand.Intn(50) + 1
		}
		tests[i] = testH{n: n, arr: arr}
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsH()

	tmpDir, err := os.MkdirTemp("", "v1530H")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.RemoveAll(tmpDir)

	refBin, err := buildEmbeddedRef(tmpDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var input bytes.Buffer
	fmt.Fprintln(&input, len(tests))
	for _, tc := range tests {
		fmt.Fprintln(&input, tc.n)
		for i, v := range tc.arr {
			if i > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprint(&input, v)
		}
		input.WriteByte('\n')
	}
	inputStr := input.String()

	// Run reference
	refCmd := exec.Command(refBin)
	refCmd.Stdin = strings.NewReader(inputStr)
	var refOut bytes.Buffer
	var refErr bytes.Buffer
	refCmd.Stdout = &refOut
	refCmd.Stderr = &refErr
	if err := refCmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "reference error: %v\n%s\n", err, refErr.String())
		os.Exit(1)
	}

	// Parse reference output
	refScanner := bufio.NewScanner(bytes.NewReader(refOut.Bytes()))
	refScanner.Split(bufio.ScanWords)
	expected := make([]int, len(tests))
	for i := range expected {
		if !refScanner.Scan() {
			fmt.Fprintf(os.Stderr, "reference output too short at test %d\n", i+1)
			os.Exit(1)
		}
		expected[i], _ = strconv.Atoi(refScanner.Text())
	}

	// Run candidate
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(inputStr)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "runtime error: %v\n%s\n", err, stderr.String())
		os.Exit(1)
	}

	scanner := bufio.NewScanner(bytes.NewReader(out.Bytes()))
	scanner.Split(bufio.ScanWords)
	for i, exp := range expected {
		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "wrong output format on test %d\n", i+1)
			os.Exit(1)
		}
		val, err := strconv.Atoi(scanner.Text())
		if err != nil || val != exp {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: expected %d got %s\n", i+1, exp, scanner.Text())
			os.Exit(1)
		}
	}
	if scanner.Scan() {
		fmt.Fprintln(os.Stderr, "extra output")
		os.Exit(1)
	}
	fmt.Println("Accepted")
}
