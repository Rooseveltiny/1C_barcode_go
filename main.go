package main

import (
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/ean"
	"github.com/boombuler/barcode/qr"
)

var (
	barcodeType    string
	barcodeLink    string
	barcodePath    string
	barcodeContent string
	barcodeWidth   int
	barcodeHeight  int
)

func init() {
	flag.StringVar(&barcodeType, "type", "ean", "set type of barcode which should be performed. I.e. ean or qrcode or code128.")
	flag.StringVar(&barcodeLink, "link", "", "barcode name of a png file. should be uuidv4 or guid.")
	flag.StringVar(&barcodePath, "path", "", "path to prepeared image in some folder of a file system.")
	flag.StringVar(&barcodeContent, "content", "", "content for barcode, for ean it should be letters, for qrcode it can be any text.")
	flag.IntVar(&barcodeWidth, "width", 100, "setting width of barcode png image.")
	flag.IntVar(&barcodeHeight, "height", 100, "setting height of barcode png image.")
	flag.Parse()
}

func PerformScale(bc *barcode.Barcode) {
	var err error
	*bc, err = barcode.Scale(*bc, barcodeWidth, barcodeHeight)

	if err != nil {
		writeDownLogMessage(fmt.Sprintf("barcode wasn't scaled probably given incorrect width or height! height is %d; width is %d",
			barcodeHeight, barcodeWidth), err)
	}

}

func PerformFileSaving(bc image.Image) {
	file, err := os.Create(fmt.Sprintf("%s/%s", barcodePath, barcodeLink))
	if err != nil {
		writeDownLogMessage("file wasn't created!", err)
	}
	defer file.Close()

	png.Encode(file, bc)
}

func PerformBarcode() {

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

	// scale given barcode
	PerformScale(&barcodeItem)
	// make a content text underneath a given barcode
	var barcodeToSave image.Image
	if barcodeItem.Metadata().CodeKind != "QR Code" {
		// to make a subtitle beneath a barcode
		barcodeToSave = subtitleBarcode(barcodeItem)
	} else {
		// to make qrcode 24 bit png image
		img := image.NewRGBA(image.Rect(0, 0, barcodeItem.Bounds().Dx(), barcodeItem.Bounds().Dy()))
		draw.Draw(img, image.Rect(0, 0, barcodeItem.Bounds().Dx(), barcodeItem.Bounds().Dy()), barcodeItem, barcodeItem.Bounds().Min, draw.Over)
		barcodeToSave = img
	}

	// save created barcode to png
	PerformFileSaving(barcodeToSave)

}

func PrepareFolder() {
	os.MkdirAll(barcodePath, 0755)
}

func main() {
	PrepareFolder()
	PerformBarcode()
}
