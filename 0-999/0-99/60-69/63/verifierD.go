package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCase struct {
	input string
}

func runBinary(bin string, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", bin)
	} else {
		cmd = exec.CommandContext(ctx, bin)
	}
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("time limit")
	}
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, out)
	}
	return strings.TrimSpace(string(out)), nil
}

// verifyOutput validates a candidate grid for CF 63D
// It checks:
// - First line is YES
// - Printed grid has correct dimensions rows=max(b,d), cols=a+c
// - Only land cells are filled; sea cells are '.'
// - Letters are in ['a'..'a'+n) and counts match the provided sizes
func verifyOutput(a, b, c, d int, sizes []int, out string) error {
	rows := b
	if d > rows {
		rows = d
	}
	w := a + c
	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	if len(lines) == 0 {
		return fmt.Errorf("empty output")
	}
	first := strings.TrimSpace(lines[0])
	if strings.ToUpper(first) != "YES" {
		return fmt.Errorf("first line must be YES, got %q", first)
	}
	if len(lines)-1 != rows {
		return fmt.Errorf("expected %d rows, got %d", rows, len(lines)-1)
	}
	grid := make([][]byte, rows)
	for y := 0; y < rows; y++ {
		row := []byte(lines[y+1])
		if len(row) != w {
			return fmt.Errorf("row %d length %d != %d", y+1, len(row), w)
		}
		grid[y] = row
	}
	// valid interval per row
	L := make([]int, rows)
	R := make([]int, rows)
	for y := 0; y < rows; y++ {
		l, r := -1, -2 // empty by default (sea only)
		if y < b && y < d {
			l, r = 0, w-1
		} else if y < b {
			l, r = 0, a-1
		} else if y < d {
			l, r = a, w-1
		}
		L[y], R[y] = l, r
	}
	cnt := make([]int, len(sizes))
	filled := 0
	for y := 0; y < rows; y++ {
		for x := 0; x < w; x++ {
			inside := x >= L[y] && x <= R[y]
			ch := grid[y][x]
			if !inside {
				if ch != '.' {
					return fmt.Errorf("row %d col %d must be '.', got %q", y+1, x+1, ch)
				}
				continue
			}
			if ch < 'a' || ch >= byte('a'+len(sizes)) {
				return fmt.Errorf("invalid letter %q at row %d col %d", ch, y+1, x+1)
			}
			cnt[int(ch-'a')]++
			filled++
		}
	}
	// counts must match and fill total area
	area := a*b + c*d
	if filled != area {
		return fmt.Errorf("filled cells %d != area %d", filled, area)
	}
	for i := range sizes {
		if cnt[i] != sizes[i] {
			return fmt.Errorf("letter %c count %d != %d", 'a'+byte(i), cnt[i], sizes[i])
		}
	}
	return nil
}

func generateCases() []testCase {
	rand.Seed(4)
	cases := make([]testCase, 100)
	for t := 0; t < 100; t++ {
		a := rand.Intn(4) + 1
		b := rand.Intn(4) + 1
		c := rand.Intn(4) + 1
		d := rand.Intn(4) + 1
		area := a*b + c*d
		n := rand.Intn(5) + 1
		if n > 26 {
			n = 26
		}
		if n > area { // ensure feasibility with positive sizes
			n = area
		}
		if n == 0 { // degenerate tiny area guard
			n = 1
		}
		sizes := make([]int, n)
		rem := area
		for i := 0; i < n; i++ {
			if i == n-1 {
				sizes[i] = rem
			} else {
				maxv := rem - (n - i - 1) // >= 1
				// pick in [1, maxv]
				if maxv <= 1 {
					sizes[i] = 1
				} else {
					sizes[i] = rand.Intn(maxv) + 1
				}
				rem -= sizes[i]
			}
		}
		var buf bytes.Buffer
		fmt.Fprintf(&buf, "%d %d %d %d %d\n", a, b, c, d, n)
		for i := 0; i < n; i++ {
			if i > 0 {
				buf.WriteByte(' ')
			}
			fmt.Fprint(&buf, sizes[i])
		}
		buf.WriteByte('\n')
		cases[t] = testCase{input: buf.String()}
	}
	return cases
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierD.go <binary>")
		os.Exit(1)
	}
	cases := generateCases()
	for i, tc := range cases {
		out, err := runBinary(os.Args[1], tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		// parse input
		flds := strings.Fields(tc.input)
		p := 0
		a := atoi(flds[p])
		p++
		b := atoi(flds[p])
		p++
		c := atoi(flds[p])
		p++
		d := atoi(flds[p])
		p++
		n := atoi(flds[p])
		p++
		sz := make([]int, n)
		for k := 0; k < n; k++ {
			sz[k] = atoi(flds[p])
			p++
		}
		if err := verifyOutput(a, b, c, d, sz, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%soutput:\n%s\n", i+1, err, tc.input, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

func atoi(s string) int {
	sign := 1
	i := 0
	if len(s) > 0 && s[0] == '-' {
		sign = -1
		i = 1
	}
	v := 0
	for ; i < len(s); i++ {
		v = v*10 + int(s[i]-'0')
	}
	return v * sign
}
