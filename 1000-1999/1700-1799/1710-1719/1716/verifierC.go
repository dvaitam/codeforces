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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func solveCase(m int, top, bottom []int) int {
	// a[0][0..m-1] = top row (top[1..m] converted to 0-indexed)
	// a[1][0..m-1] = bottom row
	a := [2][]int{make([]int, m), make([]int, m)}
	for i := 0; i < m; i++ {
		a[0][i] = top[i+1]
		a[1][i] = bottom[i+1]
	}
	// Start cell (1,1) is already occupied at time 0; treat unlock time as -1
	// so that a[0][0]+1 = 0 (minimum arrival time = 0).
	a[0][0] = -1

	// p[j][i] = minimum finish time of the U-shape suffix that starts at row j,
	// column i (0-indexed) and covers all columns i..m-1.
	// The U-shape visits: (j,i),(j,i+1),...,(j,m-1),(1-j,m-1),...,(1-j,i).
	// Recurrence:
	//   p[j][i] = max(a[1-j][i]+1,        // last cell (1-j,i) must be unlocked at finish
	//                 a[j][i]+(m-i)*2,     // first cell (j,i) unlock pushes finish time
	//                 p[j][i+1]+1)         // inner suffix finish + 1 step to (1-j,i)
	p := [2][]int{make([]int, m+1), make([]int, m+1)}
	// p[j][m] = 0 (zero-initialized: empty suffix)
	for i := m - 1; i >= 0; i-- {
		for j := 0; j < 2; j++ {
			v := a[1-j][i] + 1
			if w := a[j][i] + (m-i)*2; w > v {
				v = w
			}
			if w := p[j][i+1] + 1; w > v {
				v = w
			}
			p[j][i] = v
		}
	}

	// Iterate over all possible split columns i (0-indexed).
	// At split i: prefix covers columns 0..i-1 (zigzag), suffix covers i..m-1 (U-shape).
	// c = i&1 is the row at which the suffix starts (alternates with each column).
	// n = minimum finish time of the suffix due to prefix timing constraints.
	// ans = min over i of max(n, p[c][i]).
	ans := int(^uint(0) >> 1)
	n := 0
	for i := 0; i < m; i++ {
		c := i & 1
		cur := n
		if p[c][i] > cur {
			cur = p[c][i]
		}
		if cur < ans {
			ans = cur
		}
		// Update n: propagate finish-time constraints from cells (c,i) and (1-c,i)
		// into the next split's prefix timing.
		if w := a[c][i] + (m-i)*2; w > n {
			n = w
		}
		if w := a[1-c][i] + (m-i-1)*2 - 1 + 2; w > n {
			n = w
		}
	}
	return ans
}

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

func verifyCase(bin string, m int, top, bottom []int) error {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", m))
	for i := 1; i <= m; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(top[i]))
	}
	sb.WriteByte('\n')
	for i := 1; i <= m; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(bottom[i]))
	}
	sb.WriteByte('\n')
	expected := fmt.Sprint(solveCase(m, top, bottom))
	out, err := runCandidate(bin, sb.String())
	if err != nil {
		return err
	}
	if strings.TrimSpace(out) != expected {
		return fmt.Errorf("expected %s got %s", expected, out)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		m := rng.Intn(50) + 2
		top := make([]int, m+1)
		bottom := make([]int, m+1)
		for j := 1; j <= m; j++ {
			top[j] = rng.Intn(1000000000)
		}
		for j := 1; j <= m; j++ {
			bottom[j] = rng.Intn(1000000000)
		}
		if err := verifyCase(bin, m, top, bottom); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
