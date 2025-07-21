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

func expected(n int, j int64, x0, y0 int64, Ax, Ay []int64) (int64, int64) {
	var Cx, Cy int64
	for t := 0; t < n; t++ {
		if t%2 == 0 {
			Cx += Ax[t]
			Cy += Ay[t]
		} else {
			Cx -= Ax[t]
			Cy -= Ay[t]
		}
	}
	
	q := j / int64(n)
	r := int(j % int64(n))
	
	var mx, my int64
	if q%2 == 0 {
		mx, my = x0, y0
	} else {
		mx = -x0 + 2*Cx
		my = -y0 + 2*Cy
	}
	
	for t := 0; t < r; t++ {
		mx = -mx + 2*Ax[t]
		my = -my + 2*Ay[t]
	}
	
	return mx, my
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	
	for i := 0; i < 100; i++ {
		n := 2*rng.Intn(50) + 1
		j := int64(rng.Intn(1000000) + 1)
		
		x0 := int64(rng.Intn(2001) - 1000)
		y0 := int64(rng.Intn(2001) - 1000)
		
		Ax := make([]int64, n)
		Ay := make([]int64, n)
		for k := 0; k < n; k++ {
			Ax[k] = int64(rng.Intn(2001) - 1000)
			Ay[k] = int64(rng.Intn(2001) - 1000)
		}
		
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d %d\n", n, j))
		input.WriteString(fmt.Sprintf("%d %d\n", x0, y0))
		for k := 0; k < n; k++ {
			input.WriteString(fmt.Sprintf("%d %d\n", Ax[k], Ay[k]))
		}
		
		expectedX, expectedY := expected(n, j, x0, y0, Ax, Ay)
		got, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input.String())
			os.Exit(1)
		}
		
		parts := strings.Fields(got)
		if len(parts) != 2 {
			fmt.Fprintf(os.Stderr, "case %d failed: expected 2 coordinates, got %d\ninput:\n%s", i+1, len(parts), input.String())
			os.Exit(1)
		}
		
		gotX, parseErr := strconv.ParseInt(parts[0], 10, 64)
		if parseErr != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: cannot parse X coordinate: %v\ninput:\n%s", i+1, parseErr, input.String())
			os.Exit(1)
		}
		
		gotY, parseErr := strconv.ParseInt(parts[1], 10, 64)
		if parseErr != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: cannot parse Y coordinate: %v\ninput:\n%s", i+1, parseErr, input.String())
			os.Exit(1)
		}
		
		if gotX != expectedX || gotY != expectedY {
			fmt.Fprintf(os.Stderr, "case %d failed: expected (%d, %d) got (%d, %d)\ninput:\n%s", i+1, expectedX, expectedY, gotX, gotY, input.String())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}