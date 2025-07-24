package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	cubes := make([][10]bool, n)
	for i := 0; i < n; i++ {
		for j := 0; j < 6; j++ {
			var d int
			fmt.Fscan(in, &d)
			cubes[i][d] = true
		}
	}

	isPossible := func(num int) bool {
		d1 := num % 10
		num /= 10
		if num == 0 {
			for i := 0; i < n; i++ {
				if cubes[i][d1] {
					return true
				}
			}
			return false
		}
		d2 := num % 10
		num /= 10
		if num == 0 { // two-digit
			for i := 0; i < n; i++ {
				if !cubes[i][d2] {
					continue
				}
				for j := 0; j < n; j++ {
					if i == j {
						continue
					}
					if cubes[j][d1] {
						return true
					}
				}
			}
			return false
		}
		d3 := num % 10
		num /= 10
		if num > 0 || n < 3 { // number too big or not enough cubes
			return false
		}
		for i := 0; i < n; i++ {
			if !cubes[i][d3] {
				continue
			}
			for j := 0; j < n; j++ {
				if j == i || !cubes[j][d2] {
					continue
				}
				for k := 0; k < n; k++ {
					if k == i || k == j {
						continue
					}
					if cubes[k][d1] {
						return true
					}
				}
			}
		}
		return false
	}

	ans := 0
	for i := 1; i <= 999; i++ {
		if isPossible(i) {
			ans = i
		} else {
			break
		}
	}
	fmt.Println(ans)
}
