package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type op struct {
	t, id int
}

func chooseSpot(n int, occupied map[int]bool) int {
	if len(occupied) == 0 {
		return 1
	}
	bestPos := -1
	bestDist := -1
	for pos := 1; pos <= n; pos++ {
		if occupied[pos] {
			continue
		}
		minDist := math.MaxInt32
		for o := range occupied {
			d := int(math.Abs(float64(pos - o)))
			if d < minDist {
				minDist = d
			}
		}
		dist := 4 * minDist
		if dist > bestDist || (dist == bestDist && pos < bestPos) {
			bestDist = dist
			bestPos = pos
		}
	}
	return bestPos
}

func expectedAnswerE(n int, ops []op) []int {
	occupied := make(map[int]bool)
	carPos := make(map[int]int)
	result := make([]int, 0)
	for _, o := range ops {
		if o.t == 1 {
			pos := chooseSpot(n, occupied)
			result = append(result, pos)
			occupied[pos] = true
			carPos[o.id] = pos
		} else {
			pos := carPos[o.id]
			delete(carPos, o.id)
			delete(occupied, pos)
		}
	}
	return result
}

func generateCaseE(rng *rand.Rand) (int, []op) {
	n := rng.Intn(10) + 1
	m := rng.Intn(20) + 1
	ops := make([]op, 0, m)
	nextID := 1
	parked := make([]int, 0)
	for i := 0; i < m; i++ {
		if len(parked) == 0 || (len(parked) < n && rng.Intn(2) == 0) {
			// arrival
			id := nextID
			nextID++
			ops = append(ops, op{1, id})
			parked = append(parked, id)
		} else {
			// departure
			idx := rng.Intn(len(parked))
			id := parked[idx]
			parked = append(parked[:idx], parked[idx+1:]...)
			ops = append(ops, op{2, id})
		}
	}
	return n, ops
}

func runCaseE(bin string, n int, ops []op) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, len(ops)))
	for _, o := range ops {
		sb.WriteString(fmt.Sprintf("%d %d\n", o.t, o.id))
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(strings.TrimSpace(out.String()))
	expected := expectedAnswerE(n, ops)
	if len(fields) != len(expected) {
		return fmt.Errorf("expected %d numbers got %d", len(expected), len(fields))
	}
	for i, f := range fields {
		var val int
		fmt.Sscan(f, &val)
		if val != expected[i] {
			return fmt.Errorf("expected %v got %v", expected, fields)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	if err := runCaseE(bin, 1, []op{{1, 1}}); err != nil {
		fmt.Fprintln(os.Stderr, "deterministic case failed:", err)
		os.Exit(1)
	}

	for i := 0; i < 100; i++ {
		n, ops := generateCaseE(rng)
		if err := runCaseE(bin, n, ops); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
