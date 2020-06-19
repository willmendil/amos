package main

import (
         // "time"
        "fmt"
        "gobot.io/x/gobot"
        "gobot.io/x/gobot/drivers/gpio"
        "gobot.io/x/gobot/platforms/firmata"
)

func sq(a int ) int{
  x := float64(a)
    return int((0.0003*x*x + 0.0105*x + 0.1608)*0.25)
    // return int(0.0012*x*x - 0.0417*x + 0.5944)

}

func stepperInit(stepper *gpio.StepperDriver, firmataAdaptor *firmata.Adaptor) {

  fmt.Println("Init Stepper")
  for  {
    hit, err := firmataAdaptor.DigitalRead("4")
    if err != nil {
      fmt.Println("ERROR")
    }
    if hit == 0{
      stepper.Halt()
      stepper.Move(-30)
      return
    }else{
      stepper.Move(10)
    }
  }

}

func stepperF(msg chan delta, stepper *gpio.StepperDriver , servoY *gpio.ServoDriver) {
  stepper.SetSpeed(300)
  threshold := 30
  n := 0
  m := 130
  for s := range msg{
    p := s.dx
    //set spped
    switch {
    case s.dx > threshold && n < 0:

           if err := stepper.Move(p); err != nil {
             fmt.Println(err)
           }
         n = n +p

       case s.dx < -threshold && n > -1000:

           if err := stepper.Move(p); err != nil {
             fmt.Println(err)
           }
           n = n+p
}
           switch {
           case s.dy > 10 && m > 20:

                       for i := 0; i < sq(m); i ++{
                         servoY.Move(uint8(m))
                         m++
                         // time.Sleep(100*time.Millisecond)
                         }
                           // servoY.Move(uint8(m))
                           // m = m + sq(m)
                     case s.dy < -10 && m < 160:
                       for i := 0; i < sq(m); i ++{
                         servoY.Move(uint8(m))
                         m--
                         // time.Sleep(100*time.Millisecond)
                         }
                           // servoY.Move(uint8(m))
                           // m = m - sq(m)

}
}
}

func Arduino(msg chan delta) *gobot.Robot {
        // maxX := 160
        // minX := 60
          firmataAdaptor := firmata.NewAdaptor("/dev/ttyACM0")
          servoY := gpio.NewServoDriver(firmataAdaptor, "3")
        	stepper := gpio.NewStepperDriver(firmataAdaptor, [4]string{"7", "9", "8", "10"}, gpio.StepperModes.DualPhaseStepping, 2048)
          // stepperSwitch := gpio.NewButtonDriver(firmataAdaptor, "4")
          work := func() {
            servoY.Center()
            stepperInit(stepper, firmataAdaptor)
            go stepperF(msg,  stepper, servoY )

        }

          	robot := gobot.NewRobot("stepperBot",
          		[]gobot.Connection{firmataAdaptor},
              []gobot.Device{servoY},
          		[]gobot.Device{stepper},
          		work,
          	)

       return robot
}



// func Arduino(msg chan delta) *gobot.Robot {
//         // maxX := 160
//         // minX := 60
//         firmataAdaptor := firmata.NewAdaptor("/dev/ttyACM0")
//         servoX := gpio.NewServoDriver(firmataAdaptor, "2")
//         servoY := gpio.NewServoDriver(firmataAdaptor, "3")
//
//         fmt.Println(msg)
//         n := 100
//         m := 100
//         threshold := 20
//         work := func() {
//           servoX.Center()
//           servoY.Center()
//           for s := range msg{
//             fmt.Println(s, n, m)
//             switch {
//             case s.dx > threshold && n > 20:
//                   for i := 0; i < sq(n); i ++{
//                     servoX.Move(uint8(n))
//                     n++
//                     // time.Sleep(20*time.Millisecond)
//                     }
//
//             case s.dx < -threshold && n < 160:
//
//               for i := 0; i < sq(n); i ++{
//                 servoX.Move(uint8(n))
//                 n--
//                 // time.Sleep(20*time.Millisecond)
//                 }
//
//                   // servoX.Move(uint8(n))
//                   // n = n - sq(n)
// }
//                   switch {
//             case s.dy > 20 && m > 20:
//
//               for i := 0; i < sq(m); i ++{
//                 servoY.Move(uint8(m))
//                 m++
//                 // time.Sleep(100*time.Millisecond)
//                 }
//                   // servoY.Move(uint8(m))
//                   // m = m + sq(m)
//             case s.dy < -20 && m < 160:
//               for i := 0; i < sq(m); i ++{
//                 servoY.Move(uint8(m))
//                 m--
//                 // time.Sleep(100*time.Millisecond)
//                 }
//                   // servoY.Move(uint8(m))
//                   // m = m - sq(m)
//           }
//           }
//         }
//
//         robot := gobot.NewRobot("bot",
//                 []gobot.Connection{firmataAdaptor},
//                 []gobot.Device{servoX},
//                 []gobot.Device{servoY},
//                 work,
//         )
//
//        return robot
// }
