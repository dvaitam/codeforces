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

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func buildRef(src, out string) error {
	cmd := exec.Command("go", "build", "-o", out, src)
	return cmd.Run()
}

func generateCaseE(rng *rand.Rand) string {
	r := rng.Intn(5) + 1
	n := rng.Intn(10) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", r, n)
	nextID := 1
	active := make([]int, 0)
	coords := map[[2]int]bool{}
	for i := 0; i < n; i++ {
		typ := rng.Intn(3) + 1
		if typ == 1 || len(active) < 2 {
			typ = 1
		}
		if typ == 1 {
			x := rng.Intn(10) + r + 1
			y := rng.Intn(10) + 1
			for coords[[2]int{x, y}] {
				x = rng.Intn(10) + r + 1
				y = rng.Intn(10) + 1
			}
			coords[[2]int{x, y}] = true
			fmt.Fprintf(&sb, "1 %d %d\n", x, y)
			active = append(active, nextID)
			nextID++
		} else if typ == 2 {
			idx := rng.Intn(len(active))
			id := active[idx]
			active = append(active[:idx], active[idx+1:]...)
			fmt.Fprintf(&sb, "2 %d\n", id)
		} else {
			i1 := rng.Intn(len(active))
			j1 := rng.Intn(len(active) - 1)
			if j1 >= i1 {
				j1++
			}
			fmt.Fprintf(&sb, "3 %d %d\n", active[i1], active[j1])
		}
	}
	return sb.String()
}

func runCaseE(bin, ref, input string) error {
	expect, err := runBinary(ref, input)
	if err != nil {
		return fmt.Errorf("reference error: %v", err)
	}
	got, err := runBinary(bin, input)
	if err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, got)
	}
	if strings.TrimSpace(got) != strings.TrimSpace(expect) {
		return fmt.Errorf("expected %q got %q", strings.TrimSpace(expect), strings.TrimSpace(got))
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref := "./refE.bin"
	if err := buildRef("856E.go", ref); err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(ref)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCaseE(rng)
		if err := runCaseE(bin, ref, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
