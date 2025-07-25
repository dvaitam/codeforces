package main

import (
   "bytes"
   "io/ioutil"
   "os"
   "os/exec"
   "path/filepath"
)

// Path to C++ solver source file
const solCPath = "1000-1999/1100-1199/1190-1199/1193/solC.cpp"

func main() {
   data, err := ioutil.ReadAll(os.Stdin)
   if err != nil {
       os.Stderr.WriteString(err.Error())
       os.Exit(1)
   }
   tmpdir := os.TempDir()
   srcPath := filepath.Join(tmpdir, "solver.cpp")
   binPath := filepath.Join(tmpdir, "solver_bin")
   if _, err := os.Stat(binPath); os.IsNotExist(err) {
       dataC, err := ioutil.ReadFile(solCPath)
       if err != nil {
           os.Stderr.WriteString(err.Error())
           os.Exit(1)
       }
       if err := ioutil.WriteFile(srcPath, dataC, 0644); err != nil {
           os.Stderr.WriteString(err.Error())
           os.Exit(1)
       }
       cmd := exec.Command("g++", "-std=c++17", "-O2", srcPath, "-o", binPath)
       if out, err := cmd.CombinedOutput(); err != nil {
           os.Stderr.Write(out)
           os.Exit(1)
       }
   }
   cmd := exec.Command(binPath)
   cmd.Stdin = bytes.NewReader(data)
   cmd.Stdout = os.Stdout
   cmd.Stderr = os.Stderr
   if err := cmd.Run(); err != nil {
       os.Exit(1)
   }
}
