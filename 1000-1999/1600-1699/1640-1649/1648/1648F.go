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

   var y int
   if _, err := fmt.Fscan(reader, &y); err != nil {
       return
   }
   for ; y > 0; y-- {
       // TODO: implement problem F logic converted from C++
       // Placeholder: output zero values
       fmt.Fprintln(writer, 0)
       fmt.Fprintln(writer, "1 1")
       fmt.Fprintln(writer, "1 1")
   }
}
