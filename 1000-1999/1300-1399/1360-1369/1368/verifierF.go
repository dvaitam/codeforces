package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func runCmd(path string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func main() {
	var cand string
	if len(os.Args) == 2 {
		cand = os.Args[1]
	} else if len(os.Args) == 3 && os.Args[1] == "--" {
		cand = os.Args[2]
	} else {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		return
	}
	ref := "./refF.bin"
	if err := exec.Command("go", "build", "-o", ref, "1368F.go").Run(); err != nil {
		fmt.Println("failed to build reference solution:", err)
		return
	}
	defer os.Remove(ref)

	want, err := runCmd(ref)
	if err != nil {
		fmt.Println("reference failed:", err)
		os.Exit(1)
	}
	got, err := runCmd(cand)
	if err != nil {
		fmt.Println("candidate runtime error:", err)
		os.Exit(1)
	}
	if strings.TrimSpace(want) != strings.TrimSpace(got) {
		fmt.Println("wrong answer")
		fmt.Println("expected:\n", want)
		fmt.Println("got:\n", got)
		os.Exit(1)
	}
	fmt.Println("All tests passed.")
}
