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

func canSplit(a []int, pos1, pos2 int) bool {
	n := len(a)
	last1, last2 := -1, -1
	peak1, peak2 := a[pos1], a[pos2]
	for i := 0; i < n; i++ {
		val := a[i]
		if i == pos1 {
			last1 = val
			continue
		}
		if i == pos2 {
			last2 = val
			continue
		}
		ok1 := false
		if i < pos1 {
			if val > last1 && val < peak1 {
				ok1 = true
			}
		} else if i > pos1 {
			if val < last1 {
				ok1 = true
			}
		}
		ok2 := false
		if i < pos2 {
			if val > last2 && val < peak2 {
				ok2 = true
			}
		} else if i > pos2 {
			if val < last2 {
				ok2 = true
			}
		}
		if ok1 && !ok2 {
			last1 = val
		} else if ok2 && !ok1 {
			last2 = val
		} else if ok1 && ok2 {
			if last1 <= last2 {
				last1 = val
			} else {
				last2 = val
			}
		} else {
			return false
		}
	}
	return true
}

func solveCase(a []int) string {
	pos := make(map[int]int)
	maxVal := -1
	for i, v := range a {
		pos[v] = i
		if v > maxVal {
			maxVal = v
		}
	}
	pairs := make(map[[2]int]struct{})
	posMax := pos[maxVal]
	for val, p := range pos {
		if val == maxVal {
			continue
		}
		if canSplit(a, posMax, p) {
			pair := [2]int{val, maxVal}
			if pair[0] > pair[1] {
				pair[0], pair[1] = pair[1], pair[0]
			}
			pairs[pair] = struct{}{}
		}
	}
	return fmt.Sprintf("%d\n", len(pairs))
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(7) + 2
	arr := make([]int, n)
	used := make(map[int]bool)
	for i := 0; i < n; i++ {
		for {
			v := rng.Intn(50) + 1
			if !used[v] {
				used[v] = true
				arr[i] = v
				break
			}
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", arr[i])
	}
	sb.WriteByte('\n')
	input := sb.String()
	expected := solveCase(arr)
	return input, expected
}

func runCase(bin, input, expected string) error {
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
	outStr := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if outStr != exp {
		return fmt.Errorf("expected %q got %q", exp, outStr)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
