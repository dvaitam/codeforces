package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
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

func isIndependent(b []int) bool {
	n := len(b)
	mat := make([][]float64, n)
	for i := 0; i < n; i++ {
		mat[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			idx := (j - i + n) % n
			mat[i][j] = float64(b[idx])
		}
	}
	// Gaussian elimination
	for col, row := 0, 0; col < n && row < n; col++ {
		pivot := row
		for i := row; i < n; i++ {
			if mat[i][col] != 0 {
				pivot = i
				break
			}
		}
		if mat[pivot][col] == 0 {
			continue
		}
		mat[row], mat[pivot] = mat[pivot], mat[row]
		pv := mat[row][col]
		for j := col; j < n; j++ {
			mat[row][j] /= pv
		}
		for i := 0; i < n; i++ {
			if i == row {
				continue
			}
			factor := mat[i][col]
			for j := col; j < n; j++ {
				mat[i][j] -= factor * mat[row][j]
			}
		}
		row++
	}
	rank := 0
	for i := 0; i < n; i++ {
		zero := true
		for j := 0; j < n; j++ {
			if math.Abs(mat[i][j]) > 1e-9 {
				zero = false
				break
			}
		}
		if !zero {
			rank++
		}
	}
	return rank == n
}

func computeC(a, b []int) []int {
	n := len(a)
	c := make([]int, n)
	for i := 0; i < n; i++ {
		sum := 0
		for k := 0; k < n; k++ {
			diff := b[(k-i+n)%n] - a[k]
			sum += diff * diff
		}
		c[i] = sum
	}
	return c
}

func enumerateSolutions(b []int, c []int) [][]int {
	n := len(b)
	sol := [][]int{}
	vals := []int{-2, -1, 0, 1, 2}
	cur := make([]int, n)
	var dfs func(int)
	dfs = func(pos int) {
		if pos == n {
			if equalSlices(computeC(cur, b), c) {
				tmp := make([]int, n)
				copy(tmp, cur)
				sol = append(sol, tmp)
			}
			return
		}
		for _, v := range vals {
			cur[pos] = v
			dfs(pos + 1)
		}
	}
	dfs(0)
	sort.Slice(sol, func(i, j int) bool {
		for k := 0; k < n; k++ {
			if sol[i][k] != sol[j][k] {
				return sol[i][k] < sol[j][k]
			}
		}
		return false
	})
	return sol
}

func equalSlices(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func main() {
	if len(os.Args) == 3 && os.Args[1] == "--" {
		os.Args = append([]string{os.Args[0]}, os.Args[2])
	}
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		n := rng.Intn(3) + 1
		var b []int
		for {
			b = make([]int, n)
			for i := 0; i < n; i++ {
				b[i] = rng.Intn(4)
			}
			if isIndependent(b) {
				break
			}
		}
		a := make([]int, n)
		for i := 0; i < n; i++ {
			a[i] = rng.Intn(5) - 2
		}
		c := computeC(a, b)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", n)
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", b[i])
		}
		sb.WriteByte('\n')
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", c[i])
		}
		sb.WriteByte('\n')
		input := sb.String()
		solutions := enumerateSolutions(b, c)
		var exp strings.Builder
		fmt.Fprintf(&exp, "%d", len(solutions))
		for _, sol := range solutions {
			exp.WriteByte('\n')
			for i := 0; i < n; i++ {
				if i > 0 {
					exp.WriteByte(' ')
				}
				fmt.Fprintf(&exp, "%d", sol[i])
			}
		}
		expected := exp.String()
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", t+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expected) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\ninput:%s", t+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
