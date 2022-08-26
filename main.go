package main

import (
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/codabar"
	"github.com/boombuler/barcode/code128"
	"github.com/boombuler/barcode/code39"
	"github.com/boombuler/barcode/code93"
	"github.com/boombuler/barcode/datamatrix"
	"github.com/boombuler/barcode/ean"
	"github.com/boombuler/barcode/pdf417"
	"github.com/boombuler/barcode/qr"
	"github.com/boombuler/barcode/twooffive"
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
	flag.StringVar(&barcodeContent, "content", "", "content for barcode, for ean it should be digits, for qrcode it can be any text.")
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
		barcodeItem, errorBarcodeCreation = ean.Encode(barcodeContent)
	case "qrcode":
		barcodeItem, errorBarcodeCreation = qr.Encode(barcodeContent, qr.M, qr.Auto)
	case "codabar":
		barcodeItem, errorBarcodeCreation = codabar.Encode(barcodeContent)
	case "code128":
		barcodeItem, errorBarcodeCreation = code128.Encode(barcodeContent)
	case "code39":
		barcodeItem, errorBarcodeCreation = code39.Encode(barcodeContent, true, true)
	case "code93":
		barcodeItem, errorBarcodeCreation = code93.Encode(barcodeContent, true, true)
	case "datamatrix":
		barcodeItem, errorBarcodeCreation = datamatrix.Encode(barcodeContent)
	case "pdf417":
		barcodeItem, errorBarcodeCreation = pdf417.Encode(barcodeContent, 4)
	case "2of5":
		barcodeItem, errorBarcodeCreation = twooffive.Encode(barcodeContent, false)
	case "2of5interleaved":
		barcodeItem, errorBarcodeCreation = twooffive.Encode(barcodeContent, true)
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
