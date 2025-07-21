package main

import (
   "bufio"
   "fmt"
   "os"
)

// win determines if the current player wins with numbers a and b
// using the Euclid's game rules: subtract any positive multiple of the smaller
// from the larger, or take the remainder, until one number is zero.
func win(a, b uint64) bool {
   // reversed tracks parity of recursive negations
   var reversed bool
   for {
       if a == 0 || b == 0 {
           // base: no moves, loss for current player; apply parity
           return reversed
       }
       if a > b {
           a, b = b, a
       }
       // if we can subtract more than one multiple, win immediately
       if b/a > 1 {
           // inner returns true, so apply parity
           return !reversed
       }
       // only one move: subtract a once (b = b % a) and toggle
       b = b % a
       reversed = !reversed
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for i := 0; i < t; i++ {
       var a, b uint64
       fmt.Fscan(reader, &a, &b)
       if win(a, b) {
           fmt.Fprintln(writer, "First")
       } else {
           fmt.Fprintln(writer, "Second")
       }
   }
}
