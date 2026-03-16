package main

import (
    "bufio"
    "bytes"
    "fmt"
    "os"
    "os/exec"
    "strings"
)

func buildOracle() (string, error) {
    exe := "oracleC"
    cmd := exec.Command("go", "build", "-o", exe, "940C.go")
    if out, err := cmd.CombinedOutput(); err != nil {
        return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
    }
    return "./" + exe, nil
}

func runProgram(bin, input string) (string, error) {
    var cmd *exec.Cmd
    if strings.HasSuffix(bin, ".go") {
        cmd = exec.Command("go", "run", bin)
    } else {
        cmd = exec.Command(bin)
    }
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    var errb bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &errb
    if err := cmd.Run(); err != nil {
        return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
    }
    return strings.TrimSpace(out.String()), nil
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("usage: go run verifierC.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    oracle, err := buildOracle()
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
    defer os.Remove(oracle)

    const testcasesRaw = `2 3 cl
6 10 itgtbs
6 14 umzxql
18 15 qibalokmnqfrfhhafk
6 5 qqlqvr
6 15 znxqyl
19 12 lofymwxouqhpipqqzlv
15 15 lsxrxopvhkwftiy
16 10 jzwqrqqutsnjxgpq
12 20 czkxagxdbsub
9 19 hvdyqeihg
2 14 wy
2 2 ll
6 8 vacdca
2 1 li
5 6 xfqwa
13 19 bzhebaltuxxdj
11 16 ajorytxbiym
20 5 phcvvkdaozeqsympqkek
9 9 tnuawrevb
9 2 effdouhqw
2 8 hw
15 3 icshtzztwlivniq
1 5 b
13 14 fdqxchddafyhd
7 1 qvoojru
13 7 vygxznnqassbn
17 19 fdvzplaqdtljwljav
14 4 djgyvazobnupog
19 20 cajaljxchypgdslmwoe
12 13 diddctkumgwd
1 20 v
16 2 xwpjloezlipqpxxz
14 16 vjmhfptirnwvwc
19 19 dclfreznczcvzubejmh
11 15 fqjderyndkq
8 17 iffowhml
19 5 ooxaztmxfmqbpimiwxn
16 12 rkwxvcyxhrtgmvmu
1 11 o
17 15 ufdamgxstmgdmryzg
9 19 sgpzteatv
14 16 iqsfowgyclaprv
3 19 pvk
15 9 qoactylfyyzmivu
5 2 fpmov
10 5 ajroalbrms
15 7 vjpuepwrwjcikjk
10 13 qcqugmtqez
17 3 jbhorhqibddvzmlgk
12 3 kolfpojoewou
7 9 kfdhpgy
12 6 lezehizrummz
11 9 xtqswxkxmyw
10 18 tuvcljmpfi
12 15 pcfkmeadlflc
14 1 rkhtmrjpuelkgp
4 5 zgki
5 14 lickg
8 8 xtbkluyt
2 5 fc
14 15 yiekqsdkuywtmh
2 13 yp
16 20 krttcsqrvpmwofnm
17 15 bdosedvqfcmjozwai
4 12 hfae
14 3 kuobphcperaewq
18 2 bgraqkvqhelpaerdhd
15 7 zbtgumktumwqqyv
6 17 deugfm
7 10 knenemk
10 4 rdpijqypih
14 5 wrvdatrygggmsb
5 1 xiwxp
18 2 xyhetkbwgdeuwrfycv
15 10 gfkwiqscnnvxboj
4 9 agnk
9 18 xmsqxgnye
6 15 olmpti
20 7 spogypskjcfltuphytuv
19 5 vjgrjdazagkbkrizxvk
15 3 npwajssegeftymx
3 19 oiu
3 16 pzz
8 5 sjhgtwks
20 13 qnhugrbivhetxmndommp
13 10 ghhbrqctrvabm
14 13 hqidlqlqzpscwo
8 9 aapbeueg
11 8 rbteujydurr
3 5 nwe
2 10 qv
9 16 brlykvdtl
4 20 zllz
9 16 jqteabknu
1 12 v
18 2 vcwrqtynnnhzfftbas
12 6 jabhszhmcldt
3 8 hrg
4 1 wmcq`

    scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
    idx := 0
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line == "" {
            continue
        }
        idx++
        input := line + "\n"
        exp, err := runProgram(oracle, input)
        if err != nil {
            fmt.Fprintf(os.Stderr, "oracle error on case %d: %v\n", idx, err)
            os.Exit(1)
        }
        got, err := runProgram(bin, input)
        if err != nil {
            fmt.Fprintf(os.Stderr, "case %d: %v\n", idx, err)
            os.Exit(1)
        }
        if got != exp {
            fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx, exp, got)
            os.Exit(1)
        }
    }
    if err := scanner.Err(); err != nil {
        fmt.Fprintln(os.Stderr, "scanner error:", err)
        os.Exit(1)
    }
    fmt.Printf("All %d tests passed\n", idx)
}

