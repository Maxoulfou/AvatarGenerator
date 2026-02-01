package main

import (
	"crypto/sha256"
	"encoding/hex"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	defaultSize = 64
	maxSize     = 128
)

func main() {
	http.HandleFunc("/avatar", avatarHandler)

	addr := ":8080"
	log.Printf("avatar service listening on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func avatarHandler(w http.ResponseWriter, r *http.Request) {
	input := strings.TrimSpace(r.URL.Query().Get("input"))
	if input == "" {
		http.Error(w, "missing input query parameter", http.StatusBadRequest)
		return
	}

	size := defaultSize
	if sizeParam := r.URL.Query().Get("size"); sizeParam != "" {
		parsed, err := strconv.Atoi(sizeParam)
		if err != nil {
			http.Error(w, "invalid size", http.StatusBadRequest)
			return
		}
		size = parsed
	}
	if size != defaultSize && size != maxSize {
		http.Error(w, "size must be 64 or 128", http.StatusBadRequest)
		return
	}

	timeKey, err := resolveTimeKey(r.URL.Query().Get("timestamp"))
	if err != nil {
		http.Error(w, "invalid timestamp", http.StatusBadRequest)
		return
	}

	hash := hashInput(input, timeKey)
	img := generateAvatar(hash, size)

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("X-Avatar-Hash", hex.EncodeToString(hash))
	w.Header().Set("X-Avatar-Time-Key", timeKey)
	if err := png.Encode(w, img); err != nil {
		http.Error(w, "failed to encode image", http.StatusInternalServerError)
		return
	}
}

func resolveTimeKey(raw string) (string, error) {
	if raw == "" {
		return time.Now().UTC().Format("2006-01-02"), nil
	}
	parsed, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return "", err
	}
	return time.Unix(parsed, 0).UTC().Format("2006-01-02"), nil
}

func hashInput(input string, timeKey string) []byte {
	h := sha256.Sum256([]byte(input + ":" + timeKey))
	return h[:]
}

func generateAvatar(hash []byte, size int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	rng := newByteRNG(hash)
	background := blendColor(pickColor(rng, backgroundPalette), 0.08)
	draw.Draw(img, img.Bounds(), &image.Uniform{C: background}, image.Point{}, draw.Src)

	center := image.Point{X: size / 2, Y: size / 2}
	headRadius := int(float64(size) * (0.32 + 0.06*float64(rng.nextInt(4))))
	skin := pickColor(rng, skinPalette)
	hair := pickColor(rng, hairPalette)
	eye := pickColor(rng, eyePalette)
	mouth := pickColor(rng, mouthPalette)
	highlight := blendColor(skin, 0.2)
	accessory := pickColor(rng, accessoryPalette)
	brow := pickColor(rng, eyebrowPalette)
	blush := pickColor(rng, blushPalette)
	neck := pickColor(rng, neckPalette)
	clothing := pickColor(rng, clothingPalette)
	accent := pickColor(rng, accentPalette)
	scar := pickColor(rng, scarPalette)
	mask := pickColor(rng, maskPalette)
	lip := pickColor(rng, lipPalette)
	shadow := pickColor(rng, shadowPalette)
	frame := pickColor(rng, framePalette)
	mark := pickColor(rng, markPalette)
	hood := pickColor(rng, hoodPalette)
	irisHighlight := pickColor(rng, irisHighlightPalette)
	cape := pickColor(rng, capePalette)

	drawFilledCircle(img, center, headRadius, skin)
	drawFilledCircle(img, image.Point{X: center.X - headRadius/3, Y: center.Y + headRadius/5}, headRadius/6, highlight)
	drawBackgroundGradient(img, background, accent)
	drawHair(img, center, headRadius, hair, rng)
	drawHairStrands(img, center, headRadius, blendColor(hair, 0.2), rng)
	drawSideburns(img, center, headRadius, hair, rng)
	drawNeck(img, center, headRadius, neck)
	drawCape(img, center, headRadius, cape, rng)
	drawShoulders(img, center, headRadius, clothing, accent, rng)
	drawBackgroundAccents(img, center, headRadius, accent, rng)
	drawFrameBorder(img, frame)
	drawAccessories(img, center, headRadius, accessory, skin, rng)
	drawMask(img, center, headRadius, mask, rng)
	drawEyes(img, center, headRadius, eye, rng)
	drawIrisHighlights(img, center, headRadius, irisHighlight, rng)
	drawEyebrows(img, center, headRadius, brow, rng)
	drawNose(img, center, headRadius)
	drawBlush(img, center, headRadius, blush, rng)
	drawScar(img, center, headRadius, scar, rng)
	drawMouth(img, center, headRadius, mouth, rng)
	drawLipShine(img, center, headRadius, lip, rng)
	drawMustache(img, center, headRadius, hair, rng)
	drawChinShadow(img, center, headRadius, shadow, rng)
	drawForeheadMark(img, center, headRadius, mark, rng)
	drawHood(img, center, headRadius, hood, rng)
	applyVignette(img, center, int(float64(size)*0.48))
	applyNoise(img, rng, size/2)

	return img
}

type byteRNG struct {
	data []byte
	idx  int
}

func newByteRNG(seed []byte) *byteRNG {
	return &byteRNG{data: seed}
}

func (r *byteRNG) nextByte() byte {
	b := r.data[r.idx%len(r.data)]
	r.idx++
	return b
}

func (r *byteRNG) nextInt(max int) int {
	if max <= 0 {
		return 0
	}
	return int(r.nextByte()) % max
}

var (
	skinPalette = []color.RGBA{
		{R: 241, G: 194, B: 125, A: 255},
		{R: 224, G: 172, B: 105, A: 255},
		{R: 198, G: 134, B: 66, A: 255},
		{R: 141, G: 85, B: 36, A: 255},
		{R: 255, G: 220, B: 180, A: 255},
		{R: 205, G: 133, B: 63, A: 255},
	}
	hairPalette = []color.RGBA{
		{R: 45, G: 34, B: 30, A: 255},
		{R: 71, G: 52, B: 39, A: 255},
		{R: 120, G: 90, B: 60, A: 255},
		{R: 200, G: 160, B: 120, A: 255},
		{R: 35, G: 30, B: 50, A: 255},
	}
	eyePalette = []color.RGBA{
		{R: 36, G: 70, B: 142, A: 255},
		{R: 84, G: 52, B: 32, A: 255},
		{R: 32, G: 102, B: 66, A: 255},
		{R: 70, G: 70, B: 70, A: 255},
	}
	mouthPalette = []color.RGBA{
		{R: 141, G: 62, B: 62, A: 255},
		{R: 128, G: 50, B: 80, A: 255},
		{R: 160, G: 72, B: 92, A: 255},
	}
	accessoryPalette = []color.RGBA{
		{R: 60, G: 60, B: 60, A: 255},
		{R: 220, G: 180, B: 90, A: 255},
		{R: 180, G: 200, B: 220, A: 255},
		{R: 120, G: 160, B: 200, A: 255},
		{R: 200, G: 120, B: 140, A: 255},
	}
	eyebrowPalette = []color.RGBA{
		{R: 40, G: 32, B: 28, A: 255},
		{R: 70, G: 52, B: 38, A: 255},
		{R: 110, G: 80, B: 50, A: 255},
		{R: 160, G: 120, B: 90, A: 255},
		{R: 25, G: 25, B: 35, A: 255},
	}
	blushPalette = []color.RGBA{
		{R: 238, G: 168, B: 168, A: 220},
		{R: 230, G: 150, B: 140, A: 220},
		{R: 210, G: 120, B: 130, A: 220},
		{R: 240, G: 180, B: 190, A: 220},
	}
	neckPalette = []color.RGBA{
		{R: 236, G: 190, B: 126, A: 255},
		{R: 217, G: 168, B: 104, A: 255},
		{R: 190, G: 132, B: 70, A: 255},
		{R: 135, G: 84, B: 45, A: 255},
	}
	clothingPalette = []color.RGBA{
		{R: 52, G: 86, B: 136, A: 255},
		{R: 88, G: 120, B: 76, A: 255},
		{R: 170, G: 92, B: 92, A: 255},
		{R: 60, G: 60, B: 70, A: 255},
		{R: 120, G: 68, B: 144, A: 255},
		{R: 180, G: 132, B: 60, A: 255},
	}
	accentPalette = []color.RGBA{
		{R: 255, G: 210, B: 90, A: 255},
		{R: 210, G: 90, B: 120, A: 255},
		{R: 90, G: 170, B: 200, A: 255},
		{R: 90, G: 200, B: 140, A: 255},
		{R: 200, G: 200, B: 200, A: 255},
	}
	scarPalette = []color.RGBA{
		{R: 160, G: 90, B: 90, A: 255},
		{R: 140, G: 70, B: 70, A: 255},
		{R: 120, G: 60, B: 60, A: 255},
	}
	maskPalette = []color.RGBA{
		{R: 235, G: 235, B: 235, A: 230},
		{R: 210, G: 220, B: 230, A: 230},
		{R: 190, G: 210, B: 220, A: 230},
		{R: 220, G: 200, B: 210, A: 230},
	}
	lipPalette = []color.RGBA{
		{R: 166, G: 72, B: 98, A: 255},
		{R: 190, G: 90, B: 110, A: 255},
		{R: 140, G: 60, B: 82, A: 255},
		{R: 120, G: 45, B: 70, A: 255},
		{R: 200, G: 120, B: 140, A: 255},
	}
	shadowPalette = []color.RGBA{
		{R: 90, G: 72, B: 62, A: 120},
		{R: 110, G: 92, B: 82, A: 120},
		{R: 70, G: 58, B: 50, A: 120},
	}
	framePalette = []color.RGBA{
		{R: 30, G: 30, B: 30, A: 255},
		{R: 220, G: 210, B: 190, A: 255},
		{R: 80, G: 90, B: 120, A: 255},
		{R: 180, G: 140, B: 80, A: 255},
		{R: 90, G: 120, B: 90, A: 255},
	}
	markPalette = []color.RGBA{
		{R: 220, G: 90, B: 90, A: 200},
		{R: 90, G: 160, B: 220, A: 200},
		{R: 120, G: 200, B: 140, A: 200},
		{R: 200, G: 180, B: 100, A: 200},
	}
	hoodPalette = []color.RGBA{
		{R: 55, G: 65, B: 90, A: 220},
		{R: 90, G: 80, B: 70, A: 220},
		{R: 70, G: 90, B: 80, A: 220},
		{R: 100, G: 60, B: 80, A: 220},
	}
	irisHighlightPalette = []color.RGBA{
		{R: 255, G: 255, B: 255, A: 200},
		{R: 230, G: 240, B: 255, A: 200},
		{R: 255, G: 240, B: 230, A: 200},
	}
	capePalette = []color.RGBA{
		{R: 40, G: 60, B: 120, A: 200},
		{R: 120, G: 60, B: 40, A: 200},
		{R: 50, G: 90, B: 70, A: 200},
		{R: 100, G: 40, B: 80, A: 200},
	}
	backgroundPalette = []color.RGBA{
		{R: 232, G: 244, B: 255, A: 255},
		{R: 255, G: 240, B: 234, A: 255},
		{R: 240, G: 255, B: 244, A: 255},
		{R: 244, G: 240, B: 255, A: 255},
	}
)

func pickColor(rng *byteRNG, palette []color.RGBA) color.RGBA {
	return palette[rng.nextInt(len(palette))]
}

func blendColor(c color.RGBA, factor float64) color.RGBA {
	apply := func(v uint8) uint8 {
		return uint8(float64(v) + (255.0-float64(v))*factor)
	}
	return color.RGBA{R: apply(c.R), G: apply(c.G), B: apply(c.B), A: c.A}
}

func drawFilledCircle(img *image.RGBA, center image.Point, radius int, fill color.RGBA) {
	r2 := radius * radius
	for y := center.Y - radius; y <= center.Y+radius; y++ {
		for x := center.X - radius; x <= center.X+radius; x++ {
			dx := x - center.X
			dy := y - center.Y
			if dx*dx+dy*dy <= r2 {
				img.Set(x, y, fill)
			}
		}
	}
}

func drawHair(img *image.RGBA, center image.Point, radius int, hair color.RGBA, rng *byteRNG) {
	height := int(float64(radius) * (0.55 + 0.1*float64(rng.nextInt(3))))
	top := center.Y - radius
	for y := top; y < top+height; y++ {
		for x := center.X - radius; x <= center.X+radius; x++ {
			dx := x - center.X
			dy := y - (center.Y - radius/2)
			if dx*dx+dy*dy <= radius*radius {
				img.Set(x, y, hair)
			}
		}
	}
}

func drawAccessories(img *image.RGBA, center image.Point, radius int, accessory color.RGBA, skin color.RGBA, rng *byteRNG) {
	switch rng.nextInt(5) {
	case 0:
		drawGlasses(img, center, radius, accessory, rng)
	case 1:
		drawHat(img, center, radius, accessory, rng)
	case 2:
		drawEarrings(img, center, radius, accessory)
	case 3:
		drawFreckles(img, center, radius, blendColor(skin, 0.4), rng)
	default:
		drawBeard(img, center, radius, accessory, rng)
	}
}

func drawBackgroundGradient(img *image.RGBA, base color.RGBA, accent color.RGBA) {
	bounds := img.Bounds()
	for y := 0; y < bounds.Dy(); y++ {
		t := float64(y) / float64(bounds.Dy())
		blend := color.RGBA{
			R: uint8(float64(base.R)*(1-t) + float64(accent.R)*t),
			G: uint8(float64(base.G)*(1-t) + float64(accent.G)*t),
			B: uint8(float64(base.B)*(1-t) + float64(accent.B)*t),
			A: 255,
		}
		for x := 0; x < bounds.Dx(); x++ {
			img.Set(x, y, blend)
		}
	}
}

func drawHairStrands(img *image.RGBA, center image.Point, radius int, hair color.RGBA, rng *byteRNG) {
	count := 8 + rng.nextInt(6)
	for i := 0; i < count; i++ {
		startX := center.X - radius + rng.nextInt(radius*2)
		startY := center.Y - radius + rng.nextInt(radius/2)
		length := radius/2 + rng.nextInt(radius/2)
		for y := 0; y < length; y++ {
			img.Set(startX, startY+y, hair)
		}
	}
}

func drawSideburns(img *image.RGBA, center image.Point, radius int, hair color.RGBA, rng *byteRNG) {
	if rng.nextInt(2) == 0 {
		return
	}
	width := radius / 6
	height := radius / 2
	leftX := center.X - radius + width
	rightX := center.X + radius - width
	topY := center.Y - radius/4
	for y := topY; y < topY+height; y++ {
		for x := leftX; x < leftX+width; x++ {
			img.Set(x, y, hair)
		}
		for x := rightX - width; x < rightX; x++ {
			img.Set(x, y, hair)
		}
	}
}

func drawCape(img *image.RGBA, center image.Point, radius int, cape color.RGBA, rng *byteRNG) {
	if rng.nextInt(3) != 0 {
		return
	}
	width := radius * 2
	height := radius
	startY := center.Y + radius + radius/4
	for y := startY; y < startY+height; y++ {
		offset := (y - startY) / 2
		for x := center.X - width/2 - offset; x <= center.X+width/2+offset; x++ {
			img.Set(x, y, cape)
		}
	}
}

func drawNeck(img *image.RGBA, center image.Point, radius int, neck color.RGBA) {
	width := radius / 2
	height := radius / 2
	startX := center.X - width/2
	startY := center.Y + radius/2
	for y := startY; y < startY+height; y++ {
		for x := startX; x < startX+width; x++ {
			img.Set(x, y, neck)
		}
	}
}

func drawShoulders(img *image.RGBA, center image.Point, radius int, clothing color.RGBA, accent color.RGBA, rng *byteRNG) {
	width := radius * 2
	height := radius / 2
	startY := center.Y + radius
	for y := startY; y < startY+height; y++ {
		for x := center.X - width/2; x <= center.X+width/2; x++ {
			img.Set(x, y, clothing)
		}
	}
	if rng.nextInt(2) == 0 {
		drawChevron(img, image.Point{X: center.X, Y: startY + height/3}, width/2, height/3, accent)
	} else {
		drawStripe(img, image.Point{X: center.X, Y: startY + height/3}, width/2, height/3, accent)
	}
}

func drawBackgroundAccents(img *image.RGBA, center image.Point, radius int, accent color.RGBA, rng *byteRNG) {
	switch rng.nextInt(6) {
	case 0:
		drawOrbitRings(img, center, radius, accent)
	case 1:
		drawStars(img, rng, radius, accent)
	case 2:
		drawHexGrid(img, center, radius, accent, rng)
	case 3:
		drawCircuitTrace(img, center, radius, accent, rng)
	case 4:
		drawConstellation(img, center, radius, accent, rng)
	default:
		drawAurora(img, center, radius, accent, rng)
	}
}

func drawFrameBorder(img *image.RGBA, stroke color.RGBA) {
	bounds := img.Bounds()
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		img.Set(x, bounds.Min.Y, stroke)
		img.Set(x, bounds.Max.Y-1, stroke)
	}
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		img.Set(bounds.Min.X, y, stroke)
		img.Set(bounds.Max.X-1, y, stroke)
	}
	drawCornerTicks(img, stroke, 6)
}

func drawEyes(img *image.RGBA, center image.Point, radius int, eye color.RGBA, rng *byteRNG) {
	offsetX := radius / 2
	offsetY := radius / 5
	eyeRadius := int(float64(radius) * 0.12)
	white := color.RGBA{R: 248, G: 248, B: 248, A: 255}
	pupilRadius := int(float64(eyeRadius) * 0.6)
	eyeShift := rng.nextInt(3) - 1

	left := image.Point{X: center.X - offsetX + eyeShift, Y: center.Y - offsetY}
	right := image.Point{X: center.X + offsetX + eyeShift, Y: center.Y - offsetY}
	drawFilledCircle(img, left, eyeRadius, white)
	drawFilledCircle(img, right, eyeRadius, white)
	drawFilledCircle(img, left, pupilRadius, eye)
	drawFilledCircle(img, right, pupilRadius, eye)
}

func drawIrisHighlights(img *image.RGBA, center image.Point, radius int, highlight color.RGBA, rng *byteRNG) {
	offsetX := radius / 2
	offsetY := radius / 5
	size := radius / 12
	shift := rng.nextInt(2)
	left := image.Point{X: center.X - offsetX + shift, Y: center.Y - offsetY - shift}
	right := image.Point{X: center.X + offsetX + shift, Y: center.Y - offsetY - shift}
	drawFilledCircle(img, left, size, highlight)
	drawFilledCircle(img, right, size, highlight)
}

func drawEyebrows(img *image.RGBA, center image.Point, radius int, brow color.RGBA, rng *byteRNG) {
	width := radius / 2
	height := radius / 10
	offsetX := radius / 2
	offsetY := radius / 3
	tilt := rng.nextInt(5) - 2
	drawSlantedRect(img, image.Point{X: center.X - offsetX, Y: center.Y - offsetY}, width, height, tilt, brow)
	drawSlantedRect(img, image.Point{X: center.X + offsetX, Y: center.Y - offsetY}, width, height, -tilt, brow)
}

func drawGlasses(img *image.RGBA, center image.Point, radius int, frame color.RGBA, rng *byteRNG) {
	eyeOffsetX := radius / 2
	eyeOffsetY := radius / 5
	lensWidth := radius / 2
	lensHeight := radius / 3
	bridge := radius / 8
	thickness := 2 + rng.nextInt(2)

	left := image.Point{X: center.X - eyeOffsetX, Y: center.Y - eyeOffsetY}
	right := image.Point{X: center.X + eyeOffsetX, Y: center.Y - eyeOffsetY}

	drawRectOutline(img, left, lensWidth, lensHeight, thickness, frame)
	drawRectOutline(img, right, lensWidth, lensHeight, thickness, frame)
	for x := left.X + lensWidth/2; x < left.X+lensWidth/2+bridge; x++ {
		for t := -thickness; t <= thickness; t++ {
			img.Set(x, left.Y+t, frame)
		}
	}
}

func drawMask(img *image.RGBA, center image.Point, radius int, mask color.RGBA, rng *byteRNG) {
	if rng.nextInt(4) != 0 {
		return
	}
	width := int(float64(radius) * 1.4)
	height := radius / 2
	startY := center.Y + radius/4
	for y := startY; y < startY+height; y++ {
		for x := center.X - width/2; x <= center.X+width/2; x++ {
			img.Set(x, y, mask)
		}
	}
	stripe := blendColor(mask, 0.15)
	for x := center.X - width/2; x <= center.X+width/2; x++ {
		img.Set(x, startY+height/2, stripe)
	}
}

func drawMustache(img *image.RGBA, center image.Point, radius int, hair color.RGBA, rng *byteRNG) {
	if rng.nextInt(3) != 0 {
		return
	}
	width := radius / 2
	height := radius / 8
	startY := center.Y + radius/6
	for y := 0; y < height; y++ {
		for x := -width; x <= width; x++ {
			if x < 0 {
				img.Set(center.X+x, startY+y, hair)
			}
			if x > 0 {
				img.Set(center.X+x, startY+y, hair)
			}
		}
	}
}

func drawChinShadow(img *image.RGBA, center image.Point, radius int, shadow color.RGBA, rng *byteRNG) {
	if rng.nextInt(2) != 0 {
		return
	}
	width := radius / 2
	height := radius / 4
	startY := center.Y + radius/2
	for y := 0; y < height; y++ {
		for x := -width; x <= width; x++ {
			if x*x+y*y <= width*width {
				img.Set(center.X+x, startY+y, shadow)
			}
		}
	}
}

func drawForeheadMark(img *image.RGBA, center image.Point, radius int, mark color.RGBA, rng *byteRNG) {
	if rng.nextInt(4) != 0 {
		return
	}
	size := radius / 6
	startY := center.Y - radius/2
	drawDiamond(img, image.Point{X: center.X, Y: startY}, size, mark)
}

func drawHood(img *image.RGBA, center image.Point, radius int, hood color.RGBA, rng *byteRNG) {
	if rng.nextInt(3) != 0 {
		return
	}
	width := radius * 2
	height := radius + radius/2
	startY := center.Y - radius
	for y := startY; y < startY+height; y++ {
		for x := center.X - width/2; x <= center.X+width/2; x++ {
			dx := float64(x - center.X)
			dy := float64(y - (center.Y - radius/3))
			if (dx*dx)/(float64(width*width)/4)+(dy*dy)/(float64(height*height)/4) <= 1 {
				if img.RGBAAt(x, y).A != 0 {
					img.Set(x, y, blendColor(hood, 0.05))
				}
			}
		}
	}
}

func drawHat(img *image.RGBA, center image.Point, radius int, hat color.RGBA, rng *byteRNG) {
	height := radius/2 + rng.nextInt(radius/4)
	top := center.Y - radius - height/3
	brimHeight := radius / 10
	brimWidth := radius + radius/2
	for y := top; y < top+height; y++ {
		for x := center.X - radius; x <= center.X+radius; x++ {
			dx := x - center.X
			dy := y - (center.Y - radius)
			if dx*dx+dy*dy <= radius*radius {
				img.Set(x, y, hat)
			}
		}
	}
	for y := center.Y - radius; y < center.Y-radius+brimHeight; y++ {
		for x := center.X - brimWidth/2; x <= center.X+brimWidth/2; x++ {
			img.Set(x, y, hat)
		}
	}
}

func drawEarrings(img *image.RGBA, center image.Point, radius int, jewel color.RGBA) {
	offsetX := radius * 5 / 6
	offsetY := radius / 10
	size := radius / 8
	drawFilledCircle(img, image.Point{X: center.X - offsetX, Y: center.Y + offsetY}, size, jewel)
	drawFilledCircle(img, image.Point{X: center.X + offsetX, Y: center.Y + offsetY}, size, jewel)
}

func drawFreckles(img *image.RGBA, center image.Point, radius int, freckle color.RGBA, rng *byteRNG) {
	count := 6 + rng.nextInt(8)
	for i := 0; i < count; i++ {
		x := center.X - radius/2 + rng.nextInt(radius)
		y := center.Y + rng.nextInt(radius/3)
		img.Set(x, y, freckle)
	}
}

func drawScar(img *image.RGBA, center image.Point, radius int, scar color.RGBA, rng *byteRNG) {
	if rng.nextInt(4) != 0 {
		return
	}
	length := radius / 2
	startX := center.X - length/2
	startY := center.Y - radius/6
	angle := float64(rng.nextInt(5)-2) * 0.2
	for i := 0; i < length; i++ {
		x := startX + i
		y := startY + int(float64(i)*angle)
		img.Set(x, y, scar)
	}
}

func drawBeard(img *image.RGBA, center image.Point, radius int, beard color.RGBA, rng *byteRNG) {
	height := radius/2 + rng.nextInt(radius/4)
	startY := center.Y + radius/4
	for y := startY; y < startY+height; y++ {
		for x := center.X - radius/2; x <= center.X+radius/2; x++ {
			dx := x - center.X
			dy := y - (center.Y + radius/3)
			if dx*dx+dy*dy <= radius*radius/2 {
				img.Set(x, y, beard)
			}
		}
	}
}

func drawNose(img *image.RGBA, center image.Point, radius int) {
	noseColor := color.RGBA{R: 180, G: 120, B: 90, A: 255}
	height := int(float64(radius) * 0.25)
	for y := 0; y < height; y++ {
		width := int(float64(height-y) * 0.3)
		for x := -width; x <= width; x++ {
			img.Set(center.X+x, center.Y+y/2, noseColor)
		}
	}
}

func drawBlush(img *image.RGBA, center image.Point, radius int, blush color.RGBA, rng *byteRNG) {
	if rng.nextInt(3) == 0 {
		return
	}
	offsetX := radius / 2
	offsetY := radius / 6
	size := radius / 6
	drawFilledCircle(img, image.Point{X: center.X - offsetX, Y: center.Y + offsetY}, size, blush)
	drawFilledCircle(img, image.Point{X: center.X + offsetX, Y: center.Y + offsetY}, size, blush)
}

func drawMouth(img *image.RGBA, center image.Point, radius int, mouth color.RGBA, rng *byteRNG) {
	width := int(float64(radius) * 0.7)
	curve := float64(rng.nextInt(6)-2) / 10.0
	baseY := float64(center.Y) + float64(radius)/3.0
	thickness := int(float64(radius) * 0.08)

	for x := -width / 2; x <= width/2; x++ {
		xf := float64(x) / float64(width/2)
		y := baseY + curve*math.Pow(xf, 2)*float64(radius)*1.2
		for t := -thickness; t <= thickness; t++ {
			img.Set(center.X+x, int(y)+t, mouth)
		}
	}
}

func drawLipShine(img *image.RGBA, center image.Point, radius int, lip color.RGBA, rng *byteRNG) {
	if rng.nextInt(2) == 0 {
		return
	}
	width := radius / 3
	height := radius / 20
	startY := center.Y + radius/3
	for y := 0; y < height; y++ {
		for x := -width / 2; x <= width/2; x++ {
			img.Set(center.X+x, startY+y, lip)
		}
	}
}

func drawRectOutline(img *image.RGBA, center image.Point, width int, height int, thickness int, stroke color.RGBA) {
	left := center.X - width/2
	right := center.X + width/2
	top := center.Y - height/2
	bottom := center.Y + height/2
	for t := 0; t < thickness; t++ {
		for x := left; x <= right; x++ {
			img.Set(x, top+t, stroke)
			img.Set(x, bottom-t, stroke)
		}
		for y := top; y <= bottom; y++ {
			img.Set(left+t, y, stroke)
			img.Set(right-t, y, stroke)
		}
	}
}

func drawDiamond(img *image.RGBA, center image.Point, radius int, fill color.RGBA) {
	for y := -radius; y <= radius; y++ {
		for x := -radius; x <= radius; x++ {
			if abs(x)+abs(y) <= radius {
				img.Set(center.X+x, center.Y+y, fill)
			}
		}
	}
}

func drawSlantedRect(img *image.RGBA, center image.Point, width int, height int, slope int, fill color.RGBA) {
	left := center.X - width/2
	top := center.Y - height/2
	for y := 0; y < height; y++ {
		shift := (y * slope) / height
		for x := 0; x < width; x++ {
			img.Set(left+x+shift, top+y, fill)
		}
	}
}

func drawChevron(img *image.RGBA, center image.Point, width int, height int, fill color.RGBA) {
	for y := 0; y < height; y++ {
		offset := int(float64(y) * 0.8)
		for x := -width/2 + offset; x <= width/2-offset; x++ {
			img.Set(center.X+x, center.Y+y, fill)
		}
	}
}

func drawStripe(img *image.RGBA, center image.Point, width int, height int, fill color.RGBA) {
	for y := 0; y < height; y++ {
		if y%2 == 0 {
			for x := -width / 2; x <= width/2; x++ {
				img.Set(center.X+x, center.Y+y, fill)
			}
		}
	}
}

func drawOrbitRings(img *image.RGBA, center image.Point, radius int, accent color.RGBA) {
	ringRadius := radius + radius/2
	for angle := 0.0; angle < 2*math.Pi; angle += math.Pi / 64 {
		x := center.X + int(float64(ringRadius)*math.Cos(angle))
		y := center.Y + int(float64(ringRadius)*math.Sin(angle)*0.5)
		img.Set(x, y, accent)
	}
}

func drawStars(img *image.RGBA, rng *byteRNG, radius int, accent color.RGBA) {
	count := 12 + rng.nextInt(10)
	for i := 0; i < count; i++ {
		x := rng.nextInt(radius*2) + radius/2
		y := rng.nextInt(radius*2) + radius/2
		img.Set(x, y, accent)
	}
}

func drawCircuitTrace(img *image.RGBA, center image.Point, radius int, accent color.RGBA, rng *byteRNG) {
	count := 4 + rng.nextInt(4)
	for i := 0; i < count; i++ {
		start := image.Point{
			X: center.X - radius + rng.nextInt(radius*2),
			Y: center.Y - radius + rng.nextInt(radius*2),
		}
		length := radius/2 + rng.nextInt(radius/2)
		current := start
		for j := 0; j < length; j++ {
			img.Set(current.X, current.Y, accent)
			switch rng.nextInt(4) {
			case 0:
				current.X++
			case 1:
				current.X--
			case 2:
				current.Y++
			default:
				current.Y--
			}
			if current.X < 0 || current.Y < 0 || current.X >= img.Bounds().Dx() || current.Y >= img.Bounds().Dy() {
				break
			}
		}
	}
}

func drawConstellation(img *image.RGBA, center image.Point, radius int, accent color.RGBA, rng *byteRNG) {
	nodes := 6 + rng.nextInt(4)
	points := make([]image.Point, 0, nodes)
	for i := 0; i < nodes; i++ {
		points = append(points, image.Point{
			X: center.X - radius + rng.nextInt(radius*2),
			Y: center.Y - radius + rng.nextInt(radius*2),
		})
	}
	for i := 0; i < len(points); i++ {
		drawLine(img, points[i], points[(i+1)%len(points)], accent)
		drawFilledCircle(img, points[i], 1+rng.nextInt(2), accent)
	}
}

func drawAurora(img *image.RGBA, center image.Point, radius int, accent color.RGBA, rng *byteRNG) {
	bands := 3 + rng.nextInt(3)
	for i := 0; i < bands; i++ {
		offset := rng.nextInt(radius) - radius/2
		for x := center.X - radius; x <= center.X+radius; x++ {
			y := center.Y - radius/2 + int(math.Sin(float64(x+offset)/float64(radius))*float64(radius)/4)
			if y >= 0 && y < img.Bounds().Dy() {
				img.Set(x, y, blendColor(accent, 0.3))
				img.Set(x, y+1, accent)
			}
		}
	}
	drawGridOverlay(img, center, radius, blendColor(accent, 0.4), rng)
}

func drawGridOverlay(img *image.RGBA, center image.Point, radius int, accent color.RGBA, rng *byteRNG) {
	step := 4 + rng.nextInt(4)
	for y := center.Y - radius; y <= center.Y+radius; y += step {
		for x := center.X - radius; x <= center.X+radius; x++ {
			if x >= 0 && y >= 0 && x < img.Bounds().Dx() && y < img.Bounds().Dy() {
				img.Set(x, y, accent)
			}
		}
	}
}

func drawCornerTicks(img *image.RGBA, stroke color.RGBA, length int) {
	bounds := img.Bounds()
	for i := 0; i < length; i++ {
		img.Set(bounds.Min.X+i, bounds.Min.Y, stroke)
		img.Set(bounds.Min.X, bounds.Min.Y+i, stroke)
		img.Set(bounds.Max.X-1-i, bounds.Min.Y, stroke)
		img.Set(bounds.Max.X-1, bounds.Min.Y+i, stroke)
		img.Set(bounds.Min.X+i, bounds.Max.Y-1, stroke)
		img.Set(bounds.Min.X, bounds.Max.Y-1-i, stroke)
		img.Set(bounds.Max.X-1-i, bounds.Max.Y-1, stroke)
		img.Set(bounds.Max.X-1, bounds.Max.Y-1-i, stroke)
	}
}

func drawHexGrid(img *image.RGBA, center image.Point, radius int, accent color.RGBA, rng *byteRNG) {
	step := radius / 3
	for y := center.Y - radius; y <= center.Y+radius; y += step {
		rowShift := 0
		if ((y - center.Y) / step % 2) != 0 {
			rowShift = step / 2
		}
		for x := center.X - radius; x <= center.X+radius; x += step {
			drawHexagon(img, image.Point{X: x + rowShift, Y: y}, step/3, accent, rng)
		}
	}
}

func drawHexagon(img *image.RGBA, center image.Point, radius int, accent color.RGBA, rng *byteRNG) {
	if rng.nextInt(4) != 0 {
		return
	}
	points := make([]image.Point, 0, 6)
	for i := 0; i < 6; i++ {
		angle := float64(i) * math.Pi / 3
		points = append(points, image.Point{
			X: center.X + int(float64(radius)*math.Cos(angle)),
			Y: center.Y + int(float64(radius)*math.Sin(angle)),
		})
	}
	for i := 0; i < len(points); i++ {
		drawLine(img, points[i], points[(i+1)%len(points)], accent)
	}
}

func drawLine(img *image.RGBA, a image.Point, b image.Point, stroke color.RGBA) {
	dx := int(math.Abs(float64(b.X - a.X)))
	dy := -int(math.Abs(float64(b.Y - a.Y)))
	sx := -1
	if a.X < b.X {
		sx = 1
	}
	sy := -1
	if a.Y < b.Y {
		sy = 1
	}
	err := dx + dy
	x := a.X
	y := a.Y
	for {
		img.Set(x, y, stroke)
		if x == b.X && y == b.Y {
			break
		}
		e2 := 2 * err
		if e2 >= dy {
			err += dy
			x += sx
		}
		if e2 <= dx {
			err += dx
			y += sy
		}
	}
}

func applyVignette(img *image.RGBA, center image.Point, radius int) {
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			dx := float64(x - center.X)
			dy := float64(y - center.Y)
			dist := math.Sqrt(dx*dx + dy*dy)
			if dist > float64(radius) {
				pixel := img.RGBAAt(x, y)
				factor := math.Min((dist-float64(radius))/float64(radius), 0.6)
				img.SetRGBA(x, y, color.RGBA{
					R: uint8(float64(pixel.R) * (1 - factor)),
					G: uint8(float64(pixel.G) * (1 - factor)),
					B: uint8(float64(pixel.B) * (1 - factor)),
					A: pixel.A,
				})
			}
		}
	}
}

func applyNoise(img *image.RGBA, rng *byteRNG, intensity int) {
	if intensity <= 0 {
		return
	}
	for i := 0; i < intensity*intensity; i++ {
		x := rng.nextInt(img.Bounds().Dx())
		y := rng.nextInt(img.Bounds().Dy())
		p := img.RGBAAt(x, y)
		shift := int(rng.nextInt(5)) - 2
		img.SetRGBA(x, y, color.RGBA{
			R: clampChannel(int(p.R) + shift),
			G: clampChannel(int(p.G) + shift),
			B: clampChannel(int(p.B) + shift),
			A: p.A,
		})
	}
}

func clampChannel(value int) uint8 {
	if value < 0 {
		return 0
	}
	if value > 255 {
		return 255
	}
	return uint8(value)
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}
