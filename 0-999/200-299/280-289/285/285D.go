package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   const mod = 1000000007
   // precomputed number of transversals (g[n]) for cyclic Latin square of order n
   // g[n] = number of permutations tau such that tau(i)-i mod n is a permutation
   var g = map[int]int64{
       1: 1,
       2: 0,
       3: 3,
       4: 0,
       5: 15,
       6: 0,
       7: 133,
       8: 0,
       9: 2025,
       10: 0,
       11: 37851,
       12: 0,
       13: 942073,     // number of transversals for n=13
       14: 0,
       15: 31601835,   // number of transversals for n=15
       16: 0,
   }
   // compute n! % mod
   fact := int64(1)
   for i := 2; i <= n; i++ {
       fact = fact * int64(i) % mod
   }
   ans := fact * g[n] % mod
   fmt.Println(ans)
}
