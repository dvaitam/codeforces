package main

import (
   "bufio"
   "fmt"
   "os"
)

// Solve the problem as described in problemE.txt (CF 204E)
func main() {
   in := bufio.NewReader(os.Stdin)
   var n, k int
   if _, err := fmt.Fscan(in, &n, &k); err != nil {
       return
   }
   strs := make([]string, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &strs[i])
   }
   // Edge cases
   if k > n {
       // no substring can appear in at least k strings
       for i := 0; i < n; i++ {
           if i > 0 {
               fmt.Print(" ")
           }
           fmt.Print(0)
       }
       fmt.Println()
       return
   }
   if k == 1 {
       // every substring counts
       for i, s := range strs {
           L := int64(len(s))
           total := L * (L + 1) / 2
           if i > 0 {
               fmt.Print(" ")
           }
           fmt.Print(total)
       }
       fmt.Println()
       return
   }
   // Build Suffix Automaton over all strings with unique separators
   // Map letters 'a'-'z' to 0..25, separators to 26..26+n-1
   totalLen := 0
   codesList := make([][]int, n)
   for i, s := range strs {
       m := len(s)
       codes := make([]int, m)
       for j := 0; j < m; j++ {
           codes[j] = int(s[j]-'a')
       }
       codesList[i] = codes
       totalLen += m + 1
   }
   // SAM data structures
   maxStates := 2*totalLen + 1
   link := make([]int, maxStates)
   length := make([]int, maxStates)
   nexts := make([]map[int]int, maxStates)
   isGood := make([]bool, maxStates)
   sets := make([]map[int]struct{}, maxStates)
   size := 1
   last := 0
   link[0] = -1
   for i := range nexts {
       nexts[i] = nil
   }
   posState := make([]int, totalLen)
   posStrID := make([]int, totalLen)
   // extend function
   extend := func(c, sid int) int {
       cur := size
       size++
       length[cur] = length[last] + 1
       nexts[cur] = make(map[int]int)
       // init
       p := last
       for p >= 0 {
           if _, ok := nexts[p][c]; ok {
               break
           }
           nexts[p][c] = cur
           p = link[p]
       }
       if p == -1 {
           link[cur] = 0
       } else {
           q := nexts[p][c]
           if length[p]+1 == length[q] {
               link[cur] = q
           } else {
               clone := size
               size++
               length[clone] = length[p] + 1
               // copy transitions
               nexts[clone] = make(map[int]int)
               for k2, v2 := range nexts[q] {
                   nexts[clone][k2] = v2
               }
               link[clone] = link[q]
               for p2 := p; p2 >= 0; p2 = link[p2] {
                   if nexts[p2][c] != q {
                       break
                   }
                   nexts[p2][c] = clone
               }
               link[q] = clone
               link[cur] = clone
           }
       }
       last = cur
       posState[sid] = cur
       return cur
   }
   // Build
   idx := 0
   for sid, codes := range codesList {
       for _, c := range codes {
           extend(c, idx)
           posStrID[idx] = sid + 1
           idx++
       }
       sep := 26 + sid
       extend(sep, idx)
       posStrID[idx] = 0
       idx++
   }
   // Initial assignment of sets
   for i := 0; i < idx; i++ {
       st := posState[i]
       sid := posStrID[i]
       if sid == 0 {
           continue
       }
       if isGood[st] {
           continue
       }
       if sets[st] == nil {
           sets[st] = make(map[int]struct{})
       }
       if _, exists := sets[st][sid]; !exists {
           sets[st][sid] = struct{}{}
           if len(sets[st]) >= k {
               isGood[st] = true
               sets[st] = nil
           }
       }
   }
   // Prepare order by state length
   maxLen := 0
   for i := 0; i < size; i++ {
       if length[i] > maxLen {
           maxLen = length[i]
       }
   }
   bucket := make([]int, maxLen+1)
   for i := 0; i < size; i++ {
       bucket[length[i]]++
   }
   for i := 1; i <= maxLen; i++ {
       bucket[i] += bucket[i-1]
   }
   order := make([]int, size)
   for i := size - 1; i >= 0; i-- {
       l := length[i]
       bucket[l]--
       order[bucket[l]] = i
   }
   // Merge sets in suffix link tree
   for i := size - 1; i > 0; i-- {
       v := order[i]
       p := link[v]
       if p < 0 {
           continue
       }
       if isGood[p] {
           continue
       }
       // merge v into p
       if isGood[v] {
           isGood[p] = true
           sets[p] = nil
       } else if sets[v] != nil {
           if sets[p] == nil {
               // adopt
               if len(sets[v]) >= k {
                   isGood[p] = true
               } else {
                   sets[p] = sets[v]
               }
           } else {
               // merge smaller into larger
               if len(sets[p]) < len(sets[v]) {
                   sets[p], sets[v] = sets[v], sets[p]
               }
               for sid := range sets[v] {
                   if _, ex := sets[p][sid]; !ex {
                       sets[p][sid] = struct{}{}
                       if len(sets[p]) >= k {
                           isGood[p] = true
                           sets[p] = nil
                           break
                       }
                   }
               }
           }
       }
   }
   // For each string, compute answer
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   for idxStr, codes := range codesList {
       var p, l int
       var res int64
       for _, c := range codes {
           // extend match
           if nexts[p] != nil && nexts[p][c] != 0 {
               p = nexts[p][c]
               l++
           } else {
               for p >= 0 && (nexts[p] == nil || nexts[p][c] == 0) {
                   p = link[p]
               }
               if p < 0 {
                   p = 0
                   l = 0
               } else {
                   l = length[p] + 1
                   p = nexts[p][c]
               }
           }
           // shrink until good
           for p > 0 && !isGood[p] {
               p = link[p]
               l = length[p]
           }
           res += int64(l)
       }
       if idxStr > 0 {
           fmt.Fprint(out, ' ')
       }
       fmt.Fprint(out, res)
   }
   fmt.Fprintln(out)
}
