package main

import (
	"fmt"
)

func main() {
	var a int64
	if _, err := fmt.Scan(&a); err != nil {
		return
	}
	if a == 1 {
		fmt.Println(1)
		fmt.Println(0)
		return
	}
	vec := make([]int64, 0, 4)
	cnt := 0
	pr := false
	n := a
	// factor out 2s
	for n%2 == 0 {
		cnt++
		vec = append(vec, 2)
		if cnt > 2 {
			break
		}
		n /= 2
	}
	// factor odd numbers if needed
	if cnt <= 2 {
		for i := int64(3); i*i <= n; i += 2 {
			for n%i == 0 {
				cnt++
				vec = append(vec, i)
				if cnt > 2 {
					break
				}
				n /= i
			}
			if cnt > 2 {
				break
			}
		}
	}
	// remaining prime factor
	if cnt <= 2 && n > 2 {
		cnt++
		vec = append(vec, n)
	}
	// determine result
	ch := 0
	if cnt == 2 {
		ch = 0
	} else if cnt == 1 {
		pr = true
		ch = 1
	} else {
		ch = 1
	}
	if ch != 0 {
		if pr {
			fmt.Println(1)
			fmt.Println(0)
		} else {
			fmt.Println(1)
			ans := vec[0] * vec[1]
			fmt.Println(ans)
		}
	} else {
		fmt.Println(2)
	}
}
