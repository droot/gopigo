package main

import (
  "encoding/json"
  "fmt"
  "io/ioutil"
  "net/http"
  "strings"
  "time"

  "github.com/droot/gopigo"
  "github.com/golang/glog"
  "github.com/gorilla/mux"
  "github.com/kidoman/embd"
)

type server struct {
  bus      embd.I2CBus
  gp       *gopigo.GoPiGo
  cmdChan  chan string
  quitChan chan struct{}
  r        *mux.Router
}

func (s *server) setup() error {
  if s.bus == nil {
    return fmt.Errorf("i2bus is missing")
  }
  s.gp = gopigo.New(s.bus)

  s.cmdChan = make(chan string)
  s.quitChan = make(chan struct{})

  // launch command execution loop.
  go cmdExecutor(s.gp, s.cmdChan, s.quitChan)
  s.r = s.setupRoutes()
  return nil
}

func (s *server) setupRoutes() *mux.Router {
  r := mux.NewRouter().StrictSlash(false)
  r.PathPrefix("/public/").Handler(
    http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

  pi := r.Path("/pi").Subrouter()
  pi.Methods("POST").HandlerFunc(s.handlePICommands)
  return r
}

func (s *server) run() {
  glog.Fatal(http.ListenAndServe(":8070", s.r))
}

type cmdRequest struct {
  Content string `json:"content"`
}

func (s *server) handlePICommands(w http.ResponseWriter, r *http.Request) {
  // parse commands
  body, err := ioutil.ReadAll(r.Body)
  if err != nil {
    glog.Errorf("error reading request body :: %v", err)
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }
  glog.Infof("body received: %s", string(body))

  cmdReq := cmdRequest{}
  err = json.Unmarshal(body, &cmdReq)
  if err != nil {
    glog.Errorf("error parsing command request :: %v", err)
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }
  s.cmdChan <- cmdReq.Content
  w.WriteHeader(http.StatusOK)
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

func main() {

  bus := I2CBus()
  defer cleanupI2Bus()

  s := server{bus: bus}
  if err := s.setup(); err != nil {
    panic(err)
  }

  s.run()
}
