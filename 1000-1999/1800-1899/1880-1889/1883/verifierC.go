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

type testCaseC struct {
	n int
	k int
	a []int
}

func solveCaseC(tc testCaseC) string {
	ans := 0
	switch tc.k {
	case 2:
		ans = 1
		for _, v := range tc.a {
			if v%2 == 0 {
				ans = 0
				break
			}
		}
	case 3:
		ans = 2
		for _, v := range tc.a {
			d := (3 - v%3) % 3
			if d < ans {
				ans = d
				if ans == 0 {
					break
				}
			}
		}
	case 4:
		twos := 0
		for _, v := range tc.a {
			x := v
			for x%2 == 0 {
				twos++
				x /= 2
				if twos >= 2 {
					break
				}
			}
			if twos >= 2 {
				break
			}
		}
		if twos >= 2 {
			ans = 0
			break
		}
		best4 := int(1e9)
		min1, min2 := int(1e9), int(1e9)
		for _, v := range tc.a {
			c4 := (4 - v%4) % 4
			if c4 < best4 {
				best4 = c4
			}
			ce := 0
			if v%2 != 0 {
				ce = 1
			}
			if ce < min1 {
				min2 = min1
				min1 = ce
			} else if ce < min2 {
				min2 = ce
			}
		}
		if min2 == int(1e9) {
			min2 = 0
		}
		if best4 < min1+min2 {
			ans = best4
		} else {
			ans = min1 + min2
		}
	case 5:
		ans = 4
		for _, v := range tc.a {
			d := (5 - v%5) % 5
			if d < ans {
				ans = d
				if ans == 0 {
					break
				}
			}
		}
	}
	return fmt.Sprint(ans)
}

func runCaseC(bin string, tc testCaseC) error {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.k))
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
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
	got := strings.TrimSpace(out.String())
	exp := solveCaseC(tc)
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func randomCaseC(rng *rand.Rand) testCaseC {
	n := rng.Intn(6) + 1
	k := rng.Intn(4) + 2 // 2..5
	a := make([]int, n)
	for i := range a {
		a[i] = rng.Intn(10) + 1
	}
	return testCaseC{n: n, k: k, a: a}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testCaseC{{n: 1, k: 2, a: []int{1}}, {n: 3, k: 4, a: []int{2, 3, 4}}}
	for i := 0; i < 100; i++ {
		cases = append(cases, randomCaseC(rng))
	}
	for idx, tc := range cases {
		if err := runCaseC(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
