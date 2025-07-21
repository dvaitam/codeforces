package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s1, s2, s3 string
   if _, err := fmt.Fscan(reader, &s1, &s2, &s3); err != nil {
       return
   }
   // map notes to semitone numbers
   noteMap := map[string]int{
       "C":  0,
       "C#": 1,
       "D":  2,
       "D#": 3,
       "E":  4,
       "F":  5,
       "F#": 6,
       "G":  7,
       "G#": 8,
       "A":  9,
       "B":  10,
       "H":  11,
   }
   notes := []int{noteMap[s1], noteMap[s2], noteMap[s3]}
   // all permutations of indices 0,1,2
   perms := [6][3]int{
       {0, 1, 2}, {0, 2, 1},
       {1, 0, 2}, {1, 2, 0},
       {2, 0, 1}, {2, 1, 0},
   }
   for _, p := range perms {
       x := notes[p[0]]
       y := notes[p[1]]
       z := notes[p[2]]
       d1 := (y - x + 12) % 12
       d2 := (z - y + 12) % 12
       if d1 == 4 && d2 == 3 {
           fmt.Println("major")
           return
       }
       if d1 == 3 && d2 == 4 {
           fmt.Println("minor")
           return
       }
   }
   fmt.Println("strange")
}
