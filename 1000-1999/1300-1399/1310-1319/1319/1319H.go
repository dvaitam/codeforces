package main

import (
	"bufio"
	"fmt"
	"os"
)

func idx(x, y, z, m, k int) int {
	return x*m*k + y*k + z
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m, k int
	if _, err := fmt.Fscan(in, &n, &m, &k); err != nil {
		return
	}

	left := make([][]int, m)
	for i := 0; i < m; i++ {
		left[i] = make([]int, k)
		for j := 0; j < k; j++ {
			fmt.Fscan(in, &left[i][j])
		}
	}

	right := make([][]int, m)
	for i := 0; i < m; i++ {
		right[i] = make([]int, k)
		for j := 0; j < k; j++ {
			fmt.Fscan(in, &right[i][j])
		}
	}

	front := make([][]int, n) // from y=0 towards +y
	for i := 0; i < n; i++ {
		front[i] = make([]int, k)
		for j := 0; j < k; j++ {
			fmt.Fscan(in, &front[i][j])
		}
	}

	back := make([][]int, n) // from y=m+1 towards -y
	for i := 0; i < n; i++ {
		back[i] = make([]int, k)
		for j := 0; j < k; j++ {
			fmt.Fscan(in, &back[i][j])
		}
	}

	down := make([][]int, n) // from z=0 towards +z
	for i := 0; i < n; i++ {
		down[i] = make([]int, m)
		for j := 0; j < m; j++ {
			fmt.Fscan(in, &down[i][j])
		}
	}

	up := make([][]int, n) // from z=k+1 towards -z
	for i := 0; i < n; i++ {
		up[i] = make([]int, m)
		for j := 0; j < m; j++ {
			fmt.Fscan(in, &up[i][j])
		}
	}

	size := n * m * k
	grid := make([]int, size)
	forced := make([]bool, size)

	// process lines that must be empty
	for y := 0; y < m; y++ {
		for z := 0; z < k; z++ {
			L := left[y][z]
			R := right[y][z]
			if L == 0 || R == 0 {
				if L != 0 || R != 0 {
					fmt.Println(-1)
					return
				}
				for x := 0; x < n; x++ {
					forced[idx(x, y, z, m, k)] = true
				}
			}
		}
	}

	for x := 0; x < n; x++ {
		for z := 0; z < k; z++ {
			F := front[x][z]
			B := back[x][z]
			if F == 0 || B == 0 {
				if F != 0 || B != 0 {
					fmt.Println(-1)
					return
				}
				for y := 0; y < m; y++ {
					forced[idx(x, y, z, m, k)] = true
				}
			}
		}
	}

	for x := 0; x < n; x++ {
		for y := 0; y < m; y++ {
			D := down[x][y]
			U := up[x][y]
			if D == 0 || U == 0 {
				if D != 0 || U != 0 {
					fmt.Println(-1)
					return
				}
				for z := 0; z < k; z++ {
					forced[idx(x, y, z, m, k)] = true
				}
			}
		}
	}

	// assign from sensors with positive values
	for y := 0; y < m; y++ {
		for z := 0; z < k; z++ {
			L := left[y][z]
			if L == 0 {
				continue
			}
			ok := false
			for x := 0; x < n; x++ {
				id := idx(x, y, z, m, k)
				if forced[id] {
					continue
				}
				if grid[id] == 0 {
					grid[id] = L
					ok = true
					break
				}
				if grid[id] == L {
					ok = true
					break
				}
				// another block encountered
				fmt.Println(-1)
				return
			}
			if !ok {
				fmt.Println(-1)
				return
			}
		}
	}

	for y := 0; y < m; y++ {
		for z := 0; z < k; z++ {
			R := right[y][z]
			if R == 0 {
				continue
			}
			ok := false
			for x := n - 1; x >= 0; x-- {
				id := idx(x, y, z, m, k)
				if forced[id] {
					continue
				}
				if grid[id] == 0 {
					grid[id] = R
					ok = true
					break
				}
				if grid[id] == R {
					ok = true
					break
				}
				fmt.Println(-1)
				return
			}
			if !ok {
				fmt.Println(-1)
				return
			}
		}
	}

	for x := 0; x < n; x++ {
		for z := 0; z < k; z++ {
			F := front[x][z]
			if F == 0 {
				continue
			}
			ok := false
			for y := 0; y < m; y++ {
				id := idx(x, y, z, m, k)
				if forced[id] {
					continue
				}
				if grid[id] == 0 {
					grid[id] = F
					ok = true
					break
				}
				if grid[id] == F {
					ok = true
					break
				}
				fmt.Println(-1)
				return
			}
			if !ok {
				fmt.Println(-1)
				return
			}
		}
	}

	for x := 0; x < n; x++ {
		for z := 0; z < k; z++ {
			B := back[x][z]
			if B == 0 {
				continue
			}
			ok := false
			for y := m - 1; y >= 0; y-- {
				id := idx(x, y, z, m, k)
				if forced[id] {
					continue
				}
				if grid[id] == 0 {
					grid[id] = B
					ok = true
					break
				}
				if grid[id] == B {
					ok = true
					break
				}
				fmt.Println(-1)
				return
			}
			if !ok {
				fmt.Println(-1)
				return
			}
		}
	}

	for x := 0; x < n; x++ {
		for y := 0; y < m; y++ {
			D := down[x][y]
			if D == 0 {
				continue
			}
			ok := false
			for z := 0; z < k; z++ {
				id := idx(x, y, z, m, k)
				if forced[id] {
					continue
				}
				if grid[id] == 0 {
					grid[id] = D
					ok = true
					break
				}
				if grid[id] == D {
					ok = true
					break
				}
				fmt.Println(-1)
				return
			}
			if !ok {
				fmt.Println(-1)
				return
			}
		}
	}

	for x := 0; x < n; x++ {
		for y := 0; y < m; y++ {
			U := up[x][y]
			if U == 0 {
				continue
			}
			ok := false
			for z := k - 1; z >= 0; z-- {
				id := idx(x, y, z, m, k)
				if forced[id] {
					continue
				}
				if grid[id] == 0 {
					grid[id] = U
					ok = true
					break
				}
				if grid[id] == U {
					ok = true
					break
				}
				fmt.Println(-1)
				return
			}
			if !ok {
				fmt.Println(-1)
				return
			}
		}
	}

	// verification
	for y := 0; y < m; y++ {
		for z := 0; z < k; z++ {
			L := left[y][z]
			R := right[y][z]
			if L == 0 && R == 0 {
				for x := 0; x < n; x++ {
					if grid[idx(x, y, z, m, k)] != 0 {
						fmt.Println(-1)
						return
					}
				}
			} else {
				val := 0
				for x := 0; x < n; x++ {
					if grid[idx(x, y, z, m, k)] != 0 {
						val = grid[idx(x, y, z, m, k)]
						break
					}
				}
				if val != L {
					fmt.Println(-1)
					return
				}
				val = 0
				for x := n - 1; x >= 0; x-- {
					if grid[idx(x, y, z, m, k)] != 0 {
						val = grid[idx(x, y, z, m, k)]
						break
					}
				}
				if val != R {
					fmt.Println(-1)
					return
				}
			}
		}
	}

	for x := 0; x < n; x++ {
		for z := 0; z < k; z++ {
			F := front[x][z]
			B := back[x][z]
			if F == 0 && B == 0 {
				for y := 0; y < m; y++ {
					if grid[idx(x, y, z, m, k)] != 0 {
						fmt.Println(-1)
						return
					}
				}
			} else {
				val := 0
				for y := 0; y < m; y++ {
					if grid[idx(x, y, z, m, k)] != 0 {
						val = grid[idx(x, y, z, m, k)]
						break
					}
				}
				if val != F {
					fmt.Println(-1)
					return
				}
				val = 0
				for y := m - 1; y >= 0; y-- {
					if grid[idx(x, y, z, m, k)] != 0 {
						val = grid[idx(x, y, z, m, k)]
						break
					}
				}
				if val != B {
					fmt.Println(-1)
					return
				}
			}
		}
	}

	for x := 0; x < n; x++ {
		for y := 0; y < m; y++ {
			D := down[x][y]
			U := up[x][y]
			if D == 0 && U == 0 {
				for z := 0; z < k; z++ {
					if grid[idx(x, y, z, m, k)] != 0 {
						fmt.Println(-1)
						return
					}
				}
			} else {
				val := 0
				for z := 0; z < k; z++ {
					if grid[idx(x, y, z, m, k)] != 0 {
						val = grid[idx(x, y, z, m, k)]
						break
					}
				}
				if val != D {
					fmt.Println(-1)
					return
				}
				val = 0
				for z := k - 1; z >= 0; z-- {
					if grid[idx(x, y, z, m, k)] != 0 {
						val = grid[idx(x, y, z, m, k)]
						break
					}
				}
				if val != U {
					fmt.Println(-1)
					return
				}
			}
		}
	}

	out := bufio.NewWriter(os.Stdout)
	for x := 0; x < n; x++ {
		for y := 0; y < m; y++ {
			for z := 0; z < k; z++ {
				if z > 0 {
					fmt.Fprint(out, " ")
				}
				fmt.Fprint(out, grid[idx(x, y, z, m, k)])
			}
			if y < m-1 {
				fmt.Fprintln(out)
			}
		}
		if x < n-1 {
			fmt.Fprintln(out)
		}
	}
	out.Flush()
}
