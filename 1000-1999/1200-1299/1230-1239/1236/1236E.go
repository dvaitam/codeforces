package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	a := make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(reader, &a[i])
	}

	var ans int64
	for x := 1; x <= n; x++ {
		L := x - 1
		if L < 1 {
			L = 1
		}
		R := x + 1
		if R > n {
			R = n
		}
		if a[0] >= L && a[0] <= R {
			if L == R && L == a[0] {
				continue
			} else if a[0] == L {
				L++
			} else if a[0] == R {
				R--
			}
		}
		valid := true
		for i := 1; i < m && valid; i++ {
			L -= 2
			if L < 1 {
				L = 1
			}
			R += 2
			if R > n {
				R = n
			}
			if L > R {
				valid = false
				break
			}
			ai := a[i]
			if ai >= L && ai <= R {
				if L == R && L == ai {
					valid = false
					break
				} else if ai == L {
					L++
				} else if ai == R {
					R--
				}
			}
		}
		if !valid {
			continue
		}
		L -= 1
		if L < 1 {
			L = 1
		}
		R += 1
		if R > n {
			R = n
		}
		if L <= R {
			ans += int64(R - L + 1)
		}
	}

	fmt.Fprintln(writer, ans)
}
