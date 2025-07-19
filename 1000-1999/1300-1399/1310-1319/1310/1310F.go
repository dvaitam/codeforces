package main

import (
   "bufio"
   "fmt"
   "os"
   "math/bits"
)

const l uint64 = 0xFFFFFFFFFFFFFFFF

var p = []int{3, 5, 17, 257, 641, 65537, 6700417}
var mArr = []int{1, 2, 6, 16, 25, 256, 2588}
var I = []int{2, 2, 1, 32, 590, 16384, 3883315}

var F [64][64]uint64
var T [64][256][256]uint64

// mul2 multiplies two monomials in GF(2)[x] mod the chosen irreducible polynomial
func mul2(a, b uint64) uint64 {
   if a == 1 || b == 1 {
       return a * b
   }
   if a < b {
       a, b = b, a
   }
   x := bits.TrailingZeros64(a)
   y := bits.TrailingZeros64(b)
   // find largest power-of-two <= x
   exp := bits.Len(uint(x)) - 1
   shift := uint(exp)
   mask := uint64(1) << shift
   if uint64(y)&mask != 0 {
       // (1<<shift >>1) * 3 == 3 << (shift-1)
       return mul((uint64(1)<<shift>>1)*3, mul2(a>>shift, b>>shift))
   }
   return mul2(a>>shift, b) << shift
}

// mul multiplies two field elements by decomposing into monomials
func mul(a, b uint64) uint64 {
   var c uint64
   for ai := a; ai != 0; ai ^= uint64(1) << bits.TrailingZeros64(ai) {
       x := uint64(1) << bits.TrailingZeros64(ai)
       for bi := b; bi != 0; bi ^= uint64(1) << bits.TrailingZeros64(bi) {
           y := uint64(1) << bits.TrailingZeros64(bi)
           c ^= mul2(x, y)
       }
   }
   return c
}

// prod multiplies two general field elements using a precomputed table
func prod(a, b uint64) uint64 {
   var c uint64
   for i := 0; i < 8; i++ {
       for j := 0; j < 8; j++ {
           c ^= T[i<<3|j][(a>>(i*8))&0xFF][(b>>(j*8))&0xFF]
       }
   }
   return c
}

// qpow computes exponentiation in the field
func qpow(x, y uint64) uint64 {
   var t uint64 = 1
   for y > 0 {
       if y&1 != 0 {
           t = prod(t, x)
       }
       x = prod(x, x)
       y >>= 1
   }
   return t
}

func main() {
   // build F table for monomial products
   for x := 0; x < 64; x++ {
       for y := 0; y < 64; y++ {
           F[x][y] = mul2(1<<uint(x), 1<<uint(y))
       }
   }
   // build T table for byte-wise products
   for i := 0; i < 8; i++ {
       for j := 0; j < 8; j++ {
           idx := (i << 3) | j
           // initialize for single-bit bytes
           for x := 0; x < 8; x++ {
               for y := 0; y < 8; y++ {
                   T[idx][1<<x][1<<y] = F[(i<<3)|x][(j<<3)|y]
               }
           }
           // fill for all byte values
           for a := 1; a < 256; a++ {
               for b := 1; b < 256; b++ {
                   if a^(a&-a) != 0 {
                       T[idx][a][b] = T[idx][a&-a][b] ^ T[idx][a^(a&-a)][b]
                   } else if b^(b&-b) != 0 {
                       T[idx][a][b] = T[idx][a][b&-b] ^ T[idx][a][b^(b&-b)]
                   }
               }
           }
       }
   }

   reader := bufio.NewReader(os.Stdin)
   var t int
   fmt.Fscan(reader, &t)
   for ; t > 0; t-- {
       var a, b uint64
       fmt.Fscan(reader, &a, &b)
       var ans uint64
       var r [7]int
       fail := false
       for i := 0; i < 7; i++ {
           A := qpow(a, l/uint64(p[i]))
           B := qpow(b, l/uint64(p[i]))
           Am := qpow(A, uint64(mArr[i]))
           P := uint64(1)
           id := make(map[uint64]int)
           for j := 0; j < mArr[i]; j++ {
               id[B] = j
               B = prod(B, A)
           }
           res := -1
           for j := 0; j <= (p[i]-1)/mArr[i]; j++ {
               if v, ok := id[P]; ok {
                   res = j*mArr[i] - v
                   res %= p[i]
                   if res < 0 {
                       res += p[i]
                   }
                   break
               }
               P = prod(P, Am)
           }
           if res >= 0 {
               r[i] = res
           } else {
               fmt.Println(-1)
               fail = true
               break
           }
       }
       if fail {
           continue
       }
       for i := 0; i < 7; i++ {
           ans = (ans + uint64(I[i])*(l/uint64(p[i]))*uint64(r[i])) % l
       }
       fmt.Println(ans)
   }
}
