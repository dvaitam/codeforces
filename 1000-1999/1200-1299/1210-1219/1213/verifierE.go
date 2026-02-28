package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func buildRef() (string, error) {
	_, cur, _, _ := runtime.Caller(0)
	dir := filepath.Dir(cur)
	src := filepath.Join(dir, "1213E.go")
	bin := filepath.Join(os.TempDir(), "ref1213E.bin")
	cmd := exec.Command("go", "build", "-o", bin, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v\n%s", err, out)
	}
	return bin, nil
}

func run(bin string, input string) (string, error) {
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

func randPair(rng *rand.Rand) string {
	letters := []byte{'a', 'b', 'c'}
	b := []byte{letters[rng.Intn(3)], letters[rng.Intn(3)]}
	return string(b)
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(10) + 1
	s := randPair(rng)
	t := randPair(rng)
	return fmt.Sprintf("%d\n%s\n%s\n", n, s, t)
}

// parseInput extracts n, s, t from a test case input string.
func parseInput(input string) (n int, s, t string) {
	parts := strings.Fields(input)
	fmt.Sscan(parts[0], &n)
	s = parts[1]
	t = parts[2]
	return
}

// validate checks whether the candidate output is correct given the input.
// It uses the reference verdict only to decide whether "NO" is acceptable.
func validate(input, refOut, candOut string, testNum int) bool {
	n, s, t := parseInput(input)

	refLines := strings.SplitN(refOut, "\n", 2)
	refVerdict := strings.TrimSpace(refLines[0])

	candLines := strings.SplitN(candOut, "\n", 2)
	candVerdict := strings.TrimSpace(candLines[0])

	if candVerdict == "NO" {
		// Candidate claims no solution exists; accept only if reference also says NO.
		if refVerdict != "NO" {
			fmt.Printf("wrong answer on test %d\ninput:\n%sexpected:YES\ngot:NO\n",
				testNum, input)
			return false
		}
		return true
	}

	if candVerdict != "YES" {
		fmt.Printf("wrong answer on test %d\ninput:\n%sinvalid first line: %q\n",
			testNum, input, candVerdict)
		return false
	}

	// Reference says NO but candidate says YES — candidate is wrong.
	if refVerdict == "NO" {
		fmt.Printf("wrong answer on test %d\ninput:\n%sexpected:NO\ngot:YES\n",
			testNum, input)
		return false
	}

	// Candidate says YES — verify the string.
	if len(candLines) < 2 {
		fmt.Printf("wrong answer on test %d\ninput:\n%smissing result string after YES\n",
			testNum, input)
		return false
	}
	res := strings.TrimSpace(candLines[1])

	if len(res) != 3*n {
		fmt.Printf("wrong answer on test %d\ninput:\n%sresult length %d, expected %d\n",
			testNum, input, len(res), 3*n)
		return false
	}

	var ca, cb, cc int
	for _, ch := range res {
		switch ch {
		case 'a':
			ca++
		case 'b':
			cb++
		case 'c':
			cc++
		default:
			fmt.Printf("wrong answer on test %d\ninput:\n%sinvalid character %q in result\n",
				testNum, input, ch)
			return false
		}
	}
	if ca != n || cb != n || cc != n {
		fmt.Printf("wrong answer on test %d\ninput:\n%scharacter counts a=%d b=%d c=%d, expected %d each\n",
			testNum, input, ca, cb, cc, n)
		return false
	}

	if strings.Contains(res, s) {
		fmt.Printf("wrong answer on test %d\ninput:\n%sresult contains forbidden substring %q\n",
			testNum, input, s)
		return false
	}
	if strings.Contains(res, t) {
		fmt.Printf("wrong answer on test %d\ninput:\n%sresult contains forbidden substring %q\n",
			testNum, input, t)
		return false
	}

	return true
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}
	cand := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		input := genCase(rng)
		want, err := run(ref, input)
		if err != nil {
			fmt.Printf("reference failed on test %d: %v\n", i, err)
			os.Exit(1)
		}
		got, err := run(cand, input)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n", i, err)
			os.Exit(1)
		}
		if !validate(input, want, got, i) {
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
