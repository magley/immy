package util

import (
	"fmt"
	"os"
    "encoding/base64"
)

func getFullPath(filename string) (string) {
    return "/files/" + filename
}

func SaveFile(filename string, fileBytes string) (error) {
    f, err := os.Create(getFullPath(filename))
    if err != nil {
        return err
    }
    defer f.Close()

    data, err := base64.StdEncoding.DecodeString(fileBytes)
    if err != nil {
        return err
    }

    _, err = f.Write(data)
    if err != nil {
        fmt.Print(err.Error())
        return err
    }

    return nil
}