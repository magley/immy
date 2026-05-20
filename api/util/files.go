package util

import (
	"bytes"
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

func SaveFile(filename string, fileBytes string) (error) {
    fullPath := getFullPath(filename)

    err := os.MkdirAll(filepath.Dir(fullPath), 0666)
    if err != nil {
        fmt.Print(err.Error())
        return err
    }

    f, err := os.Create(fullPath)
    if err != nil {
        fmt.Print(err.Error())
        return err
    }
    defer f.Close()

    data, err := base64.StdEncoding.DecodeString(fileBytes)
    if err != nil {
        fmt.Print(err.Error())
        return err
    }

    _, err = f.Write(data)
    if err != nil {
        fmt.Print(err.Error())
        return err
    }

    return nil
}

func SaveImage(filename string, fileBytes string) (error) {
    thumbnailBytes, err := createThumbnailBytes(fileBytes, 6, draw.CatmullRom)
    if err != nil {
        fmt.Print(err.Error())
        return err
    }

    fnameBase := strings.TrimSuffix(filename, filepath.Ext(filename))
    filenameThumb := fmt.Sprintf("%s-thumb.jpg", fnameBase)

    err = SaveFile(filename, fileBytes)
    if err != nil {
        fmt.Print(err.Error())
        return err
    }

    err = SaveFile(filenameThumb, thumbnailBytes)
    if err != nil {
        fmt.Print(err.Error())
        return err
    }

    return nil
}

func GetPostImageFilename(boardCode string, sourceFilename string) string {
    ext := filepath.Ext(sourceFilename)
    return fmt.Sprintf("%s/%d%s", boardCode, time.Now().UnixMilli(), ext)
}

func createThumbnailBytes(srcBytes string, scaleDown int, scale draw.Scaler) (string, error) {
    data, err := base64.StdEncoding.DecodeString(srcBytes)
    if err != nil {
        fmt.Print(err.Error())
        return "", err
    }

    img, _, err := image.Decode(bytes.NewReader(data))
    if err != nil {
        fmt.Print(err.Error())
        return "", err
    }

    thumbRect := image.Rect(0, 0, img.Bounds().Size().X / scaleDown, img.Bounds().Size().Y / scaleDown)
    thumbnail := createThumbnail(img, thumbRect, scale)

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