package main

import (
    "bufio"
    "bytes"
    "fmt"
    "os"
    "os/exec"
    "strconv"
    "strings"
)

func runExe(path string, input []byte) (string, error) {
    var cmd *exec.Cmd
    if strings.HasSuffix(path, ".go") {
        cmd = exec.Command("go", "run", path)
    } else {
        cmd = exec.Command(path)
    }
    cmd.Stdin = bytes.NewReader(input)
    var out bytes.Buffer
    var errb bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &errb
    if err := cmd.Run(); err != nil {
        return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
    }
    return strings.TrimSpace(out.String()), nil
}

func buildRef() (string, error) {
    ref := "./refB.bin"
    cmd := exec.Command("go", "build", "-o", ref, "1347B.go")
    if out, err := cmd.CombinedOutput(); err != nil {
        return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
    }
    return ref, nil
}

func readTests() ([]string, error) {
    f, err := os.Open("testcasesB.txt")
    if err != nil {
        return nil, err
    }
    defer f.Close()
    scan := bufio.NewScanner(f)
    scan.Split(bufio.ScanWords)
    if !scan.Scan() {
        return nil, fmt.Errorf("empty test file")
    }
    t, _ := strconv.Atoi(scan.Text())
    tests := make([]string, 0, t)
    for i := 0; i < t; i++ {
        if !scan.Scan() {
            return nil, fmt.Errorf("invalid test %d", i+1)
        }
        a1 := scan.Text()
        if !scan.Scan() {
            return nil, fmt.Errorf("invalid test %d", i+1)
        }
        b1 := scan.Text()
        if !scan.Scan() {
            return nil, fmt.Errorf("invalid test %d", i+1)
        }
        a2 := scan.Text()
        if !scan.Scan() {
            return nil, fmt.Errorf("invalid test %d", i+1)
        }
        b2 := scan.Text()
        tests = append(tests, fmt.Sprintf("%s %s\n%s %s\n", a1, b1, a2, b2))
    }
    return tests, nil
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("Usage: go run verifierB.go /path/to/binary")
        return
    }
    bin := os.Args[1]
    ref, err := buildRef()
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
    defer os.Remove(ref)

    tests, err := readTests()
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }

    for i, input := range tests {
        exp, err := runExe(ref, []byte(input))
        if err != nil {
            fmt.Fprintf(os.Stderr, "reference runtime error on test %d: %v\n", i+1, err)
            os.Exit(1)
        }
        got, err := runExe(bin, []byte(input))
        if err != nil {
            fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\n", i+1, err)
            os.Exit(1)
        }
        if strings.TrimSpace(exp) != strings.TrimSpace(got) {
            fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%sGot:\n%s\n", i+1, input, exp, got)
            os.Exit(1)
        }
    }
    fmt.Println("All tests passed")
}
