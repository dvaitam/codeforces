package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func runProg(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func runRef(input string) (string, error) {
	_, self, _, _ := runtime.Caller(0)
	dir := filepath.Dir(self)
	ref := filepath.Join(dir, "1548D1.go")
	return runProg(ref, input)
}

func collinear(a, b, c [2]int) bool {
	return (b[0]-a[0])*(c[1]-a[1]) == (c[0]-a[0])*(b[1]-a[1])
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(4) + 3
	points := make([][2]int, 0, n)
	for len(points) < n {
		x := rng.Intn(10) * 2
		y := rng.Intn(10) * 2
		dup := false
		for _, p := range points {
			if p[0] == x && p[1] == y {
				dup = true
				break
			}
		}
		if dup {
			continue
		}
		good := true
		for i := 0; i < len(points) && good; i++ {
			for j := i + 1; j < len(points) && good; j++ {
				if collinear(points[i], points[j], [2]int{x, y}) {
					good = false
				}
			}
		}
		if !good {
			continue
		}
		points = append(points, [2]int{x, y})
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for _, p := range points {
		fmt.Fprintf(&sb, "%d %d\n", p[0], p[1])
	}
	return sb.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	if bin == "--" && len(os.Args) >= 3 {
		bin = os.Args[2]
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in := genCase(rng)
		expect, err := runRef(in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runProg(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("test %d failed\ninput:\n%sexpected:%s\nactual:%s\n", i+1, in, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
