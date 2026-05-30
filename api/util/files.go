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
	"os/exec"
	"path/filepath"
	"strconv"
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
        if mimeType == "video/webm" {
            return saveFileAsWebm(filename, fileBytes)
        } else {

        }
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

func saveFileAsWebm(filename string, fileBytes []byte) (*SaveFileResult, error) {
    fnameBase := strings.TrimSuffix(filename, filepath.Ext(filename))
    filenameThumb := fmt.Sprintf("%s-thumb.jpg", fnameBase)

    // 1) Save file

    nBytesImg, err := WriteFile(filename, fileBytes)
    if err != nil {
        fmt.Print(err.Error())
        return nil, err
    }

    // 2) Get video WxH using ffprobe

    var stderr bytes.Buffer
    cmd := exec.Command("/usr/local/bin/ffprobe", "-v", "error", "-show_entries", "stream=width,height", "-of", "default=noprint_wrappers=1:nokey=1", getFullPath(filename))
    cmd.Stderr = &stderr

    stdout, err := cmd.Output()
    if err != nil {
        fmt.Println("Error:", err.Error())
        fmt.Println("Stderr:", stderr.String())
        return nil, err
    }
    widthHeightStrings := strings.Split(string(stdout[:]), "\n")[:2]

    if len(widthHeightStrings) != 2 {
        return nil, errors.New(fmt.Sprintf("Invalid probe output:%v", widthHeightStrings))
    }

    videoWidth, err := strconv.ParseUint(widthHeightStrings[0], 10, 64)
    if err != nil {
        return nil, err
    }
    videoHeight, err := strconv.ParseUint(widthHeightStrings[1], 10, 64)
    if err != nil {
        return nil, err
    }

    // 3) Save video thumbnail using ffmpeg

    cmd = exec.Command("/usr/local/bin/ffmpeg", "-i",  getFullPath(filename), "-ss", "00:00:01.000", "-vframes", "1", getFullPath(filenameThumb))
    cmd.Stderr = &stderr
    _, err = cmd.Output()
    if err != nil {
        fmt.Println("Error:", err.Error())
        fmt.Println("Stderr:", stderr.String())
        return nil, err
    }

    // 4) Get thumbnail WxH

    imgThumb, err := os.Open(getFullPath(filenameThumb))
    if err != nil {
        fmt.Println(err.Error())
        return nil, err
    }

    imgThumbStat, err := imgThumb.Stat()
    if err != nil {
        fmt.Println(err.Error())
        return nil, err
    }

    // Result

    return &SaveFileResult{
        SizeImageBytes: nBytesImg,
        SizeThumbnailBytes: uint(imgThumbStat.Size()),
        ImageWidth: uint(videoWidth),
        ImageHeight: uint(videoHeight),
    }, nil
}


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