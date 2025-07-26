package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type Point struct{ x, y int64 }

func solve(pts []Point) int64 {
	n := len(pts)
	sort.Slice(pts, func(i, j int) bool {
		if pts[i].y != pts[j].y {
			return pts[i].y < pts[j].y
		}
		return pts[i].x < pts[j].x
	})
	var ans int64
	for s := 0; s < n; s++ {
		m := n - s - 1
		if m < 4 {
			break
		}
		O := pts[s]
		Q := make([]Point, m)
		copy(Q, pts[s+1:])
		sort.Slice(Q, func(i, j int) bool {
			return (Q[i].x-O.x)*(Q[j].y-O.y)-(Q[i].y-O.y)*(Q[j].x-O.x) > 0
		})
		dp3 := make([][]int, m)
		for i := range dp3 {
			dp3[i] = make([]int, m)
		}
		for k := 1; k < m; k++ {
			for j := k + 1; j < m; j++ {
				cnt := 0
				xk, yk := Q[k].x, Q[k].y
				xj, yj := Q[j].x, Q[j].y
				for i := 0; i < k; i++ {
					if (xk-Q[i].x)*(yj-Q[i].y)-(yk-Q[i].y)*(xj-Q[i].x) > 0 {
						cnt++
					}
				}
				dp3[k][j] = cnt
			}
		}
		dp4 := make([][]int, m)
		for i := range dp4 {
			dp4[i] = make([]int, m)
		}
		for u := 2; u < m; u++ {
			xu, yu := Q[u].x, Q[u].y
			for t := u + 1; t < m; t++ {
				xt, yt := Q[t].x, Q[t].y
				cnt := 0
				for k := 1; k < u; k++ {
					if (xu-Q[k].x)*(yt-Q[k].y)-(yu-Q[k].y)*(xt-Q[k].x) > 0 {
						cnt += dp3[k][u]
					}
				}
				dp4[u][t] = cnt
			}
		}
		for u := 2; u < m; u++ {
			xu, yu := Q[u].x, Q[u].y
			for t := u + 1; t < m; t++ {
				if dp4[u][t] > 0 {
					if (Q[t].x-xu)*(O.y-yu)-(Q[t].y-yu)*(O.x-xu) > 0 {
						ans += int64(dp4[u][t])
					}
				}
			}
		}
	}
	return ans
}

func genCase(rng *rand.Rand) (string, int64) {
	n := rng.Intn(4) + 5 // 5..8
	pts := make([]Point, n)
	for {
		ok := true
		for i := 0; i < n; i++ {
			pts[i] = Point{int64(rng.Intn(11) - 5), int64(rng.Intn(11) - 5)}
		}
		// ensure no duplicates or collinear triple
		dup := false
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				if pts[i] == pts[j] {
					dup = true
				}
			}
		}
		if dup {
			continue
		}
		for i := 0; i < n && ok; i++ {
			for j := i + 1; j < n && ok; j++ {
				for k := j + 1; k < n && ok; k++ {
					if (pts[j].x-pts[i].x)*(pts[k].y-pts[i].y)-(pts[j].y-pts[i].y)*(pts[k].x-pts[i].x) == 0 {
						ok = false
					}
				}
			}
		}
		if ok {
			break
		}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for _, p := range pts {
		fmt.Fprintf(&sb, "%d %d\n", p.x, p.y)
	}
	return sb.String(), solve(pts)
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if out != fmt.Sprint(exp) {
			fmt.Fprintf(os.Stderr, "case %d: expected %d got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
