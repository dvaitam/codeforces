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

func run(prog, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(prog, ".go") {
		cmd = exec.Command("go", "run", prog)
	} else {
		cmd = exec.Command(prog)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand, n int) string {
	perm := rng.Perm(n)
	for i := range perm {
		perm[i]++
	}
	mat := make([][]int, n)
	for i := 0; i < n; i++ {
		mat[i] = make([]int, n)
		for j := 0; j < n; j++ {
			if perm[i] < perm[j] {
				mat[i][j] = perm[i]
			} else {
				mat[i][j] = perm[j]
			}
		}
	}
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(n))
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(mat[i][j]))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

// validate checks that the output permutation is consistent with the matrix.
func validate(input, output string) string {
	lines := strings.Fields(input)
	if len(lines) == 0 {
		return "empty input"
	}
	n, _ := strconv.Atoi(lines[0])
	mat := make([][]int, n)
	for i := 0; i < n; i++ {
		mat[i] = make([]int, n)
		for j := 0; j < n; j++ {
			mat[i][j], _ = strconv.Atoi(lines[1+i*n+j])
		}
	}

	fields := strings.Fields(output)
	if len(fields) != n {
		return fmt.Sprintf("expected %d values, got %d", n, len(fields))
	}
	p := make([]int, n)
	seen := make([]bool, n+1)
	for i, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil || v < 1 || v > n {
			return fmt.Sprintf("value %q out of range", f)
		}
		if seen[v] {
			return fmt.Sprintf("duplicate value %d", v)
		}
		seen[v] = true
		p[i] = v
	}

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if i == j {
				continue
			}
			want := p[i]
			if p[j] < want {
				want = p[j]
			}
			if mat[i][j] != want {
				return fmt.Sprintf("p[%d]=%d p[%d]=%d but matrix[%d][%d]=%d (want %d)", i, p[i], j, p[j], i, j, mat[i][j], want)
			}
		}
	}
	return ""
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 120; t++ {
		n := rng.Intn(50) + 2
		input := genCase(rng, n)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", t+1, err, input)
			os.Exit(1)
		}
		if msg := validate(input, got); msg != "" {
			fmt.Fprintf(os.Stderr, "case %d failed: %s\ngot: %s\ninput:\n%s", t+1, msg, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
