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
	if a < 0 {
		return -a
	}
	return a
}

func divisors(x int) []int {
	res := []int{}
	for i := 1; i*i <= x; i++ {
		if x%i == 0 {
			res = append(res, i)
			if i*i != x {
				res = append(res, x/i)
			}
		}
	}
	return res
}

func check(cp1 int, anchor int, x []int, allowed []bool, periods []int) (bool, int, int, int, int) {
	r := x[anchor] % cp1
	anchor2 := -1
	g2 := 0
	count2 := 0
	for _, v := range x {
		if v%cp1 != r {
			if anchor2 == -1 {
				anchor2 = v
			} else {
				g2 = gcd(g2, v-anchor2)
			}
			count2++
		}
	}
	if count2 == 0 {
		cp2 := periods[0]
		return true, x[anchor], cp1, 1, cp2
	}
	if count2 == 1 {
		cp2 := periods[0]
		return true, x[anchor], cp1, anchor2, cp2
	}
	divs := divisors(g2)
	for _, d := range divs {
		if d <= 1000000 && allowed[d] {
			return true, x[anchor], cp1, anchor2, d
		}
	}
	return false, 0, 0, 0, 0
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var k, n int
	if _, err := fmt.Fscan(in, &k, &n); err != nil {
		return
	}
	periods := make([]int, k)
	allowed := make([]bool, 1000001)
	for i := 0; i < k; i++ {
		fmt.Fscan(in, &periods[i])
		if periods[i] <= 1000000 {
			allowed[periods[i]] = true
		}
	}
	x := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &x[i])
	}

	g := 0
	for i := 1; i < n; i++ {
		g = gcd(g, x[i]-x[0])
	}
	for _, p := range periods {
		if g%p == 0 {
			fmt.Println("YES")
			fmt.Printf("%d %d\n", x[0], p)
			fmt.Printf("1 %d\n", periods[0])
			return
		}
	}

	idx := []int{0}
	if n > 1 {
		idx = append(idx, 1)
	}
	if n > 2 {
		idx = append(idx, 2)
	}

	for i := 0; i < len(idx); i++ {
		for j := i + 1; j < len(idx); j++ {
			d := x[idx[j]] - x[idx[i]]
			if d < 0 {
				d = -d
			}
			divs := divisors(d)
			for _, cp1 := range divs {
				if cp1 <= 1000000 && allowed[cp1] {
					ok, s1, c1, s2, c2 := check(cp1, idx[i], x, allowed, periods)
					if ok {
						fmt.Println("YES")
						fmt.Printf("%d %d\n", s1, c1)
						fmt.Printf("%d %d\n", s2, c2)
						return
					}
					ok, s1, c1, s2, c2 = check(cp1, idx[j], x, allowed, periods)
					if ok {
						fmt.Println("YES")
						fmt.Printf("%d %d\n", s1, c1)
						fmt.Printf("%d %d\n", s2, c2)
						return
					}
				}
			}
		}
	}

	fmt.Println("NO")
}
