package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

var mod int64
var almostRoot [18]int64
var almostRootPow [18]int64

// multMod computes (a * b) % mod avoiding overflow
func multMod(a, b int64) int64 {
   // split b to high and low 16 bits
   res := (a * (b >> 16)) % mod
   res = (res << 16) % mod
   res = (res + a*(b & ((1<<16)-1))) % mod
   return res
}

func main() {
   // compute mod = 5^17
   mod = 1
   for i := 0; i < 17; i++ {
       mod *= 5
   }
   // precompute almost roots: 2^(5^i) mod mod
   almostRoot[0] = 2
   almostRoot[1] = 16
   almostRootPow[0] = 1
   almostRootPow[1] = 4
   for i := 2; i <= 17; i++ {
       x := multMod(almostRoot[i-1], almostRoot[i-1])
       x = multMod(x, x)
       x = multMod(x, almostRoot[i-1])
       almostRoot[i] = x
       almostRootPow[i] = almostRootPow[i-1] * 5
   }

   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   // read number of queries
   line, err := reader.ReadString('\n')
   if err != nil {
       return
   }
   n, err := strconv.Atoi(line[:len(line)-1])
   if err != nil {
       return
   }
   for qi := 0; qi < n; qi++ {
       // read ai
       line, err = reader.ReadString('\n')
       if err != nil {
           return
       }
       ai, err := strconv.ParseInt(line[:len(line)-1], 10, 64)
       if err != nil {
           return
       }
       // align a to multiple of 2^17 above ai*1e6
       a := ai * 1000000
       a += 1 << 17
       a &= ^((1 << 17) - 1)
       if a%5 == 0 {
           a += 1 << 17
       }
       // initial power and value r = 2^60 mod mod
       p := int64(60)
       r := (1 << 60) % mod
       curmod := int64(5)
       // refine r to match a modulo increasing powers of 5
       for modp := 1; modp <= 17; modp++ {
           for (r-a)%curmod != 0 {
               r = multMod(r, almostRoot[modp-1])
               p += almostRootPow[modp-1]
           }
           curmod *= 5
       }
       fmt.Fprintln(writer, p)
   }
}
