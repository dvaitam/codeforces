package main

import (
   "bufio"
   "fmt"
   "os"
)

// Entry holds count of files containing a substring and a representative index
type Entry struct {
   count int
   rep   int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   fmt.Fscan(reader, &n)
   names := make([]string, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &names[i])
   }

   // map from substring to Entry
   subs := make(map[string]Entry, n*30)
   for idx, s := range names {
       L := len(s)
       for i := 0; i < L; i++ {
           for j := i + 1; j <= L; j++ {
               key := s[i:j]
               if e, ok := subs[key]; ok {
                   e.count++
                   subs[key] = e
               } else {
                   subs[key] = Entry{count: 1, rep: idx}
               }
           }
       }
   }

   var q int
   fmt.Fscan(reader, &q)
   for j := 0; j < q; j++ {
       var t string
       fmt.Fscan(reader, &t)
       if e, ok := subs[t]; ok {
           fmt.Fprintf(writer, "%d %s\n", e.count, names[e.rep])
       } else {
           fmt.Fprintln(writer, "-")
       }
   }
}
