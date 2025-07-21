package main

import "fmt"

func main() {
	var n, radix int
	if _, err := fmt.Scan(&n); err != nil {
		return
	}
	if _, err := fmt.Scan(&radix); err != nil {
		return
	}
	sum := 0
	for i := 0; i < n; i++ {
		var s string
		fmt.Scan(&s)
		v := 0
		for _, c := range s {
			var d int
			if c >= '0' && c <= '9' {
				d = int(c - '0')
			} else if c >= 'A' && c <= 'Z' {
				d = int(c - 'A' + 10)
			}
			v = v*radix + d
		}
		sum += v
	}
	if sum == 0 {
		fmt.Println("0")
		return
	}
	var res []rune
	for sum > 0 {
		rem := sum % radix
		var c rune
		if rem < 10 {
			c = rune('0' + rem)
		} else {
			c = rune('A' + rem - 10)
		}
		res = append(res, c)
		sum /= radix
	}
	for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
		res[i], res[j] = res[j], res[i]
	}
	fmt.Println(string(res))
}
