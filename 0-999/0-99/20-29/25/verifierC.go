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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func expected(n int, dist [][]int, roads [][]int) []int64 {
	d := make([][]int, n)
	for i := 0; i < n; i++ {
		d[i] = make([]int, n)
		copy(d[i], dist[i])
	}
	
	result := make([]int64, len(roads))
	
	for idx, road := range roads {
		u, v, c := road[0]-1, road[1]-1, road[2]
		
		du := make([]int, n)
		dv := make([]int, n)
		for i := 0; i < n; i++ {
			du[i] = d[i][u]
			dv[i] = d[i][v]
		}
		
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				alt := min(du[i]+c+dv[j], dv[i]+c+du[j])
				d[i][j] = min(d[i][j], alt)
			}
		}
		
		var sum int64
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				sum += int64(d[i][j])
			}
		}
		result[idx] = sum
	}
	
	return result
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

func generateDistanceMatrix(n int, rng *rand.Rand) [][]int {
	d := make([][]int, n)
	for i := 0; i < n; i++ {
		d[i] = make([]int, n)
		for j := 0; j < n; j++ {
			if i == j {
				d[i][j] = 0
			} else {
				d[i][j] = rng.Intn(100) + 1
			}
		}
	}
	
	for k := 0; k < n; k++ {
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				d[i][j] = min(d[i][j], d[i][k]+d[k][j])
			}
		}
	}
	
	return d
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	
	for i := 0; i < 100; i++ {
		n := rng.Intn(8) + 2
		k := rng.Intn(10) + 1
		
		d := generateDistanceMatrix(n, rng)
		roads := make([][]int, k)
		for j := 0; j < k; j++ {
			u := rng.Intn(n) + 1
			v := rng.Intn(n) + 1
			for u == v {
				v = rng.Intn(n) + 1
			}
			c := rng.Intn(100) + 1
			roads[j] = []int{u, v, c}
		}
		
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d\n", n))
		for row := 0; row < n; row++ {
			for col := 0; col < n; col++ {
				if col > 0 {
					input.WriteString(" ")
				}
				input.WriteString(strconv.Itoa(d[row][col]))
			}
			input.WriteString("\n")
		}
		input.WriteString(fmt.Sprintf("%d\n", k))
		for j := 0; j < k; j++ {
			input.WriteString(fmt.Sprintf("%d %d %d\n", roads[j][0], roads[j][1], roads[j][2]))
		}
		
		expectedOut := expected(n, d, roads)
		got, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input.String())
			os.Exit(1)
		}
		
		gotParts := strings.Fields(got)
		if len(gotParts) != len(expectedOut) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d values, got %d\ninput:\n%s", i+1, len(expectedOut), len(gotParts), input.String())
			os.Exit(1)
		}
		
		for j, expectedVal := range expectedOut {
			gotVal, parseErr := strconv.ParseInt(gotParts[j], 10, 64)
			if parseErr != nil {
				fmt.Fprintf(os.Stderr, "case %d failed: cannot parse output value %s\ninput:\n%s", i+1, gotParts[j], input.String())
				os.Exit(1)
			}
			if gotVal != expectedVal {
				fmt.Fprintf(os.Stderr, "case %d failed: at position %d expected %d got %d\ninput:\n%s", i+1, j, expectedVal, gotVal, input.String())
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}