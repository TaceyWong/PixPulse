package converter

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"

	"PixPulse/internal/utils"

	"golang.org/x/image/bmp"
)

// Converter handles the conversion processes
type Converter struct{}

// NewConverter creates a new Converter
func NewConverter() *Converter {
	return &Converter{}
}

// Convert processes the request based on the mode
func (c *Converter) Convert(req ConvertRequest) (*ConvertResponse, error) {
	if _, err := os.Stat(req.InputPath); os.IsNotExist(err) {
		return &ConvertResponse{Success: false, Error: "Input file does not exist"}, nil
	}

	var err error
	var output string

	switch req.Mode {
	case "color":
		output, err = c.convertColor(req)
	case "bw":
		output, err = c.convertBW(req)
	case "render":
		output, err = c.renderSVG(req)
	default:
		return &ConvertResponse{Success: false, Error: "Invalid mode"}, nil
	}

	if err != nil {
		return &ConvertResponse{Success: false, Error: fmt.Sprintf("Conversion error: %v\n%s", err, output)}, nil
	}

	return &ConvertResponse{Success: true}, nil
}

// convertColor uses vtracer: vtracer input.png -o output.svg
func (c *Converter) convertColor(req ConvertRequest) (string, error) {
	bin := utils.GetBinPath("vtracer")
	return RunCommand(bin, "--input", req.InputPath, "--output", req.OutputPath)
}

// convertBW uses potrace: png -> bmp -> potrace -> svg
func (c *Converter) convertBW(req ConvertRequest) (string, error) {
	// 1. Convert to BMP
	bmpPath := filepath.Join(os.TempDir(), "pixpulse_temp.bmp")
	err := c.toBMP(req.InputPath, bmpPath)
	if err != nil {
		return "", fmt.Errorf("failed to convert to BMP: %w", err)
	}
	defer os.Remove(bmpPath)

	// 2. Run potrace
	bin := utils.GetBinPath("potrace")
	// -s means SVG output
	return RunCommand(bin, bmpPath, "-s", "-o", req.OutputPath)
}

// renderSVG uses resvg: resvg input.svg output.png
func (c *Converter) renderSVG(req ConvertRequest) (string, error) {
	bin := utils.GetBinPath("resvg")
	return RunCommand(bin, req.InputPath, req.OutputPath)
}

// toBMP converts any image format to a monochrome-friendly BMP
func (c *Converter) toBMP(input, output string) error {
	f, err := os.Open(input)
	if err != nil {
		return err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return err
	}

	// Create a new RGBA image with a white background
	// This ensures transparent areas don't become black
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)
	draw.Draw(newImg, bounds, &image.Uniform{color.White}, image.Point{}, draw.Src)
	draw.Draw(newImg, bounds, img, bounds.Min, draw.Over)

	out, err := os.Create(output)
	if err != nil {
		return err
	}
	defer out.Close()

	return bmp.Encode(out, newImg)
}
