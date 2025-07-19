package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   n     int
   a     []int
   tk    []bool
   seq1  []int
   seq2  []int
   writer *bufio.Writer
)

// Print prints the sequence with spaces and newline
func Print(v []int) {
   for i, x := range v {
       if i > 0 {
           writer.WriteByte(' ')
       }
       writer.WriteString(fmt.Sprint(x))
   }
   writer.WriteByte('\n')
}

// Check returns true if v is empty? false. If size<=2 always true, else checks constant difference
func Check(v []int) bool {
   m := len(v)
   if m == 0 {
       return false
   }
   if m <= 2 {
       return true
   }
   d := v[1] - v[0]
   for i := 2; i < m; i++ {
       if v[i]-v[i-1] != d {
           return false
       }
   }
   return true
}

// Get tries to build one arithmetic progression starting from a[l] with diff a[r]-a[l]
// and checks if the remainder is also an arithmetic progression. Prints and returns true on success.
func Get(l, r int) bool {
   seq1 = seq1[:0]
   // reset taken flags
   for i := 0; i < n; i++ {
       tk[i] = false
   }
   nd := a[l]
   delt := a[r] - a[l]
   ind := -1
   for i := 0; i < n; i++ {
       if a[i] == nd {
           seq1 = append(seq1, a[i])
           ind = i
           tk[i] = true
           nd += delt
       }
   }
   seq2 = seq2[:0]
   for i := 0; i < n; i++ {
       if !tk[i] {
           seq2 = append(seq2, a[i])
       }
   }
   if Check(seq2) {
       Print(seq1)
       Print(seq2)
       return true
   }
   // rollback last element
   if ind >= 0 {
       tk[ind] = false
       seq1 = seq1[:len(seq1)-1]
   }
   seq2 = seq2[:0]
   for i := 0; i < n; i++ {
       if !tk[i] {
           seq2 = append(seq2, a[i])
       }
   }
   if Check(seq2) {
       Print(seq1)
       Print(seq2)
       return true
   }
   return false
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer = bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   fmt.Fscan(reader, &n)
   a = make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   tk = make([]bool, n)
   seq1 = make([]int, 0, n)
   seq2 = make([]int, 0, n)
   if n == 2 {
       Print([]int{a[0]})
       Print([]int{a[1]})
       return
   }
   if Get(0, 1) || Get(0, 2) || Get(1, 2) {
       return
   }
   writer.WriteString("No solution\n")
}
