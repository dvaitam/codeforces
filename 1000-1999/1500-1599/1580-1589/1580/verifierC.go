package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const solverB = 450

func solverNextInt(data []byte, idx *int) int {
	n := len(data)
	for *idx < n {
		c := data[*idx]
		if c >= '0' && c <= '9' {
			break
		}
		*idx++
	}
	v := 0
	for *idx < n {
		c := data[*idx]
		if c < '0' || c > '9' {
			break
		}
		v = v*10 + int(c-'0')
		*idx++
	}
	return v
}

func solve1580C(input string) string {
	data := []byte(input)
	idx := 0

	n := solverNextInt(data, &idx)
	m := solverNextInt(data, &idx)

	x := make([]int, n+1)
	y := make([]int, n+1)
	per := make([]int, n+1)

	hasSmall := make([]bool, solverB+1)
	smallPeriods := make([]int, 0, solverB)

	for i := 1; i <= n; i++ {
		xi := solverNextInt(data, &idx)
		yi := solverNextInt(data, &idx)
		p := xi + yi
		x[i] = xi
		y[i] = yi
		per[i] = p
		if p <= solverB && !hasSmall[p] {
			hasSmall[p] = true
			smallPeriods = append(smallPeriods, p)
		}
	}

	sdiff := make([][]int, solverB+1)
	for _, p := range smallPeriods {
		sdiff[p] = make([]int, p)
	}

	rem := make([]int, solverB+1)
	curr := make([]int, solverB+1)
	smallTotal := 0

	start := make([]int, n+1)
	diff := make([]int, m+3)
	curLarge := 0

	out := make([]byte, 0, m*8)

	for day := 1; day <= m; day++ {
		op := solverNextInt(data, &idx)
		k := solverNextInt(data, &idx)

		for i := 0; i < len(smallPeriods); i++ {
			p := smallPeriods[i]
			nr := rem[p] + 1
			if nr == p {
				rem[p] = 0
				old := curr[p]
				nv := sdiff[p][0]
				curr[p] = nv
				smallTotal += nv - old
			} else {
				rem[p] = nr
				d := sdiff[p][nr]
				curr[p] += d
				smallTotal += d
			}
		}

		pk := per[k]
		if op == 1 {
			start[k] = day
			if pk <= solverB {
				d := sdiff[pk]
				l := (day + x[k]) % pk
				length := y[k]
				end := l + length
				if end <= pk {
					d[l]++
					if end < pk {
						d[end]--
					}
				} else {
					d[l]++
					d[0]++
					d[end-pk]--
				}
				r := rem[pk]
				if (end <= pk && r >= l && r < end) || (end > pk && (r >= l || r < end-pk)) {
					curr[pk]++
					smallTotal++
				}
			} else {
				step := pk
				yy := y[k]
				for l := day + x[k]; l <= m; l += step {
					r := l + yy
					if r > m+1 {
						r = m + 1
					}
					diff[l]++
					diff[r]--
				}
			}
		} else {
			s := start[k]
			start[k] = 0
			if pk <= solverB {
				d := sdiff[pk]
				l := (s + x[k]) % pk
				length := y[k]
				end := l + length
				if end <= pk {
					d[l]--
					if end < pk {
						d[end]++
					}
				} else {
					d[l]--
					d[0]--
					d[end-pk]++
				}
				r := rem[pk]
				if (end <= pk && r >= l && r < end) || (end > pk && (r >= l || r < end-pk)) {
					curr[pk]--
					smallTotal--
				}
			} else {
				step := pk
				yy := y[k]
				for l := s + x[k]; l <= m; l += step {
					r := l + yy
					if r > m+1 {
						r = m + 1
					}
					L := l
					if L < day {
						L = day
					}
					if L < r {
						diff[L]--
						diff[r]++
					}
				}
			}
		}

		curLarge += diff[day]
		out = strconv.AppendInt(out, int64(curLarge+smallTotal), 10)
		out = append(out, '\n')
	}

	return strings.TrimSpace(string(out))
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(r *rand.Rand) string {
	n := r.Intn(3) + 2 // 2..4
	m := r.Intn(10) + 1
	x := make([]int, n+1)
	y := make([]int, n+1)
	for i := 1; i <= n; i++ {
		x[i] = r.Intn(3) + 1
		y[i] = r.Intn(3) + 1
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 1; i <= n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", x[i], y[i]))
	}
	active := make([]bool, n+1)
	for i := 0; i < m; i++ {
		op := 1
		if anyActive := func() bool {
			for _, a := range active[1:] {
				if a {
					return true
				}
			}
			return false
		}(); anyActive && r.Intn(2) == 0 {
			op = 2
		}
		var k int
		if op == 1 {
			for {
				k = r.Intn(n) + 1
				if !active[k] {
					break
				}
			}
			active[k] = true
		} else {
			for {
				k = r.Intn(n) + 1
				if active[k] {
					break
				}
			}
			active[k] = false
		}
		sb.WriteString(fmt.Sprintf("%d %d\n", op, k))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		input := genCase(r)
		want := solve1580C(input)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, input)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i, want, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All 100 tests passed")
}
