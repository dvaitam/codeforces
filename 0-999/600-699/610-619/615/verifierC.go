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

func computeMinSegments(s1, s2 string) int {
	b1 := append([]byte{0}, []byte(s1)...)
	b2 := append([]byte{0}, []byte(s2)...)
	l1 := len(s1)
	l2 := len(s2)
	i := 1
	count := 0
	for i <= l2 {
		best := 0
		for j := 1; j <= l1; j++ {
			if b1[j] != b2[i] {
				continue
			}
			k := j
			for k+1 <= l1 && (k+1-j+i) <= l2 && b1[k+1] == b2[k+1-j+i] {
				k++
			}
			if k-j+1 > best {
				best = k - j + 1
			}
			k2 := j - 1
			for k2 >= 1 && (j-k2+i) <= l2 && b1[k2] == b2[j-k2+i] {
				k2--
			}
			if j-k2 > best {
				best = j - k2
			}
		}
		if best == 0 {
			return -1
		}
		count++
		i += best
	}
	return count
}

func buildFromSegments(s1 string, pairs [][2]int) string {
	var sb strings.Builder
	for _, p := range pairs {
		l, r := p[0], p[1]
		if l >= 1 && r >= 1 && l <= len(s1) && r <= len(s1) {
			if l <= r {
				sb.WriteString(s1[l-1 : r])
			} else {
				for i := l - 1; i >= r-1; i-- {
					sb.WriteByte(s1[i])
				}
			}
		}
	}
	return sb.String()
}

func genReachableCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	letters := "abcde"
	s1b := make([]byte, n)
	for i := range s1b {
		s1b[i] = letters[rng.Intn(len(letters))]
	}
	s1 := string(s1b)
	var s2 strings.Builder
	pieces := rng.Intn(4) + 1
	for i := 0; i < pieces; i++ {
		l := rng.Intn(n)
		r := rng.Intn(n)
		if l <= r {
			s2.WriteString(s1[l : r+1])
		} else {
			for j := l; j >= r; j-- {
				s2.WriteByte(s1[j])
			}
		}
	}
	return s1, s2.String()
}

func genUnreachableCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	letters := "abcde"
	s1b := make([]byte, n)
	for i := range s1b {
		s1b[i] = letters[rng.Intn(len(letters))]
	}
	s1 := string(s1b)
	s2 := s1 + "z"
	return s1, s2
}

func runCase(bin string, s1, s2 string) error {
	input := fmt.Sprintf("%s\n%s\n", s1, s2)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(out.String())
	if len(fields) == 0 {
		return fmt.Errorf("no output")
	}
	if fields[0] == "-1" {
		exp := computeMinSegments(s1, s2)
		if exp != -1 {
			return fmt.Errorf("expected %d segments but got -1", exp)
		}
		return nil
	}
	var k int
	if _, err := fmt.Sscan(fields[0], &k); err != nil {
		return fmt.Errorf("cannot parse k: %v", err)
	}
	if len(fields) != 1+2*k {
		return fmt.Errorf("expected %d coordinates, got %d", 2*k, len(fields)-1)
	}
	pairs := make([][2]int, k)
	idx := 1
	for i := 0; i < k; i++ {
		var l, r int
		fmt.Sscan(fields[idx], &l)
		fmt.Sscan(fields[idx+1], &r)
		pairs[i] = [2]int{l, r}
		idx += 2
	}
	built := buildFromSegments(s1, pairs)
	if built != s2 {
		return fmt.Errorf("segments do not form target")
	}
	exp := computeMinSegments(s1, s2)
	if exp != k {
		return fmt.Errorf("expected %d segments got %d", exp, k)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	s1 := "abac"
	s2 := "aba"
	cases := [][2]string{{s1, s2}}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(cases) < 80 {
		a, b := genReachableCase(rng)
		cases = append(cases, [2]string{a, b})
	}
	for len(cases) < 100 {
		a, b := genUnreachableCase(rng)
		cases = append(cases, [2]string{a, b})
	}

	for i, tc := range cases {
		if err := runCase(bin, tc[0], tc[1]); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s\n%s\n", i+1, err, tc[0], tc[1])
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
