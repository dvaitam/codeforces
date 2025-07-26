package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func expectedA(n int, pos []int) int {
	sort.Ints(pos)
	m := n / 2
	black := 0
	white := 0
	for i := 0; i < m; i++ {
		black += abs(pos[i] - (2*i + 1))
		white += abs(pos[i] - (2 * (i + 1)))
	}
	if black < white {
		return black
	}
	return white
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func generateCaseA(rng *rand.Rand) (int, []int) {
	n := rng.Intn(50)*2 + 2 // even between 2 and 100
	m := n / 2
	perm := rng.Perm(n)
	pos := make([]int, m)
	for i := 0; i < m; i++ {
		pos[i] = perm[i] + 1
	}
	return n, pos
}

func runCaseA(bin string, n int, pos []int) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, p := range pos {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(p))
	}
	sb.WriteByte('\n')

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	expected := expectedA(n, append([]int(nil), pos...))
	if got != expected {
		return fmt.Errorf("expected %d got %d", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, pos := generateCaseA(rng)
		if err := runCaseA(bin, n, pos); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%d\n%v\n", i+1, err, n, pos)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
