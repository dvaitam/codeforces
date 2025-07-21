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
   var pattern string
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   if _, err := fmt.Fscan(reader, &pattern); err != nil {
       return
   }
   pokemons := []string{"vaporeon", "jolteon", "flareon", "espeon", "umbreon", "leafeon", "glaceon", "sylveon"}
   for _, p := range pokemons {
       if len(p) != n {
           continue
       }
       match := true
       for i := 0; i < n; i++ {
           if pattern[i] != '.' && pattern[i] != p[i] {
               match = false
               break
           }
       }
       if match {
           fmt.Fprintln(writer, p)
           return
       }
   }
}
