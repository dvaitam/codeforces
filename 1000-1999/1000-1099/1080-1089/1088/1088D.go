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

   // ask function: query with c, d and return response
   ask := func(c, d int) int {
       fmt.Fprintf(writer, "%d %d\n", c, d)
       writer.Flush()
       var res int
       if _, err := fmt.Fscan(reader, &res); err != nil {
           os.Exit(0)
       }
       return res
   }

   var a, b int
   // initial response for (0,0)
   p := ask(0, 0)
   // build bits from highest to lowest
   for i := 29; i >= 0; i-- {
       bit := 1 << i
       // test flipping both bits
       t := ask(a|bit, b|bit)
       if t == p {
           // bits are equal: either both 0 or both 1
           // distinguish by flipping only a's bit
           u := ask(a|bit, b)
           if u == -1 {
               // both bits are 1
               a |= bit
               b |= bit
           }
           // p remains unchanged
       } else {
           // bits differ
           if t == 1 {
               // a_i = 1, b_i = 0
               a |= bit
           } else {
               // a_i = 0, b_i = 1
               b |= bit
           }
           // update p with current a,b
           p = ask(a, b)
       }
   }
   // output final answer
   fmt.Fprintf(writer, "! %d %d\n", a, b)
}
