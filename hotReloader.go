package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

func runCommand(cmd *exec.Cmd, watchDir string) {
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Dir = watchDir
    err := cmd.Start()
    if err != nil {
        fmt.Println("Error Running Command:", cmd.String())
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Command success:", cmd.String())
}

func buildAndStartModule( watchDir string, builtModuleName string ) *exec.Cmd {
    buildCmd := exec.Command("go", "build", watchDir)
    runCommand(buildCmd, watchDir)
    buildCmd.Wait()

    runCmd := exec.Command("./" + builtModuleName)
    runCommand(runCmd, watchDir)
    return runCmd
}

func main() {
    numArgs := len(os.Args[1:])

    if numArgs < 2 {
        fmt.Println("")
    }
        
    watchDir := os.Args[1] 
    builtModuleName := os.Args[2]
    pollWaitTimeMs := time.Duration(2000)
    if numArgs >= 3  {
        enteredTime, err := strconv.Atoi(os.Args[3])
        if err == nil {
            pollWaitTimeMs = time.Duration(enteredTime)
        } else {
            fmt.Println("Error: Couldn't parse entered time, defaulting to 2000ms.", err)
        }
    }

    cmd := buildAndStartModule( watchDir, builtModuleName )

    fileModTimes := make(map[string]time.Time)

    for {
        filepath.Walk(watchDir, func(path string, info os.FileInfo, err error) error {
            if err != nil {
                fmt.Println("Error walking directory:", err)
                return err
            }
            if !info.IsDir() && filepath.Ext(info.Name()) == ".go" {
                if modTime, ok := fileModTimes[path]; ok {
                    if modTime.Before(info.ModTime()) {
                        fmt.Println("File modified:", path)
                        cmd.Process.Kill()
                        fmt.Println("Module stopped.")
                        cmd.Wait()
                        fmt.Println("Restarting the module...")
                        cmd = buildAndStartModule(watchDir, builtModuleName)
                    }
                }
                fileModTimes[path] = info.ModTime()
            }
            return nil
        })

        time.Sleep(pollWaitTimeMs * time.Millisecond)
    }
}
