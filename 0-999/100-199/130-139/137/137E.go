package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   n := len(s)
   // Compute prefix A: vowels +1, consonants -2
   A := make([]int, n+1)
   isVowel := func(c byte) bool {
       switch c {
       case 'a', 'e', 'i', 'o', 'u', 'A', 'E', 'I', 'O', 'U':
           return true
       }
       return false
   }
   for i := 1; i <= n; i++ {
       if isVowel(s[i-1]) {
           A[i] = A[i-1] + 1
       } else {
           A[i] = A[i-1] - 2
       }
   }
   // Compress A
   comp := make([]int, n+1)
   vals := make([]int, n+1)
   copy(vals, A)
   sort.Ints(vals)
   uniq := vals[:1]
   for i := 1; i <= n; i++ {
       if vals[i] != uniq[len(uniq)-1] {
           uniq = append(uniq, vals[i])
       }
   }
   m := len(uniq)
   for i := 0; i <= n; i++ {
       // 1-based
       idx := sort.SearchInts(uniq, A[i])
       comp[i] = idx + 1
   }
   // Fenwick tree for suffix min via reversed idx
   size := m
   INF := n + 5
   tree := make([]int, size+1)
   for i := 1; i <= size; i++ {
       tree[i] = INF
   }
   // BIT functions (1..size)
   update := func(pos, v int) {
       for i := pos; i <= size; i += i & -i {
           if v < tree[i] {
               tree[i] = v
           }
       }
   }
   query := func(pos int) int {
       res := INF
       for i := pos; i > 0; i -= i & -i {
           if tree[i] < res {
               res = tree[i]
           }
       }
       return res
   }
   // Insert A[0]
   // reversed index: ridx = m - comp[i] + 1
   ridx0 := size - comp[0] + 1
   update(ridx0, 0)
   maxlen := 0
   count := 0
   for r := 1; r <= n; r++ {
       cr := comp[r]
       rcr := size - cr + 1
       // want minimal l with A[l] >= A[r], i.e., comp[l] >= cr -> reversed idx <= rcr
       lmin := query(rcr)
       if lmin < INF {
           length := r - lmin
           if length > maxlen {
               maxlen = length
               count = 1
           } else if length == maxlen {
               count++
           }
       }
       // insert current r
       posr := size - comp[r] + 1
       update(posr, r)
   }
   if count > 0 {
       fmt.Printf("%d %d", maxlen, count)
   } else {
       fmt.Print("No solution")
   }
}
