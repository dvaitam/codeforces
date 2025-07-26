package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type complaint struct{ x, y int }

type pair struct{ u, v int }

func exists(n, p, M int, comps []complaint, l, r []int, inter []pair) bool {
	// brute force over f and subset of stations
	for f := 1; f <= M; f++ {
		for mask := 0; mask < (1 << p); mask++ {
			valid := true
			// check power constraints
			for j := 0; j < p; j++ {
				if mask&(1<<j) != 0 {
					if f < l[j] || f > r[j] {
						valid = false
						break
					}
				}
			}
			if !valid {
				continue
			}
			// interference
			for _, pr := range inter {
				if mask&(1<<pr.u) != 0 && mask&(1<<pr.v) != 0 {
					valid = false
					break
				}
			}
			if !valid {
				continue
			}
			// complaints
			for _, c := range comps {
				if mask&(1<<c.x) == 0 && mask&(1<<c.y) == 0 {
					valid = false
					break
				}
			}
			if valid {
				return true
			}
		}
	}
	return false
}

func verifySolution(out string, n, p, M int, comps []complaint, l, r []int, inter []pair, expectExist bool) error {
	fields := strings.Fields(out)
	if !expectExist {
		if len(fields) != 1 || fields[0] != "-1" {
			return fmt.Errorf("expected -1")
		}
		return nil
	}
	if len(fields) < 2 {
		return fmt.Errorf("incomplete output")
	}
	var k, f int
	if _, err := fmt.Sscanf(fields[0], "%d", &k); err != nil {
		return fmt.Errorf("bad k")
	}
	if _, err := fmt.Sscanf(fields[1], "%d", &f); err != nil {
		return fmt.Errorf("bad f")
	}
	if len(fields) != 2+k {
		return fmt.Errorf("wrong number of station indices")
	}
	used := make(map[int]bool)
	for i := 0; i < k; i++ {
		var x int
		fmt.Sscanf(fields[2+i], "%d", &x)
		if x < 1 || x > p {
			return fmt.Errorf("station index out of range")
		}
		if used[x] {
			return fmt.Errorf("duplicate station")
		}
		used[x] = true
		if f < l[x-1] || f > r[x-1] {
			return fmt.Errorf("power constraint")
		}
	}
	// interference
	for _, pr := range inter {
		if used[pr.u+1] && used[pr.v+1] {
			return fmt.Errorf("interference")
		}
	}
	// complaints
	for _, c := range comps {
		if !used[c.x+1] && !used[c.y+1] {
			return fmt.Errorf("complaint not satisfied")
		}
	}
	return nil
}

func run(binary, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, binary)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]

	rand.Seed(47)
	const T = 20
	for tc := 0; tc < T; tc++ {
		n := rand.Intn(3) + 1 // complaints
		p := rand.Intn(3) + 2 // stations
		M := rand.Intn(4) + 2
		m := rand.Intn(2) + 1
		comps := make([]complaint, n)
		for i := 0; i < n; i++ {
			x := rand.Intn(p)
			y := rand.Intn(p)
			for y == x {
				y = rand.Intn(p)
			}
			if x > y {
				x, y = y, x
			}
			comps[i] = complaint{x, y}
		}
		l := make([]int, p)
		r := make([]int, p)
		for i := 0; i < p; i++ {
			a := rand.Intn(M) + 1
			b := rand.Intn(M) + 1
			if a > b {
				a, b = b, a
			}
			l[i] = a
			r[i] = b
		}
		inter := make([]pair, 0, m)
		seen := make(map[[2]int]bool)
		for len(inter) < m {
			u := rand.Intn(p)
			v := rand.Intn(p)
			if u == v {
				continue
			}
			if u > v {
				u, v = v, u
			}
			key := [2]int{u, v}
			if seen[key] {
				continue
			}
			seen[key] = true
			inter = append(inter, pair{u, v})
		}
		// build input string
		input := fmt.Sprintf("%d %d %d %d\n", p, n, M, m)
		for _, c := range comps {
			input += fmt.Sprintf("%d %d\n", c.x+1, c.y+1)
		}
		for i := 0; i < p; i++ {
			input += fmt.Sprintf("%d %d\n", l[i], r[i])
		}
		for _, pr := range inter {
			input += fmt.Sprintf("%d %d\n", pr.u+1, pr.v+1)
		}
		exist := exists(n, p, M, comps, l, r, inter)
		out, err := run(binary, input)
		if err != nil {
			fmt.Printf("test %d: execution failed: %v\ninput:%s\n", tc+1, err, input)
			os.Exit(1)
		}
		if err := verifySolution(out, n, p, M, comps, l, r, inter, exist); err != nil {
			fmt.Printf("test %d failed: %v\ninput:%s\noutput:%s\n", tc+1, err, input, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
