package main

import (
 // "fmt"
        "gobot.io/x/gobot"
        // "gobot.io/x/gobot/api"
)


func main() {

  msg := make(chan delta)
  master := gobot.NewMaster()
  // server := api.NewAPI(master)
  // server.Port = "8081"
  // server.AddHandler(api.BasicAuth("abc", "abc"))
  // server.Start()
  master.AddRobot(Arduino(msg))
  // master.AddRobot(Controller())
  master.AddRobot(Camera(msg))



  master.Start()
}
