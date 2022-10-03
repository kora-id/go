// Copyright 2018 Kora ID. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package uploader

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/nfnt/resize"
)

// Image data type for uploading images
type Image struct {
	Name        string
	ContentType string
	Extention   string
	Data        []byte
	Size        int
}

// MakeThumbnail generate thumbnail from the image instances
func (i *Image) MakeThumbnail() (*Image, error) {
	w, _ := strconv.Atoi(os.Getenv("THUMB_SIZE"))
	h, _ := strconv.Atoi(os.Getenv("THUMB_SIZE"))

	if i.ContentType == "image/jpeg" {
		return thumbJPEG(i, w, h)
	}

	return thumbPNG(i, w, h)
}

// FileName return filename from the image
func (i *Image) FileName() string {
	return fmt.Sprintf("%s%s", i.Name, i.Extention)
}

// NewImage making new instance image
func NewImage(name string, fb []byte) *Image {
	ct, fx := readFileExtention(fb)

	return &Image{
		Name:        name,
		ContentType: ct,
		Extention:   fx,
		Data:        fb,
		Size:        len(fb),
	}
}

// ReadFileBase64 reading base64 string into byte
func ReadFileBase64(s string) (fb []byte, e error) {
	fb, e = base64.StdEncoding.DecodeString(strings.Split(s, "base64,")[1])

	if e == nil {
		if _, extention := readFileExtention(fb); extention == "" {
			e = errors.New("invalid file")
		}
	}

	return
}

func readFileExtention(fb []byte) (ct string, fx string) {
	ct = http.DetectContentType(fb)
	fx = getExtention(ct)

	return ct, fx
}

func thumbJPEG(i *Image, w int, h int) (*Image, error) {
	img, _, _ := image.Decode(bytes.NewReader(i.Data))
	thumbnail := resize.Thumbnail(uint(w), uint(h), img, resize.Lanczos3)

	data := new(bytes.Buffer)
	err := jpeg.Encode(data, thumbnail, &jpeg.Options{
		Quality: 100,
	})

	if err != nil {
		return nil, err
	}

	bs := data.Bytes()

	t := &Image{
		Name:        fmt.Sprintf("%s_thumb", i.Name),
		ContentType: i.ContentType,
		Extention:   i.Extention,
		Data:        bs,
		Size:        len(bs),
	}

	return t, nil
}

func thumbPNG(i *Image, w int, h int) (*Image, error) {
	img, _, _ := image.Decode(bytes.NewReader(i.Data))
	thumbnail := resize.Thumbnail(uint(w), uint(h), img, resize.Lanczos3)

	data := new(bytes.Buffer)
	err := png.Encode(data, thumbnail)

	if err != nil {
		return nil, err
	}

	bs := data.Bytes()

	t := &Image{
		Name:        fmt.Sprintf("%s_thumb", i.Name),
		ContentType: i.ContentType,
		Extention:   i.Extention,
		Data:        bs,
		Size:        len(bs),
	}

	return t, nil
}
