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

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func genCase(rng *rand.Rand) string {
	t := rng.Intn(3) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		n := rng.Intn(6) + 1
		sb.WriteString(fmt.Sprintf("%d\n", n))
		// Generate a random permutation a of 1..n
		a := make([]int, n)
		for j := 0; j < n; j++ {
			a[j] = j + 1
		}
		// Fisher-Yates shuffle
		for j := n - 1; j > 0; j-- {
			k := rng.Intn(j + 1)
			a[j], a[k] = a[k], a[j]
		}
		// Compute b[i] = floor((i+1)/a[i]) to ensure a is a valid restoration
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			bi := (j + 1) / a[j]
			sb.WriteString(fmt.Sprintf("%d", bi))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		input := genCase(rng)
		candOut, cErr := runBinary(candidate, input)
		if cErr != nil {
			fmt.Fprintf(os.Stderr, "case %d: candidate error: %v\ninput:\n%s", t+1, cErr, input)
			os.Exit(1)
		}
		// Parse input
		lines := strings.Split(strings.TrimSpace(input), "\n")
		if len(lines) < 2 {
			fmt.Fprintf(os.Stderr, "case %d: malformed input generated\n", t+1)
			os.Exit(1)
		}
		tt, _ := strconv.Atoi(strings.TrimSpace(lines[0]))
		ns := make([]int, tt)
		idx := 1
		for i := 0; i < tt; i++ {
			n, _ := strconv.Atoi(strings.TrimSpace(lines[idx]))
			ns[i] = n
			idx++
			// skip b line
			idx++
		}
		// Parse candidate output numbers across whitespace
		toks := strings.Fields(candOut)
		pos := 0
		for cas := 0; cas < tt; cas++ {
			n := ns[cas]
			if pos+n > len(toks) {
				fmt.Fprintf(os.Stderr, "case %d: insufficient numbers in output for case %d (need %d more)\ninput:\n%s\noutput:\n%s\n", t+1, cas+1, n-(len(toks)-pos), input, candOut)
				os.Exit(1)
			}
			a := make([]int, n)
			used := make([]bool, n+1)
			ok := true
			for i := 0; i < n; i++ {
				v, err := strconv.Atoi(toks[pos+i])
				if err != nil {
					ok = false
					break
				}
				if v < 1 || v > n || used[v] {
					ok = false
					break
				}
				used[v] = true
				a[i] = v
			}
			if !ok {
				fmt.Fprintf(os.Stderr, "case %d: invalid permutation in output for case %d\ninput:\n%s\noutput:\n%s\n", t+1, cas+1, input, candOut)
				os.Exit(1)
			}
			// Validate b[i] from a
			// Reconstruct b line from input for this case
			// Find the line index of b for this case
			// We recompute index to avoid complexity: reparse lines
			// lines structure: [0]=t, then repeating blocks of [n, b]
			// Compute b for this case by scanning from start
			idx2 := 1
			for k := 0; k < cas; k++ {
				// skip n line
				idx2++
				// skip b line
				idx2++
			}
			// Now idx2 at n line for this case
			idx2++ // move to b line
			bFields := strings.Fields(lines[idx2])
			if len(bFields) != n {
				fmt.Fprintf(os.Stderr, "case %d: malformed b line for case %d\ninput:\n%s\n", t+1, cas+1, input)
				os.Exit(1)
			}
			for i := 0; i < n; i++ {
				bi, _ := strconv.Atoi(bFields[i])
				ai := a[i]
				expected := (i + 1) / ai
				if expected != bi {
					fmt.Fprintf(os.Stderr, "case %d failed for subcase %d at position %d\ninput:\n%s\nexpected b=%d from a=%d at i=%d\n", t+1, cas+1, i+1, input, bi, ai, i+1)
					os.Exit(1)
				}
			}
			pos += n
		}
		// Ensure no extra tokens are provided; allow trailing newlines
		if pos != len(toks) {
			// Not fatal, but we can flag it if desired; ignore for flexibility
		}
	}
	fmt.Println("All tests passed")
}
