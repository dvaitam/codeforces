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
	oracle := filepath.Join(dir, "oracleD_bin")
	cmd := exec.Command("go", "build", "-o", oracle, "1292D.go")
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
	n := rand.Intn(7) + 1
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		if i > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(fmt.Sprintf("%d", rand.Intn(20)))
	}
	b.WriteByte('\n')
	return b.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
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
