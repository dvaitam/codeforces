package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func buildOracle() string {
	dir, _ := os.Getwd()
	oracle := filepath.Join(dir, "oracleB_bin")
	cmd := exec.Command("go", "build", "-o", oracle, "1292B.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		panic(fmt.Sprintf("failed to build oracle: %v\n%s", err, out))
	}
	return oracle
}

func runProg(bin, input string) (string, error) {
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
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func genCase() string {
	x0 := rand.Int63n(5) + 1
	y0 := rand.Int63n(5) + 1
	ax := rand.Int63n(3) + 1
	ay := rand.Int63n(3) + 1
	bx := rand.Int63n(4)
	by := rand.Int63n(4)
	xs := rand.Int63n(15)
	ys := rand.Int63n(15)
	t := rand.Int63n(30) + 1
	return fmt.Sprintf("%d %d %d %d %d %d\n%d %d %d\n", x0, y0, ax, ay, bx, by, xs, ys, t)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}
	rand.Seed(time.Now().UnixNano())
	bin := os.Args[1]
	oracle := buildOracle()
	defer os.Remove(oracle)
	for i := 0; i < 100; i++ {
		input := genCase()
		exp, err := runProg(oracle, input)
		if err != nil {
			fmt.Printf("oracle failed on case %d: %v\n", i+1, err)
			return
		}
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Printf("binary failed on case %d: %v\n", i+1, err)
			return
		}
		if exp != got {
			fmt.Printf("mismatch on case %d\ninput:\n%sexpected:%s\nactual:%s\n", i+1, input, exp, got)
			return
		}
	}
	fmt.Println("all tests passed")
}
