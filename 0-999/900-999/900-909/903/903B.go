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

   var h1, a1, c1 int
   fmt.Fscan(reader, &h1, &a1, &c1)
   var h2, a2 int
   fmt.Fscan(reader, &h2, &a2)

   Ineed := h2 / a1
   if h2%a1 != 0 {
       Ineed++
   }
   HeNeed := h1 / a2
   if h1%a2 != 0 {
       HeNeed++
   }
   actions := make([]string, 0)

   // Heal until hero can survive enough strikes
   for Ineed > HeNeed {
       actions = append(actions, "HEAL")
       // Hero heals
       h1 += c1
       // Monster strikes
       h1 -= a2
       // Recompute needed hero turns
       HeNeed = h1 / a2
       if h1%a2 != 0 {
           HeNeed++
       }
   }
   // Then strike Ineed times
   for i := 0; i < Ineed; i++ {
       actions = append(actions, "STRIKE")
   }

   // Output
   fmt.Fprintln(writer, len(actions))
   for _, op := range actions {
       fmt.Fprintln(writer, op)
   }
}
