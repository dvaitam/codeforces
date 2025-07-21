package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var grid [4]string
   for i := 0; i < 4; i++ {
       if _, err := fmt.Fscan(in, &grid[i]); err != nil {
           return
       }
   }

   for i := 0; i < 3; i++ {
       for j := 0; j < 3; j++ {
           cnt := 0
           if grid[i][j] == '#' {
               cnt++
           }
           if grid[i][j+1] == '#' {
               cnt++
           }
           if grid[i+1][j] == '#' {
               cnt++
           }
           if grid[i+1][j+1] == '#' {
               cnt++
           }
           if cnt != 2 {
               fmt.Fprintln(out, "YES")
               return
           }
       }
   }
   fmt.Fprintln(out, "NO")
}
