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

func expectedC(n, k int, l int64, arr []int64) int64 {
	m := n * k
	b := append([]int64(nil), arr...)
	sort.Slice(b, func(i, j int) bool { return b[i] < b[j] })
	limit := b[0] + l
	good := 0
	for good < m && b[good] <= limit {
		good++
	}
	if good < n {
		return 0
	}
	pos := good - 1
	right := m - 1
	var result int64
	for i := 0; i < n; i++ {
		result += b[pos]
		pos--
		for j := 0; j < k-1; j++ {
			if right > pos {
				right--
			} else {
				pos--
			}
		}
	}
	return result
}

func generateCaseC(rng *rand.Rand) (int, int, int64, []int64) {
	for {
		n := rng.Intn(5) + 1
		k := rng.Intn(5) + 1
		if n*k <= 20 {
			l := int64(rng.Intn(10))
			m := n * k
			arr := make([]int64, m)
			for i := 0; i < m; i++ {
				arr[i] = int64(rng.Intn(20)) + 1
			}
			return n, k, l, arr
		}
	}
}

func runCaseC(bin string, n, k int, l int64, arr []int64) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, k, l))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
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
	var got int64
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	expected := expectedC(n, k, l, arr)
	if got != expected {
		return fmt.Errorf("expected %d got %d", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, k, l, arr := generateCaseC(rng)
		if err := runCaseC(bin, n, k, l, arr); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
