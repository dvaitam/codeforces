package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   A := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &A[i])
   }
   B := make([]int, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &B[i])
   }
   sort.Ints(A)
   sort.Ints(B)
   // Build blocks of consecutive monsters
   type block struct{ a, c, length int }
   blocks := make([]block, 0, n)
   for i := 0; i < n; i++ {
       start := A[i]
       j := i
       for j+1 < n && A[j+1] == A[j]+1 {
           j++
       }
       blocks = append(blocks, block{a: start, c: A[j], length: j - i + 1})
       i = j
   }
   total := 0
   infL := -2000000000
   infR := 2000000000
   // For each block, compute its interval and solve
   for i, b := range blocks {
       L := infL
       R := infR
       if i > 0 {
           L = blocks[i-1].c + 1
       }
       if i+1 < len(blocks) {
           R = blocks[i+1].a - 1
       }
       // find special cells in [L, R]
       l := sort.Search(m, func(j int) bool { return B[j] >= L })
       r := sort.Search(m, func(j int) bool { return B[j] > R }) - 1
       if l >= m || r < l {
           continue
       }
       // two pointers on B[l..r]
       best := 0
       right := l
       maxDist := b.length - 1
       for left := l; left <= r; left++ {
           for right <= r && B[right]-B[left] <= maxDist {
               right++
           }
           cnt := right - left
           if cnt > best {
               best = cnt
           }
       }
       total += best
   }
   fmt.Fprintln(writer, total)
}
