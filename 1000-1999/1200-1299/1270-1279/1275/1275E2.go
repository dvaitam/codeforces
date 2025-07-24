package main

import (
    "bufio"
    "fmt"
    "hash/crc32"
    "os"
)

var table = crc32.MakeTable(crc32.IEEE)

func update(crc uint32, data []byte) uint32 {
    for _, b := range data {
        crc = table[byte(crc)^b] ^ (crc >> 8)
    }
    return crc
}

func computeCRC(data []byte) uint32 {
    crc := uint32(0xffffffff)
    crc = update(crc, data)
    return crc ^ 0xffffffff
}

func computePatch(data []byte, oldCRC uint32, i, j int, x []byte) ([]byte, bool) {
    if i+4 > len(data) || j+4 > len(data) || i < 0 || j < 0 || i > j {
        return nil, false
    }
    crc := uint32(0xffffffff)
    crc = update(crc, data[:i])
    crc = update(crc, x)
    crc = update(crc, data[i+4:j])
    pre := crc

    crcZero := update(pre, []byte{0, 0, 0, 0})
    crcZero = update(crcZero, data[j+4:])
    crcZero ^= 0xffffffff

    var diffCols [32]uint32
    for b := 0; b < 32; b++ {
        patch := []byte{0, 0, 0, 0}
        patch[b/8] = 1 << uint(b%8)
        c := update(pre, patch)
        c = update(c, data[j+4:])
        c ^= 0xffffffff
        diffCols[b] = c ^ crcZero
    }

    delta := crcZero ^ oldCRC
    rows := make([]uint64, 32)
    for r := 0; r < 32; r++ {
        var row uint64
        for c := 0; c < 32; c++ {
            if (diffCols[c]>>uint(r))&1 == 1 {
                row |= 1 << uint(c)
            }
        }
        if (delta>>uint(r))&1 == 1 {
            row |= 1 << 32
        }
        rows[r] = row
    }

    col := 0
    var pos [32]int
    for k := 0; k < 32; k++ {
        pos[k] = -1
    }
    for r := 0; r < 32 && col < 32; {
        pivot := -1
        for i2 := r; i2 < 32; i2++ {
            if (rows[i2]>>uint(col))&1 == 1 {
                pivot = i2
                break
            }
        }
        if pivot == -1 {
            col++
            continue
        }
        rows[r], rows[pivot] = rows[pivot], rows[r]
        pos[r] = col
        for i2 := 0; i2 < 32; i2++ {
            if i2 != r && ((rows[i2]>>uint(col))&1) == 1 {
                rows[i2] ^= rows[r]
            }
        }
        r++
        col++
    }
    for i := 0; i < 32; i++ {
        if rows[i]&0xffffffff == 0 && (rows[i]>>32)&1 == 1 {
            return nil, false
        }
    }

    sol := uint32(0)
    for i := 31; i >= 0; i-- {
        if pos[i] == -1 {
            continue
        }
        c := pos[i]
        val := (rows[i] >> 32) & 1
        for j2 := c + 1; j2 < 32; j2++ {
            if (rows[i]>>uint(j2))&1 == 1 {
                if (sol>>uint(j2))&1 == 1 {
                    val ^= 1
                }
            }
        }
        if val == 1 {
            sol |= 1 << uint(c)
        }
    }

    patchRes := make([]byte, 4)
    for k := 0; k < 4; k++ {
        patchRes[k] = byte((sol >> uint(8*k)) & 0xff)
    }

    tmp := append([]byte(nil), data...)
    copy(tmp[i:i+4], x)
    copy(tmp[j:j+4], patchRes)
    if computeCRC(tmp) != oldCRC {
        return nil, false
    }
    return patchRes, true
}

func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    var n, q int
    if _, err := fmt.Fscan(in, &n, &q); err != nil {
        return
    }
    data := make([]byte, n)
    for i := 0; i < n; i++ {
        var v int
        fmt.Fscan(in, &v)
        data[i] = byte(v)
    }

    oldCRC := computeCRC(data)
    for ; q > 0; q-- {
        var i, j int
        var x0, x1, x2, x3 int
        fmt.Fscan(in, &i, &j, &x0, &x1, &x2, &x3)
        patch, ok := computePatch(data, oldCRC, i, j, []byte{byte(x0), byte(x1), byte(x2), byte(x3)})
        if !ok {
            fmt.Fprintln(out, "No solution.")
        } else {
            fmt.Fprintf(out, "%d %d %d %d\n", patch[0], patch[1], patch[2], patch[3])
        }
    }
}

