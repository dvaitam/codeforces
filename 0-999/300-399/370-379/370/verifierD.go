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

func buildRef() (string, error) {
	refSrc := os.Getenv("REFERENCE_SOURCE_PATH")
	if refSrc == "" {
		return "", fmt.Errorf("REFERENCE_SOURCE_PATH not set")
	}
	exe, err := os.CreateTemp("", "ref370D-*")
	if err != nil {
		return "", err
	}
	exe.Close()
	path := exe.Name()
	cmd := exec.Command("go", "build", "-o", path, refSrc)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return path, nil
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

func genCase(rng *rand.Rand) string {
	n := rng.Intn(4) + 2
	m := rng.Intn(4) + 2
	grid := make([][]byte, n)
	hasW := false
	for i := 0; i < n; i++ {
		row := make([]byte, m)
		for j := 0; j < m; j++ {
			if rng.Intn(4) == 0 {
				row[j] = 'w'
				hasW = true
			} else {
				row[j] = '.'
			}
		}
		grid[i] = row
	}
	if !hasW {
		grid[0][0] = 'w'
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		sb.Write(grid[i])
		sb.WriteByte('\n')
	}
	return sb.String()
}

// parseGrid parses the input to get n, m and the original grid
func parseInput(input string) (int, int, [][]byte) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	var n, m int
	fmt.Sscanf(lines[0], "%d %d", &n, &m)
	grid := make([][]byte, n)
	for i := 0; i < n; i++ {
		grid[i] = []byte(lines[i+1])
	}
	return n, m, grid
}

// validateOutput checks if the candidate output is valid for the given input.
// Returns nil if valid, an error describing the problem otherwise.
func validateOutput(input, output string, refIsNeg1 bool) error {
	n, m, inGrid := parseInput(input)

	// Find bounding box and minimum d from input
	rmin, rmax := n, -1
	cmin, cmax := m, -1
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if inGrid[i][j] == 'w' {
				if i < rmin {
					rmin = i
				}
				if i > rmax {
					rmax = i
				}
				if j < cmin {
					cmin = j
				}
				if j > cmax {
					cmax = j
				}
			}
		}
	}

	outTrimmed := strings.TrimSpace(output)

	if outTrimmed == "-1" {
		if !refIsNeg1 {
			return fmt.Errorf("candidate says -1 but reference found a valid frame")
		}
		return nil
	}

	// Don't blindly trust refIsNeg1 — the reference may be wrong.
	// Instead, validate the candidate's output structurally below.

	// Parse output grid
	outLines := strings.Split(outTrimmed, "\n")
	if len(outLines) != n {
		return fmt.Errorf("expected %d lines, got %d", n, len(outLines))
	}
	outGrid := make([][]byte, n)
	for i := 0; i < n; i++ {
		outGrid[i] = []byte(outLines[i])
		if len(outGrid[i]) != m {
			return fmt.Errorf("line %d: expected %d chars, got %d", i, m, len(outGrid[i]))
		}
	}

	// Check: all original 'w' preserved, '.' either stays '.' or becomes '+'
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			orig := inGrid[i][j]
			out := outGrid[i][j]
			if orig == 'w' {
				if out != 'w' {
					return fmt.Errorf("pixel (%d,%d) was 'w' in input but '%c' in output", i, j, out)
				}
			} else { // orig == '.'
				if out != '.' && out != '+' {
					return fmt.Errorf("pixel (%d,%d) was '.' in input but '%c' in output (expected '.' or '+')", i, j, out)
				}
			}
		}
	}

	// Collect all frame pixels (w and +)
	fr, fc := n, -1
	lr, lc := m, -1
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if outGrid[i][j] == 'w' || outGrid[i][j] == '+' {
				if i < fr {
					fr = i
				}
				if i > fc {
					fc = i
				}
				if j < lr {
					lr = j
				}
				if j > lc {
					lc = j
				}
			}
		}
	}

	// Must be a square
	h := fc - fr + 1
	w := lc - lr + 1
	if h != w {
		return fmt.Errorf("frame bounding box is %dx%d, not square", h, w)
	}
	d := h

	// Check minimum size: d must equal max(rmax-rmin+1, cmax-cmin+1)
	minD := rmax - rmin + 1
	if cmax-cmin+1 > minD {
		minD = cmax - cmin + 1
	}
	if d != minD {
		return fmt.Errorf("frame size %d != minimum required %d", d, minD)
	}

	// All frame border pixels must be 'w' or '+', all interior must be '.'
	for i := fr; i <= fc; i++ {
		for j := lr; j <= lc; j++ {
			isBorder := (i == fr || i == fc || j == lr || j == lc)
			ch := outGrid[i][j]
			if isBorder {
				if ch != 'w' && ch != '+' {
					return fmt.Errorf("border pixel (%d,%d) is '%c', expected 'w' or '+'", i, j, ch)
				}
			} else {
				if ch != '.' {
					return fmt.Errorf("interior pixel (%d,%d) is '%c', expected '.'", i, j, ch)
				}
			}
		}
	}

	// No frame pixels outside the square
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if outGrid[i][j] == '+' || outGrid[i][j] == 'w' {
				if i < fr || i > fc || j < lr || j > lc {
					return fmt.Errorf("frame pixel at (%d,%d) outside bounding box", i, j)
				}
				// Also ensure it's on border
				if i != fr && i != fc && j != lr && j != lc {
					return fmt.Errorf("frame pixel at (%d,%d) not on border of square", i, j)
				}
			}
		}
	}

	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierD <candidate-binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		input := genCase(rng)
		refOut, err := run(ref, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on case %d: %v\n", t+1, err)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", t+1, err)
			os.Exit(1)
		}
		refIsNeg1 := strings.TrimSpace(refOut) == "-1"
		if err := validateOutput(input, got, refIsNeg1); err != nil {
			fmt.Fprintf(os.Stderr, "case %d validation failed: %v\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n",
				t+1, err, input, refOut, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
