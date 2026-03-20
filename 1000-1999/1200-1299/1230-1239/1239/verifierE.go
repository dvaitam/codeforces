package main

import (
	"bytes"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"sort"
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
	cmd.Stderr = &out
	cmd.Stdout = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

// computeOptimalCost computes the optimal (minimum) cost for the given array.
// Cost = max(sum of row1, sum of row2) where we split sorted 2n elements into
// two rows of n, row1 ascending, row2 descending, and cost = a[0]+sum(row2) or a[n-1]+sum(row1) etc.
// Actually cost = max(v[0]+sum_row2, v[n-1_of_row1]+sum_row1) where the two smallest go to different rows.
// We use the same DP approach from the reference solution.
func computeOptimalCost(a []int) int {
	n := len(a) / 2
	sort.Ints(a)

	chk := make([][]*big.Int, 2*n)
	for i := 0; i < 2*n; i++ {
		chk[i] = make([]*big.Int, n+1)
		for j := 0; j <= n; j++ {
			chk[i][j] = new(big.Int)
		}
	}
	chk[1][0].SetBit(chk[1][0], 0, 1)

	for i := 2; i < 2*n; i++ {
		for j := 0; j <= i-1 && j <= n-1; j++ {
			chk[i][j].Set(chk[i-1][j])
			if j > 0 {
				shifted := new(big.Int).Lsh(chk[i-1][j-1], uint(a[i]))
				chk[i][j].Or(chk[i][j], shifted)
			}
		}
	}

	tot := 0
	for i := 2; i < 2*n; i++ {
		tot += a[i]
	}
	finalDP := chk[2*n-1][n-1]
	bestCost := math.MaxInt64
	for s := 0; s <= tot; s++ {
		if finalDP.Bit(s) == 1 {
			// cost = a[0] + a[1] + max(s, tot-s)
			mx := s
			if tot-s > mx {
				mx = tot - s
			}
			cost := a[0] + a[1] + mx
			if cost < bestCost {
				bestCost = cost
			}
		}
	}
	return bestCost
}

// checkAnswer verifies the candidate output for a given input.
func checkAnswer(input, output string) error {
	inR := strings.NewReader(input)
	var n int
	fmt.Fscan(inR, &n)
	a := make([]int, 2*n)
	for i := 0; i < 2*n; i++ {
		fmt.Fscan(inR, &a[i])
	}

	optCost := computeOptimalCost(a)

	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) < 2 {
		return fmt.Errorf("expected 2 lines, got %d", len(lines))
	}

	row1Fields := strings.Fields(lines[0])
	row2Fields := strings.Fields(lines[1])
	if len(row1Fields) != n {
		return fmt.Errorf("row1 has %d elements, expected %d", len(row1Fields), n)
	}
	if len(row2Fields) != n {
		return fmt.Errorf("row2 has %d elements, expected %d", len(row2Fields), n)
	}

	row1 := make([]int, n)
	row2 := make([]int, n)
	for i := 0; i < n; i++ {
		v, err := strconv.Atoi(row1Fields[i])
		if err != nil {
			return fmt.Errorf("parse row1[%d]: %v", i, err)
		}
		row1[i] = v
	}
	for i := 0; i < n; i++ {
		v, err := strconv.Atoi(row2Fields[i])
		if err != nil {
			return fmt.Errorf("parse row2[%d]: %v", i, err)
		}
		row2[i] = v
	}

	// Check row1 is non-decreasing
	for i := 1; i < n; i++ {
		if row1[i] < row1[i-1] {
			return fmt.Errorf("row1 not non-decreasing at index %d: %d > %d", i, row1[i-1], row1[i])
		}
	}
	// Check row2 is non-increasing
	for i := 1; i < n; i++ {
		if row2[i] > row2[i-1] {
			return fmt.Errorf("row2 not non-increasing at index %d: %d < %d", i, row2[i-1], row2[i])
		}
	}

	// Check it's a valid partition of a
	sortedA := make([]int, 2*n)
	copy(sortedA, a)
	sort.Ints(sortedA)

	all := make([]int, 0, 2*n)
	all = append(all, row1...)
	all = append(all, row2...)
	sort.Ints(all)

	for i := 0; i < 2*n; i++ {
		if all[i] != sortedA[i] {
			return fmt.Errorf("output is not a valid partition of input elements")
		}
	}

	// Check cost: max path = max over j of (sum(row1[0..j]) + sum(row2[j..n-1]))
	// With row1 ascending, row2 descending, the max is:
	// max(row1[0] + sum(row2), sum(row1) + row2[n-1])
	sum1 := 0
	for _, v := range row1 {
		sum1 += v
	}
	sum2 := 0
	for _, v := range row2 {
		sum2 += v
	}
	path1 := row1[0] + sum2
	path2 := sum1 + row2[n-1]
	cost := path1
	if path2 > cost {
		cost = path2
	}
	if cost != optCost {
		return fmt.Errorf("cost %d != optimal %d", cost, optCost)
	}

	return nil
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(10) + 2
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", n)
	for i := 0; i < 2*n; i++ {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(&b, "%d", rng.Intn(100))
	}
	b.WriteByte('\n')
	return b.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := genCase(rng)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if err := checkAnswer(input, got); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%soutput:\n%s\n", i+1, err, input, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
