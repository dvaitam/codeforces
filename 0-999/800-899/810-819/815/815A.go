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

   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   a := make([][]int64, n)
   for i := 0; i < n; i++ {
       a[i] = make([]int64, m)
       for j := 0; j < m; j++ {
           fmt.Fscan(reader, &a[i][j])
       }
   }
   // compute minimal operations for columns
   r := make([]int64, n)
   c := make([]int64, m)
   // minimal value in first row
   rmn := a[0][0]
   for j := 0; j < m; j++ {
       if a[0][j] < rmn {
           rmn = a[0][j]
       }
   }
   for j := 0; j < m; j++ {
       c[j] = a[0][j] - rmn
   }
   // minimal value in first column
   cmn := a[0][0]
   for i := 0; i < n; i++ {
       if a[i][0] < cmn {
           cmn = a[i][0]
       }
   }
   for i := 0; i < n; i++ {
       r[i] = a[i][0] - cmn
   }
   // base difference
   df := a[0][0] - r[0] - c[0]
   // check consistency
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           if a[i][j]-r[i]-c[j] != df {
               fmt.Fprintln(writer, -1)
               return
           }
       }
   }
   // distribute extra operations
   if n < m {
       for i := 0; i < n; i++ {
           r[i] += df
       }
   } else {
       for j := 0; j < m; j++ {
           c[j] += df
       }
   }
   // total operations
   var ans int64
   for i := 0; i < n; i++ {
       ans += r[i]
   }
   for j := 0; j < m; j++ {
       ans += c[j]
   }
   fmt.Fprintln(writer, ans)
   // print operations
   for i := 0; i < n; i++ {
       for k := int64(0); k < r[i]; k++ {
           fmt.Fprintf(writer, "row %d\n", i+1)
       }
   }
   for j := 0; j < m; j++ {
       for k := int64(0); k < c[j]; k++ {
           fmt.Fprintf(writer, "col %d\n", j+1)
       }
   }
}
