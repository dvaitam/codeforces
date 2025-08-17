package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func buildRef() (string, error) {
	ref := "./refC.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1187C.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	return ref, nil
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(8) + 2
	m := rng.Intn(8) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < m; i++ {
		t := rng.Intn(2)
		l := rng.Intn(n-1) + 1
		r := l + rng.Intn(n-l) + 1
		if r > n {
			r = n
		}
		if l == r {
			if r < n {
				r++
			} else {
				l--
			}
		}
		sb.WriteString(fmt.Sprintf("%d %d %d\n", t, l, r))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	cand := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := generateCase(rng)
		exp, err := run(ref, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on case %d: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		out, err := run(cand, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		n, _, t, l, r, err := parseInput(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse input on case %d: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		expAns, _, err := parseOutput(exp, n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on case %d: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		outAns, arr, err := parseOutput(out, n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		if expAns == "NO" {
			if outAns != "NO" {
				fmt.Fprintf(os.Stderr, "case %d failed: expected NO, got %s\ninput:%s", i+1, outAns, input)
				os.Exit(1)
			}
			continue
		}
		if outAns != "YES" {
			fmt.Fprintf(os.Stderr, "case %d failed: expected YES, got %s\ninput:%s", i+1, outAns, input)
			os.Exit(1)
		}
		if err := checkConstraints(arr, t, l, r); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

func parseInput(input string) (int, int, []int, []int, []int, error) {
	reader := strings.NewReader(input)
	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return 0, 0, nil, nil, nil, err
	}
	t := make([]int, m)
	l := make([]int, m)
	r := make([]int, m)
	for i := 0; i < m; i++ {
		if _, err := fmt.Fscan(reader, &t[i], &l[i], &r[i]); err != nil {
			return 0, 0, nil, nil, nil, err
		}
		l[i]--
		r[i]--
	}
	return n, m, t, l, r, nil
}

func parseOutput(out string, n int) (string, []int, error) {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return "", nil, fmt.Errorf("empty output")
	}
	res := strings.ToUpper(fields[0])
	if res != "YES" && res != "NO" {
		return "", nil, fmt.Errorf("first token must be YES or NO")
	}
	if res == "NO" {
		return res, nil, nil
	}
	if len(fields) != n+1 {
		return "", nil, fmt.Errorf("expected %d numbers, got %d", n, len(fields)-1)
	}
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		v, err := strconv.Atoi(fields[i+1])
		if err != nil {
			return "", nil, fmt.Errorf("invalid integer: %v", err)
		}
		if v < 1 || v > 1_000_000_000 {
			return "", nil, fmt.Errorf("value out of range")
		}
		arr[i] = v
	}
	return res, arr, nil
}

func checkConstraints(a []int, t, l, r []int) error {
	for i := 0; i < len(t); i++ {
		if t[i] == 1 {
			for j := l[i] + 1; j <= r[i]; j++ {
				if a[j] < a[j-1] {
					return fmt.Errorf("segment %d expected sorted", i+1)
				}
			}
		} else {
			ok := false
			for j := l[i] + 1; j <= r[i]; j++ {
				if a[j] < a[j-1] {
					ok = true
					break
				}
			}
			if !ok {
				return fmt.Errorf("segment %d expected not sorted", i+1)
			}
		}
	}
	return nil
}
