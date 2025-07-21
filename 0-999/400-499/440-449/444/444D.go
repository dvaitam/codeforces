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

   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   n := len(s)
   // map substrings of length <=4 to their start positions
   substrPos := make(map[string][]int)
   // build positions
   for i := 0; i < n; i++ {
       var buf [4]byte
       for l := 1; l <= 4 && i+l <= n; l++ {
           buf[l-1] = s[i+l-1]
           key := string(buf[:l])
           substrPos[key] = append(substrPos[key], i)
       }
   }
   var q int
   fmt.Fscan(reader, &q)
   cache := make(map[string]int)
   for i := 0; i < q; i++ {
       var a, b string
       fmt.Fscan(reader, &a, &b)
       key := a + "#" + b
       if ans, ok := cache[key]; ok {
           fmt.Fprintln(writer, ans)
           continue
       }
       // same pattern: single occurrence suffices
       if a == b {
           if list, ok := substrPos[a]; !ok || len(list) == 0 {
               cache[key] = -1
               fmt.Fprintln(writer, -1)
           } else {
               cache[key] = len(a)
               fmt.Fprintln(writer, len(a))
           }
           continue
       }
       A, okA := substrPos[a]
       B, okB := substrPos[b]
       if !okA || !okB {
           cache[key] = -1
           fmt.Fprintln(writer, -1)
           continue
       }
       la, lb := len(a), len(b)
       // choose smaller occurrences list to iterate
       var S, L []int
       var lenS, lenL int
       if len(A) < len(B) {
           S, L = A, B
           lenS, lenL = la, lb
       } else {
           S, L = B, A
           lenS, lenL = lb, la
       }
       ans := n + 1
       // merge-like scan using binary search
       for _, x := range S {
           idx := sort.Search(len(L), func(i int) bool { return L[i] >= x })
           // check right neighbor
           if idx < len(L) {
               y := L[idx]
               start := x
               // compute end = max(x+lenS-1, y+lenL-1)
               end := x + lenS - 1
               if y+lenL-1 > end {
                   end = y + lenL - 1
               }
               if cur := end - start + 1; cur < ans {
                   ans = cur
               }
           }
           // check left neighbor
           if idx > 0 {
               y := L[idx-1]
               start := y
               end := x + lenS - 1
               if y+lenL-1 > end {
                   end = y + lenL - 1
               }
               if cur := end - start + 1; cur < ans {
                   ans = cur
               }
           }
       }
       if ans == n+1 {
           ans = -1
       }
       cache[key] = ans
       fmt.Fprintln(writer, ans)
   }
}
