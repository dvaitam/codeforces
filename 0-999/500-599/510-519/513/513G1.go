package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   var k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   p := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &p[i])
   }
   // Generate all permutations of 1..n
   perms := make([][]int, 0)
   perm := make([]int, n)
   used := make([]bool, n+1)
   var dfs func(int)
   dfs = func(d int) {
       if d == n {
           cp := make([]int, n)
           copy(cp, perm)
           perms = append(perms, cp)
           return
       }
       for v := 1; v <= n; v++ {
           if !used[v] {
               used[v] = true
               perm[d] = v
               dfs(d + 1)
               used[v] = false
           }
       }
   }
   dfs(0)
   S := len(perms)
   // Map permutation to index
   idxMap := make(map[string]int, S)
   for i, pr := range perms {
       var sb strings.Builder
       for j, v := range pr {
           if j > 0 {
               sb.WriteByte(',')
           }
           sb.WriteString(strconv.Itoa(v))
       }
       idxMap[sb.String()] = i
   }
   // Find start index
   var sb0 strings.Builder
   for i, v := range p {
       if i > 0 {
           sb0.WriteByte(',')
       }
       sb0.WriteString(strconv.Itoa(v))
   }
   start := idxMap[sb0.String()]
   // Precompute inversion counts
   invCount := make([]int, S)
   for i, pr := range perms {
       cnt := 0
       for a := 0; a < n; a++ {
           for b := a + 1; b < n; b++ {
               if pr[a] > pr[b] {
                   cnt++
               }
           }
       }
       invCount[i] = cnt
   }
   // Precompute intervals
   intervals := make([][2]int, 0, n*(n+1)/2)
   for l := 0; l < n; l++ {
       for r := l; r < n; r++ {
           intervals = append(intervals, [2]int{l, r})
       }
   }
   M := len(intervals)
   // Precompute transitions
   neighbors := make([][]int, S)
   for i, pr := range perms {
       neighbors[i] = make([]int, M)
       for t, lr := range intervals {
           l, r := lr[0], lr[1]
           tmp := make([]int, n)
           copy(tmp, pr)
           // reverse [l..r]
           for x, y := l, r; x < y; x, y = x+1, y-1 {
               tmp[x], tmp[y] = tmp[y], tmp[x]
           }
           var sb strings.Builder
           for j, v := range tmp {
               if j > 0 {
                   sb.WriteByte(',')
               }
               sb.WriteString(strconv.Itoa(v))
           }
           neighbors[i][t] = idxMap[sb.String()]
       }
   }
   // Probability distribution
   cur := make([]float64, S)
   cur[start] = 1.0
   // Perform k operations
   for step := 0; step < k; step++ {
       next := make([]float64, S)
       for i, prob := range cur {
           if prob == 0 {
               continue
           }
           pseg := prob / float64(M)
           for _, j := range neighbors[i] {
               next[j] += pseg
           }
       }
       cur = next
   }
   // Compute expected inversions
   var exp float64
   for i, prob := range cur {
       exp += prob * float64(invCount[i])
   }
   // Output
   fmt.Printf("%.10f", exp)
}
