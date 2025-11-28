package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func runProg(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return out.String() + errBuf.String(), fmt.Errorf("runtime error: %v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(3) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	counts := make([]int, n)
	for i := 0; i < n; i++ {
		c := rng.Intn(3) + 1
		counts[i] = c
		fmt.Fprintf(&sb, "%d", c)
		
		currentVal := 0
		for j := 0; j < c; j++ {
			currentVal += rng.Intn(5) + 1 // Ensure strict increase
			fmt.Fprintf(&sb, " %d", currentVal)
		}
		sb.WriteByte('\n')
	}
	totalBuilds := 1
	for _, c := range counts {
		totalBuilds *= c
	}
	
	// Ensure we don't ban all builds (problem guarantee)
	maxM := 3 // Arbitrary upper bound for M
	if maxM >= totalBuilds {
		maxM = totalBuilds - 1
	}
	
	var m int
	if maxM > 0 {
		m = rng.Intn(maxM + 1)
	} else {
		m = 0
	}

	fmt.Fprintf(&sb, "%d\n", m)
	seen := make(map[string]bool)
	for len(seen) < m {
		var build strings.Builder
		for i := 0; i < n; i++ {
			idx := rng.Intn(counts[i]) + 1
			if i > 0 {
				build.WriteByte(' ')
			}
			fmt.Fprintf(&build, "%d", idx)
		}
		s := build.String()
		if seen[s] {
			continue
		}
		seen[s] = true
		sb.WriteString(s)
		sb.WriteByte('\n')
	}
	return sb.String()
}

type Key [10]int

func solveBruteForce(n int, a [][]int, banned map[Key]bool) (int64, Key) {
	var maxStrength int64 = -1
	var bestBuild Key

	// Helper for recursion
	var indices Key
	var backtrack func(depth int)
	backtrack = func(depth int) {
		if depth == n {
			// Check if banned
			if banned[indices] {
				return
			}
			// Calc strength
			var currentStrength int64
			for i := 0; i < n; i++ {
				// indices is 1-based
				val := a[i][indices[i]-1]
				currentStrength += int64(val)
			}
			if currentStrength > maxStrength {
				maxStrength = currentStrength
				bestBuild = indices
			}
			return
		}
		// Try all items for slot depth
		count := len(a[depth])
		for i := 1; i <= count; i++ {
			indices[depth] = i
			backtrack(depth + 1)
		}
	}
	backtrack(0)
	return maxStrength, bestBuild
}

func verify(input string, output string) error {
	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(bufio.ScanWords)

	scanInt := func() int {
		scanner.Scan()
		val, _ := strconv.Atoi(scanner.Text())
		return val
	}

	if !scanner.Scan() {
		return fmt.Errorf("empty input")
	}
	n, _ := strconv.Atoi(scanner.Text())

	a := make([][]int, n)
	for i := 0; i < n; i++ {
		c := scanInt()
		a[i] = make([]int, c)
		for j := 0; j < c; j++ {
			a[i][j] = scanInt()
		}
	}

	m := scanInt()
	banned := make(map[Key]bool)
	for i := 0; i < m; i++ {
		var k Key
		for j := 0; j < n; j++ {
			k[j] = scanInt()
		}
		banned[k] = true
	}

	// Solve brute force
	maxStrength, _ := solveBruteForce(n, a, banned)

	// Parse output
	outScanner := bufio.NewScanner(strings.NewReader(output))
	outScanner.Split(bufio.ScanWords)
	
	var candidateIndices Key
	for i := 0; i < n; i++ {
		if !outScanner.Scan() {
			return fmt.Errorf("output too short, expected %d integers", n)
		}
		val, err := strconv.Atoi(outScanner.Text())
		if err != nil {
			return fmt.Errorf("output parse error: %v", err)
		}
		candidateIndices[i] = val
	}

	// Validate candidate
	for i := 0; i < n; i++ {
		if candidateIndices[i] < 1 || candidateIndices[i] > len(a[i]) {
			return fmt.Errorf("index %d out of bounds for slot %d (size %d)", candidateIndices[i], i+1, len(a[i]))
		}
	}

	if banned[candidateIndices] {
		return fmt.Errorf("candidate build is banned: %v", candidateIndices)
	}

	var candidateStrength int64
	for i := 0; i < n; i++ {
		candidateStrength += int64(a[i][candidateIndices[i]-1])
	}

	if candidateStrength != maxStrength {
		return fmt.Errorf("suboptimal strength: got %d, expected %d", candidateStrength, maxStrength)
	}

	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	cand := os.Args[1]
	
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := genCase(rng)
		out, err := runProg(cand, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if err := verify(input, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s\noutput:\n%s\n", i+1, err, input, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
