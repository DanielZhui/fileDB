# fileDB 📁🔑

fileDB is a lightweight, file-based key-value storage system implemented in Go. It provides a simple interface for storing and retrieving string data using a persistent file storage mechanism. 💾

> refer: https://github.com/avinassh/go-caskdb

## ✨ Features

- 📂 File-based persistent storage
- 🔑 Simple key-value operations (Get and Set)
- 🚀 Efficient data encoding and decoding
- ⏱️ Automatic timestamp recording for each entry

## 🛠️ Installation

To use fileDB in your Go project, you can clone this repository or import it in your Go module:

```bash
go get github.com/DanielZhui/fileDB
```

## 🚀 Usage

Here's a quick example of how to use fileDB:

```go
package main

import (
    "fmt"
    "github.com/DanielZhui/fileDB"
)

func main() {
    // Initialize the disk store
    ds, err := fileDB.InitDiskStore("./test.db")
    if err != nil {
        panic(err)
    }

    // Set some key-value pairs
    ds.Set("hello", "world")
    ds.Set("foo", "bar")

    // Retrieve a value
    value := ds.Get("foo")
    fmt.Println(value) // Output: bar
}
```

## 📚 API

### InitDiskStore(fileName string) (*DiskStore, error)

Initializes a new DiskStore or loads an existing one from the specified file.

### (d *DiskStore) Set(key string, value string)

Stores a key-value pair in the database.

### (d *DiskStore) Get(key string) string

Retrieves the value associated with the given key. Returns an empty string if the key is not found.

## 🏗️ Data Structure

Each entry in the file is stored in the following format:

- Header (12 bytes):
  - ⏱️ Timestamp (4 bytes)
  - 📏 Key Size (4 bytes)
  - 📏 Value Size (4 bytes)
- 🔑 Key (variable length)
- 📄 Value (variable length)

## ⚠️ Limitations and Future Improvements

- 🔒 Currently not thread-safe
- 🗑️ No delete or update operations
- 💾 All key information is stored in memory, which may not be suitable for large datasets
- 🔐 No data compression or integrity checks

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request. 👨‍💻👩‍💻

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.