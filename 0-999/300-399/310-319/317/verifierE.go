package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func runRef(input string) (string, error) {
	cmd := exec.Command("go", "run", "317E.go")
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rand.Seed(5)
	for i := 0; i < 100; i++ {
		vx := rand.Intn(7) - 3
		vy := rand.Intn(7) - 3
		sx := rand.Intn(7) - 3
		sy := rand.Intn(7) - 3
		for vx == sx && vy == sy {
			sx = rand.Intn(7) - 3
			sy = rand.Intn(7) - 3
		}
		m := rand.Intn(5)
		trees := make(map[[2]int]bool)
		for len(trees) < m {
			tx := rand.Intn(7) - 3
			ty := rand.Intn(7) - 3
			if (tx == vx && ty == vy) || (tx == sx && ty == sy) {
				continue
			}
			trees[[2]int{tx, ty}] = true
		}
		input := fmt.Sprintf("%d %d %d %d %d\n", vx, vy, sx, sy, m)
		for t := range trees {
			input += fmt.Sprintf("%d %d\n", t[0], t[1])
		}
		exp, err := runRef(input)
		if err != nil {
			fmt.Println("reference run error:", err)
			return
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("binary run error on test %d: %v\n", i+1, err)
			return
		}
		if exp != got {
			fmt.Printf("mismatch on test %d\ninput:\n%sexpected:\n%s\n got:\n%s\n", i+1, input, exp, got)
			return
		}
	}
	fmt.Println("all tests passed")
}
