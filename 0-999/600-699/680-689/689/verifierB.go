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

func expected(n int, a []int) []int {
	dist := make([]int, n+1)
	for i := 2; i <= n; i++ {
		dist[i] = -1
	}
	q := make([]int, 0, n)
	q = append(q, 1)
	head := 0
	for head < len(q) {
		v := q[head]
		head++
		d := dist[v]
		if v > 1 && dist[v-1] == -1 {
			dist[v-1] = d + 1
			q = append(q, v-1)
		}
		if v < n && dist[v+1] == -1 {
			dist[v+1] = d + 1
			q = append(q, v+1)
		}
		if dist[a[v]] == -1 {
			dist[a[v]] = d + 1
			q = append(q, a[v])
		}
	}
	return dist[1:]
}

func runCase(bin string, n int, a []int) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range a[1:] {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteString("\n")
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	gotFields := strings.Fields(out.String())
	if len(gotFields) != n {
		return fmt.Errorf("expected %d numbers, got %d", n, len(gotFields))
	}
	exp := expected(n, a)
	for i := 0; i < n; i++ {
		var val int
		fmt.Sscan(gotFields[i], &val)
		if val != exp[i] {
			return fmt.Errorf("index %d expected %d got %d", i+1, exp[i], val)
		}
	}
	return nil
}

func randomCase(rng *rand.Rand) (int, []int) {
	n := rng.Intn(20) + 1
	a := make([]int, n+1)
	cur := 1
	for i := 1; i <= n; i++ {
		minVal := cur
		if minVal < i {
			minVal = i
		}
		val := rng.Intn(n-minVal+1) + minVal
		a[i] = val
		cur = val
	}
	return n, a
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, a := randomCase(rng)
		if err := runCase(bin, n, a); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
