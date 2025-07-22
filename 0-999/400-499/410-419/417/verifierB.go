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

func expected(n int, pairs [][2]int) string {
	const maxK = 100000
	mex := make([]int, maxK+1)
	for i := 0; i < n; i++ {
		x := pairs[i][0]
		k := pairs[i][1]
		cur := mex[k]
		if x > cur {
			return "NO"
		}
		if x == cur {
			mex[k] = cur + 1
		}
	}
	return "YES"
}

func runCase(bin string, n int, pairs [][2]int) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", pairs[i][0], pairs[i][1]))
	}
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expect := expected(n, pairs)
	if strings.ToUpper(got) != expect {
		return fmt.Errorf("expected %s got %s", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(15) + 1
		pairs := make([][2]int, n)
		tmpMex := make(map[int]int)
		for j := 0; j < n; j++ {
			k := rng.Intn(10) + 1
			cur := tmpMex[k]
			var x int
			if rng.Intn(3) == 0 {
				x = cur + rng.Intn(3) + 1 // invalid
			} else {
				if cur > 0 {
					x = rng.Intn(cur + 1)
				} else {
					if rng.Intn(2) == 0 {
						x = 0
					} else {
						x = 1
					}
				}
				if x == cur {
					tmpMex[k] = cur + 1
				}
			}
			if x < 0 {
				x = 0
			}
			pairs[j] = [2]int{x, k}
		}
		if err := runCase(bin, n, pairs); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			fmt.Fprintf(os.Stderr, "input:\n%d\n", n)
			for _, p := range pairs {
				fmt.Fprintf(os.Stderr, "%d %d\n", p[0], p[1])
			}
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
