package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 1000000007

func main() {
   in := bufio.NewReader(os.Stdin)
   var k int
   if _, err := fmt.Fscan(in, &k); err != nil {
       return
   }
   // T[d]: number of descending sequences starting at a node with d levels below (empty included)
   T := make([]int, k)
   T[0] = 1
   // precompute powers of 2 up to k
   pow2 := make([]int, k+2)
   pow2[0] = 1
   for i := 1; i < len(pow2); i++ {
       pow2[i] = pow2[i-1] * 2 % mod
   }
   // compute T[d]
   for d := 1; d < k; d++ {
       sum := 1 // empty sequence
       for i := 1; i <= d; i++ {
           // choose descendant at depth i (2^i nodes) and then sequence in subtree of size d-i
           sum += pow2[i] * T[d-i] % mod
           if sum >= mod {
               sum -= mod
           }
       }
       T[d] = sum
   }
   inv2 := (mod + 1) / 2
   var ans int
   // sum over d = 0..k-1, level L = k-d, nodes at level = 2^(L-1) = pow2[k-d-1]
   for d := 0; d < k; d++ {
       t := T[d]
       // count of paths with minimal-depth at this node
       // count = (t^2 + 2*t - 1) / 2 mod mod
       cnt := (int((int64(t)*int64(t)%mod + int64(2*t)%mod - 1 + mod) % mod) * inv2) % mod
       nodes := pow2[k-d-1]
       ans = (ans + int((int64(nodes)*int64(cnt))%mod)) % mod
   }
   // output
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   fmt.Fprint(w, ans)
}
