package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var s string
	if _, err := fmt.Fscan(in, &s); err != nil {
		return
	}
	n := len(s)
	hasX := false
	for i := 0; i < n; i++ {
		if s[i] == 'X' {
			hasX = true
			break
		}
	}
	endings := [][2]int{{0, 0}, {2, 5}, {5, 0}, {7, 5}}
	total := 0

	if n == 1 {
		loops := 10
		if !hasX {
			loops = 1
		}
		for d := 0; d < loops; d++ {
			valid := true
			ch := s[0]
			if ch == 'X' {
				if d != 0 {
					valid = false
				}
			} else if ch == '_' {
				// must be 0 -> always valid
			} else {
				if int(ch-'0') != 0 {
					valid = false
				}
			}
			if valid {
				total++
			}
		}
		fmt.Println(total)
		return
	}

	loops := 10
	if !hasX {
		loops = 1
	}
	for _, e := range endings {
		a, b := e[0], e[1]
		for d := 0; d < loops; d++ {
			valid := true
			ways := 1
			for i := 0; i < n; i++ {
				want := -1
				if i == n-2 {
					want = a
				}
				if i == n-1 {
					want = b
				}
				ch := s[i]
				if ch == 'X' {
					if want != -1 {
						if d != want {
							valid = false
							break
						}
						if i == 0 && want == 0 {
							valid = false
							break
						}
					} else {
						if i == 0 && d == 0 {
							valid = false
							break
						}
					}
				} else if ch == '_' {
					if want != -1 {
						if i == 0 && want == 0 {
							valid = false
							break
						}
					} else {
						if i == 0 {
							ways *= 9
						} else {
							ways *= 10
						}
					}
				} else {
					digit := int(ch - '0')
					if want != -1 && digit != want {
						valid = false
						break
					}
					if i == 0 && digit == 0 {
						valid = false
						break
					}
				}
			}
			if valid {
				total += ways
			}
		}
	}
	fmt.Println(total)
}
