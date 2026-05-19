package util

import (
	"fmt"
	"os"
)

func getFullPath(filename string) (string) {
    return "/files/" + filename
}

func SaveFile(filename string, fileBytes string) (error) {
    f, err := os.Create(getFullPath(filename))
    if err != nil {
        fmt.Print(err.Error())
        return err
    }
    defer f.Close()

    
    data := []byte(fileBytes)

    _, err = f.Write(data)
    if err != nil {
        fmt.Print(err.Error())
        return err
    }

    return nil
}