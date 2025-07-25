package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   s := make([]int, m)
   h := make([]int, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(in, &s[i], &h[i])
       s[i]--
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   for r := 0; r < n; r++ {
       leftIdx, rightIdx := -1, -1
       for i := 0; i < m; i++ {
           if s[i] < r {
               if leftIdx == -1 || s[i] < s[leftIdx] {
                   leftIdx = i
               }
           }
           if s[i] > r {
               if rightIdx == -1 || s[i] > s[rightIdx] {
                   rightIdx = i
               }
           }
       }
       ok := true
       if leftIdx != -1 {
           sum, need := 0, 0
           for j := s[leftIdx] + 1; j <= r; j++ {
               v := a[j]
               if v < 0 {
                   req := -v - sum
                   if req > need {
                       need = req
                   }
               }
               sum += v
           }
           if h[leftIdx] < need {
               ok = false
           }
       }
       if !ok {
           continue
       }
       if rightIdx != -1 {
           sum, need := 0, 0
           for j := s[rightIdx] - 1; j >= r; j-- {
               v := a[j]
               if v < 0 {
                   req := -v - sum
                   if req > need {
                       need = req
                   }
               }
               sum += v
           }
           if h[rightIdx] < need {
               ok = false
           }
       }
       if !ok {
           continue
       }
       lefts := make([]int, 0)
       mids := make([]int, 0)
       rights := make([]int, 0)
       for i := 0; i < m; i++ {
           if s[i] < r {
               lefts = append(lefts, i)
           } else if s[i] > r {
               rights = append(rights, i)
           } else {
               mids = append(mids, i)
           }
       }
       // sort lefts asc and rights desc
       for i := 0; i < len(lefts); i++ {
           for j := i + 1; j < len(lefts); j++ {
               if s[lefts[j]] < s[lefts[i]] {
                   lefts[i], lefts[j] = lefts[j], lefts[i]
               }
           }
       }
       for i := 0; i < len(rights); i++ {
           for j := i + 1; j < len(rights); j++ {
               if s[rights[j]] > s[rights[i]] {
                   rights[i], rights[j] = rights[j], rights[i]
               }
           }
       }
       order := make([]int, 0, m)
       order = append(order, lefts...)
       order = append(order, mids...)
       order = append(order, rights...)
       w := bufio.NewWriter(os.Stdout)
       fmt.Fprintln(w, r+1)
       for i, idx := range order {
           if i > 0 {
               w.WriteByte(' ')
           }
           fmt.Fprintf(w, "%d", idx+1)
       }
       w.WriteByte('\n')
       w.Flush()
       return
   }
   fmt.Println(-1)
}
