package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func buildCube(n int) [][][]int {
	a := make([][][]int, n)
	for x := 0; x < n; x++ {
		a[x] = make([][]int, n)
		for y := 0; y < n; y++ {
			a[x][y] = make([]int, n)
		}
	}
	cnt := 1
	idx := 0
	for x := 0; x < n; x++ {
		if x%2 == 0 {
			for y := 0; y < n; y++ {
				if idx%2 == 0 {
					for z := 0; z < n; z++ {
						a[x][y][z] = cnt
						cnt++
					}
				} else {
					for z := n - 1; z >= 0; z-- {
						a[x][y][z] = cnt
						cnt++
					}
				}
				idx++
			}
		} else {
			for y := n - 1; y >= 0; y-- {
				if idx%2 == 0 {
					for z := 0; z < n; z++ {
						a[x][y][z] = cnt
						cnt++
					}
				} else {
					for z := n - 1; z >= 0; z-- {
						a[x][y][z] = cnt
						cnt++
					}
				}
				idx++
			}
		}
	}
	return a
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(4) + 1
	cube := buildCube(n)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for z := 0; z < n; z++ {
		for x := 0; x < n; x++ {
			for y := 0; y < n; y++ {
				if y > 0 {
					sb.WriteByte(' ')
				}
				fmt.Fprintf(&sb, "%d", cube[x][y][z])
			}
			sb.WriteByte('\n')
		}
		if z != n-1 {
			sb.WriteByte('\n')
		}
	}
	return fmt.Sprintf("%d\n", n), sb.String()
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	if strings.TrimSpace(out.String()) != strings.TrimSpace(expected) {
		return fmt.Errorf("output mismatch")
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
