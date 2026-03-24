package main

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

// magnets: 0 = N, 1 = S, 2 = demagnetized (-)
func genTestCase(rng *rand.Rand) (int, []int) {
	n := rng.Intn(8) + 4 // n >= 4 to guarantee at least 2 non-demagnetized and 1 demagnetized
	types := make([]int, n)
	// Ensure at least 2 non-demagnetized and 1 demagnetized
	// Place at least one N, one S (or two of same type), and one demagnetized
	for i := range types {
		types[i] = rng.Intn(3)
	}
	// Count non-demagnetized
	nonDemag := 0
	demag := 0
	for _, t := range types {
		if t < 2 {
			nonDemag++
		} else {
			demag++
		}
	}
	// Fix: need at least 2 non-demagnetized and at least 1 demagnetized
	for nonDemag < 2 {
		idx := rng.Intn(n)
		if types[idx] == 2 {
			types[idx] = rng.Intn(2)
			nonDemag++
			demag--
		}
	}
	for demag < 1 {
		idx := rng.Intn(n)
		if types[idx] < 2 {
			types[idx] = 2
			nonDemag--
			demag++
		}
	}
	// Make sure we still have >=2 non-demagnetized
	for nonDemag < 2 {
		idx := rng.Intn(n)
		if types[idx] == 2 && demag > 1 {
			types[idx] = rng.Intn(2)
			nonDemag++
			demag--
		}
	}
	return n, types
}

func computeForce(types []int, left, right []int) int {
	n1, s1 := 0, 0
	for _, idx := range left {
		if types[idx] == 0 {
			n1++
		} else if types[idx] == 1 {
			s1++
		}
	}
	n2, s2 := 0, 0
	for _, idx := range right {
		if types[idx] == 0 {
			n2++
		} else if types[idx] == 1 {
			s2++
		}
	}
	return n1*n2 + s1*s2 - n1*s2 - n2*s1
}

func runInteractive(bin string, t int, testCases []struct {
	n     int
	types []int
}) error {
	cmd := exec.Command(bin)
	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("stdin pipe: %v", err)
	}
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("stdout pipe: %v", err)
	}
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start: %v", err)
	}

	writer := bufio.NewWriter(stdinPipe)
	reader := bufio.NewReader(stdoutPipe)

	// Send number of test cases
	fmt.Fprintf(writer, "%d\n", t)
	writer.Flush()

	for tc := 0; tc < t; tc++ {
		n := testCases[tc].n
		types := testCases[tc].types
		maxQueries := n + intLog2(n)

		// Send n
		fmt.Fprintf(writer, "%d\n", n)
		writer.Flush()

		queryCount := 0
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					return fmt.Errorf("tc %d: unexpected EOF from candidate", tc+1)
				}
				return fmt.Errorf("tc %d: read error: %v", tc+1, err)
			}
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			if line[0] == '!' {
				// Answer line
				parts := strings.Fields(line)
				if len(parts) < 2 {
					return fmt.Errorf("tc %d: invalid answer format: %s", tc+1, line)
				}
				cnt, err := strconv.Atoi(parts[1])
				if err != nil {
					return fmt.Errorf("tc %d: invalid count in answer: %s", tc+1, line)
				}
				if len(parts) != cnt+2 {
					return fmt.Errorf("tc %d: answer count mismatch: declared %d but got %d values", tc+1, cnt, len(parts)-2)
				}
				// Collect answered demagnetized indices
				answered := make(map[int]bool)
				for i := 2; i < len(parts); i++ {
					v, err := strconv.Atoi(parts[i])
					if err != nil || v < 1 || v > n {
						return fmt.Errorf("tc %d: invalid index in answer: %s", tc+1, parts[i])
					}
					answered[v] = true
				}
				// Check correctness
				expected := make(map[int]bool)
				for i, t := range types {
					if t == 2 {
						expected[i+1] = true
					}
				}
				if len(answered) != len(expected) {
					return fmt.Errorf("tc %d: wrong answer: expected %d demagnetized, got %d", tc+1, len(expected), len(answered))
				}
				for k := range expected {
					if !answered[k] {
						return fmt.Errorf("tc %d: wrong answer: missing demagnetized magnet %d", tc+1, k)
					}
				}
				break
			} else if line[0] == '?' {
				// Query line
				queryCount++
				if queryCount > maxQueries {
					stdinPipe.Close()
					cmd.Wait()
					return fmt.Errorf("tc %d: too many queries (%d > %d)", tc+1, queryCount, maxQueries)
				}
				// Parse: ? la ra followed by la values then ra values
				// Actually the format from the solution is:
				// ? la ra\n
				// left values\n
				// right values\n
				parts := strings.Fields(line)
				if len(parts) < 3 {
					return fmt.Errorf("tc %d: invalid query format: %s", tc+1, line)
				}
				la, _ := strconv.Atoi(parts[1])
				ra, _ := strconv.Atoi(parts[2])

				// Read left indices
				leftLine, err := reader.ReadString('\n')
				if err != nil {
					return fmt.Errorf("tc %d: read left error: %v", tc+1, err)
				}
				leftParts := strings.Fields(strings.TrimSpace(leftLine))
				if len(leftParts) != la {
					return fmt.Errorf("tc %d: expected %d left values, got %d", tc+1, la, len(leftParts))
				}
				left := make([]int, la)
				for i, p := range leftParts {
					v, _ := strconv.Atoi(p)
					left[i] = v - 1 // 0-indexed
				}

				// Read right indices
				rightLine, err := reader.ReadString('\n')
				if err != nil {
					return fmt.Errorf("tc %d: read right error: %v", tc+1, err)
				}
				rightParts := strings.Fields(strings.TrimSpace(rightLine))
				if len(rightParts) != ra {
					return fmt.Errorf("tc %d: expected %d right values, got %d", tc+1, ra, len(rightParts))
				}
				right := make([]int, ra)
				for i, p := range rightParts {
					v, _ := strconv.Atoi(p)
					right[i] = v - 1 // 0-indexed
				}

				force := computeForce(types, left, right)

				// Check that |force| <= n
				absForce := force
				if absForce < 0 {
					absForce = -absForce
				}
				if absForce > n {
					stdinPipe.Close()
					cmd.Wait()
					return fmt.Errorf("tc %d: machine broke! force=%d, n=%d", tc+1, force, n)
				}

				fmt.Fprintf(writer, "%d\n", force)
				writer.Flush()
			} else {
				return fmt.Errorf("tc %d: unexpected line from candidate: %s", tc+1, line)
			}
		}
	}

	stdinPipe.Close()
	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("candidate exited with error: %v", err)
	}
	return nil
}

func intLog2(n int) int {
	r := 0
	v := n
	for v > 1 {
		v >>= 1
		r++
	}
	return r
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]
	rng := rand.New(rand.NewSource(1491))

	// Run multiple batches of test cases
	for batch := 0; batch < 10; batch++ {
		t := 5
		testCases := make([]struct {
			n     int
			types []int
		}, t)
		for i := 0; i < t; i++ {
			n, types := genTestCase(rng)
			testCases[i].n = n
			testCases[i].types = types
		}

		err := runInteractive(candidate, t, testCases)
		if err != nil {
			fmt.Fprintf(os.Stderr, "batch %d: %v\n", batch+1, err)
			os.Exit(1)
		}
	}

	// Also test some edge cases: all same type non-demagnetized
	edgeCases := []struct {
		n     int
		types []int
	}{
		{4, []int{0, 0, 2, 2}},   // NN--
		{4, []int{1, 1, 2, 2}},   // SS--
		{5, []int{0, 1, 2, 0, 1}}, // N-SNS -> N S - N S
		{4, []int{0, 1, 0, 2}},   // NSNC
	}
	err := runInteractive(candidate, len(edgeCases), edgeCases)
	if err != nil {
		fmt.Fprintf(os.Stderr, "edge cases: %v\n", err)
		os.Exit(1)
	}

	_ = sort.Ints // suppress unused import
	fmt.Println("All tests passed")
}
