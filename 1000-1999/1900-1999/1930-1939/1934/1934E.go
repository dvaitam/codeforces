package main

import (
	"bufio"
	"fmt"
	"os"
)

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func solve(n int, out *bufio.Writer) {
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
			// r == 2*l-1
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
	// output
	fmt.Fprintln(out, len(ans))
	for _, v := range ans {
		fmt.Fprintf(out, "%d %d %d\n", v[0], v[1], v[2])
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var t, n int
	fmt.Fscan(in, &t)
	for t > 0 {
		t--
		fmt.Fscan(in, &n)
		solve(n, out)
	}
}
