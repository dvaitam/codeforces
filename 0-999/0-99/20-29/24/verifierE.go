package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

const INF = 1e20

type Particle struct {
	x, v float64
}

func expected(particles []Particle) float64 {
	n := len(particles)
	sort.Slice(particles, func(i, j int) bool { return particles[i].x < particles[j].x })
	
	ans := INF
	for i := 0; i < n; {
		j := i
		for j < n && particles[j].v > 0 {
			j++
		}
		k := j
		for k < n && particles[k].v < 0 {
			k++
		}
		
		if i == j || j == k {
			i = k
			continue
		}
		
		lo, hi := 0.0, 1e9
		for rep := 0; rep < 60; rep++ {
			mid := (lo + hi) / 2.0
			maxLeft := -INF
			for idx := i; idx < j; idx++ {
				pos := particles[idx].x + particles[idx].v*mid
				if pos > maxLeft {
					maxLeft = pos
				}
			}
			ok := true
			for idx := j; idx < k; idx++ {
				pos := particles[idx].x + particles[idx].v*mid
				if pos <= maxLeft {
					ok = false
					break
				}
			}
			if ok {
				lo = mid
			} else {
				hi = mid
			}
		}
		
		if hi < ans {
			ans = hi
		}
		i = k
	}
	
	if ans >= INF {
		return -1.0
	}
	return ans
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	
	for i := 0; i < 100; i++ {
		n := rng.Intn(100) + 1
		particles := make([]Particle, n)
		
		used := make(map[int]bool)
		for j := 0; j < n; j++ {
			x := rng.Intn(2001) - 1000
			for used[x] {
				x = rng.Intn(2001) - 1000
			}
			used[x] = true
			
			v := rng.Intn(1999) + 1
			if rng.Intn(2) == 0 {
				v = -v
			}
			
			particles[j] = Particle{float64(x), float64(v)}
		}
		
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d\n", n))
		for _, p := range particles {
			input.WriteString(fmt.Sprintf("%d %d\n", int(p.x), int(p.v)))
		}
		
		expectedOut := expected(particles)
		got, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input.String())
			os.Exit(1)
		}
		
		if expectedOut == -1.0 {
			if got != "-1" {
				fmt.Fprintf(os.Stderr, "case %d failed: expected -1 got %s\ninput:\n%s", i+1, got, input.String())
				os.Exit(1)
			}
		} else {
			gotFloat, parseErr := strconv.ParseFloat(got, 64)
			if parseErr != nil {
				fmt.Fprintf(os.Stderr, "case %d failed: cannot parse output: %v\ninput:\n%s", i+1, parseErr, input.String())
				os.Exit(1)
			}
			
			relErr := math.Abs(gotFloat-expectedOut) / math.Max(1.0, math.Abs(expectedOut))
			if relErr > 1e-8 && math.Abs(gotFloat-expectedOut) > 1e-8 {
				fmt.Fprintf(os.Stderr, "case %d failed: expected %.12f got %.12f (rel err %.12f)\ninput:\n%s", i+1, expectedOut, gotFloat, relErr, input.String())
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}