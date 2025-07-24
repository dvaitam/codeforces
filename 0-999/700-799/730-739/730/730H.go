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

   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   names := make([]string, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &names[i])
   }
   dels := make([]int, m)
   toDel := make([]bool, n)
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &dels[i])
       dels[i]--
       toDel[dels[i]] = true
   }
   // All target filenames must have the same length
   L := len(names[dels[0]])
   for i := 1; i < m; i++ {
       if len(names[dels[i]]) != L {
           fmt.Fprintln(writer, "No")
           return
       }
   }
   // Build pattern
   pattern := []byte(names[dels[0]])
   for i := 1; i < m; i++ {
       cur := []byte(names[dels[i]])
       for j := 0; j < L; j++ {
           if pattern[j] != cur[j] {
               pattern[j] = '?'
           }
       }
   }
   // Verify no other filename matches the pattern
   for i := 0; i < n; i++ {
       if toDel[i] {
           continue
       }
       if len(names[i]) != L {
           continue
       }
       ok := true
       for j := 0; j < L; j++ {
           if pattern[j] != '?' && pattern[j] != names[i][j] {
               ok = false
               break
           }
       }
       if ok {
           fmt.Fprintln(writer, "No")
           return
       }
   }
   fmt.Fprintln(writer, "Yes")
   fmt.Fprintln(writer, string(pattern))
}
