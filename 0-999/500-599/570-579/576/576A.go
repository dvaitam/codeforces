package main

import (
   "fmt"
   "math"
)

// isPrime checks if x is a prime number
func isPrime(x int) bool {
	if x <= 1 {
		return false
	}
	if x == 2 {
		return true
	}
	if x%2 == 0 {
		return false
	}
	r := int(math.Sqrt(float64(x)))
	for i := 3; i <= r; i += 2 {
		if x%i == 0 {
			return false
		}
	}
	return true
}

func main() {
   var n int
   fmt.Scan(&n)

	var res []int
	for i := 2; i <= n; i++ {
		if isPrime(i) {
			p := i
			for p <= n {
				res = append(res, p)
				if p > n/i {
					break
				}
				p *= i
			}
		}
	}

   fmt.Println(len(res))
   for _, v := range res {
       fmt.Printf("%d ", v)
   }
   fmt.Println()
}
