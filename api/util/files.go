package util

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/image/draw"
)

type SaveFileResult struct {
    SizeImageBytes uint
    SizeThumbnailBytes uint
    ImageWidth uint
    ImageHeight uint
}

// TODO: Sandboxing?
func getFullPath(filename string) (string) {
    return "/files/" + filename
}

func GetPostImageFilename(boardCode string, sourceFilename string) string {
    ext := filepath.Ext(sourceFilename)
    return fmt.Sprintf("%s/%d%s", boardCode, time.Now().UnixMilli(), ext)
}

func GetFileHashB64(fileBytes string) (string) {
    hashed := md5.Sum([]byte(fileBytes))
    return base64.StdEncoding.EncodeToString([]byte(hashed[:]))
}

// Returns file size in bytes on success.
// `filename` is relative to the files directory used by the server.
func WriteFile(filename string, fileBytes []byte) (uint, error) {
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

    writtenBytes, err := f.Write(fileBytes)
    if err != nil {
        fmt.Print(err.Error())
        return 0, err
    }

    return uint(writtenBytes), nil
}

func SaveFile(filename string, fileBytesB64 string) (*SaveFileResult, error) {
    fileBytes, err := base64.StdEncoding.DecodeString(fileBytesB64)
    if err != nil {
        return nil, err
    }

    mimeType := http.DetectContentType(fileBytes[:512])

    if strings.HasPrefix(mimeType, "image/") {
        return saveFileAsImage(filename, fileBytes)
    } else {
    }

    return nil, errors.New("Unexpected MIME type: " + mimeType)
}

func saveFileAsImage(filename string, fileBytes []byte) (*SaveFileResult, error) {
    img, _, err := image.Decode(bytes.NewReader(fileBytes))
    if err != nil {
        return nil, err
    }

    fileBytesThumbnail, err := createThumbnail(img, 6, draw.CatmullRom)
    if err != nil {
        return nil, err
    }

    fnameBase := strings.TrimSuffix(filename, filepath.Ext(filename))
    filenameThumb := fmt.Sprintf("%s-thumb.jpg", fnameBase)

    nBytesImg, err := WriteFile(filename, fileBytes)
    if err != nil {
        fmt.Print(err.Error())
        return nil, err
    }

    nBytesThumbnail, err := WriteFile(filenameThumb, fileBytesThumbnail)
    if err != nil {
        fmt.Print(err.Error())
        return nil, err
    }

    return &SaveFileResult{
        SizeImageBytes: nBytesImg,
        SizeThumbnailBytes: nBytesThumbnail,
        ImageWidth: uint(img.Bounds().Size().X),
        ImageHeight: uint(img.Bounds().Size().Y),
    }, nil
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
func createThumbnail(srcImage image.Image, scaleDown int, scale draw.Scaler) ([]byte, error) {
    thumbRect := image.Rect(0, 0, srcImage.Bounds().Size().X / scaleDown, srcImage.Bounds().Size().Y / scaleDown)
    thumbnail := scaleImage(srcImage, thumbRect, scale)

    var buf bytes.Buffer
    if err := jpeg.Encode(&buf, thumbnail, &jpeg.Options{Quality: 50}); err != nil {
        fmt.Print(err.Error())
        return nil, err
    }

    thumbBytes := buf.Bytes()
    return thumbBytes, nil
}

func scaleImage(src image.Image, rect image.Rectangle, scale draw.Scaler) image.Image {
    dst := image.NewRGBA(rect)
    scale.Scale(dst, rect, src, src.Bounds(), draw.Over, nil)
    return dst
}