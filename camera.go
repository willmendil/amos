package main

import (

  "fmt"
	"log"
	"net/http"
  "image"
  "image/color"
  "bufio"
  "os"
	_ "net/http/pprof"

	"github.com/hybridgroup/mjpeg"

  "gobot.io/x/gobot"
  "gocv.io/x/gocv"
)

var (
	deviceID int
	err      error
	webcam   *gocv.VideoCapture
	stream   *mjpeg.Stream
  width int
  height int
)

type delta struct {
  dx int
  dy int
}

func mjpegCapture() {
	img := gocv.NewMat()
	defer img.Close()

	for {
		if ok := webcam.Read(&img); !ok {
			fmt.Printf("Device closed: %v\n", deviceID)
			return
		}
		if img.Empty() {
			continue
		}


		buf, _ := gocv.IMEncode(".jpg", img)
		stream.UpdateJPEG(buf)
	}
}

func getCenterFace(faceCoord image.Rectangle) delta {

  bwidth := faceCoord.Max.X - faceCoord.Min.X
  x := faceCoord.Min.X + (bwidth/2)
  bheight := faceCoord.Max.Y - faceCoord.Min.Y
  y := faceCoord.Min.Y + (bheight/2)

  shift := delta{ dx:(width/2)-x , dy: (height/2)-y }

  return shift
}

func captureFace(tunnelChan chan delta) {
  xmlFile := "gobotAssets/haarcascade_frontalface_alt.xml"


	img := gocv.NewMat()
	defer img.Close()

  // color for the rect when faces detected
  blue := color.RGBA{0, 0, 255, 0}

  // load classifier to recognize faces
  classifier := gocv.NewCascadeClassifier()
  defer classifier.Close()

  if !classifier.Load(xmlFile) {
      fmt.Printf("Error reading cascade file: %v\n", xmlFile)
      return
  }

	for {
		if ok := webcam.Read(&img); !ok {
			fmt.Printf("Device closed: %v\n", deviceID)
			return
		}
		if img.Empty() {
			continue
		}

    // detect faces
    rects := classifier.DetectMultiScale(img)
    fmt.Printf("found %d faces\n", len(rects))


    if len(rects)> 0 {
      tunnelChan <- getCenterFace(rects[0])
    }


    // draw a rectangle around each face on the original image,
    // along with text identifying as "Human"
    for _, r := range rects {
        gocv.Rectangle(&img, r, blue, 3)
        s := getCenterFace(r)
        text := fmt.Sprintf("dx %v - dy %v Human", s.dx, s.dy)
        size := gocv.GetTextSize(text, gocv.FontHersheyPlain, 1.2, 2)
        pt := image.Pt(r.Min.X+(r.Min.X/2)-(size.X/2), r.Min.Y-2)
        gocv.PutText(&img, text, pt, gocv.FontHersheyPlain, 1.2, blue, 2)
    }



    // delta.dx = int()

    // gocv.Resize(img, &img, image.Point{}, 0.3, 0.3, 0)

		buf, _ := gocv.IMEncode(".jpg", img)
		stream.UpdateJPEG(buf)
	}
}


// readDescriptions reads the descriptions from a file
// and returns a slice of its lines.
func readDescriptions(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func caffe() {
  model := "gobotAssets/bvlc_googlenet.caffemodel"
	config := "gobotAssets/bvlc_googlenet.prototxt"
	descr := "gobotAssets/classification_classes_ILSVRC2012.txt"
	descriptions, err := readDescriptions(descr)
	if err != nil {
		fmt.Printf("Error reading descriptions file: %v\n", descr)
		return
	}

  backend := gocv.NetBackendDefault
	if len(os.Args) > 5 {
		backend = gocv.ParseNetBackend("openvino")
	}

	target := gocv.NetTargetCPU
	if len(os.Args) > 6 {
		target = gocv.ParseNetTarget("fp16")
	}

  img := gocv.NewMat()
  defer img.Close()

  // open DNN classifier
	net := gocv.ReadNet(model, config)
	if net.Empty() {
		fmt.Printf("Error reading network model from : %v %v\n", model, config)
		return
	}
	defer net.Close()

  net.SetPreferableBackend(gocv.NetBackendType(backend))
	net.SetPreferableTarget(gocv.NetTargetType(target))

	status := "Ready"
	statusColor := color.RGBA{0, 255, 0, 0}
	fmt.Printf("Start reading device: %v\n", deviceID)

  for {
		if ok := webcam.Read(&img); !ok {
			fmt.Printf("Device closed: %v\n", deviceID)
			return
		}
		if img.Empty() {
			continue
		}

		// convert image Mat to 224x224 blob that the classifier can analyze
		blob := gocv.BlobFromImage(img, 1.0, image.Pt(224, 224), gocv.NewScalar(104, 117, 123, 0), false, false)

		// feed the blob into the classifier
		net.SetInput(blob, "")

		// run a forward pass thru the network
		prob := net.Forward("")

		// reshape the results into a 1x1000 matrix
		probMat := prob.Reshape(1, 1)

		// determine the most probable classification
		_, maxVal, _, maxLoc := gocv.MinMaxLoc(probMat)

		// display classification
		status = fmt.Sprintf("description: %v, maxVal: %v\n", descriptions[maxLoc.X], maxVal)
		gocv.PutText(&img, status, image.Pt(10, 20), gocv.FontHersheyPlain, 1.2, statusColor, 2)

		blob.Close()
		prob.Close()
		probMat.Close()
    buf, _ := gocv.IMEncode(".jpg", img)
		stream.UpdateJPEG(buf)
	}

}


func streamFunc(tunnelChan chan delta){

  	// parse args
  	deviceID := 0
  	host := "0.0.0.0:8082"

  	// open webcam
  	webcam, err = gocv.OpenVideoCapture(deviceID)
  	if err != nil {
  		fmt.Printf("Error opening capture device: %v\n", deviceID)
  		return
  	}
  	defer webcam.Close()

    width = 320
    height = int(float64(width) *0.75)

    webcam.Set(3, float64(width) )
    webcam.Set(4, float64(height) )
    webcam.Set(5, 5 )
  	// create the mjpeg stream
  	stream = mjpeg.NewStream()

  	// start capturing
  	// go mjpegCapture()
    go captureFace(tunnelChan)
    // go caffe()

  	fmt.Println("Capturing. Point your browser to " + host)

  	// start http server
  	http.Handle("/", stream)
  	log.Fatal(http.ListenAndServe(host, nil))
  }

func Camera(tunnelChan chan delta) *gobot.Robot {
  // window := opencv.NewWindowDriver()
    // camera := opencv.NewCameraDriver(0)

    work := func() {
             streamFunc(tunnelChan)
    }

    robot := gobot.NewRobot("cameraBot",
            // []gobot.Device{camera},
            work,
    )


	return robot
}
