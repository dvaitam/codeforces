package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   // Tree1 backward strings
   s1 := []string{""}
   count1 := make(map[string]int)
   count1[""] = 1
   // Tree2 forward strings from root
   t2 := []string{""}
   parent2 := []int{-1}
   count2 := make(map[string]int)
   count2[""] = 1
   var total int64 = 1

   for i := 0; i < n; i++ {
       var t, v int
       var cs string
       fmt.Fscan(reader, &t, &v, &cs)
       c := cs
       v-- // zero-based
       if t == 1 {
           // add to tree1
           newID := len(s1)
           // backward path: char + parent string
           s := c + s1[v]
           s1 = append(s1, s)
           count1[s]++
           if cnt2, ok := count2[s]; ok {
               total += int64(cnt2)
           }
       } else {
           // add to tree2
           newID := len(t2)
           parent2 = append(parent2, v)
           // forward path: parent string + char
           tstr := t2[v] + c
           t2 = append(t2, tstr)
           // for all ancestors j of newID, including itself
           for j := newID; j >= 0; j = parent2[j] {
               prefixLen := len(t2[j])
               s := tstr[prefixLen:]
               count2[s]++
               if cnt1, ok := count1[s]; ok {
                   total += int64(cnt1)
               }
               if j == 0 {
                   break
               }
           }
       }
       fmt.Fprintln(writer, total)
   }
}
