package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 998244353

// DSU with parity for XOR constraints
type DSU struct {
   p, d []int
}

func NewDSU(n int) *DSU {
   p := make([]int, n)
   d := make([]int, n)
   for i := 0; i < n; i++ {
       p[i] = i
       d[i] = 0
   }
   return &DSU{p: p, d: d}
}

// find root of x, return root and parity from x to root
func (u *DSU) find(x int) (int, int) {
   if u.p[x] == x {
       return x, 0
   }
   r, dr := u.find(u.p[x])
   u.d[x] ^= dr
   u.p[x] = r
   return r, u.d[x]
}

// unite x and y with constraint x XOR y = val, return false if conflict
func (u *DSU) unite(x, y, val int) bool {
   rx, dx := u.find(x)
   ry, dy := u.find(y)
   if rx == ry {
       // check consistency: dx XOR dy == val
       return (dx ^ dy) == val
   }
   // attach rx under ry: set d[rx] so that
   // dx ^ d_rx ^ dy == val => d_rx = val ^ dx ^ dy
   u.p[rx] = ry
   u.d[rx] = val ^ dx ^ dy
   return true
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var s string
   fmt.Fscan(in, &s)
   m := len(s)
   // parse s: -1 for '?', 0 or 1
   sp := make([]int, m)
   for i, ch := range s {
       if ch == '?' {
           sp[i] = -1
       } else if ch == '0' {
           sp[i] = 0
       } else {
           sp[i] = 1
       }
   }
   // precompute powers of 2
   pw2 := make([]int, m*2+5)
   pw2[0] = 1
   for i := 1; i < len(pw2); i++ {
       pw2[i] = pw2[i-1] * 2 % MOD
   }
   ans := 0
   // iterate A length L from 1 to m-1, B length is m
   for L := 1; L < m; L++ {
       // JA = number of a-vars, JB = number of b-vars, plus const node
       JA := (L + 1) / 2
       JB := (m + 1) / 2
       N := JA + JB + 1
       constID := JA + JB
       dsu := NewDSU(N)
       ok := true
       // a_0 = 1
       if !dsu.unite(0, constID, 1) {
           ok = false
       }
       // b_0 = 1
       if ok {
           if !dsu.unite(JA+0, constID, 1) {
               ok = false
           }
       }
       // constraints for positions
       for i := 0; ok && i < m; i++ {
           v := sp[i]
           if v < 0 {
               continue
           }
           // determine if in A padded region or A region
           if i < m-L {
               // A[i]=0, so B[i] = v
               k := i
               if m-1-i < i {
                   k = m-1-i
               }
               bid := JA + k
               if !dsu.unite(bid, constID, v) {
                   ok = false
                   break
               }
           } else {
               // i in [m-L, m-1]
               j := i - (m - L)
               l := j
               if L-1-j < j {
                   l = L-1-j
               }
               aid := l
               k := i
               if m-1-i < i {
                   k = m-1-i
               }
               bid := JA + k
               // a[aid] XOR b[bid] = v
               if !dsu.unite(aid, bid, v) {
                   ok = false
                   break
               }
           }
       }
       if !ok {
           continue
       }
       // count free components (not connected to constID)
       seen := make(map[int]bool)
       constRoot, _ := dsu.find(constID)
       cnt := 0
       for x := 0; x < JA+JB; x++ {
           rx, _ := dsu.find(x)
           if rx == constRoot {
               continue
           }
           if !seen[rx] {
               seen[rx] = true
               cnt++
           }
       }
       ans = (ans + pw2[cnt]) % MOD
   }
   fmt.Fprintln(out, ans)
}
