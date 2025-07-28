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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	if errBuf.Len() > 0 {
		return "", fmt.Errorf(errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func solve(n int) string {
	l, r, t := n/2+1, n, 1
	type tri [3]int
	var ans []tri
	var lastv, lastt, lastl int
	for {
		if r <= 11 {
			switch r {
			case 3:
				ans = append(ans, tri{t, 2 * t, 3 * t})
			case 4:
				ans = append(ans, tri{t, 3 * t, 4 * t})
			case 5:
				ans = append(ans, tri{3 * t, 4 * t, 5 * t})
			case 6:
				ans = append(ans, tri{3 * t, 4 * t, 5 * t})
			case 7:
				ans = append(ans, tri{t, 3 * t, 4 * t})
				ans = append(ans, tri{5 * t, 6 * t, 7 * t})
			case 8:
				ans = append(ans, tri{t, 5 * t, 7 * t})
				ans = append(ans, tri{2 * t, 6 * t, 8 * t})
			case 9:
				ans = append(ans, tri{t, 5 * t, 6 * t})
				ans = append(ans, tri{7 * t, 8 * t, 9 * t})
			case 10:
				ans = append(ans, tri{2 * t, 6 * t, 10 * t})
				ans = append(ans, tri{7 * t, 8 * t, 9 * t})
			case 11:
				ans = append(ans, tri{t, 10 * t, 11 * t})
				ans = append(ans, tri{7 * t, 8 * t, 9 * t})
			}
			break
		}
		skip := false
		if l%4 == 3 && r%4 == 1 {
			ans = append(ans, tri{l * t, r * t, 2 * t})
			l++
			skip = true
		}
		for l%4 > 1 {
			l--
		}
		i := (l/4)*4 + 1
		for i+2 <= r {
			ans = append(ans, tri{i * t, (i + 1) * t, (i + 2) * t})
			i += 4
		}
		if i <= r {
			if i+1 <= r {
				ans = append(ans, tri{t, i * t, (i + 1) * t})
			} else {
				if !skip {
					if lastv != 0 {
						if gcd(i*t, lastv) < lastl*lastt {
							d := gcd(i*t, lastv)
							ans = append(ans, tri{d, i * t, lastv})
							lastv = 0
							lastt = 0
						} else {
							ans = append(ans, tri{lastt, 2 * lastt, lastv})
							lastv = i * t
							lastt = t
							lastl = l
						}
					} else {
						lastv = i * t
						lastt = t
						lastl = l
					}
				}
			}
		}
		l = (l-1)/4 + 1
		r = r / 4
		t *= 4
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(ans)))
	for _, v := range ans {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", v[0], v[1], v[2]))
	}
	return strings.TrimSpace(sb.String())
}

func genTests() []int {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]int, 100)
	for i := range cases {
		cases[i] = rng.Intn(2000) + 3
	}
	cases = append(cases, 3, 4, 5, 6, 7)
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := genTests()
	for idx, n := range cases {
		input := fmt.Sprintf("1\n%d\n", n)
		expect := solve(n)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("case %d failed: n=%d expected\n%s\ngot\n%s\n", idx+1, n, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
