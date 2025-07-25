package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func runCandidate(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expected(h, w int, grid []string, queries [][4]int) []int {
	hPref := make([][]int, h+1)
	vPref := make([][]int, h+1)
	for i := 0; i <= h; i++ {
		hPref[i] = make([]int, w+1)
		vPref[i] = make([]int, w+1)
	}
	for i := 1; i <= h; i++ {
		for j := 1; j <= w; j++ {
			if j < w && grid[i-1][j-1] == '.' && grid[i-1][j] == '.' {
				hPref[i][j] = 1
			}
			if i < h && grid[i-1][j-1] == '.' && grid[i][j-1] == '.' {
				vPref[i][j] = 1
			}
			hPref[i][j] += hPref[i-1][j] + hPref[i][j-1] - hPref[i-1][j-1]
			vPref[i][j] += vPref[i-1][j] + vPref[i][j-1] - vPref[i-1][j-1]
		}
	}
	res := make([]int, len(queries))
	for idx, q := range queries {
		r1, c1, r2, c2 := q[0], q[1], q[2], q[3]
		horizontal := hPref[r2][c2-1] - hPref[r1-1][c2-1] - hPref[r2][c1-1] + hPref[r1-1][c1-1]
		vertical := vPref[r2-1][c2] - vPref[r1-1][c2] - vPref[r2-1][c1-1] + vPref[r1-1][c1-1]
		res[idx] = horizontal + vertical
	}
	return res
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(3))
	for t := 0; t < 100; t++ {
		h := rng.Intn(4) + 2
		w := rng.Intn(4) + 2
		grid := make([]string, h)
		for i := 0; i < h; i++ {
			var sb strings.Builder
			for j := 0; j < w; j++ {
				if rng.Intn(4) == 0 {
					sb.WriteByte('#')
				} else {
					sb.WriteByte('.')
				}
			}
			grid[i] = sb.String()
		}
		q := rng.Intn(5) + 1
		queries := make([][4]int, q)
		for i := 0; i < q; i++ {
			r1 := rng.Intn(h) + 1
			r2 := rng.Intn(h-r1+1) + r1
			c1 := rng.Intn(w) + 1
			c2 := rng.Intn(w-c1+1) + c1
			queries[i] = [4]int{r1, c1, r2, c2}
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", h, w)
		for _, row := range grid {
			sb.WriteString(row)
			sb.WriteByte('\n')
		}
		fmt.Fprintf(&sb, "%d\n", q)
		for _, qu := range queries {
			fmt.Fprintf(&sb, "%d %d %d %d\n", qu[0], qu[1], qu[2], qu[3])
		}
		input := sb.String()
		want := expected(h, w, grid, queries)
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", t+1, err)
			os.Exit(1)
		}
		fields := strings.Fields(strings.TrimSpace(out))
		if len(fields) != len(want) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d numbers got %d\n", t+1, len(want), len(fields))
			os.Exit(1)
		}
		for i, f := range fields {
			val, err := strconv.Atoi(f)
			if err != nil || val != want[i] {
				fmt.Fprintf(os.Stderr, "case %d failed: expected %v got %v\n", t+1, want, fields)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
