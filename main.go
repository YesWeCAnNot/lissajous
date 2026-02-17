package main

import (
	"image"
	"image/color"
	"image/gif"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
)

var allColors = []color.Color{
	color.RGBA{0xFF, 0xFF, 0xFF, 0xFF},
	color.RGBA{0x00, 0x00, 0x00, 0xFF},
	color.RGBA{0xFF, 0x00, 0x00, 0xFF},
	color.RGBA{0x00, 0xFF, 0x00, 0xFF},
	color.RGBA{0x00, 0x00, 0xFF, 0xFF},
	color.RGBA{0xFF, 0xFF, 0x00, 0xFF},
	color.RGBA{0xFF, 0x00, 0xFF, 0xFF},
	color.RGBA{0x00, 0xFF, 0xFF, 0xFF},
	color.RGBA{0xFF, 0xA5, 0x00, 0xFF},
	color.RGBA{0x80, 0x00, 0x80, 0xFF},
}

const (
	bgIndex = 0 // first color in palette
)

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	cylcesParam := r.URL.Query().Get("cycles")
	if cylcesParam == "" {
		cylcesParam = "5"
	}
	cycles, err := strconv.Atoi(cylcesParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sizeParam := r.URL.Query().Get("size")
	if sizeParam == "" {
		sizeParam = "400"
	}

	size, err := strconv.Atoi(sizeParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	lissajous(w, float64(cycles), size)

}

func lissajous(out io.Writer, cycle float64, siz int) {
	const (
		res     = 0.001
		nframes = 64
		delay   = 1
	)
	var cycles float64 = cycle
	var size int = siz

	bgIndex := rand.Intn(len(allColors))
	bgColor := allColors[bgIndex]

	var lineColor color.Color
	for {
		candidate := allColors[rand.Intn(len(allColors))]
		if candidate != bgColor {
			lineColor = candidate
			break
		}
	}

	palette := []color.Color{bgColor, lineColor}

	freq := rand.Float64() * 2.0
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0 // фаза
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*float64(size)+0.5), size+int(y*float64(size)+0.5), 1)
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim)
}
