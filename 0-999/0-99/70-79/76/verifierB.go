package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

func solve(data string) string {
	in := bufio.NewReader(strings.NewReader(data))
	var n, m int
	var y0, y1 int
	if _, err := fmt.Fscan(in, &n, &m, &y0, &y1); err != nil {
		return ""
	}
	mice := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &mice[i])
	}
	cheese := make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &cheese[i])
	}
	if m == 0 {
		return fmt.Sprintf("%d\n", n)
	}
	t := make([]int, m)
	inf := int64(1) << 60
	d2 := make([]int64, m)
	c := make([]int, m)
	for j := 0; j < m; j++ {
		d2[j] = inf
	}
	for _, x := range mice {
		k := sort.SearchInts(cheese, x)
		if k > 0 && k < m {
			dl := absInt(x - cheese[k-1])
			dr := cheese[k] - x
			if dr < 0 {
				dr = -dr
			}
			if dl == dr {
				t[k-1]++
				continue
			}
		}
		var j int
		if k == 0 {
			j = 0
		} else if k == m {
			j = m - 1
		} else {
			dl := absInt(x - cheese[k-1])
			dr := cheese[k] - x
			if dr < 0 {
				dr = -dr
			}
			if dl < dr {
				j = k - 1
			} else {
				j = k
			}
		}
		dist2 := int64(absInt(x-cheese[j])) * 2
		if c[j] == 0 || dist2 < d2[j] {
			d2[j] = dist2
			c[j] = 1
		} else if dist2 == d2[j] {
			c[j]++
		}
	}
	dpPrev := [2]int64{0, -inf}
	dpCur := [2]int64{0, 0}
	for j := 0; j < m-1; j++ {
		dpCur[0], dpCur[1] = -inf, -inf
		for s := 0; s < 2; s++ {
			base := dpPrev[s]
			if base < 0 {
				continue
			}
			for ns := 0; ns < 2; ns++ {
				var Lcnt, Rcnt int
				if s == 1 && j > 0 {
					Lcnt = t[j-1]
				}
				if ns == 1 {
					Rcnt = t[j]
				}
				minDist := inf
				if c[j] > 0 {
					minDist = d2[j]
				}
				if Lcnt > 0 {
					dl := int64(cheese[j] - cheese[j-1])
					if dl < minDist {
						minDist = dl
					}
				}
				if Rcnt > 0 {
					dr := int64(cheese[j+1] - cheese[j])
					if dr < minDist {
						minDist = dr
					}
				}
				fed := int64(0)
				if c[j] > 0 && d2[j] == minDist {
					fed += int64(c[j])
				}
				if Lcnt > 0 && int64(cheese[j]-cheese[j-1]) == minDist {
					fed += int64(Lcnt)
				}
				if Rcnt > 0 && int64(cheese[j+1]-cheese[j]) == minDist {
					fed += int64(Rcnt)
				}
				val := base + fed
				if val > dpCur[ns] {
					dpCur[ns] = val
				}
			}
		}
		dpPrev = dpCur
	}
	best := int64(0)
	for s := 0; s < 2; s++ {
		base := dpPrev[s]
		if base < 0 {
			continue
		}
		j := m - 1
		var Lcnt int
		if s == 1 && j > 0 {
			Lcnt = t[j-1]
		}
		minDist := inf
		if c[j] > 0 {
			minDist = d2[j]
		}
		if Lcnt > 0 {
			dl := int64(cheese[j] - cheese[j-1])
			if dl < minDist {
				minDist = dl
			}
		}
		fed := int64(0)
		if c[j] > 0 && d2[j] == minDist {
			fed += int64(c[j])
		}
		if Lcnt > 0 && int64(cheese[j]-cheese[j-1]) == minDist {
			fed += int64(Lcnt)
		}
		total := base + fed
		if total > best {
			best = total
		}
	}
	hungry := int64(n) - best
	if hungry < 0 {
		hungry = 0
	}
	return fmt.Sprintf("%d\n", hungry)
}

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(9) + 1
	m := rng.Intn(9)
	y0 := rng.Intn(20)
	y1 := rng.Intn(20)
	if y1 == y0 {
		y1++
	}
	mice := make([]int, n)
	for i := range mice {
		mice[i] = rng.Intn(50)
	}
	sort.Ints(mice)
	cheese := make([]int, m)
	for i := range cheese {
		cheese[i] = rng.Intn(50)
	}
	sort.Ints(cheese)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d %d\n", n, m, y0, y1))
	for i, v := range mice {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	for i, v := range cheese {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	input := sb.String()
	expected := solve(input)
	return input, expected
}

func runCase(exe string, in, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	if strings.TrimSpace(out.String()) != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %q got %q", strings.TrimSpace(expected), strings.TrimSpace(out.String()))
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(exe, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
