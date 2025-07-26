package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
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

func solveC(n int, x1, y1, x2, y2 float64, data [][4]float64) string {
	l, r := 0.0, 1e18
	for _, m := range data {
		x, y, vx, vy := m[0], m[1], m[2], m[3]
		u, v := 0.0, 1e18
		if vx == 0 {
			if x <= x1 || x >= x2 {
				return "-1"
			}
		} else {
			u = (x1 - x) / vx
			v = (x2 - x) / vx
			if u > v {
				u, v = v, u
			}
		}
		l = math.Max(l, u)
		r = math.Min(r, v)
		if vy == 0 {
			if y <= y1 || y >= y2 {
				return "-1"
			}
		} else {
			u = (y1 - y) / vy
			v = (y2 - y) / vy
			if u > v {
				u, v = v, u
			}
		}
		l = math.Max(l, u)
		r = math.Min(r, v)
	}
	if l >= r {
		return "-1"
	}
	return fmt.Sprintf("%.6f", l)
}

func genCase(rng *rand.Rand) (int, float64, float64, float64, float64, [][4]float64) {
	n := rng.Intn(5) + 1
	x1 := float64(rng.Intn(10))
	y1 := float64(rng.Intn(10))
	x2 := x1 + float64(rng.Intn(10)+1)
	y2 := y1 + float64(rng.Intn(10)+1)
	data := make([][4]float64, n)
	for i := 0; i < n; i++ {
		data[i][0] = float64(rng.Intn(15))
		data[i][1] = float64(rng.Intn(15))
		data[i][2] = float64(rng.Intn(5) - 2)
		data[i][3] = float64(rng.Intn(5) - 2)
		if data[i][2] == 0 && (data[i][0] <= x1 || data[i][0] >= x2) {
			data[i][0] = (x1 + x2) / 2
		}
		if data[i][3] == 0 && (data[i][1] <= y1 || data[i][1] >= y2) {
			data[i][1] = (y1 + y2) / 2
		}
	}
	return n, x1, y1, x2, y2, data
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, x1, y1, x2, y2, data := genCase(rng)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", n)
		fmt.Fprintf(&sb, "%d %d %d %d\n", int(x1), int(y1), int(x2), int(y2))
		for _, m := range data {
			fmt.Fprintf(&sb, "%d %d %d %d\n", int(m[0]), int(m[1]), int(m[2]), int(m[3]))
		}
		expect := solveC(n, x1, y1, x2, y2, data)
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		got = strings.TrimSpace(got)
		var gv float64
		fmt.Sscan(got, &gv)
		var ev float64
		fmt.Sscan(expect, &ev)
		if math.Abs(gv-ev) > 1e-4 {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected %s got %s\n", i+1, sb.String(), expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
