package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// bruteCanWin checks (by exhaustive search) whether any partition exists.
func bruteCanWin(a []int, x int) bool {
	n := len(a)
	for mask := 1; mask < (1 << n); mask++ {
		red := 0
		for i := 0; i < n; i++ {
			if mask>>i&1 == 1 {
				red++
			}
		}
		if red < 2 || red > n-2 {
			continue
		}
		g := 0
		andVal := (1 << 30) - 1
		for i := 0; i < n; i++ {
			if mask>>i&1 == 1 {
				g = gcd(g, a[i])
			} else {
				andVal &= a[i]
			}
		}
		if g > andVal+x {
			return true
		}
	}
	return false
}

func multisetEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	cnt := make(map[int]int)
	for _, v := range a {
		cnt[v]++
	}
	for _, v := range b {
		cnt[v]--
		if cnt[v] < 0 {
			return false
		}
	}
	return true
}

func parseNums(line string) ([]int, error) {
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return nil, fmt.Errorf("empty line")
	}
	var count int
	if _, err := fmt.Sscan(fields[0], &count); err != nil {
		return nil, fmt.Errorf("bad count: %v", err)
	}
	if len(fields)-1 != count {
		return nil, fmt.Errorf("expected %d numbers, got %d", count, len(fields)-1)
	}
	nums := make([]int, count)
	for i := 0; i < count; i++ {
		if _, err := fmt.Sscan(fields[i+1], &nums[i]); err != nil {
			return nil, fmt.Errorf("bad number: %v", err)
		}
	}
	return nums, nil
}

func checkOutput(output string, a []int, x int) error {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	verdict := strings.ToUpper(strings.TrimSpace(lines[0]))

	canWin := bruteCanWin(a, x)

	if verdict == "NO" {
		if canWin {
			return fmt.Errorf("said NO but a solution exists")
		}
		return nil
	}

	if verdict != "YES" {
		return fmt.Errorf("unexpected verdict %q", verdict)
	}
	if !canWin {
		return fmt.Errorf("said YES but no solution exists")
	}
	if len(lines) < 3 {
		return fmt.Errorf("YES answer needs 3 lines, got %d", len(lines))
	}

	red, err := parseNums(lines[1])
	if err != nil {
		return fmt.Errorf("parsing red line: %v", err)
	}
	blue, err := parseNums(lines[2])
	if err != nil {
		return fmt.Errorf("parsing blue line: %v", err)
	}

	n := len(a)
	if len(red) < 2 || len(red) > n-2 {
		return fmt.Errorf("red count %d not in [2, %d]", len(red), n-2)
	}
	if !multisetEqual(append(red, blue...), a) {
		return fmt.Errorf("partition is not a rearrangement of the input")
	}

	g := 0
	for _, v := range red {
		g = gcd(g, v)
	}
	andVal := (1 << 30) - 1
	for _, v := range blue {
		andVal &= v
	}
	if g <= andVal+x {
		return fmt.Errorf("GCD(%d) not > AND+x(%d+%d=%d)", g, andVal, x, andVal+x)
	}
	return nil
}

func runExe(path, input string) (string, error) {
	cmd := exec.Command(path)
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
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

func generateCase(rng *rand.Rand) (string, []int, int) {
	n := rng.Intn(7) + 4
	x := rng.Intn(20)
	a := make([]int, n)
	for i := range a {
		a[i] = rng.Intn(50) + 1
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	fmt.Fprintf(&sb, "%d %d\n", n, x)
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	return sb.String(), a, x
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, a, x := generateCase(rng)
		got, err := runExe(exe, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if err := checkOutput(got, a, x); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%sgot:\n%s\n", i+1, err, input, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
