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

func countAC(s string) int {
	cnt := 0
	for i := 0; i+1 < len(s); i++ {
		if s[i] == 'A' && s[i+1] == 'C' {
			cnt++
		}
	}
	return cnt
}

func gen(k, x, n, m int) (string, string, bool) {
	letters := []byte{'A', 'X', 'C'}
	for fs1 := 0; fs1 < 3; fs1++ {
		for ls1 := 0; ls1 < 3; ls1++ {
			for fs2 := 0; fs2 < 3; fs2++ {
				for ls2 := 0; ls2 < 3; ls2++ {
					for num1 := 0; num1 <= n/2; num1++ {
						for num2 := 0; num2 <= m/2; num2++ {
							if n == 1 && fs1 != ls1 {
								continue
							}
							if m == 1 && fs2 != ls2 {
								continue
							}
							if n > 1 && n%2 == 0 && num1 == n/2 && (fs1 != 0 || ls1 != 2) {
								continue
							}
							if m > 1 && m%2 == 0 && num2 == m/2 && (fs2 != 0 || ls2 != 2) {
								continue
							}
							if n > 1 && n%2 == 1 && num1 == n/2 && (fs1 != 0 && ls1 != 2) {
								continue
							}
							if m > 1 && m%2 == 1 && num2 == m/2 && (fs2 != 0 && ls2 != 2) {
								continue
							}
							if n == 2 && fs1 == 0 && ls1 == 2 && num1 != 1 {
								continue
							}
							if m == 2 && fs2 == 0 && ls2 == 2 && num2 != 1 {
								continue
							}
							// simulate counts
							c1 := num1
							c2 := num2
							of, ol := fs1, ls1
							nf, nl := fs2, ls2
							ok := true
							if k == 1 {
								if c1 != x {
									ok = false
								}
							} else if k == 2 {
								if c2 != x {
									ok = false
								}
							} else {
								for it := 2; it < k && ok; it++ {
									next := c1 + c2
									if ol == 0 && nf == 2 {
										next++
									}
									c1, c2 = c2, next
									of, ol, nf, nl = nf, nl, of, nl
									if c2 > x {
										ok = false
									}
								}
								if ok && c2 != x {
									ok = false
								}
							}
							if ok {
								// build strings
								var sb1, sb2 strings.Builder
								sb1.WriteByte(letters[fs1])
								rem := num1
								for i := 1; i < n-1; i++ {
									if rem > 0 && sb1.String()[i-1] == 'A' {
										sb1.WriteByte('C')
										rem--
									} else {
										sb1.WriteByte('A')
									}
								}
								if n > 1 {
									sb1.WriteByte(letters[ls1])
								}
								rem = num2
								sb2.WriteByte(letters[fs2])
								for i := 1; i < m-1; i++ {
									if rem > 0 && sb2.String()[i-1] == 'A' {
										sb2.WriteByte('C')
										rem--
									} else {
										sb2.WriteByte('A')
									}
								}
								if m > 1 {
									sb2.WriteByte(letters[ls2])
								}
								return sb1.String(), sb2.String(), true
							}
						}
					}
				}
			}
		}
	}
	return "", "", false
}

func expectedValid(k, x, n, m int, s1, s2 string) bool {
	if len(s1) != n || len(s2) != m {
		return false
	}
	c1 := countAC(s1)
	c2 := countAC(s2)
	of, ol := s1[0], s1[len(s1)-1]
	nf, nl := s2[0], s2[len(s2)-1]
	for it := 2; it < k; it++ {
		next := c1 + c2
		if ol == 'A' && nf == 'C' {
			next++
		}
		c1, c2 = c2, next
		of, ol, nf, nl = nf, nl, of, nl
		if c2 > x {
			return false
		}
	}
	return c2 == x
}

func runCase(exe string, k, x, n, m int) error {
	input := fmt.Sprintf("%d %d %d %d\n", k, x, n, m)
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) == 1 && strings.TrimSpace(lines[0]) == "Happy new year!" {
		if _, _, ok := gen(k, x, n, m); ok {
			return fmt.Errorf("solution claims impossible but there is a pair")
		}
		return nil
	}
	if len(lines) != 2 {
		return fmt.Errorf("expected two lines, got %d", len(lines))
	}
	s1 := strings.TrimSpace(lines[0])
	s2 := strings.TrimSpace(lines[1])
	if !expectedValid(k, x, n, m, s1, s2) {
		return fmt.Errorf("output strings do not form valid solution")
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	params := [][4]int{{3, 0, 1, 1}, {3, 1, 2, 2}, {4, 2, 3, 3}}
	for i := 0; i < 100; i++ {
		k := rng.Intn(5) + 3
		x := rng.Intn(5)
		n := rng.Intn(3) + 1
		m := rng.Intn(3) + 1
		params = append(params, [4]int{k, x, n, m})
	}
	for idx, p := range params {
		if err := runCase(exe, p[0], p[1], p[2], p[3]); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
