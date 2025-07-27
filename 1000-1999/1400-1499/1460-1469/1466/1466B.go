package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

func main() {
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   // Instead of ReadBytes for each int, use bufio.Scanner for simplicity
   // Reinitialize scanner
   scanner := bufio.NewScanner(os.Stdin)
   scanner.Split(bufio.ScanWords)
   var nextInt func() int
   nextInt = func() int {
       if !scanner.Scan() {
           return 0
       }
       v, _ := strconv.Atoi(scanner.Text())
       return v
   }

   t := nextInt()
   for tc := 0; tc < t; tc++ {
       n := nextInt()
       // process sequence
       last := 0
       count := 0
       for i := 0; i < n; i++ {
           x := nextInt()
           if x > last {
               last = x
               count++
           } else if x+1 > last {
               last = x + 1
               count++
           }
       }
       fmt.Fprintln(writer, count)
   }
}
