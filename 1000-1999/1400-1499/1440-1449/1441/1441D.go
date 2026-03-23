package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	scanInt := func() int {
		scanner.Scan()
		res, _ := strconv.Atoi(scanner.Text())
		return res
	}

	if !scanner.Scan() {
		return
	}
	t, _ := strconv.Atoi(scanner.Text())

	for i := 0; i < t; i++ {
		n := scanInt()
		color := make([]int, n+1)
		for j := 1; j <= n; j++ {
			color[j] = scanInt()
		}
		adj := make([][]int, n+1)
		for j := 0; j < n-1; j++ {
			u := scanInt()
			v := scanInt()
			adj[u] = append(adj[u], v)
			adj[v] = append(adj[v], u)
		}

		ans := 0
		const negInf = -1000000000

		var dfs func(u, p int) (int, int)
		dfs = func(u, p int) (int, int) {
			dp0, dp1 := negInf, negInf
			if color[u] == 1 {
				dp0 = 0
			} else if color[u] == 2 {
				dp1 = 0
			}

			for _, v := range adj[u] {
				if v == p {
					continue
				}
				v0, v1 := dfs(v, u)

				if color[u] == 1 {
					tmpV1 := v1
					if tmpV1 != negInf {
						tmpV1++
					}
					bestV := v0
					if tmpV1 > bestV {
						bestV = tmpV1
					}
					if dp0+bestV > ans {
						ans = dp0 + bestV
					}
					if bestV > dp0 {
						dp0 = bestV
					}
				} else if color[u] == 2 {
					tmpV0 := v0
					if tmpV0 != negInf {
						tmpV0++
					}
					bestV := v1
					if tmpV0 > bestV {
						bestV = tmpV0
					}
					if dp1+bestV > ans {
						ans = dp1 + bestV
					}
					if bestV > dp1 {
						dp1 = bestV
					}
				} else {
					if dp0 != negInf && v1 != negInf {
						if dp0+v1+1 > ans {
							ans = dp0 + v1 + 1
						}
					}
					if dp1 != negInf && v0 != negInf {
						if dp1+v0+1 > ans {
							ans = dp1 + v0 + 1
						}
					}
					if dp0 != negInf && v0 != negInf {
						if dp0+v0 > ans {
							ans = dp0 + v0
						}
					}
					if dp1 != negInf && v1 != negInf {
						if dp1+v1 > ans {
							ans = dp1 + v1
						}
					}
					if v0 > dp0 {
						dp0 = v0
					}
					if v1 > dp1 {
						dp1 = v1
					}
				}
			}
			return dp0, dp1
		}

		dfs(1, 0)
		fmt.Fprintln(out, (ans+3)/2)
	}
}
