package main

import (
  "bufio"
  "flag"
  "fmt"
  "io"
  "os"
  "strings"

  "github.com/droot/gopigo"
  "github.com/kidoman/embd"

  _ "github.com/kidoman/embd/host/all"
  "time"
)

var (
  cmdFilePath = flag.String("cmd_file", "", "file containing GoPiGo commands")
)

func main() {
  flag.Parse()

  if err := embd.InitI2C(); err != nil {
    panic(err)
  }
  defer embd.CloseI2C()

  bus := embd.NewI2CBus(1)

  // create GoPiGo instance.
  gp := gopigo.New(bus)

  cmdChan := make(chan string)
  quitChan := make(chan struct{})

  // launch command execution loop.
  go cmdExecutor(gp, cmdChan, quitChan)

  var reader io.Reader
  if *cmdFilePath != "" {
    if f, err := os.Open(*cmdFilePath); err != nil {
      fmt.Errorf("error opening file :: %v", err)
      return
    } else {
      reader = f
      defer f.Close()
    }
  } else {
    reader = os.Stdin
  }

  go cmdReader(reader, cmdChan, quitChan)

  <-quitChan
}

func cmdReader(reader io.Reader, cmds chan string, quitChan chan struct{}) {
  scanner := bufio.NewScanner(reader)
  for scanner.Scan() {
    cmd := scanner.Text()
    if cmd == "quit" {
      close(quitChan)
      break
    }
    // valid command, send it for execution.
    cmds <- scanner.Text()
  }
  if scanner.Err() != nil {
    fmt.Printf("error scanning from stdin :: %v\n", scanner.Err())
  }
}

func cmdExecutor(gp *gopigo.GoPiGo, cmds chan string, quitChan chan struct{}) {
  var err error
  for {
    select {
    case <-quitChan:
      return
    case cmd := <-cmds:
      cmdParts := strings.Split(cmd, " ")
      if len(cmdParts) <= 0 {
        fmt.Println("invalid command")
        continue
      }
      switch cmdParts[0] {
      case "f":
        var dist int
        fmt.Sscanf(cmd, "f %d", &dist)
        err = gp.Fwd(dist)
      case "b":
        var dist int
        fmt.Sscanf(cmd, "b %d", &dist)
        err = gp.Bwd(dist)
      case "sleep":
        var secs int
        fmt.Sscanf(cmd, "sleep %d", &secs)
        time.Sleep(time.Duration(secs) * time.Second)
      case "s":
        err = gp.Stop()
      case "l":
        err = gp.Left()
      case "lr":
        err = gp.LeftRotate()
      case "r":
        err = gp.Right()
      case "rr":
        err = gp.RightRotate()
      case "tl":
        var degrees float64
        fmt.Sscanf(cmd, "tl %f", &degrees)
        err = gp.TurnLeft(degrees)
      case "tr":
        var degrees float64
        fmt.Sscanf(cmd, "tr %f", &degrees)
        err = gp.TurnRight(degrees)
      case "volt":
        volt, err := gp.BatteryVoltage()
        if err != nil {
          fmt.Println("error reading battery voltage: ", err)
        }
        fmt.Println("battery voltage: ", volt)
      default:
        fmt.Println("invalid cmd received:", cmd, ":")
      }
      if err != nil {
        fmt.Errorf("error in executing command: ", cmd)
      }
    }
  }
}
