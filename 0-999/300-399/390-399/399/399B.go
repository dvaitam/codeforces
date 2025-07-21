package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   // Read n
   line, err := reader.ReadString('\n')
   if err != nil && err.Error() != "EOF" {
       fmt.Fprintln(os.Stderr, err)
       return
   }
   fields := strings.Fields(line)
   if len(fields) < 1 {
       return
   }
   n, _ := strconv.Atoi(fields[0])
   // Read stack string
   sline, err := reader.ReadString('\n')
   if err != nil && err.Error() != "EOF" {
       fmt.Fprintln(os.Stderr, err)
       return
   }
   s := strings.TrimSpace(sline)
   if len(s) != n {
       // invalid input length
       return
   }
   // Use byte slice for mutability
   b := []byte(s)
   ops := 0
   for {
       // Check for any blue balls
       hasB := false
       for i := 0; i < n; i++ {
           if b[i] == 'B' {
               hasB = true
               break
           }
       }
       if !hasB {
           break
       }
       // Count leading reds
       k := 0
       for k < n && b[k] == 'R' {
           k++
       }
       // Build new stack
       nb := make([]byte, n)
       // Push k blues
       for i := 0; i < k; i++ {
           nb[i] = 'B'
       }
       // Change the next blue to red
       nb[k] = 'R'
       // Copy remaining balls beyond k
       for i := k + 1; i < n; i++ {
           nb[i] = b[i]
       }
       b = nb
       ops++
   }
   fmt.Println(ops)
}
