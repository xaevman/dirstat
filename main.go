package main

import (
    "flag"
    "fmt"
    "os"
    "path/filepath"
    "strings"
)

const (
    TB = 1000000000000
    GB = 1000000000
    MB = 1000000
    KB = 1000
)

// command line options
var (
    targetDir string
)

var (
    totalSize      uint64
    fileCount      uint64
    directoryCount uint64
    currentDir     string
    sizeMap        map[string]uint64 = make(map[string]uint64)
    errList        []error           = make([]error, 0)
)

func main() {
    flag.StringVar(&targetDir, "targetDir", ".", "The target directory to map out.")
    flag.Parse()

    walkPath()
    printSummary()
}

func walkPath() {
    err := filepath.Walk(targetDir, onWalkDir)
    if err != nil {
        panic(err)
    }
}

func getSizeStr(size uint64) string {
    if size / TB >= 1 {
        return fmt.Sprintf("%.2f TB (%d Bytes)", float64(size) / TB, size)
    }

    if size / GB >= 1 {
        return fmt.Sprintf("%.2f GB  (%d Bytes)", float64(size) / GB, size)
    }

    if size / MB >= 1 {
        return fmt.Sprintf("%.2f MB  (%d Bytes)", float64(size) / MB, size)
    }

    if size / KB >= 1 {
        return fmt.Sprintf("%.2f KB  (%d Bytes)", float64(size) / KB, size)
    }

    return fmt.Sprintf("%d Bytes", size)
}

func printSummary() {
    fmt.Printf(
        "\nWalk complete. (files: %d, directories: %d, errors: %d)\n",
        fileCount,
        directoryCount,
        len(errList),
    )

    for i := range errList {
        fmt.Println(errList[i])
    }

    fmt.Printf(
        "\n*** Total Size: %s\n",
        getSizeStr(totalSize),
    )
}

func onWalkDir(path string, info os.FileInfo, err error) error {
    if err != nil {
        errList = append(errList, err)
        return nil
    }

    sPath := strings.Replace(path, targetDir, "", -1)
    pathParts := strings.Split(sPath, string(os.PathSeparator))
    subDir := sPath
    if len(pathParts) > 1 {
        subDir = pathParts[1]
    }
    
    _, ok := sizeMap[subDir]
    if !ok {
        fmt.Printf(" (size: %s)\n", getSizeStr(sizeMap[currentDir]))

        currentDir = subDir
        sizeMap[currentDir] = 0
        fmt.Printf("Entering subdirectory %s...", currentDir)
    }
    
    if info.IsDir() {
        directoryCount++
        return nil
    }

    fileSize := uint64(info.Size())
    sizeMap[subDir] += fileSize
    totalSize += fileSize
    fileCount++

    return nil
}
