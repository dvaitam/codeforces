package main

import (
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
	"time"
)

func buildOracle() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "1427A.go")
	bin := filepath.Join(os.TempDir(), "oracle1427A.bin")
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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand) string {
	q := rng.Intn(4) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", q))
	for i := 0; i < q; i++ {
		n := rng.Intn(8) + 1
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(rng.Intn(11) - 5))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

// parseInput returns the array for each test case.
func parseInput(s string) [][]int {
	fields := strings.Fields(s)
	idx := 0
	t, _ := strconv.Atoi(fields[idx])
	idx++
	var cases [][]int
	for i := 0; i < t; i++ {
		n, _ := strconv.Atoi(fields[idx])
		idx++
		a := make([]int, n)
		for j := 0; j < n; j++ {
			a[j], _ = strconv.Atoi(fields[idx])
			idx++
		}
		cases = append(cases, a)
	}
	return cases
}

// parseOutput splits a multi-testcase output into per-case (verdict, arrayLine) pairs.
func parseOutput(s string) [][]string {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	var result [][]string
	i := 0
	for i < len(lines) {
		verdict := strings.TrimSpace(lines[i])
		i++
		if strings.ToUpper(verdict) == "YES" && i < len(lines) {
			result = append(result, []string{"YES", strings.TrimSpace(lines[i])})
			i++
		} else {
			result = append(result, []string{strings.ToUpper(verdict)})
		}
	}
	return result
}

func sortedCopy(a []int) []int {
	b := make([]int, len(a))
	copy(b, a)
	sort.Ints(b)
	return b
}

// check validates user output against oracle output and problem constraints.
// For NO: user must also say NO (oracle is authoritative).
// For YES: verify user's array is a valid rearrangement with all prefix sums nonzero.
func check(input, oracleOut, userOut string) error {
	cases := parseInput(input)
	oracleParsed := parseOutput(oracleOut)
	userParsed := parseOutput(userOut)

	if len(userParsed) != len(cases) {
		return fmt.Errorf("expected %d answers, got %d", len(cases), len(userParsed))
	}

	for i, a := range cases {
		oVerdict := oracleParsed[i][0]
		uCase := userParsed[i]
		uVerdict := uCase[0]

		if oVerdict == "NO" {
			if uVerdict != "NO" {
				return fmt.Errorf("test %d: oracle says NO but user says %s", i+1, uVerdict)
			}
			continue
		}

		// oracle says YES
		if uVerdict == "NO" {
			return fmt.Errorf("test %d: oracle says YES but user says NO", i+1)
		}

		// user says YES — validate the array
		if len(uCase) < 2 {
			return fmt.Errorf("test %d: user said YES but provided no array", i+1)
		}
		bFields := strings.Fields(uCase[1])
		if len(bFields) != len(a) {
			return fmt.Errorf("test %d: expected %d elements, got %d", i+1, len(a), len(bFields))
		}
		b := make([]int, len(a))
		for j, f := range bFields {
			v, err := strconv.Atoi(f)
			if err != nil {
				return fmt.Errorf("test %d: invalid element %q", i+1, f)
			}
			b[j] = v
		}
		// Check it's a rearrangement
		if sa, sb := sortedCopy(a), sortedCopy(b); fmt.Sprint(sa) != fmt.Sprint(sb) {
			return fmt.Errorf("test %d: array is not a rearrangement of input", i+1)
		}
		// Check all prefix sums nonzero
		sum := 0
		for j, v := range b {
			sum += v
			if sum == 0 {
				return fmt.Errorf("test %d: prefix sum is zero at position %d", i+1, j+1)
			}
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	user := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]string, 0, 100)
	for len(cases) < 100 {
		cases = append(cases, genCase(rng))
	}
	for i, in := range cases {
		oracleOut, err := run(oracle, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		userOut, err := run(user, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if err := check(in, oracleOut, userOut); err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%soracle:\n%s\ngot:\n%s\n", i+1, err, in, oracleOut, userOut)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
