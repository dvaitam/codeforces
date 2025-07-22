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

// BIT for prefix max
type BIT struct {
	n    int
	tree []int
}

func NewBIT(n int) *BIT {
	return &BIT{n: n, tree: make([]int, n+1)}
}

func (b *BIT) Update(i, v int) {
	for ; i <= b.n; i += i & -i {
		if b.tree[i] < v {
			b.tree[i] = v
		}
	}
}

func (b *BIT) Query(i int) int {
	m := 0
	for ; i > 0; i -= i & -i {
		if b.tree[i] > m {
			m = b.tree[i]
		}
	}
	return m
}

func expected(a []int) string {
	n := len(a)
	maxVal := 0
	for _, v := range a {
		if v > maxVal {
			maxVal = v
		}
	}
	bit1 := NewBIT(maxVal + 2)
	dp1 := make([]int, n)
	L := 0
	for i := 0; i < n; i++ {
		v := a[i]
		best := 0
		if v > 1 {
			best = bit1.Query(v - 1)
		}
		dp1[i] = best + 1
		if dp1[i] > L {
			L = dp1[i]
		}
		bit1.Update(v, dp1[i])
	}
	bit2 := NewBIT(maxVal + 2)
	dp2 := make([]int, n)
	for i := n - 1; i >= 0; i-- {
		v := a[i]
		ra := maxVal - v + 1
		best := 0
		if ra > 1 {
			best = bit2.Query(ra - 1)
		}
		dp2[i] = best + 1
		bit2.Update(ra, dp2[i])
	}
	cnt := make([]int, L+1)
	for i := 0; i < n; i++ {
		if dp1[i]+dp2[i]-1 == L {
			cnt[dp1[i]]++
		}
	}
	res := make([]byte, n)
	for i := 0; i < n; i++ {
		if dp1[i]+dp2[i]-1 < L {
			res[i] = '1'
		} else if cnt[dp1[i]] == 1 {
			res[i] = '3'
		} else {
			res[i] = '2'
		}
	}
	return string(res)
}

func verifyCase(bin string, a []int) error {
	n := len(a)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", a[i]))
	}
	sb.WriteByte('\n')
	got, err := run(bin, sb.String())
	if err != nil {
		return err
	}
	want := expected(a)
	if strings.TrimSpace(got) != want {
		return fmt.Errorf("expected %s got %s", want, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	const cases = 100
	for i := 0; i < cases; i++ {
		n := rng.Intn(20) + 1
		a := make([]int, n)
		for j := 0; j < n; j++ {
			a[j] = rng.Intn(20) + 1
		}
		if err := verifyCase(bin, a); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", cases)
}
