package main

import (
   "bufio"
   "fmt"
   "os"
)

// extgcd returns g = gcd(a,b) and x,y such that a*x + b*y = g
func extgcd(a, b int) (g, x, y int) {
   if b == 0 {
       return a, 1, 0
   }
   g, x1, y1 := extgcd(b, a%b)
   return g, y1, x1 - (a/b)*y1
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n, k int
   fmt.Fscan(in, &n, &k)
   q := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &q[i])
       q[i]--
   }
   s := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &s[i])
       s[i]--
   }
   seen := make([]bool, n)
   x0, M := 0, 1
   for i := 0; i < n; i++ {
       if seen[i] {
           continue
       }
       // build cycle of q starting at i
       cyc := []int{}
       for u := i; !seen[u]; u = q[u] {
           seen[u] = true
           cyc = append(cyc, u)
       }
       c := len(cyc)
       // find shift t such that applying q^t to cyc[0] gives s[cyc[0]]
       t := -1
       for j, v := range cyc {
           if s[cyc[0]] == v {
               t = j
               break
           }
       }
       if t < 0 {
           fmt.Fprintln(out, "NO")
           return
       }
       // verify cycle matches the same shift
       for j, v := range cyc {
           if s[v] != cyc[(j+t)%c] {
               fmt.Fprintln(out, "NO")
               return
           }
       }
       // merge CRT: x ≡ x0 mod M, x ≡ t mod c
       a, m := t, c
       g, p, _ := extgcd(M, m)
       if (a-x0)%g != 0 {
           fmt.Fprintln(out, "NO")
           return
       }
       // lcm of M and m
       lcm := M / g * m
       // compute k0 = (a - x0)/g * p mod (m/g)
       mul := (a - x0) / g
       mod := m / g
       k0 := (mul * p) % mod
       if k0 < 0 {
           k0 += mod
       }
       x0 = x0 + M*k0
       M = lcm
       x0 %= M
       if x0 < 0 {
           x0 += M
       }
   }
   // x0 is smallest non-negative exponent with q^x0 = s
   if x0 == 0 {
       // s is identity, which occurs at start
       fmt.Fprintln(out, "NO")
       return
   }
   // minimal moves to reach: d = min(x0, M-x0)
   d := x0
   if M-x0 < d {
       d = M - x0
   }
   // check reachability in exactly k moves, avoiding earlier occurrence
   if k >= d && (k-d)%2 == 0 {
       fmt.Fprintln(out, "YES")
   } else {
       fmt.Fprintln(out, "NO")
   }
}
