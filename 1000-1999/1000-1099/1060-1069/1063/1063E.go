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
   fmt.Fscan(reader, &n)
   goPos := make([]int, n)
   for i := 0; i < n; i++ {
       var a int
       fmt.Fscan(reader, &a)
       a--
       goPos[a] = i
   }
   // check if identity
   dame := false
   for i := 0; i < n; i++ {
       if goPos[i] != i {
           dame = true
           break
       }
   }
   // template line of dots
   base := make([]byte, n)
   for i := range base {
       base[i] = '.'
   }
   if !dame {
       fmt.Fprintln(writer, n)
       line := string(base)
       for i := 0; i < n; i++ {
           fmt.Fprintln(writer, line)
       }
       return
   }
   // set empty at position 0
   goPos[0] = -1
   // n-1 moves
   fmt.Fprintln(writer, n-1)
   for h := 0; h < n; h++ {
       // find leftmost and rightmost misplaced
       dL, dR := -1, -1
       for i := 0; i < n; i++ {
           if goPos[i] != i {
               dL = i
               break
           }
       }
       for i := n - 1; i >= 0; i-- {
           if goPos[i] != i {
               dR = i
               break
           }
       }
       // generate line
       gen := make([]byte, n)
       copy(gen, base)
       if dL == dR {
           // nothing
       } else if goPos[dL] == -1 {
           t := goPos[dR]
           gen[dL] = '/'
           gen[dR] = '/'
           gen[t] = '/'
           // update mapping
           goPos[dL] = goPos[t]
           goPos[t] = goPos[dR]
           goPos[dR] = -1
       } else {
           t := goPos[dL]
           gen[dR] = '\\'
           gen[dL] = '\\'
           gen[t] = '\\'
           goPos[dR] = goPos[t]
           goPos[t] = goPos[dL]
           goPos[dL] = -1
       }
       fmt.Fprintln(writer, string(gen))
   }
}
