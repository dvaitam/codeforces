package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

type stateB struct {
	n, k int
	req  []int
}

func solveB(n, k int, req []int) []string {
	p := make([][2]int, k+1)
	for i := 1; i <= k; i++ {
		p[i][0] = 0
		p[i][1] = -1
	}
	xc := (k + 1) / 2
	res := make([]string, n)
	for i := 0; i < n; i++ {
		m := req[i]
		if m > k {
			res[i] = "-1"
			continue
		}
		best := math.MaxInt32
		bx, bl, br := 0, 0, 0
		for x := 1; x <= k; x++ {
			if p[x][0] > p[x][1] {
				d := abs(x-xc)*m + (m/2)*((m+1)/2)
				if d < best {
					best = d
					bx = x
					bl = (k-m)/2 + 1
					br = bl + m - 1
				}
			} else {
				if p[x][0] > m {
					d := abs(x-xc)*m + (xc-p[x][0])*m + m*(m+1)/2
					if d < best {
						best = d
						bx = x
						bl = p[x][0] - m
						br = p[x][0] - 1
					}
				}
				if p[x][1] <= k-m {
					d := abs(x-xc)*m + (p[x][1]-xc)*m + m*(m+1)/2
					if d < best {
						best = d
						bx = x
						bl = p[x][1] + 1
						br = p[x][1] + m
					}
				}
			}
		}
		if bx == 0 {
			res[i] = "-1"
			continue
		}
		res[i] = fmt.Sprintf("%d %d %d", bx, bl, br)
		if p[bx][0] > p[bx][1] {
			p[bx][0] = bl
			p[bx][1] = br
		} else if p[bx][0] == br+1 {
			p[bx][0] = bl
		} else if p[bx][1] == bl-1 {
			p[bx][1] = br
		}
	}
	return res
}

func generateCaseB(rng *rand.Rand) (string, []string) {
	k := rng.Intn(5)*2 + 3 // odd between 3..11
	n := rng.Intn(5) + 1
	req := make([]int, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i := 0; i < n; i++ {
		req[i] = rng.Intn(k*2) + 1
		sb.WriteString(strconv.Itoa(req[i]))
		if i+1 < n {
			sb.WriteByte(' ')
		}
	}
	sb.WriteByte('\n')
	return sb.String(), solveB(n, k, req)
}

func runCaseB(bin, input string, expected []string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	scanner := bufio.NewScanner(bytes.NewReader(out.Bytes()))
	for i := 0; i < len(expected); i++ {
		if !scanner.Scan() {
			return fmt.Errorf("not enough output for case, expected %d lines", len(expected))
		}
		line := strings.TrimSpace(scanner.Text())
		if line != expected[i] {
			return fmt.Errorf("line %d: expected '%s' got '%s'", i+1, expected[i], line)
		}
	}
	if scanner.Scan() {
		return fmt.Errorf("extra output: %s", scanner.Text())
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		in, exp := generateCaseB(rng)
		if err := runCaseB(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
