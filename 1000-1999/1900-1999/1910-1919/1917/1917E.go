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

   var T int
   fmt.Fscan(reader, &T)
   for t := 0; t < T; t++ {
       if !work(reader, writer) {
           fmt.Fprintln(writer, "No")
       }
   }
}

func work(r *bufio.Reader, w *bufio.Writer) bool {
   var n, k int
   fmt.Fscan(r, &n, &k)
   if k%2 == 1 {
       return false
   }
   // initialize matrix 1-indexed
   a := make([][]int, n+1)
   for i := 1; i <= n; i++ {
       a[i] = make([]int, n+1)
   }
   if k%4 == 0 {
       fmt.Fprintln(w, "Yes")
       rem := k
       for i := 1; i+1 <= n && rem > 0; i += 2 {
           for j := 1; j+1 <= n && rem > 0; j += 2 {
               for x := i; x < i+2; x++ {
                   for y := j; y < j+2; y++ {
                       a[x][y] = 1
                   }
               }
               rem -= 4
           }
       }
       printMatrix(w, a, n)
       return true
   }
   if k >= 6 && k <= n*n-10 && k%4 == 2 {
       fmt.Fprintln(w, "Yes")
       rem := k - 6
       for i := 1; i+1 <= n && rem > 0; i += 2 {
           startJ := 1
           if i <= 3 {
               startJ = 4
           }
           for j := startJ; j+1 <= n && rem > 0; j += 2 {
               for x := i; x < i+2; x++ {
                   for y := j; y < j+2; y++ {
                       a[x][y] = 1
                   }
               }
               rem -= 4
           }
       }
       a[1][1] = 1
       a[1][2] = 1
       a[2][1] = 1
       a[2][3] = 1
       a[3][2] = 1
       a[3][3] = 1
       printMatrix(w, a, n)
       return true
   }
   if k == n {
       fmt.Fprintln(w, "Yes")
       for i := 1; i <= n; i++ {
           for j := 1; j <= n; j++ {
               if j == 1 {
                   fmt.Fprint(w, "1 ")
               } else {
                   fmt.Fprint(w, "0 ")
               }
           }
           fmt.Fprintln(w)
       }
       return true
   }
   if k == n*n-6 {
       fmt.Fprintln(w, "Yes")
       a[1][1] = 1
       a[1][2] = 1
       a[2][1] = 1
       a[2][3] = 1
       a[3][1] = 1
       a[3][2] = 1
       a[3][3] = 1
       a[3][4] = 1
       a[4][1] = 1
       a[4][4] = 1
       for i := 1; i <= n; i++ {
           for j := 1; j <= n; j++ {
               if i > 4 || j > 4 {
                   fmt.Fprint(w, "1 ")
               } else {
                   fmt.Fprintf(w, "%d ", a[i][j])
               }
           }
           fmt.Fprintln(w)
       }
       return true
   }
   return false
}

func printMatrix(w *bufio.Writer, a [][]int, n int) {
   for i := 1; i <= n; i++ {
       for j := 1; j <= n; j++ {
           fmt.Fprintf(w, "%d ", a[i][j])
       }
       fmt.Fprintln(w)
   }
}
