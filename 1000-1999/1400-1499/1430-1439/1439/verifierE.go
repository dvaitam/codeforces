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
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func compileRef() (string, error) {
	ref := "1439E.go"
	tmp := "refE_bin"
	cmd := exec.Command("go", "build", "-o", tmp, ref)
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return tmp, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/bin")
		return
	}
	bin := os.Args[1]
	ref, err := compileRef()
	if err != nil {
		fmt.Println("compile reference:", err)
		return
	}
	defer os.Remove(ref)
	rand.Seed(1)
	for tcase := 0; tcase < 100; tcase++ {
		m := rand.Intn(3) + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", m))
		for i := 0; i < m; i++ {
			x1 := rand.Intn(8)
			y1 := rand.Intn(8)
			x2 := rand.Intn(8)
			y2 := rand.Intn(8)
			x1 &= ^y1
			x2 &= ^y2 // ensure good cells (x&y==0)
			sb.WriteString(fmt.Sprintf("%d %d %d %d\n", x1, y1, x2, y2))
		}
		inp := sb.String()
		outRef, _ := run(ref, inp)
		outCand, err := run(bin, inp)
		if err != nil {
			fmt.Printf("test %d exec err %v\n", tcase, err)
			return
		}
		if strings.TrimSpace(outRef) != strings.TrimSpace(outCand) {
			fmt.Printf("test %d mismatch expected %s got %s\n", tcase, strings.TrimSpace(outRef), strings.TrimSpace(outCand))
			return
		}
	}
	fmt.Println("All tests passed")
}
