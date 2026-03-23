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

// Embedded correct solver for 13D (CF-accepted)

type Point64 struct {
	x  int64
	y  int64
	sx int64
}

func orient64(a, b, c Point64) int64 {
	return (b.x-a.x)*(c.y-a.y) - (b.y-a.y)*(c.x-a.x)
}

func solveD(input string) string {
	data := []byte(input)
	idx := 0
	nextInt := func() int64 {
		for idx < len(data) && (data[idx] == ' ' || data[idx] == '\n' || data[idx] == '\r' || data[idx] == '\t') {
			idx++
		}
		sign := int64(1)
		if idx < len(data) && data[idx] == '-' {
			sign = -1
			idx++
		}
		var v int64
		for idx < len(data) && data[idx] >= '0' && data[idx] <= '9' {
			v = v*10 + int64(data[idx]-'0')
			idx++
		}
		return sign * v
	}

	n := int(nextInt())
	m := int(nextInt())

	const K int64 = 2000000001

	red := make([]Point64, n)
	for i := 0; i < n; i++ {
		x := nextInt()
		y := nextInt()
		red[i] = Point64{x: x, y: y, sx: x + K*y}
	}

	blue := make([]Point64, m)
	for i := 0; i < m; i++ {
		x := nextInt()
		y := nextInt()
		blue[i] = Point64{x: x, y: y, sx: x + K*y}
	}

	sort.Slice(red, func(i, j int) bool {
		return red[i].sx < red[j].sx
	})

	low := make([][]int, n)
	for i := 0; i < n; i++ {
		low[i] = make([]int, n)
	}

	for i := 0; i < n; i++ {
		pi := red[i]
		si := pi.sx
		for j := i + 1; j < n; j++ {
			pj := red[j]
			sj := pj.sx
			cnt := 0
			for t := 0; t < m; t++ {
				b := blue[t]
				if b.sx > si && b.sx < sj && orient64(pi, pj, b) < 0 {
					cnt++
				}
			}
			low[i][j] = cnt
		}
	}

	var ans int64
	for i := 0; i < n; i++ {
		pi := red[i]
		lowi := low[i]
		for j := i + 1; j < n; j++ {
			pj := red[j]
			lowj := low[j]
			lij := lowi[j]
			for k := j + 1; k < n; k++ {
				var inside int
				if orient64(pi, pj, red[k]) > 0 {
					inside = lij + lowj[k] - lowi[k]
				} else {
					inside = lowi[k] - lij - lowj[k]
				}
				if inside == 0 {
					ans++
				}
			}
		}
	}

	return fmt.Sprintf("%d", ans)
}

type PointSmall struct{ x, y int }

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

func generateCaseD(rng *rand.Rand) string {
	n := rng.Intn(4) + 3
	m := rng.Intn(4)
	reds := make([]PointSmall, n)
	blues := make([]PointSmall, m)
	for i := range reds {
		reds[i] = PointSmall{rng.Intn(11) - 5, rng.Intn(11) - 5}
	}
	for i := range blues {
		blues[i] = PointSmall{rng.Intn(11) - 5, rng.Intn(11) - 5}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for _, p := range reds {
		sb.WriteString(fmt.Sprintf("%d %d\n", p.x, p.y))
	}
	for _, p := range blues {
		sb.WriteString(fmt.Sprintf("%d %d\n", p.x, p.y))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in := generateCaseD(rng)
		exp := solveD(in)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		got := strings.TrimSpace(out)
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
