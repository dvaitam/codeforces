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
   var k float64
   if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
       return
   }
   skills := make(map[string]int)
   for i := 0; i < n; i++ {
       var name string
       var exp int
       fmt.Fscan(reader, &name, &exp)
       level := int(k*float64(exp) + 1e-9)
       if level >= 100 {
           skills[name] = level
       }
   }
   for i := 0; i < m; i++ {
       var name string
       fmt.Fscan(reader, &name)
       if _, exists := skills[name]; !exists {
           skills[name] = 0
       }
   }
   names := make([]string, 0, len(skills))
   for name := range skills {
       names = append(names, name)
   }
   sort.Strings(names)
   fmt.Fprintln(writer, len(names))
   for _, name := range names {
       fmt.Fprintf(writer, "%s %d\n", name, skills[name])
   }
}
