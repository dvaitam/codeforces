package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const mod = 998244353

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierJ.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	inBytes, err := readAllStdin()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to read stdin:", err)
		os.Exit(1)
	}

	// Two modes:
	// 1) stdin provided: verify exactly that test case.
	// 2) empty stdin: run built-in randomized tests.
	if len(strings.Fields(string(inBytes))) == 0 {
		if err := runRandomTests(candidate); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println("Accepted")
		return
	}

	if err := verifySingleCase(candidate, inBytes); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("Accepted")
}

func verifySingleCase(candidate string, inBytes []byte) error {
	n, a, err := parseInput(inBytes)
	if err != nil {
		return fmt.Errorf("invalid input: %v", err)
	}
	want := solveExpected(n, a)
	candOut, err := runProgram(candidate, inBytes)
	if err != nil {
		return fmt.Errorf("candidate runtime error: %v", err)
	}
	if err := compareAnswer(want, candOut); err != nil {
		return fmt.Errorf("%v\nexpected:\n%d\ncandidate output:\n%s", err, want, candOut)
	}
	return nil
}

func readAllStdin() ([]byte, error) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(os.Stdin)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func parseInput(input []byte) (int, []int, error) {
	fields := strings.Fields(string(input))
	if len(fields) < 1 {
		return 0, nil, fmt.Errorf("missing n")
	}
	n, err := strconv.Atoi(fields[0])
	if err != nil || n < 2 {
		return 0, nil, fmt.Errorf("invalid n")
	}
	if len(fields) != n+1 {
		return 0, nil, fmt.Errorf("expected %d permutation values, got %d", n, len(fields)-1)
	}
	a := make([]int, n+1)
	seen := make([]bool, n+1)
	for i := 1; i <= n; i++ {
		v, err := strconv.Atoi(fields[i])
		if err != nil || v < 1 || v > n {
			return 0, nil, fmt.Errorf("invalid a[%d]", i)
		}
		if seen[v] {
			return 0, nil, fmt.Errorf("a is not a permutation")
		}
		seen[v] = true
		a[i] = v
	}
	if a[1] != 1 {
		return 0, nil, fmt.Errorf("a1 must be 1")
	}
	return n, a, nil
}

func add(x *int, y int) {
	*x += y
	if *x >= mod {
		*x -= mod
	}
}

func solveExpected(n int, a []int) int {
	f := make([]int, n+1)
	f[1] = 1
	for i := 1; i < n; i++ {
		p, s := 1, 0
		for j := i; j <= n; j++ {
			x := int((int64(f[j]) * int64(p)) % mod)
			f[j] = x
			p = int((int64(p) * 2) % mod)
			add(&f[j], s)
			if j < n && a[j+1] < a[j] {
				s = 0
			}
			add(&s, x)
		}
	}
	return f[n]
}

func genCase(rng *rand.Rand) []byte {
	n := rng.Intn(60) + 2
	perm := make([]int, 0, n)
	perm = append(perm, 1)
	rest := make([]int, 0, n-1)
	for i := 2; i <= n; i++ {
		rest = append(rest, i)
	}
	rng.Shuffle(len(rest), func(i, j int) {
		rest[i], rest[j] = rest[j], rest[i]
	})
	perm = append(perm, rest...)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(perm[i]))
	}
	sb.WriteByte('\n')
	return []byte(sb.String())
}

func runRandomTests(candidate string) error {
	rng := rand.New(rand.NewSource(1))
	for tc := 1; tc <= 100; tc++ {
		in := genCase(rng)
		if err := verifySingleCase(candidate, in); err != nil {
			return fmt.Errorf("case %d failed: %v\ninput:\n%s", tc, err, string(in))
		}
	}
	return nil
}

func runProgram(bin string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	return out.String(), cmd.Run()
}

func compareAnswer(expected int, got string) error {
	fields := strings.Fields(got)
	if len(fields) == 0 {
		return fmt.Errorf("empty candidate output")
	}
	v, err := strconv.Atoi(fields[0])
	if err != nil {
		return fmt.Errorf("first output token is not an integer")
	}
	if v != expected {
		return fmt.Errorf("wrong answer: expected %d, got %d", expected, v)
	}
	return nil
}
