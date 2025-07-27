package main

import (
   "bufio"
   "fmt"
   "math/bits"
   "os"
)

const MOD = 1000000007

// max carry value is 6, so dimension is 7 (0..6)
const C = 7
const B6 = 6
const B12 = 12

type matrix [C][C]int

var T [2]matrix
var M6 [1<<B6]matrix
var M12 [1<<B12]matrix

// multiply matrices: a * b
func matMul(a, b *matrix) (c matrix) {
   for i := 0; i < C; i++ {
       for k := 0; k < C; k++ {
           var sum int64
           for j := 0; j < C; j++ {
               sum += int64(a[i][j]) * int64(b[j][k])
           }
           c[i][k] = int(sum % MOD)
       }
   }
   return
}

// multiply vector v by matrix m: v * m
func vecMul(v *[C]int, m *matrix) {
   var res [C]int
   for k := 0; k < C; k++ {
       var sum int64
       for j := 0; j < C; j++ {
           sum += int64(v[j]) * int64(m[j][k])
       }
       res[k] = int(sum % MOD)
   }
   *v = res
}

func initMatrices() {
   // build T[0] and T[1]
   for c := 0; c < C; c++ {
       for a := 0; a < 8; a++ {
           s := c + a
           b := s & 1
           c2 := s >> 1
           if c2 < C {
               T[b][c][c2]++
           }
       }
   }
   // build M6: blocks of 6 bits
   var I matrix
   for i := 0; i < C; i++ {
       I[i][i] = 1
   }
   for mask := 0; mask < 1<<B6; mask++ {
       M6[mask] = I
       for j := 0; j < B6; j++ {
           bit := (mask >> j) & 1
           M6[mask] = matMul(&M6[mask], &T[bit])
       }
   }
   // build M12 by combining two M6
   for mask := 0; mask < 1<<B12; mask++ {
       low := mask & ((1 << B6) - 1)
       high := (mask >> B6) & ((1 << B6) - 1)
       M12[mask] = matMul(&M6[low], &M6[high])
   }
}

func main() {
   initMatrices()
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var t int
   fmt.Fscan(reader, &t)
   for i := 0; i < t; i++ {
       var m uint64
       fmt.Fscan(reader, &m)
       // compute number of blocks: bits needed = bitlen(m) + 3
       bl := bits.Len64(m)
       totalBits := bl + 3
       blocks := (totalBits + B12 - 1) / B12
       var v [C]int
       v[0] = 1
       for b := 0; b < blocks; b++ {
           // extract next B12 bits
           shift := uint(b * B12)
           mask := int((m >> shift) & ((1<<B12) - 1))
           vecMul(&v, &M12[mask])
       }
       // result is ways with final carry zero
       fmt.Fprintln(writer, v[0])
   }
}
