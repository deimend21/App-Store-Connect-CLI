package asc

import (
	"bytes"
	"encoding/binary"
	"hash/crc32"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"strings"
	"testing"
)

type declaredAlphaImage struct{ image.Image }

func (declaredAlphaImage) Opaque() bool { return false }

func TestAppStoreImageFormatAlphaMetadata(t *testing.T) {
	opaque := image.NewNRGBA(image.Rect(0, 0, 2, 2))
	for i := 0; i < len(opaque.Pix); i += 4 {
		opaque.Pix[i+3] = 255
	}
	palette := image.NewPaletted(image.Rect(0, 0, 2, 2), color.Palette{color.NRGBA{A: 0}, color.White})
	for _, tc := range []struct {
		name         string
		img          image.Image
		transparency bool
		wantAlpha    bool
		jpeg         bool
	}{
		{name: "RGB PNG", img: opaque},
		{name: "JPEG", img: opaque, jpeg: true},
		{name: "opaque alpha channel", img: declaredAlphaImage{opaque}, wantAlpha: true},
		{name: "transparent alpha channel", img: image.NewNRGBA(image.Rect(0, 0, 2, 2)), wantAlpha: true},
		{name: "indexed transparency", img: palette, wantAlpha: true},
		{name: "truecolor transparency", img: opaque, transparency: true, wantAlpha: true},
		{name: "grayscale transparency", img: image.NewGray(image.Rect(0, 0, 2, 2)), transparency: true, wantAlpha: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var encoded bytes.Buffer
			if tc.jpeg {
				if err := jpeg.Encode(&encoded, tc.img, nil); err != nil {
					t.Fatal(err)
				}
			} else if err := png.Encode(&encoded, tc.img); err != nil {
				t.Fatal(err)
			}
			content := encoded.Bytes()
			if tc.transparency {
				length := 6
				if _, ok := tc.img.(*image.Gray); ok {
					length = 2
				}
				chunk := make([]byte, length+12)
				binary.BigEndian.PutUint32(chunk, uint32(length))
				copy(chunk[4:8], "tRNS")
				binary.BigEndian.PutUint32(chunk[8+length:], crc32.ChecksumIEEE(chunk[4:8+length]))
				content = append(append(append([]byte{}, content[:33]...), chunk...), content[33:]...)
			}
			format, err := ReadAppStoreImageFormatFrom(bytes.NewReader(content))
			if tc.wantAlpha {
				if err == nil || !strings.Contains(err.Error(), "alpha channel or transparency") {
					t.Fatalf("format %q error %v", format, err)
				}
			} else if err != nil {
				t.Fatal(err)
			}
		})
	}
}
