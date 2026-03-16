package main

import (
    "bufio"
    "bytes"
    "fmt"
    "os"
    "os/exec"
    "strconv"
)

const testcasesDRaw = `100
4 1 11 21 31
4 2 12 22 32
4 3 13 23 33
4 4 14 24 34
4 5 15 25 35
4 6 16 26 36
4 7 17 27 37
4 8 18 28 38
4 9 19 29 39
4 10 20 30 40
4 11 21 31 41
4 12 22 32 42
4 13 23 33 43
4 14 24 34 44
4 15 25 35 45
4 16 26 36 46
4 17 27 37 47
4 18 28 38 48
4 19 29 39 49
4 20 30 40 50
4 21 31 41 51
4 22 32 42 52
4 23 33 43 53
4 24 34 44 54
4 25 35 45 55
4 26 36 46 56
4 27 37 47 57
4 28 38 48 58
4 29 39 49 59
4 30 40 50 60
4 31 41 51 61
4 32 42 52 62
4 33 43 53 63
4 34 44 54 64
4 35 45 55 65
4 36 46 56 66
4 37 47 57 67
4 38 48 58 68
4 39 49 59 69
4 40 50 60 70
4 41 51 61 71
4 42 52 62 72
4 43 53 63 73
4 44 54 64 74
4 45 55 65 75
4 46 56 66 76
4 47 57 67 77
4 48 58 68 78
4 49 59 69 79
4 50 60 70 80
4 51 61 71 81
4 52 62 72 82
4 53 63 73 83
4 54 64 74 84
4 55 65 75 85
4 56 66 76 86
4 57 67 77 87
4 58 68 78 88
4 59 69 79 89
4 60 70 80 90
4 61 71 81 91
4 62 72 82 92
4 63 73 83 93
4 64 74 84 94
4 65 75 85 95
4 66 76 86 96
4 67 77 87 97
4 68 78 88 98
4 69 79 89 99
4 70 80 90 100
4 71 81 91 101
4 72 82 92 102
4 73 83 93 103
4 74 84 94 104
4 75 85 95 105
4 76 86 96 106
4 77 87 97 107
4 78 88 98 108
4 79 89 99 109
4 80 90 100 110
4 81 91 101 111
4 82 92 102 112
4 83 93 103 113
4 84 94 104 114
4 85 95 105 115
4 86 96 106 116
4 87 97 107 117
4 88 98 108 118
4 89 99 109 119
4 90 100 110 120
4 91 101 111 121
4 92 102 112 122
4 93 103 113 123
4 94 104 114 124
4 95 105 115 125
4 96 106 116 126
4 97 107 117 127
4 98 108 118 128
4 99 109 119 129
4 100 110 120 130`

func main() {
    if len(os.Args) != 2 {
        fmt.Println("Usage: go run verifierD.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    data := []byte(testcasesDRaw)
    scan := bufio.NewScanner(bytes.NewReader(data))
    scan.Split(bufio.ScanWords)
    if !scan.Scan() {
        fmt.Println("invalid test file")
        os.Exit(1)
    }
    t, _ := strconv.Atoi(scan.Text())
    expected := make([]int, t)
    for i := 0; i < t; i++ {
        if !scan.Scan() {
            fmt.Println("bad test file")
            os.Exit(1)
        }
        n, _ := strconv.Atoi(scan.Text())
        minv := int(1<<31 - 1)
        maxv := -minv - 1
        for j := 0; j < n; j++ {
            scan.Scan()
            x, _ := strconv.Atoi(scan.Text())
            if x < minv { minv = x }
            if x > maxv { maxv = x }
        }
        expected[i] = maxv - minv
    }

    cmd := exec.Command(bin)
    cmd.Stdin = bytes.NewReader(data)
    var out bytes.Buffer
    cmd.Stdout = &out
    if err := cmd.Run(); err != nil {
        fmt.Println("execution failed:", err)
        os.Exit(1)
    }

    outScan := bufio.NewScanner(bytes.NewReader(out.Bytes()))
    outScan.Split(bufio.ScanWords)
    for i := 0; i < t; i++ {
        if !outScan.Scan() {
            fmt.Printf("missing output for test %d\n", i+1)
            os.Exit(1)
        }
        got, _ := strconv.Atoi(outScan.Text())
        if got != expected[i] {
            fmt.Printf("test %d failed: expected %d got %d\n", i+1, expected[i], got)
            os.Exit(1)
        }
    }
    if outScan.Scan() {
        fmt.Println("extra output detected")
        os.Exit(1)
    }
    fmt.Println("All tests passed")
}
