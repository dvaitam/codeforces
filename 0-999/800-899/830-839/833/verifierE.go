package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type cloud struct {
	l, r int64
	c    int64
}

func earliest(k int64, clouds []cloud, C int64) int64 {
	n := len(clouds)
	best := int64(1 << 62)
	for mask := 0; mask < (1 << n); mask++ {
		if bits.OnesCount(uint(mask)) > 2 {
			continue
		}
		var cost int64
		for i := 0; i < n; i++ {
			if mask&(1<<i) != 0 {
				cost += clouds[i].c
			}
		}
		if cost > C {
			continue
		}
		var segs [][2]int64
		for i := 0; i < n; i++ {
			if mask&(1<<i) == 0 {
				segs = append(segs, [2]int64{clouds[i].l, clouds[i].r})
			}
		}
		sort.Slice(segs, func(i, j int) bool { return segs[i][0] < segs[j][0] })
		merged := make([][2]int64, 0, len(segs))
		for _, s := range segs {
			if len(merged) == 0 || s[0] > merged[len(merged)-1][1] {
				merged = append(merged, s)
			} else if s[1] > merged[len(merged)-1][1] {
				merged[len(merged)-1][1] = s[1]
			}
		}
		t := computeEarliest(k, merged)
		if t < best {
			best = t
		}
	}
	return best
}

func computeEarliest(k int64, segs [][2]int64) int64 {
	var cum int64
	prev := int64(0)
	for _, s := range segs {
		if prev < s[0] {
			sunny := s[0] - prev
			if cum+sunny >= k {
				return prev + (k - cum)
			}
			cum += sunny
		}
		if s[1] > prev {
			prev = s[1]
		}
	}
	return prev + (k - cum)
}

func solveCase(n int, C int64, clouds []cloud, ks []int64) string {
	var sb strings.Builder
	for i, k := range ks {
		if i > 0 {
			sb.WriteByte('\n')
		}
		sb.WriteString(fmt.Sprintf("%d", earliest(k, clouds, C)))
	}
	return sb.String()
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(4)
	C := int64(rng.Intn(20))
	clouds := make([]cloud, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, C))
	for i := 0; i < n; i++ {
		l := int64(rng.Intn(10))
		r := l + int64(rng.Intn(10)+1)
		c := int64(rng.Intn(10))
		clouds[i] = cloud{l, r, c}
		sb.WriteString(fmt.Sprintf("%d %d %d\n", l, r, c))
	}
	m := rng.Intn(3) + 1
	sb.WriteString(fmt.Sprintf("%d\n", m))
	ks := make([]int64, m)
	for i := 0; i < m; i++ {
		ks[i] = int64(rng.Intn(20) + 1)
		sb.WriteString(fmt.Sprintf("%d\n", ks[i]))
	}
	expected := solveCase(n, C, clouds, ks)
	if !strings.HasSuffix(expected, "\n") {
		expected += "\n"
	}
	return sb.String(), expected
}

func runExe(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, expected := generateCase(rng)
		got, err := runExe(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != strings.TrimSpace(expected) {
			fmt.Fprintf(os.Stderr, "case %d failed\nexpected:\n%s\ngot:\n%s\ninput:%s", i+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
