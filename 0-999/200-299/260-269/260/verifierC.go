package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func simulateCase(rng *rand.Rand) (int, int, []int64) {
	n := rng.Intn(20) + 2
	orig := make([]int64, n)
	for i := range orig {
		orig[i] = int64(rng.Intn(10))
	}
	idx := rng.Intn(n)
	orig[idx] += int64(rng.Intn(5) + 1) // ensure >0 balls removed
	balls := orig[idx]
	final := make([]int64, n)
	copy(final, orig)
	final[idx] = 0
	pos := (idx + 1) % n
	for t := int64(0); t < balls; t++ {
		final[pos]++
		pos = (pos + 1) % n
	}
	x := (idx+int(balls))%n + 1
	return n, x, final
}

func expectedAnswer(n, x int, a []int64) []int64 {
	m := a[0]
	for i := 1; i < n; i++ {
		if a[i] < m {
			m = a[i]
		}
	}
	ks := []int64{m}
	if m > 0 {
		ks = append(ks, m-1)
	}
	for _, k := range ks {
		d := make([]int64, n)
		ok := true
		for i := 0; i < n; i++ {
			d[i] = a[i] - k
			if d[i] < 0 {
				ok = false
				break
			}
		}
		if !ok {
			continue
		}
		rCount := 0
		pos := x - 1
		for rCount < n && d[pos] > 0 {
			rCount++
			pos = (pos - 1 + n) % n
		}
		r := rCount % n
		if k == 0 && r == 0 {
			continue
		}
		i0 := ((x-1-r)%n + n) % n
		seg := make([]bool, n)
		for t := 1; t <= r; t++ {
			j := (i0 + t) % n
			seg[j] = true
		}
		init := make([]int64, n)
		base := k*int64(n) + int64(r)
		init[i0] = base
		valid := base > 0
		for j := 0; j < n && valid; j++ {
			if j == i0 {
				continue
			}
			if seg[j] {
				init[j] = a[j] - k - 1
			} else {
				init[j] = a[j] - k
			}
			if init[j] < 0 {
				valid = false
			}
		}
		if valid {
			return init
		}
	}
	return nil
}

func runCase(bin string, n, x int, arr, expect []int64) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, x))
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
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\nstderr: %s", err, stderr.String())
	}
	scanner := bufio.NewScanner(strings.NewReader(out.String()))
	var got []int64
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		for _, f := range fields {
			var v int64
			fmt.Sscan(f, &v)
			got = append(got, v)
		}
	}
	if len(got) != len(expect) {
		return fmt.Errorf("expected %d numbers got %d", len(expect), len(got))
	}
	for i := range expect {
		if got[i] != expect[i] {
			return fmt.Errorf("mismatch at %d expected %d got %d", i+1, expect[i], got[i])
		}
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
		n, x, arr := simulateCase(rng)
		expect := expectedAnswer(n, x, arr)
		if expect == nil {
			fmt.Fprintf(os.Stderr, "case %d generation failed\n", i+1)
			os.Exit(1)
		}
		if err := runCase(bin, n, x, arr, expect); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
