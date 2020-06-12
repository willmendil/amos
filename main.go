package main

import (

        "gobot.io/x/gobot"
        "gobot.io/x/gobot/api"
)


func main() {
        master := gobot.NewMaster()
        api.NewAPI(master).Start()

        master.AddRobot(Controller())
				master.AddRobot(Camera())

        master.Start()
}
