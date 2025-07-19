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
   if _, err := fmt.Fscan(reader, &n); err != nil {
       fmt.Fprintln(writer, -1)
       return
   }
   // Precompute sequences up to n+1
   seq := make([][]int, 0, n+2)
   seq = append(seq, []int{0})
   seq = append(seq, []int{1})
   seq = append(seq, []int{0, 1})
   seq = append(seq, []int{1, 0, 1})
   for len(seq) < n+2 {
       if !build(&seq) {
           fmt.Fprintln(writer, -1)
           return
       }
   }
   // Output first polynomial: degree n, coefficients seq[n+1]
   fmt.Fprintln(writer, n)
   for i, v := range seq[n+1] {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, v)
   }
   writer.WriteByte('\n')
   // Output second polynomial: degree n-1, coefficients seq[n]
   fmt.Fprintln(writer, n-1)
   for i, v := range seq[n] {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, v)
   }
   writer.WriteByte('\n')
}

// build extends the sequence with the next row, returns false on failure
func build(seq *[][]int) bool {
   v := *seq
   a := v[len(v)-2]
   b := v[len(v)-1]
   // shift b right and prepend 0
   newB := make([]int, len(b)+1)
   copy(newB[1:], b)
   mult := 1
   for i := 0; i < len(a); i++ {
       if mult == -1 && ((a[i] == -1 && newB[i] == 1) || (a[i] == 1 && newB[i] == -1)) {
           return false
       }
       if (a[i] == 1 && newB[i] == 1) || (a[i] == -1 && newB[i] == -1) {
           mult = -1
       }
   }
   for i := 0; i < len(a); i++ {
       newB[i] += mult * a[i]
   }
   *seq = append(*seq, newB)
   return true
}
