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

func isomorphic(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	ma := make(map[byte]byte)
	mb := make(map[byte]byte)
	for i := 0; i < len(a); i++ {
		x := a[i]
		y := b[i]
		if v, ok := ma[x]; ok {
			if v != y {
				return false
			}
		} else {
			ma[x] = y
		}
		if v, ok := mb[y]; ok {
			if v != x {
				return false
			}
		} else {
			mb[y] = x
		}
	}
	return true
}

func expectedF(s string, queries [][3]int) []string {
	res := make([]string, len(queries))
	for idx, q := range queries {
		x, y, l := q[0], q[1], q[2]
		a := s[x : x+l]
		b := s[y : y+l]
		if isomorphic(a, b) {
			res[idx] = "YES"
		} else {
			res[idx] = "NO"
		}
	}
	return res
}

func generateCaseF(rng *rand.Rand) (string, [][3]int) {
	n := rng.Intn(10) + 1
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = byte('a' + rng.Intn(3))
	}
	m := rng.Intn(10) + 1
	qs := make([][3]int, m)
	for i := 0; i < m; i++ {
		if n == 1 {
			qs[i] = [3]int{0, 0, 1}
			continue
		}
		x := rng.Intn(n)
		y := rng.Intn(n)
		l := rng.Intn(n-max(x, y)) + 1
		qs[i] = [3]int{x, y, l}
	}
	return string(b), qs
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func runCaseF(bin string, s string, qs [][3]int) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", len(s), len(qs)))
	sb.WriteString(s)
	sb.WriteByte('\n')
	for _, q := range qs {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", q[0]+1, q[1]+1, q[2]))
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	gotLines := strings.Split(strings.TrimSpace(out.String()), "\n")
	expectedLines := expectedF(s, qs)
	if len(gotLines) != len(expectedLines) {
		return fmt.Errorf("wrong number of lines: expected %d got %d", len(expectedLines), len(gotLines))
	}
	for i := range gotLines {
		if strings.TrimSpace(gotLines[i]) != expectedLines[i] {
			return fmt.Errorf("line %d: expected %s got %s", i+1, expectedLines[i], strings.TrimSpace(gotLines[i]))
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		s, qs := generateCaseF(rng)
		if err := runCaseF(bin, s, qs); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
