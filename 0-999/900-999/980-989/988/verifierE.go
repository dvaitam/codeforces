package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func minMoves(s string, x, y byte) int {
	arr := []byte(s)
	n := len(arr)
	const INF = int(1e9)
	best := INF
	for j := 0; j < n; j++ {
		if arr[j] != y {
			continue
		}
		costY := n - 1 - j
		arr1 := append(append([]byte{}, arr[:j]...), arr[j+1:]...)
		arr1 = append(arr1, y)
		for i := 0; i < n-1; i++ {
			if arr1[i] != x {
				continue
			}
			costX := n - 2 - i
			arr2 := append(append([]byte{}, arr1[:i]...), arr1[i+1:n-1]...)
			arr2 = append(arr2, x, y)
			idx := -1
			for k := 0; k < n-2; k++ {
				if arr2[k] != '0' {
					idx = k
					break
				}
			}
			if idx == -1 {
				if n == 2 {
					if arr2[0] == '0' {
						continue
					}
					idx = 0
				} else {
					continue
				}
			}
			total := costY + costX + idx
			if total < best {
				best = total
			}
		}
	}
	return best
}

func solve(s string) string {
	const INF = int(1e9)
	ans := INF
	pairs := [][2]byte{{'0', '0'}, {'2', '5'}, {'5', '0'}, {'7', '5'}}
	for _, p := range pairs {
		mv := minMoves(s, p[0], p[1])
		if mv < ans {
			ans = mv
		}
	}
	if ans == INF {
		return "-1\n"
	}
	return fmt.Sprintf("%d\n", ans)
}

func runCase(bin string, s string) error {
	input := s + "\n"
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	expected := strings.TrimSpace(solve(s))
	got := strings.TrimSpace(out.String())
	if expected != got {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases := []string{"00", "72", "5071"}
	rng := rand.New(rand.NewSource(4))
	for i := 0; i < 100; i++ {
		length := rng.Intn(18) + 1
		b := make([]byte, length)
		for j := 0; j < length; j++ {
			b[j] = byte('0' + rng.Intn(10))
		}
		if b[0] == '0' {
			b[0] = '1'
		}
		cases = append(cases, string(b))
	}

	for i, s := range cases {
		if err := runCase(bin, s); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput: %s\n", i+1, err, s)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
