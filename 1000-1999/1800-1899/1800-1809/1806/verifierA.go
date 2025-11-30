package main

import (
    "bufio"
    "bytes"
    "encoding/base64"
    "fmt"
    "os"
    "os/exec"
    "strconv"
    "strings"
)

// Base64-encoded contents of testcasesA.txt.
const testcasesA64 = "MiAzIC05IC0yCjYgNSAyIC0xCjUgMSA4IC00CjYgLTYgLTEgLTYKLTcgOSAtMiA3CjkgLTYgLTEgLTcKLTggMCA1IDcKLTcgMSAzIDAKOSAxMCAtNCA3CjUgNCA2IC0yCi05IDcgLTEwIC04CjIgMTAgLTEwIDkKNSAwIC0zIDAKLTggLTQgOCAtMwotMyAtNiA3IDQKLTggLTggMCA2CjUgLTcgLTEgNwotMSAtNyA3IDAKNyAtNCA5IDcKOCAtMSA0IC04CjkgMiAwIDgKLTMgLTEgLTUgLTQKLTUgLTkgOSAtMgo1IC04IC04IC02Ci02IC05IC04IDcKMiA2IC0yIDYKLTMgLTQgOCAzCjggLTIgNCA1CjEwIDEgLTggMAo5IC03IDUgOAoxMCAwIC00IC0zCi0xMCAtMiAtNyAtMwoxIC01IDAgMwotOSAtNyAtNiAtMwotOSA4IDEwIDcKOSAtOCAtMTAgLTcKMTAgLTQgOSA4Ci03IDIgLTggMQotNyAtOSA5IC0xMAotNCAtNSAtNyA1Ci00IC05IC0xMCA3CjMgOSAtNyAtMgotOCAtMyAtOCAxMAotMSAxIDMgLTUKLTkgNiA0IC05CjkgLTcgMiAtNAotMiAxIDUgOAotNSAtNCAtOSAtNQotNSAwIDYgLTIKLTcgOSA0IC01Ci0xMCA1IDMgOAo2IC0xIDEwIDEKMiAtMiAtNiA3Ci0xMCA0IC04IDAKLTkgNyAtMiAtNgotMyA1IDEgOQotMSAxIDggMTAKOSAtNiAtMSAyCjMgMTAgLTggLTEwCjkgLTQgMCAtNQotMyAtMyAxMCA0CjIgOCAzIC05CjIgOCAzIC05Ci01IDQgLTggLTIKLTUgNCA2IDUKNyA5IC0xMCAtOQo1IDAgLTEgNAotOSAzIC00IDcKMTAgLTggLTYgLTEwCjIgMyAwIC0xMAotNCAtMTAgLTEwIDYKOSAtNyAtNCAtNwo5IDEwIC00IC0xCi0yIC01IC03IDUKMiAxMCAtOCAtMTAKLTIgNCAtNyAtMgotNiAxMCA2IDEwCjEwIDEgLTcgLTYKLTIgLTEwIC05IC05Ci00IC0yIDcgMAoxIDggLTkgOQoxMCA1IDEwIDQKMTAgMyAxIDcKLTUgLTQgMiA4Ci0xIC0xMCAtNiAtNgotMiAwIDAgMQotOCAwIDkgLTkKLTkgLTIgLTUgLTYKOCAtMSAxIDIKNyAtNiAtMSAtNwo1IC0zIC05IC0xCi01IDYgLTggLTEKMiAwIC0xIDMKLTcgLTcgNyA1CjUgMCAwIC03CjUgLTcgNSAzCi05IC0xIDAgLTYKLTUgMTAgOCAyCjEwIC04IC04IC04Ci00IC0zIC05IDIK"

type testCase struct {
    a, b, c, d int64
}

// Embedded solution logic from 1806A.go.
func solve(tc testCase) string {
    if tc.d < tc.b || tc.c > tc.a+tc.d-tc.b {
        return "-1"
    }
    steps := (tc.d - tc.b) + (tc.a+tc.d-tc.b-tc.c)
    return strconv.FormatInt(steps, 10)
}

func parseTestcases() ([]testCase, error) {
    raw, err := base64.StdEncoding.DecodeString(testcasesA64)
    if err != nil {
        return nil, err
    }
    sc := bufio.NewScanner(bytes.NewReader(raw))
    sc.Split(bufio.ScanWords)
    tests := make([]testCase, 0)
    for {
        var tc testCase
        fields := []*int64{&tc.a, &tc.b, &tc.c, &tc.d}
        for i := 0; i < 4; i++ {
            if !sc.Scan() {
                if i == 0 {
                    return tests, nil
                }
                return nil, fmt.Errorf("incomplete testcase data")
            }
            val, err := strconv.ParseInt(sc.Text(), 10, 64)
            if err != nil {
                return nil, fmt.Errorf("parse int: %v", err)
            }
            *fields[i] = val
        }
        tests = append(tests, tc)
    }
}

func buildIfGo(path string) (string, func(), error) {
    if strings.HasSuffix(path, ".go") {
        tmp, err := os.CreateTemp("", "bin*")
        if err != nil {
            return "", nil, err
        }
        tmp.Close()
        if out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput(); err != nil {
            os.Remove(tmp.Name())
            return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
        }
        return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
    }
    return path, func() {}, nil
}

func runCandidate(bin, input string) (string, error) {
    cmd := exec.Command(bin)
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    var errBuf bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &errBuf
    if err := cmd.Run(); err != nil {
        return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
    }
    return strings.TrimSpace(out.String()), nil
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("usage: go run verifierA.go /path/to/binary")
        os.Exit(1)
    }

    cases, err := parseTestcases()
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }

    bin, cleanup, err := buildIfGo(os.Args[1])
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
    defer cleanup()

    for idx, tc := range cases {
        input := fmt.Sprintf("1\n%d %d %d %d\n", tc.a, tc.b, tc.c, tc.d)
        want := solve(tc)
        got, err := runCandidate(bin, input)
        if err != nil {
            fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
            os.Exit(1)
        }
        if got != want {
            fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", idx+1, input, want, got)
            os.Exit(1)
        }
    }
    fmt.Printf("All %d tests passed\n", len(cases))
}

