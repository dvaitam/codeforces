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

func solveA(n, q int, ops [][2]int) string {
	notifications := make([]int, 0)
	read := make([]bool, q+5)
	appQueues := make([][]int, n+1)
	appPtr := make([]int, n+1)
	globalPtr := 0
	unread := 0
	var sb strings.Builder
	for _, op := range ops {
		t := op[0]
		x := op[1]
		switch t {
		case 1:
			id := len(notifications)
			notifications = append(notifications, x)
			if len(read) <= id {
				read = append(read, false)
			}
			appQueues[x] = append(appQueues[x], id)
			unread++
		case 2:
			qlist := appQueues[x]
			ptr := appPtr[x]
			for ptr < len(qlist) {
				id := qlist[ptr]
				if !read[id] {
					read[id] = true
					unread--
				}
				ptr++
			}
			appPtr[x] = ptr
		case 3:
			tVal := x
			for globalPtr < tVal {
				if !read[globalPtr] {
					read[globalPtr] = true
					unread--
				}
				globalPtr++
			}
		}
		sb.WriteString(fmt.Sprintf("%d\n", unread))
	}
	return sb.String()
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	q := rng.Intn(20) + 1
	ops := make([][2]int, q)
	type1Count := 0
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, q))
	for i := 0; i < q; i++ {
		typ := rng.Intn(3) + 1
		var x int
		switch typ {
		case 1:
			x = rng.Intn(n) + 1
			type1Count++
		case 2:
			x = rng.Intn(n) + 1
		case 3:
			if type1Count == 0 {
				typ = 1
				x = rng.Intn(n) + 1
				type1Count++
			} else {
				x = rng.Intn(type1Count) + 1
			}
		}
		ops[i] = [2]int{typ, x}
		sb.WriteString(fmt.Sprintf("%d %d\n", typ, x))
	}
	expected := solveA(n, q, ops)
	return sb.String(), expected
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
	outStr := strings.ReplaceAll(strings.TrimSpace(out.String()), "\r", "")
	expStr := strings.ReplaceAll(strings.TrimSpace(expected), "\r", "")
	if outStr != expStr {
		return fmt.Errorf("expected %q got %q", expStr, outStr)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
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
