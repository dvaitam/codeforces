package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Embedded reference solver for 488B.
const embeddedSolver488B = `package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	a := make([]int, n)
	minA := 501
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
		if a[i] < minA {
			minA = a[i]
		}
	}

	if n == 0 {
		fmt.Fprintln(out, "YES")
		fmt.Fprintln(out, 1)
		fmt.Fprintln(out, 1)
		fmt.Fprintln(out, 3)
		fmt.Fprintln(out, 3)
		return
	}

	for base := 1; base <= minA; base++ {
		for b := base; b <= 2*base; b++ {
			cand := [4]int{base, b, 4*base - b, 3 * base}
			used := [4]bool{}
			ok := true

			for _, v := range a {
				found := false
				for i := 0; i < 4; i++ {
					if !used[i] && cand[i] == v {
						used[i] = true
						found = true
						break
					}
				}
				if !found {
					ok = false
					break
				}
			}

			if ok {
				fmt.Fprintln(out, "YES")
				for i := 0; i < 4; i++ {
					if !used[i] {
						fmt.Fprintln(out, cand[i])
					}
				}
				return
			}
		}
	}

	fmt.Fprintln(out, "NO")
}
`

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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

func buildReference() (string, error) {
	tmpSrc, err := os.CreateTemp("", "488B-src-*.go")
	if err != nil {
		return "", err
	}
	srcPath := tmpSrc.Name()
	if _, err := tmpSrc.WriteString(embeddedSolver488B); err != nil {
		tmpSrc.Close()
		os.Remove(srcPath)
		return "", err
	}
	tmpSrc.Close()
	defer os.Remove(srcPath)

	tmp, err := os.CreateTemp("", "488B-ref-*")
	if err != nil {
		return "", err
	}
	binPath := tmp.Name()
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", binPath, srcPath)
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(binPath)
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return binPath, nil
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(5)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		v := rng.Intn(500) + 1
		sb.WriteString(fmt.Sprintf("%d\n", v))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)
	var cases []string
	// simple deterministic cases
	cases = append(cases, "0\n")
	cases = append(cases, "4\n1\n1\n3\n3\n")
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCase(rng))
	}
	for i, in := range cases {
		exp, err := run(refBin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if err := verify(in, out, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d wrong answer: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

func chk(b []int) bool {
	dif := (b[3] - b[0]) * 4
	mid := (b[1] + b[2]) * 2
	if mid != dif {
		return false
	}
	sum := 0
	for _, v := range b {
		sum += v
	}
	return sum == dif
}

func verify(input, output, expected string) error {
	tokensIn := strings.Fields(input)
	if len(tokensIn) == 0 {
		return fmt.Errorf("empty input")
	}
	n, err := strconv.Atoi(tokensIn[0])
	if err != nil {
		return fmt.Errorf("bad input: %v", err)
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i], err = strconv.Atoi(tokensIn[i+1])
		if err != nil {
			return fmt.Errorf("bad input number: %v", err)
		}
	}

	tokensOut := strings.Fields(output)
	if len(tokensOut) == 0 {
		return fmt.Errorf("empty output")
	}

	tokensExp := strings.Fields(expected)
	if len(tokensExp) == 0 {
		return fmt.Errorf("empty reference output")
	}

	if tokensExp[0] == "NO" {
		if tokensOut[0] != "NO" {
			return fmt.Errorf("expected NO but got %s", tokensOut[0])
		}
		if len(tokensOut) > 1 {
			return fmt.Errorf("unexpected extra output after NO")
		}
		return nil
	}

	if tokensOut[0] != "YES" {
		return fmt.Errorf("expected YES but got %s", tokensOut[0])
	}

	need := 4 - n
	if len(tokensOut)-1 != need {
		return fmt.Errorf("expected %d numbers after YES, got %d", need, len(tokensOut)-1)
	}

	b := make([]int, 0, 4)
	b = append(b, a...)
	for i := 1; i < len(tokensOut); i++ {
		v, err := strconv.Atoi(tokensOut[i])
		if err != nil {
			return fmt.Errorf("invalid number %q", tokensOut[i])
		}
		b = append(b, v)
	}
	if len(b) != 4 {
		return fmt.Errorf("total numbers != 4")
	}
	sort.Ints(b)
	if !chk(b) {
		return fmt.Errorf("output numbers invalid: %v", b)
	}
	return nil
}
