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

const MX = 300000

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func expected(arr []int) string {
	n := len(arr)
	have := make([][]bool, 10)
	for i := range have {
		have[i] = make([]bool, MX+1)
	}
	for _, v := range arr {
		if v == 1 {
			return "1\n"
		}
		have[0][v] = true
	}
	g := arr[0]
	for i := 1; i < n; i++ {
		g = gcd(g, arr[i])
	}
	if g > 1 {
		return "-1\n"
	}
	for i := 2; i <= MX; i++ {
		if !have[0][i] {
			continue
		}
		for j := i + i; j <= MX; j += i {
			have[0][j] = false
		}
	}
	ini := make([]int, 0)
	arrs := make([]int, 0)
	for i := 1; i <= MX; i++ {
		if have[0][i] {
			ini = append(ini, i)
		}
	}
	arrs = append(arrs, ini...)
	for x := 1; x < 10; x++ {
		nextArr := make([]int, 0)
		seen := have[x]
		prevArr := arrs
		for _, v1 := range prevArr {
			for _, v2 := range ini {
				g := gcd(v1, v2)
				if g == 1 {
					return fmt.Sprintf("%d\n", x+1)
				}
				if !seen[g] {
					seen[g] = true
					nextArr = append(nextArr, g)
				}
			}
		}
		arrs = nextArr
	}
	return "-1\n"
}

func runCase(bin, input, want string) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if strings.TrimSpace(want) != got {
		return fmt.Errorf("expected %q got %q", want, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	tests := [][]int{{2, 4, 6}, {6, 10, 15}}
	for i := 0; i < 100; i++ {
		n := rng.Intn(6) + 1
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = rng.Intn(50) + 1
		}
		tests = append(tests, arr)
	}

	for idx, tc := range tests {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", len(tc)))
		for i, v := range tc {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		want := expected(tc)
		if err := runCase(bin, sb.String(), strings.TrimSpace(want)); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", idx+1, err, sb.String())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
