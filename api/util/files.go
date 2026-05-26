package util

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/image/draw"
)

func getFullPath(filename string) (string) {
    return "/files/" + filename
}

func SaveFile(filename string, fileBytes string) (uint, error) {
    fullPath := getFullPath(filename)

    err := os.MkdirAll(filepath.Dir(fullPath), 0666)
    if err != nil {
        fmt.Print(err.Error())
        return 0, err
    }

    f, err := os.Create(fullPath)
    if err != nil {
        fmt.Print(err.Error())
        return 0, err
    }
    defer f.Close()

    data, err := base64.StdEncoding.DecodeString(fileBytes)
    if err != nil {
        fmt.Print(err.Error())
        return 0, err
    }

    writtenBytes, err := f.Write(data)
    if err != nil {
        fmt.Print(err.Error())
        return 0, err
    }

    return uint(writtenBytes), nil
}

// Returns (byteSizeImage, byteSizeThumb, ImgW, ImgH, Error)
func SaveImage(filename string, fileBytes string) (uint, uint, int, int, error) {
    img, err := imageFromByteString(fileBytes)
    if err != nil {
        return 0, 0, 0, 0, err
    }

    thumbnailBytes, err := createThumbnailBytes(img, 6, draw.CatmullRom)
    if err != nil {
        fmt.Print(err.Error())
        return 0, 0, 0, 0, err
    }

    fnameBase := strings.TrimSuffix(filename, filepath.Ext(filename))
    filenameThumb := fmt.Sprintf("%s-thumb.jpg", fnameBase)

    bytes, err := SaveFile(filename, fileBytes)
    if err != nil {
        fmt.Print(err.Error())
        return 0, 0, 0, 0, err
    }

    bytesThumb, err := SaveFile(filenameThumb, thumbnailBytes)
    if err != nil {
        fmt.Print(err.Error())
        return 0, 0, 0, 0, err
    }

    return bytes, bytesThumb, img.Bounds().Size().X, img.Bounds().Size().Y, nil
}

func GetPostImageFilename(boardCode string, sourceFilename string) string {
    ext := filepath.Ext(sourceFilename)
    return fmt.Sprintf("%s/%d%s", boardCode, time.Now().UnixMilli(), ext)
}

func GetFileHashB64(fileBytes string) (string) {
    hashed := md5.Sum([]byte(fileBytes))
    return base64.StdEncoding.EncodeToString([]byte(hashed[:]))
}

func createThumbnailBytes(srcImage image.Image, scaleDown int, scale draw.Scaler) (string, error) {
    thumbRect := image.Rect(0, 0, srcImage.Bounds().Size().X / scaleDown, srcImage.Bounds().Size().Y / scaleDown)
    thumbnail := createThumbnail(srcImage, thumbRect, scale)

    var buf bytes.Buffer
    if err := jpeg.Encode(&buf, thumbnail, &jpeg.Options{Quality: 50}); err != nil {
        fmt.Print(err.Error())
        return "", err
    }

    thumbBytes := buf.Bytes()
    thumbBase64 := base64.StdEncoding.EncodeToString(thumbBytes)

    return thumbBase64, nil
}

func createThumbnail(src image.Image, rect image.Rectangle, scale draw.Scaler) image.Image {
    dst := image.NewRGBA(rect)
    scale.Scale(dst, rect, src, src.Bounds(), draw.Over, nil)
    return dst
}

func imageFromByteString(imgByteString string) (image.Image, error) {
    data, err := base64.StdEncoding.DecodeString(imgByteString)
    if err != nil {
        return nil, err
    }

    img, _, err := image.Decode(bytes.NewReader(data))
    if err != nil {
        return nil, err
    }

    return img, nil
}