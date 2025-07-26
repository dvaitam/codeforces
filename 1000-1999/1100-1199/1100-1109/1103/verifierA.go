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

func runCandidate(bin, input string) (string, error) {
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

func verifyCase(s, out string) error {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != len(s) {
		return fmt.Errorf("expected %d lines, got %d", len(s), len(lines))
	}
	grid := [4][4]int{}
	for idx, ch := range s {
		var r, c int
		if _, err := fmt.Sscanf(lines[idx], "%d %d", &r, &c); err != nil {
			return fmt.Errorf("line %d parse error: %v", idx+1, err)
		}
		if r < 1 || r > 4 || c < 1 || c > 4 {
			return fmt.Errorf("line %d out of bounds", idx+1)
		}
		r--
		c--
		if ch == '0' {
			if r+1 >= 4 {
				return fmt.Errorf("tile %d out of bounds", idx+1)
			}
			if grid[r][c] != 0 || grid[r+1][c] != 0 {
				return fmt.Errorf("tile %d overlaps", idx+1)
			}
			grid[r][c] = 1
			grid[r+1][c] = 1
		} else {
			if c+1 >= 4 {
				return fmt.Errorf("tile %d out of bounds", idx+1)
			}
			if grid[r][c] != 0 || grid[r][c+1] != 0 {
				return fmt.Errorf("tile %d overlaps", idx+1)
			}
			grid[r][c] = 1
			grid[r][c+1] = 1
		}
		// clear full rows
		for i := 0; i < 4; i++ {
			full := true
			for j := 0; j < 4; j++ {
				if grid[i][j] == 0 {
					full = false
					break
				}
			}
			if full {
				for j := 0; j < 4; j++ {
					grid[i][j] = 0
				}
			}
		}
		// clear full cols
		for j := 0; j < 4; j++ {
			full := true
			for i := 0; i < 4; i++ {
				if grid[i][j] == 0 {
					full = false
					break
				}
			}
			if full {
				for i := 0; i < 4; i++ {
					grid[i][j] = 0
				}
			}
		}
	}
	return nil
}

func genCase(r *rand.Rand) string {
	n := r.Intn(1000) + 1
	var sb strings.Builder
	for i := 0; i < n; i++ {
		if r.Intn(2) == 0 {
			sb.WriteByte('0')
		} else {
			sb.WriteByte('1')
		}
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []string{"0", "1"}
	cases = append(cases, strings.Repeat("0", 1000))
	cases = append(cases, strings.Repeat("1", 1000))
	for i := len(cases); i < 100; i++ {
		cases = append(cases, genCase(rng))
	}
	for idx, s := range cases {
		out, err := runCandidate(bin, s+"\n")
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n", idx+1, err)
			os.Exit(1)
		}
		if err := verifyCase(s, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput: %s\noutput:\n%s\n", idx+1, err, s, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
