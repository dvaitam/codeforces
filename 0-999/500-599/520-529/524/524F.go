package main

import (
   "bufio"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   line, _, err := reader.ReadLine()
   if err != nil {
       return
   }
   s := string(line)
   n := len(s)
   // Prefix sums for 2 rotations
   P := make([]int, 2*n+1)
   for i := 0; i < 2*n; i++ {
       v := 1
       if s[i%n] == ')' {
           v = -1
       }
       P[i+1] = P[i] + v
   }
   // Sliding window minimum over P[1..]
   A := P[1:] // length 2n
   M := make([]int, n)
   dq := make([]int, 0, 2*n)
   for i := 0; i < 2*n; i++ {
       for len(dq) > 0 && A[dq[len(dq)-1]] >= A[i] {
           dq = dq[:len(dq)-1]
       }
       dq = append(dq, i)
       if i >= n {
           if dq[0] == i-n {
               dq = dq[1:]
           }
       }
       if i >= n-1 {
           k := i - (n - 1)
           M[k] = A[dq[0]]
       }
   }
   neededOpen := make([]int, n)
   bestOpen := int(1e9)
   for i := 0; i < n; i++ {
       t := P[i] - M[i]
       if t < 0 {
           t = 0
       }
       neededOpen[i] = t
       if t < bestOpen {
           bestOpen = t
       }
   }
   totalSum := P[n]
   neededClose := bestOpen + totalSum
   valid := make([]bool, n)
   for i := 0; i < n; i++ {
       if neededOpen[i] == bestOpen {
           valid[i] = true
       }
   }
   // Build SA for cyclic shifts of s
   // Map '('->1, ')'->2
   sInt := make([]int, n)
   for i, ch := range s {
       if ch == '(' {
           sInt[i] = 1
       } else {
           sInt[i] = 2
       }
   }
   sa := buildSA(sInt, 3)
   start := 0
   for _, pos := range sa {
       if pos < n && valid[pos] {
           start = pos
           break
       }
   }
   // Output result
   w := bufio.NewWriter(os.Stdout)
   for i := 0; i < bestOpen; i++ {
       w.WriteByte('(')
   }
   for i := 0; i < n; i++ {
       w.WriteByte(line[(start+i)%n])
   }
   for i := 0; i < neededClose; i++ {
       w.WriteByte(')')
   }
   w.WriteByte('\n')
   w.Flush()
}

// buildSA builds the suffix array of all cyclic shifts of s of length n
// K is the alphabet size (max(s)+1)
func buildSA(s []int, K int) []int {
   n := len(s)
   p := make([]int, n)
   c := make([]int, n)
   cnt := make([]int, K)
   for i := 0; i < n; i++ {
       cnt[s[i]]++
   }
   for i := 1; i < K; i++ {
       cnt[i] += cnt[i-1]
   }
   for i := n - 1; i >= 0; i-- {
       cnt[s[i]]--
       p[cnt[s[i]]] = i
   }
   c[p[0]] = 0
   classes := 1
   for i := 1; i < n; i++ {
       if s[p[i]] != s[p[i-1]] {
           classes++
       }
       c[p[i]] = classes - 1
   }
   pn := make([]int, n)
   cn := make([]int, n)
   for h := 0; (1 << h) < n; h++ {
       shift := 1 << h
       for i := 0; i < n; i++ {
           pn[i] = p[i] - shift
           if pn[i] < 0 {
               pn[i] += n
           }
       }
       cnt = make([]int, classes)
       for i := 0; i < n; i++ {
           cnt[c[pn[i]]]++
       }
       for i := 1; i < classes; i++ {
           cnt[i] += cnt[i-1]
       }
       for i := n - 1; i >= 0; i-- {
           cnt[c[pn[i]]]--
           p[cnt[c[pn[i]]]] = pn[i]
       }
       cn[p[0]] = 0
       classes = 1
       for i := 1; i < n; i++ {
           cur := [2]int{c[p[i]], c[(p[i]+shift)%n]}
           prev := [2]int{c[p[i-1]], c[(p[i-1]+shift)%n]}
           if cur != prev {
               classes++
           }
           cn[p[i]] = classes - 1
       }
       c, cn = cn, c
       if classes == n {
           break
       }
   }
   return p
}
