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

func run(bin string, input string) (string, error) {
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

func bfsCount(f []string) int {
	n := len(f)
	visited := make([]bool, n)
	queue := []int{0}
	visited[0] = true
	count := make([]int, n)
	for head := 0; head < len(queue); head++ {
		u := queue[head]
		for v := 0; v < n; v++ {
			if f[u][v] == '1' {
				count[v]++
				if v != n-1 && !visited[v] {
					visited[v] = true
					queue = append(queue, v)
				}
			}
		}
	}
	return count[n-1]
}

func runCase(bin string, f []string) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(f)))
	for _, row := range f {
		sb.WriteString(row)
		sb.WriteByte('\n')
	}
	out, err := run(bin, sb.String())
	if err != nil {
		return err
	}
	got, err := strconv.Atoi(strings.TrimSpace(out))
	if err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	exp := bfsCount(f)
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func randomMatrix(rng *rand.Rand, n int) []string {
	rows := make([]string, n)
	for i := 0; i < n; i++ {
		b := make([]byte, n)
		for j := 0; j < n; j++ {
			if i == j {
				b[j] = '0'
			} else {
				if rng.Intn(2) == 0 {
					b[j] = '0'
				} else {
					b[j] = '1'
				}
			}
		}
		rows[i] = string(b)
	}
	// ensure symmetry
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if rows[i][j] != rows[j][i] {
				if rng.Intn(2) == 0 {
					rows[i] = rows[i][:j] + string(rows[j][i]) + rows[i][j+1:]
				} else {
					rows[j] = rows[j][:i] + string(rows[i][j]) + rows[j][i+1:]
				}
			}
		}
	}
	return rows
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(6) + 2 // between 2 and 7
		f := randomMatrix(rng, n)
		if err := runCase(bin, f); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
