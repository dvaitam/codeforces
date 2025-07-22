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

const keyboard = "1234567890-=QWERTYUIOP[]\\ASDFGHJKL;'ZXCVBNM,./"

func expected(dir string, typed string) string {
	pos := make(map[rune]int)
	for i, r := range keyboard {
		pos[r] = i
	}
	var b strings.Builder
	for _, r := range typed {
		idx := pos[r]
		if dir == "R" {
			b.WriteByte(keyboard[idx-1])
		} else {
			b.WriteByte(keyboard[idx+1])
		}
	}
	return b.String()
}

func genCase(rng *rand.Rand) (string, string) {
	dir := "L"
	if rng.Intn(2) == 0 {
		dir = "R"
	}
	length := rng.Intn(20) + 1
	var sb strings.Builder
	for i := 0; i < length; i++ {
		sb.WriteByte(keyboard[rng.Intn(len(keyboard))])
	}
	return dir, sb.String()
}

func runCase(bin, dir, typed string) error {
	in := fmt.Sprintf("%s\n%s\n", dir, typed)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := expected(dir, typed)
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		dir, typed := genCase(rng)
		if err := runCase(bin, dir, typed); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s\n%s\n", i+1, err, dir, typed)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
