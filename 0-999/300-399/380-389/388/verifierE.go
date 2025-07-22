package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type meteor struct {
	vx, vy int64
	ox, oy int64
	den    int64
}

func abs64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func solve(input string) (int, error) {
	scan := bufio.NewScanner(strings.NewReader(input))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		return 0, fmt.Errorf("empty input")
	}
	n, err := strconv.Atoi(scan.Text())
	if err != nil {
		return 0, err
	}
	meteors := make([]meteor, n)
	for i := 0; i < n; i++ {
		vals := make([]int64, 6)
		for j := 0; j < 6; j++ {
			if !scan.Scan() {
				return 0, fmt.Errorf("not enough values")
			}
			v, err := strconv.ParseInt(scan.Text(), 10, 64)
			if err != nil {
				return 0, err
			}
			vals[j] = v
		}
		t1, x1, y1, t2, x2, y2 := vals[0], vals[1], vals[2], vals[3], vals[4], vals[5]
		di := t2 - t1
		meteors[i].den = di
		meteors[i].vx = x2 - x1
		meteors[i].vy = y2 - y1
		meteors[i].ox = x1*t2 - x2*t1
		meteors[i].oy = y1*t2 - y2*t1
	}
	adj := make([][]bool, n)
	for i := range adj {
		adj[i] = make([]bool, n)
	}
	concurrency := 1
	for i := 0; i < n; i++ {
		counts := make(map[[2]int64]int)
		for j := i + 1; j < n; j++ {
			dvx := meteors[i].vx*meteors[j].den - meteors[j].vx*meteors[i].den
			dvy := meteors[i].vy*meteors[j].den - meteors[j].vy*meteors[i].den
			dox := meteors[j].ox*meteors[i].den - meteors[i].ox*meteors[j].den
			doy := meteors[j].oy*meteors[i].den - meteors[i].oy*meteors[j].den
			if dvx == 0 && dvy == 0 {
				continue
			}
			if dvx*doy != dvy*dox {
				continue
			}
			adj[i][j] = true
			adj[j][i] = true
			var tn, td int64
			if dvx != 0 {
				tn = dox
				td = dvx
			} else {
				tn = doy
				td = dvy
			}
			if td < 0 {
				tn = -tn
				td = -td
			}
			g := gcd(abs64(tn), td)
			tn /= g
			td /= g
			key := [2]int64{tn, td}
			counts[key]++
		}
		for _, c := range counts {
			if c+1 > concurrency {
				concurrency = c + 1
			}
		}
	}
	if concurrency < 3 {
		found := false
		for i := 0; i < n && !found; i++ {
			for j := i + 1; j < n && !found; j++ {
				if !adj[i][j] {
					continue
				}
				for k := j + 1; k < n; k++ {
					if adj[i][k] && adj[j][k] {
						found = true
						break
					}
				}
			}
		}
		if found && concurrency < 3 {
			concurrency = 3
		}
	}
	if concurrency < 1 {
		concurrency = 1
	}
	return concurrency, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesE.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		input := line + "\n"
		// compute expected using internal solver
		expectedVal, err := solve(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "internal solver error on case %d: %v\n", idx, err)
			os.Exit(1)
		}
		expected := fmt.Sprintf("%d", expectedVal)

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != expected {
			fmt.Printf("test %d failed\nexpected: %s\n got: %s\n", idx, expected, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
