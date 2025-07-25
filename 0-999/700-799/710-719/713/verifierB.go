package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type rect [4]int

type testCase struct {
	n     int
	rects [2]rect
}

func overlap(a, b rect) bool {
	return !(a[2] < b[0] || b[2] < a[0] || a[3] < b[1] || b[3] < a[1])
}

func randomRect(rng *rand.Rand, n int) rect {
	x1 := rng.Intn(n) + 1
	x2 := rng.Intn(n) + 1
	if x1 > x2 {
		x1, x2 = x2, x1
	}
	y1 := rng.Intn(n) + 1
	y2 := rng.Intn(n) + 1
	if y1 > y2 {
		y1, y2 = y2, y1
	}
	return rect{x1, y1, x2, y2}
}

func generateCase(rng *rand.Rand) testCase {
	n := rng.Intn(20) + 2
	r1 := randomRect(rng, n)
	var r2 rect
	for {
		r2 = randomRect(rng, n)
		if !overlap(r1, r2) {
			break
		}
	}
	return testCase{n: n, rects: [2]rect{r1, r2}}
}

func runCase(bin string, tc testCase) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		return err
	}
	fmt.Fprintf(stdin, "%d\n", tc.n)
	reader := bufio.NewReader(stdout)
	queryCnt := 0
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if ctx.Err() == context.DeadlineExceeded {
				return fmt.Errorf("time limit")
			}
			return fmt.Errorf("read error: %v stderr:%s", err, stderr.String())
		}
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "?") {
			var x1, y1, x2, y2 int
			if _, err := fmt.Sscanf(line, "? %d %d %d %d", &x1, &y1, &x2, &y2); err != nil {
				return fmt.Errorf("bad query: %q", line)
			}
			cnt := 0
			for _, r := range tc.rects {
				if x1 <= r[0] && y1 <= r[1] && x2 >= r[2] && y2 >= r[3] {
					cnt++
				}
			}
			fmt.Fprintf(stdin, "%d\n", cnt)
			queryCnt++
			if queryCnt > 200 {
				return fmt.Errorf("too many queries")
			}
		} else if strings.HasPrefix(line, "!") {
			var a rect
			var b rect
			if _, err := fmt.Sscanf(line, "! %d %d %d %d %d %d %d %d", &a[0], &a[1], &a[2], &a[3], &b[0], &b[1], &b[2], &b[3]); err != nil {
				return fmt.Errorf("bad answer: %q", line)
			}
			ok := (a == tc.rects[0] && b == tc.rects[1]) || (a == tc.rects[1] && b == tc.rects[0])
			if !ok {
				return fmt.Errorf("wrong answer: expected %v %v got %v %v", tc.rects[0], tc.rects[1], a, b)
			}
			stdin.Close()
			return cmd.Wait()
		}
	}
}

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		if out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput(); err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCase(rng)
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\nrect1:%v rect2:%v n:%d\n", i+1, err, tc.rects[0], tc.rects[1], tc.n)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
