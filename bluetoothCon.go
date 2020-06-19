package main

import (
        "fmt"
        "time"

        "gobot.io/x/gobot"
        "gobot.io/x/gobot/platforms/ble"
)

func main() {
        bleAdaptor := ble.NewClientAdaptor("5C:BA:37:0F:1A:E0")
        battery := ble.NewBatteryDriver(bleAdaptor)

        work := func() {
                gobot.Every(5*time.Second, func() {
                        fmt.Println("Battery level:", battery.GetBatteryLevel())
                })
        }

        robot := gobot.NewRobot("bleBot",
                []gobot.Connection{bleAdaptor},
                []gobot.Device{battery},
                work,
        )

        robot.Start()
}
