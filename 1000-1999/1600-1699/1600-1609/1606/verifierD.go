package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func genCase(r *rand.Rand) string {
	n := r.Intn(4) + 2
	m := r.Intn(4) + 2
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(r.Intn(20) + 1))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

// bruteCheck checks if there exists any valid coloring+cut for the given matrix.
// Returns true if a valid answer exists.
func bruteCheck(n, m int, a [][]int) bool {
	// Try all 2^n - 2 colorings (at least one R and one B)
	for mask := 1; mask < (1<<n)-1; mask++ {
		for k := 1; k < m; k++ {
			// Check: every red cell in left (cols 0..k-1) > every blue cell in left
			// and every blue cell in right (cols k..m-1) > every red cell in right
			redMinLeft := math.MaxInt64
			blueMaxLeft := math.MinInt64
			blueMinRight := math.MaxInt64
			redMaxRight := math.MinInt64
			for i := 0; i < n; i++ {
				isRed := (mask>>i)&1 == 1
				for j := 0; j < m; j++ {
					if j < k {
						if isRed {
							if a[i][j] < redMinLeft {
								redMinLeft = a[i][j]
							}
						} else {
							if a[i][j] > blueMaxLeft {
								blueMaxLeft = a[i][j]
							}
						}
					} else {
						if isRed {
							if a[i][j] > redMaxRight {
								redMaxRight = a[i][j]
							}
						} else {
							if a[i][j] < blueMinRight {
								blueMinRight = a[i][j]
							}
						}
					}
				}
			}
			if redMinLeft > blueMaxLeft && blueMinRight > redMaxRight {
				return true
			}
		}
	}
	return false
}

// validateAnswer checks if the candidate's answer is valid for the given matrix.
func validateAnswer(n, m int, a [][]int, colors string, k int) error {
	if len(colors) != n {
		return fmt.Errorf("color string length %d != n %d", len(colors), n)
	}
	hasR, hasB := false, false
	for _, c := range colors {
		if c == 'R' {
			hasR = true
		} else if c == 'B' {
			hasB = true
		} else {
			return fmt.Errorf("invalid color character %c", c)
		}
	}
	if !hasR || !hasB {
		return fmt.Errorf("must have at least one R and one B")
	}
	if k < 1 || k >= m {
		return fmt.Errorf("k=%d out of range [1, %d)", k, m)
	}

	redMinLeft := math.MaxInt64
	blueMaxLeft := math.MinInt64
	blueMinRight := math.MaxInt64
	redMaxRight := math.MinInt64

	for i := 0; i < n; i++ {
		isRed := colors[i] == 'R'
		for j := 0; j < m; j++ {
			if j < k {
				if isRed {
					if a[i][j] < redMinLeft {
						redMinLeft = a[i][j]
					}
				} else {
					if a[i][j] > blueMaxLeft {
						blueMaxLeft = a[i][j]
					}
				}
			} else {
				if isRed {
					if a[i][j] > redMaxRight {
						redMaxRight = a[i][j]
					}
				} else {
					if a[i][j] < blueMinRight {
						blueMinRight = a[i][j]
					}
				}
			}
		}
	}
	if redMinLeft <= blueMaxLeft {
		return fmt.Errorf("left: min red %d <= max blue %d", redMinLeft, blueMaxLeft)
	}
	if blueMinRight <= redMaxRight {
		return fmt.Errorf("right: min blue %d <= max red %d", blueMinRight, redMaxRight)
	}
	return nil
}

func parseInput(input string) (int, int, [][]int) {
	fields := strings.Fields(input)
	idx := 0
	// skip t (always 1)
	idx++
	n, _ := strconv.Atoi(fields[idx]); idx++
	m, _ := strconv.Atoi(fields[idx]); idx++
	a := make([][]int, n)
	for i := 0; i < n; i++ {
		a[i] = make([]int, m)
		for j := 0; j < m; j++ {
			a[i][j], _ = strconv.Atoi(fields[idx]); idx++
		}
	}
	return n, m, a
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 200; i++ {
		input := genCase(rng)
		n, m, a := parseInput(input)
		canSolve := bruteCheck(n, m, a)

		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, input)
			os.Exit(1)
		}

		lines := strings.SplitN(got, "\n", 2)
		verdict := strings.TrimSpace(lines[0])

		if !canSolve {
			if verdict != "NO" {
				fmt.Fprintf(os.Stderr, "case %d: expected NO but got %s\ninput:\n%s", i, verdict, input)
				os.Exit(1)
			}
		} else {
			if verdict != "YES" {
				fmt.Fprintf(os.Stderr, "case %d: expected YES but got %s\ninput:\n%s", i, verdict, input)
				os.Exit(1)
			}
			if len(lines) < 2 {
				fmt.Fprintf(os.Stderr, "case %d: YES but no answer line\ninput:\n%s", i, input)
				os.Exit(1)
			}
			parts := strings.Fields(lines[1])
			if len(parts) != 2 {
				fmt.Fprintf(os.Stderr, "case %d: expected 2 fields on answer line, got %d\ninput:\n%s", i, len(parts), input)
				os.Exit(1)
			}
			colors := parts[0]
			k, _ := strconv.Atoi(parts[1])
			if err := validateAnswer(n, m, a, colors, k); err != nil {
				fmt.Fprintf(os.Stderr, "case %d: invalid answer: %v\ninput:\n%s\noutput:\n%s", i, err, input, got)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All 200 tests passed")
}
