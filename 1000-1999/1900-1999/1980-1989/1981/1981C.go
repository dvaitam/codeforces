package main

import (
	"bufio"
	"fmt"
	"os"
)

func get(x int64) int {
	for i := 30; i >= 0; i-- {
		if (x>>i)&1 == 1 {
			return i
		}
	}
	return -1
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	const inf = int64(1e9)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int64, n+2)
		b := make([]int64, n+2)
		for i := 1; i <= n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		flag := true
		l := 1
		for l <= n {
			if a[l] != -1 {
				b[l] = a[l]
				l++
			} else {
				r := l
				for r <= n && a[r] == -1 {
					r++
				}
				// segment [l, r)
				if l == 1 {
					if r <= n {
						b[r] = a[r]
					} else {
						b[r] = 1
					}
					for i := r - 1; i >= 1; i-- {
						if 2*b[i+1] <= inf {
							b[i] = 2 * b[i+1]
						} else {
							b[i] = b[i+1] / 2
						}
					}
				} else if r == n+1 {
					for i := l; i < r; i++ {
						if b[i-1]*2 <= inf {
							b[i] = b[i-1] * 2
						} else {
							b[i] = b[i-1] / 2
						}
					}
				} else {
					L := a[l-1]
					R := a[r]
					can := false
					length := r - l + 1
					for suff := 0; suff <= 30 && suff <= length; suff++ {
						cur := L >> suff
						if cur == 0 {
							continue
						}
						tmp := (r - l) - suff + 1
						v1 := get(cur)
						v2 := get(R)
						if v1 > v2 {
							continue
						}
						need := v2 - v1
						value := R >> need
						if cur != value {
							continue
						}
						if need <= tmp && (tmp-need)%2 == 0 {
							can = true
							// decreasing part
							for i := l; i < l+suff; i++ {
								b[i] = b[i-1] / 2
							}
							// bitwise build
							x := need
							for i := l + suff; i < l+suff+need; i++ {
								x--
								bit := (R >> x) & 1
								b[i] = b[i-1]*2 + bit
							}
							// fill rest
							start := l + suff + need
							for i := start; i < r; i++ {
								if i%2 == start%2 {
									b[i] = b[i-1] * 2
								} else {
									b[i] = b[i-1] / 2
								}
							}
							break
						}
					}
					if !can {
						flag = false
					}
				}
				l = r
			}
		}
		// validate
		for i := 1; i < n; i++ {
			if !(b[i] == b[i+1]/2 || b[i+1] == b[i]/2) {
				flag = false
				break
			}
		}
		if !flag {
			fmt.Fprintln(writer, -1)
		} else {
			for i := 1; i <= n; i++ {
				if i > 1 {
					writer.WriteByte(' ')
				}
				fmt.Fprint(writer, b[i])
			}
			writer.WriteByte('\n')
		}
	}
}
