package main

import (
    "flag"
    "fmt"
    "os"
    "path/filepath"
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
)

func main() {
    flag.StringVar(&targetDir, "targetDir", ".", "The target directory to map out.")
    flag.Parse()

    walkPath()
    printSummary()
}

func walkPath() {
    filepath.Walk(targetDir, onWalkDir)
}

func getSizeStr() string {
    if totalSize / TB >= 1 {
        return fmt.Sprintf("%.2f TB (%d Bytes)", float64(totalSize) / TB, totalSize)
    }

    if totalSize / GB >= 1 {
        return fmt.Sprintf("%.2f GB  (%d Bytes)", float64(totalSize) / GB, totalSize)
    }

    if totalSize / MB >= 1 {
        return fmt.Sprintf("%.2f MB  (%d Bytes)", float64(totalSize) / MB, totalSize)
    }

    if totalSize / KB >= 1 {
        return fmt.Sprintf("%.2f KB  (%d Bytes)", float64(totalSize) / KB, totalSize)
    }

    return fmt.Sprintf("%d Bytes", totalSize)
}

func printSummary() {
    fmt.Printf(
        "\n*** Walk complete (files: %d, directories: %d, size: %s\n",
        fileCount,
        directoryCount,
        getSizeStr(),
    )
}

func onWalkDir(path string, info os.FileInfo, err error) error {
    if err != nil {
        return err
    }

    if info.IsDir() {
        fmt.Printf("Walking directory %s...\n", path)
        directoryCount++
        return nil
    }

    totalSize += uint64(info.Size())
    fileCount++
    return nil
}
