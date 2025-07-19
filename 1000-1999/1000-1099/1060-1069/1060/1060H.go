package main

import (
   "bufio"
   "fmt"
   "os"
)

const (
   N    = 13
   MUL  = 2999
   ZERO = 4999
   ONE  = 4997
)

var (
   d, p int
   C     [N][N]int
   m     [N][N]int
   b     [N]int
)

// powerMod computes a^n mod mod
func powerMod(a int64, n int, mod int64) int64 {
   var res int64 = 1
   a %= mod
   for n > 0 {
       if n&1 == 1 {
           res = res * a % mod
       }
       a = a * a % mod
       n >>= 1
   }
   return res
}

func add(x, y, z int) {
   fmt.Printf("+ %d %d %d\n", x, y, z)
}

func Pow(x, y int) {
   fmt.Printf("^ %d %d\n", x, y)
}

func fin(x int) {
   fmt.Printf("f %d\n", x)
}

func zero(x int) {
   add(ZERO-1, ZERO, x)
}

func mulCell(x, c, y int) {
   zero(y)
   add(x, ZERO, MUL)
   for c > 0 {
       if c&1 == 1 {
           add(y, MUL, y)
       }
       add(MUL, MUL, MUL)
       c >>= 1
   }
}

func initDevice() {
   c := p - 1
   n := 5000
   for c > 0 {
       if c&1 == 1 {
           add(n-1, n, n-1)
           add(n-2, n, n-2)
       }
       add(n, n, n)
       c >>= 1
   }
}

func getPow(x, y int) {
   if d == 2 {
       Pow(x, y)
       return
   }
   add(x, ZERO, 10)
   for i := 0; i < d; i++ {
       add(10+i, ONE, 11+i)
   }
   for i := 0; i <= d; i++ {
       Pow(10+i, 21+i)
   }
   for i := 0; i <= d; i++ {
       mulCell(21+i, b[i], 32+i)
   }
   for i := 0; i < d; i++ {
       add(32, 33+i, 32)
   }
   add(32, ZERO, y)
}

func gauss() {
   // build binomial C
   for i := 0; i <= d; i++ {
       C[i][0] = 1
       for j := 1; j <= i; j++ {
           C[i][j] = (C[i-1][j-1] + C[i-1][j]) % p
       }
   }
   // build matrix m
   for i := 0; i <= d; i++ {
       for j := 0; j <= d; j++ {
           m[i][j] = int(powerMod(int64(j), d-i, int64(p)) * int64(C[d][i]) % int64(p))
       }
   }
   b[2] = 1
   for k := 0; k <= d; k++ {
       // pivot
       i := k
       for ; i <= d; i++ {
           if m[i][k] != 0 {
               break
           }
       }
       if i != k {
           b[i], b[k] = b[k], b[i]
           for j := 0; j <= d; j++ {
               m[i][j], m[k][j] = m[k][j], m[i][j]
           }
       }
       inv := powerMod(int64(m[k][k]), p-2, int64(p))
       for j := 0; j <= d; j++ {
           m[k][j] = int(int64(m[k][j]) * inv % int64(p))
       }
       b[k] = int(int64(b[k]) * inv % int64(p))
       for ii := 0; ii <= d; ii++ {
           if ii == k {
               continue
           }
           coe := m[ii][k]
           for j := 0; j <= d; j++ {
               m[ii][j] = (m[ii][j] - coe*m[k][j]) % p
               if m[ii][j] < 0 {
                   m[ii][j] += p
               }
           }
           b[ii] = (b[ii] - coe*b[k]) % p
           if b[ii] < 0 {
               b[ii] += p
           }
       }
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   if _, err := fmt.Fscan(reader, &d, &p); err != nil {
       return
   }
   initDevice()
   if d > 2 {
       gauss()
   }
   add(1, 2, 3)
   getPow(1, 4)
   getPow(2, 5)
   getPow(3, 6)
   mulCell(4, p-1, 7)
   mulCell(5, p-1, 8)
   add(6, 7, 6)
   add(6, 8, 6)
   mulCell(6, (p+1)/2, 9)
   fin(9)
}
