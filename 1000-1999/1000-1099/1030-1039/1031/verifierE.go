package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// canSolve uses BFS on bitmask state to determine if the array can be zeroed
// using triple-flip operations. Only feasible for small n (n <= 22).
func canSolve(a []int) bool {
	n := len(a)
	var start uint32
	for i, v := range a {
		if v == 1 {
			start |= 1 << uint(i)
		}
	}
	if start == 0 {
		return true
	}
	visited := make(map[uint32]bool)
	visited[start] = true
	queue := []uint32{start}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for d := 1; d < n; d++ {
			for i := 0; i+2*d < n; i++ {
				next := cur ^ (1 << uint(i)) ^ (1 << uint(i+d)) ^ (1 << uint(i+2*d))
				if next == 0 {
					return true
				}
				if !visited[next] {
					visited[next] = true
					queue = append(queue, next)
				}
			}
		}
	}
	return false
}

func generateCase(rng *rand.Rand) (string, bool) {
	n := rng.Intn(15) + 3 // keep small for BFS
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = rng.Intn(2)
	}
	ok := canSolve(a)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(a[i]))
	}
	sb.WriteByte('\n')
	return sb.String(), ok
}

func parseOutput(out string) (bool, [][3]int, error) {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) == 0 {
		return false, nil, fmt.Errorf("empty output")
	}
	t := strings.ToUpper(strings.TrimSpace(lines[0]))
	if t == "NO" {
		return false, nil, nil
	}
	if t != "YES" {
		return false, nil, fmt.Errorf("expected YES/NO")
	}
	if len(lines) < 2 {
		return false, nil, fmt.Errorf("missing count line")
	}
	m, err := strconv.Atoi(strings.TrimSpace(lines[1]))
	if err != nil {
		return false, nil, fmt.Errorf("bad count")
	}
	ops := make([][3]int, m)
	if len(lines) != 2+m {
		return false, nil, fmt.Errorf("wrong number of lines")
	}
	for i := 0; i < m; i++ {
		fields := strings.Fields(lines[2+i])
		if len(fields) != 3 {
			return false, nil, fmt.Errorf("bad op line")
		}
		x, _ := strconv.Atoi(fields[0])
		y, _ := strconv.Atoi(fields[1])
		z, _ := strconv.Atoi(fields[2])
		ops[i] = [3]int{x, y, z}
	}
	return true, ops, nil
}

func validateOps(a []int, ops [][3]int, n int) error {
	arr := make([]int, len(a))
	copy(arr, a)
	maxOps := n/3 + 12
	if len(ops) > maxOps {
		return fmt.Errorf("too many operations: %d > %d", len(ops), maxOps)
	}
	for i, op := range ops {
		x, y, z := op[0]-1, op[1]-1, op[2]-1
		if x < 0 || z >= n || x >= y || y >= z {
			return fmt.Errorf("op %d: invalid indices", i+1)
		}
		if y-x != z-y {
			return fmt.Errorf("op %d: not arithmetic progression", i+1)
		}
		arr[x] ^= 1
		arr[y] ^= 1
		arr[z] ^= 1
	}
	for i, v := range arr {
		if v != 0 {
			return fmt.Errorf("position %d is not zero after operations", i+1)
		}
	}
	return nil
}

func runCase(bin, input string, expectOk bool) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\nstderr: %s", err, errBuf.String())
	}
	ok, ops, err := parseOutput(out.String())
	if err != nil {
		return err
	}
	if !ok && expectOk {
		return fmt.Errorf("expected YES but got NO")
	}
	if ok && !expectOk {
		// Candidate says YES but BFS says NO – validate ops anyway
		// (should fail validation if truly impossible)
		parts := strings.Fields(strings.TrimSpace(input))
		n, _ := strconv.Atoi(parts[0])
		a := make([]int, n)
		for i := 0; i < n; i++ {
			a[i], _ = strconv.Atoi(parts[i+1])
		}
		if err := validateOps(a, ops, n); err != nil {
			return fmt.Errorf("candidate says YES but operations invalid: %v", err)
		}
		// Operations are valid, so candidate found a solution BFS missed (shouldn't happen)
	}
	if ok {
		parts := strings.Fields(strings.TrimSpace(input))
		n, _ := strconv.Atoi(parts[0])
		a := make([]int, n)
		for i := 0; i < n; i++ {
			a[i], _ = strconv.Atoi(parts[i+1])
		}
		if err := validateOps(a, ops, n); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(42))
	for i := 0; i < 100; i++ {
		in, ok := generateCase(rng)
		if err := runCase(bin, in, ok); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
