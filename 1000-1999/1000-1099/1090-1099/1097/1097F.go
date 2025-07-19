package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   const maxn = 7000
   const chunks = 220
   // precompute bit counts for 16-bit numbers
   bi := make([]int, 1<<16)
   for i := 1; i < 1<<16; i++ {
       bi[i] = bi[i>>1] + (i & 1)
   }
   // compute Mobius function and primes
   miu := make([]int, maxn+1)
   isPrime := make([]bool, maxn+1)
   prime := make([]int, maxn+1)
   tot := 0
   miu[1] = 1
   for i := 2; i <= maxn; i++ {
       isPrime[i] = true
   }
   for i := 2; i <= maxn; i++ {
       if isPrime[i] {
           tot++
           prime[tot] = i
           miu[i] = -1
       }
       for j := 1; j <= tot; j++ {
           p := prime[j]
           if i*p > maxn {
               break
           }
           isPrime[i*p] = false
           if i%p == 0 {
               miu[i*p] = 0
               break
           } else {
               miu[i*p] = -miu[i]
           }
       }
   }
   // convert -1 to 1
   for i := 1; i <= maxn; i++ {
       if miu[i] == -1 {
           miu[i] = 1
       }
   }
   // build g and h bitsets
   var g [maxn+1][chunks]uint32
   var h [maxn+1][chunks]uint32
   for i := 1; i <= maxn; i++ {
       for j := 1; i*j <= maxn; j++ {
           if miu[j] != 0 {
               idx := (i * j) >> 5
               pos := uint((i*j)&31)
               g[i][idx] |= 1 << pos
           }
       }
   }
   for i := 1; i <= maxn; i++ {
       for j := i; j <= maxn; j += i {
           idx := i >> 5
           pos := uint(i & 31)
           h[j][idx] |= 1 << pos
       }
   }
   // I/O
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n, q int
   fmt.Fscan(reader, &n, &q)
   f := make([][chunks]uint32, n+1)
   for q > 0 {
       q--
       var ord, x, y, z int
       fmt.Fscan(reader, &ord, &x, &y)
       switch ord {
       case 1:
           f[x] = h[y]
       case 2:
           fmt.Fscan(reader, &z)
           for k := 0; k < chunks; k++ {
               f[x][k] = f[y][k] ^ f[z][k]
           }
       case 3:
           fmt.Fscan(reader, &z)
           for k := 0; k < chunks; k++ {
               f[x][k] = f[y][k] & f[z][k]
           }
       case 4:
           cnt := 0
           for k := 0; k < chunks; k++ {
               v := g[y][k] & f[x][k]
               cnt += bi[int(v&0xFFFF)] + bi[int(v>>16)]
           }
           if cnt&1 == 1 {
               writer.WriteByte('1')
           } else {
               writer.WriteByte('0')
           }
       }
   }
}
