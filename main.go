package main

import (
	"flag"
	"fmt"
	"image/png"
	"os"
	"time"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/ean"
	"github.com/boombuler/barcode/qr"
)

var (
	barcodeType    string
	barcodeLink    string
	barcodeContent string
	barcodeWidth   int
	barcodeHeight  int
)

func init() {
	flag.StringVar(&barcodeType, "type", "ean", "set type of barcode which should be performed. I.e. ean or qrcode or code128.")
	flag.StringVar(&barcodeLink, "link", "", "barcode name of a png file. should be uuidv4 or guid.")
	flag.StringVar(&barcodeContent, "content", "", "content for barcode, for ean it should be letters, for qrcode it can be any text.")
	flag.IntVar(&barcodeWidth, "width", 50, "setting width of barcode png image.")
	flag.IntVar(&barcodeHeight, "height", 50, "setting height of barcode png image.")
	flag.Parse()
}

func PerformBarcode(folderName string) {

	var barcodeItem barcode.Barcode
	var errorBarcodeCreation error

	switch barcodeType {
	case "ean":
		func() {
			barcodeItem, errorBarcodeCreation = ean.Encode(barcodeContent)
		}()
	case "qrcode":
		func() {
			barcodeItem, errorBarcodeCreation = qr.Encode(barcodeContent, qr.M, qr.Auto)
		}()
	default:
		writeDownLogMessage("unknown barcode type!", nil)
	}

	if errorBarcodeCreation != nil {
		writeDownLogMessage("error initializing barcode. Check given args.", errorBarcodeCreation)
	}

	file, err := os.Create(fmt.Sprintf("%s/%s.png", folderName, barcodeLink))
	if err != nil {
		writeDownLogMessage("file wasn't created!", err)
	}
	defer file.Close()

	barcodeItem, err = barcode.Scale(barcodeItem, barcodeWidth, barcodeHeight)
	if err != nil {
		writeDownLogMessage(fmt.Sprintf("barcode wasn't scaled probably given incorrect width or height! height is %d; width is %d",
			barcodeHeight, barcodeWidth), err)
	}
	png.Encode(file, barcodeItem)

}

func PrepareFolder() string {
	var t = time.Now()
	timeString := t.Format("2006-1-2")
	path := fmt.Sprintf("%s/%s", "temp", timeString)
	os.MkdirAll(path, 0755)
	return path
}

func main() {
	PerformBarcode(PrepareFolder())
}
