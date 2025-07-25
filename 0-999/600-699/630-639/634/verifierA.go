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

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func isPossible(a, b []int) bool {
	aa := make([]int, 0, len(a))
	for _, v := range a {
		if v != 0 {
			aa = append(aa, v)
		}
	}
	bb := make([]int, 0, len(b))
	for _, v := range b {
		if v != 0 {
			bb = append(bb, v)
		}
	}
	if len(aa) != len(bb) {
		return false
	}
	m := len(aa)
	if m == 0 {
		return true
	}
	pos := make(map[int]int, m)
	for i, v := range aa {
		pos[v] = i
	}
	shift := -1
	for i, v := range bb {
		idx, ok := pos[v]
		if !ok {
			return false
		}
		d := (i - idx + m) % m
		if shift == -1 {
			shift = d
		} else if shift != d {
			return false
		}
	}
	return true
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(19) + 2 // 2..20
	perm := rng.Perm(n)
	a := make([]int, n)
	copy(a, perm)

	b := make([]int, n)
	if rng.Intn(2) == 0 {
		// generate a valid rotation after removing zero
		aa := make([]int, 0, n-1)
		for _, v := range a {
			if v != 0 {
				aa = append(aa, v)
			}
		}
		if len(aa) > 0 {
			shift := rng.Intn(len(aa))
			rotated := make([]int, len(aa))
			for i := range aa {
				rotated[i] = aa[(i+shift)%len(aa)]
			}
			zeroIdx := rng.Intn(n)
			j := 0
			for i := 0; i < n; i++ {
				if i == zeroIdx {
					b[i] = 0
				} else {
					b[i] = rotated[j]
					j++
				}
			}
		} else {
			copy(b, a)
		}
	} else {
		// completely random permutation
		copy(b, rng.Perm(n))
	}

	exp := "NO"
	if isPossible(a, b) {
		exp = "YES"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	for i, v := range b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')

	return sb.String(), exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
