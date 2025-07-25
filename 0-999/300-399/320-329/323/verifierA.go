package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func checkPainting(k int, out string) bool {
	out = strings.TrimSpace(out)
	if k%2 == 1 {
		return out == "-1"
	}
	letters := make([]byte, 0, k*k*k)
	for i := 0; i < len(out); i++ {
		c := out[i]
		if c == 'w' || c == 'b' {
			letters = append(letters, c)
		}
	}
	if len(letters) != k*k*k {
		return false
	}
	idx := 0
	cube := make([][][]byte, k)
	for z := 0; z < k; z++ {
		layer := make([][]byte, k)
		for y := 0; y < k; y++ {
			row := make([]byte, k)
			copy(row, letters[idx:idx+k])
			idx += k
			layer[y] = row
		}
		cube[z] = layer
	}
	dirs := [][3]int{{1, 0, 0}, {-1, 0, 0}, {0, 1, 0}, {0, -1, 0}, {0, 0, 1}, {0, 0, -1}}
	for x := 0; x < k; x++ {
		for y := 0; y < k; y++ {
			for z := 0; z < k; z++ {
				c := cube[x][y][z]
				cnt := 0
				for _, d := range dirs {
					nx, ny, nz := x+d[0], y+d[1], z+d[2]
					if nx >= 0 && nx < k && ny >= 0 && ny < k && nz >= 0 && nz < k {
						if cube[nx][ny][nz] == c {
							cnt++
						}
					}
				}
				if cnt != 2 {
					return false
				}
			}
		}
	}
	return true
}

func runCase(binary string, k int) error {
	in := fmt.Sprintf("%d\n", k)
	cmd := exec.Command(binary)
	cmd.Stdin = strings.NewReader(in)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v", err)
	}
	if !checkPainting(k, buf.String()) {
		return fmt.Errorf("wrong answer")
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for k := 1; k <= 100; k++ {
		if err := runCase(bin, k); err != nil {
			fmt.Printf("test %d failed: %v\n", k, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
