package main

import (
        "path"
        "runtime"
        "sync/atomic"
        "time"

        "gobot.io/x/gobot"
        "gobot.io/x/gobot/platforms/opencv"
        "gocv.io/x/gocv"
)

var img atomic.Value

func Camera() *gobot.Robot {
  _, currentfile, _, _ := runtime.Caller(0)
   cascade := path.Join(path.Dir(currentfile), "gobotAssets/haarcascade_frontalface_alt.xml")

   // window := opencv.NewWindowDriver()
   camera := opencv.NewCameraDriver(0)

   work := func() {
           mat := gocv.NewMat()
           img.Store(mat)

           camera.On(opencv.Frame, func(data interface{}) {
                   i := data.(gocv.Mat)
                   img.Store(i)
           })

           gobot.Every(10*time.Millisecond, func() {
                   i := img.Load().(gocv.Mat)
                   if i.Empty() {
                           return
                   }
                   faces := opencv.DetectObjects(cascade, i)
                   opencv.DrawRectangles(i, faces, 0, 255, 0, 5)
                   // window.ShowImage(i)
                   // window.WaitKey(1)
           })
   }

   robot := gobot.NewRobot("faceBot",
           []gobot.Connection{},
           []gobot.Device{camera},
           work,
   )
        return robot
}
