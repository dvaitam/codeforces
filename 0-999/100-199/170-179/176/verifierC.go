package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func runBinary(bin, input string) (string, error) {
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
	err := cmd.Run()
	if err != nil {
		return out.String(), fmt.Errorf("%v: %s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expected(n, m, x1, y1, x2, y2 int) string {
	d := abs(x1-x2) + abs(y1-y2)
	if n == 1 || m == 1 {
		if d <= 4 {
			return "First"
		}
		return "Second"
	}
	return "First"
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	if bin == "--" && len(os.Args) >= 3 {
		bin = os.Args[2]
	}
	rand.Seed(2)
	for tcase := 0; tcase < 100; tcase++ {
		n := rand.Intn(4) + 1
		m := rand.Intn(4) + 1
		x1 := rand.Intn(n) + 1
		y1 := rand.Intn(m) + 1
		x2 := rand.Intn(n) + 1
		y2 := rand.Intn(m) + 1
		if x1 == x2 && y1 == y2 {
			x2 = (x2 % n) + 1
		}
		input := fmt.Sprintf("%d %d %d %d %d %d\n", n, m, x1, y1, x2, y2)
		exp := expected(n, m, x1, y1, x2, y2)
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", tcase+1, err)
			return
		}
		out = strings.TrimSpace(out)
		if out != exp {
			fmt.Printf("test %d failed: expected %s got %s\ninput:%s", tcase+1, exp, out, input)
			return
		}
	}
	fmt.Println("All tests passed.")
}
