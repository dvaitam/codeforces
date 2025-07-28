package main

import "fmt"

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func main() {
	var n, m int
	fmt.Scan(&n, &m)

	s := make([]byte, n*m)
	for i := 0; i < n; i++ {
		var row string
		fmt.Scan(&row)
		copy(s[i*m:], []byte(row))
	}

	a := make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Scan(&a[i])
	}

	L := make([]int, n*m)
	R := make([]int, n*m)
	for i := range L {
		L[i] = -1
		R[i] = -1
	}

	var dfs func(x, y, t, ty int)
	dfs = func(x, y, t, ty int) {
		pos := x*m + y
		arr := L
		if ty != 0 {
			arr = R
		}
		if arr[pos] != -1 {
			return
		}
		arr[pos] = t
		if x+1 < n {
			dfs(x+1, y, t, ty)
		}
		if y > 0 && s[x*m+y-1] == '#' {
			dfs(x, y-1, t, ty)
		}
		if x > 0 && s[(x-1)*m+y] == '#' {
			dfs(x-1, y, t, ty)
		}
		if y+1 < m && s[x*m+y+1] == '#' {
			dfs(x, y+1, t, ty)
		}
	}

	for j := 0; j < m; j++ {
		for i := 0; i < n; i++ {
			if s[i*m+j] == '#' {
				dfs(i, j, j, 0)
			}
		}
	}

	for j := m - 1; j >= 0; j-- {
		for i := 0; i < n; i++ {
			if s[i*m+j] == '#' {
				dfs(i, j, j, 1)
			}
		}
	}

	f := make([]int, m+2)
	for i := 0; i <= m+1; i++ {
		f[i] = m + 1
	}

	for j := 0; j < m; j++ {
		if a[j] > 0 {
			remaining := a[j]
			for i := n - 1; i >= 0; i-- {
				if s[i*m+j] == '#' {
					remaining--
					if remaining == 0 {
						pos := i*m + j
						l := L[pos]
						r := R[pos]
						f[l] = min(f[l], r+1)
						break
					}
				}
			}
		}
	}

	for i := m; i >= 0; i-- {
		f[i] = min(f[i], f[i+1])
	}

	res := 0
	nw := 0
	for nw < m+1 {
		nw = f[nw]
		res++
	}

	fmt.Printf("%d\n", res-1)
}
