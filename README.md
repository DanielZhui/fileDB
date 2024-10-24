# FileDB 📁

FileDB is a simple key-value store implementation in Go that persists data to disk. It provides basic CRUD operations and is designed for lightweight storage needs. 🚀

## ✨ Features

- 💾 Persistent storage: Data is stored on disk and can be retrieved across program restarts.
- 🔄 Basic CRUD operations: Set, Get, Update, and Delete operations are supported.
- 🗂️ Key directory: Maintains an in-memory index of keys for fast lookups.
- ➕ Append-only write: New data is appended to the file, improving write performance.

## 🛠️ Installation

To use FileDB in your Go project, you can install it using `go get`:

```bash
go get github.com/DanielZhui/fileDB@v0.0.1
```

## 🚀 Usage

Here's a basic example of how to use FileDB:

```go
package main

import (
	"fmt"
	"log"

	"github.com/DanielZhui/fileDB"
)

func main() {
	filePath := "./test.db"
	ds, err := fileDB.InitDiskStore(filePath)
	if err != nil {
		log.Fatalf("Failed to initialize disk store: %v", err)
	}

	// Set a key-value pair
	err = ds.Set("hello", "world")
	if err != nil {
		log.Printf("Failed to set 'hello': %v", err)
	}

	// Get a value
	value, err := ds.Get("hello")
	if err != nil {
		log.Printf("Failed to get 'hello': %v", err)
	} else {
		fmt.Println(value) // Output: world
	}

	// Update a value
	err = ds.Update("hello", "new world")
	if err != nil {
		log.Printf("Failed to update 'hello': %v", err)
	}

	// Delete a key
	err = ds.Delete("hello")
	if err != nil {
		log.Printf("Failed to delete 'hello': %v", err)
	}

	// List all keys (reloads the key directory from the file)
	ds.List(filePath)
}
```

## 📚 API

- 🆕 `InitDiskStore(fileName string) (*DiskStore, error)`: Initialize a new DiskStore instance.
- ✍️ `Set(key string, value string) error`: Set a key-value pair.
- 🔍 `Get(key string) (string, error)`: Retrieve the value for a given key.
- 🔄 `Update(key string, value string) error`: Update the value for an existing key.
- 🗑️ `Delete(key string) error`: Delete a key-value pair.
- 📋 `List(filePath string)`: Reload the key directory from the file.

## 🏗️ Data Structure

FileDB uses a simple file format to store data:

- 📊 Header (12 bytes): timestamp (4 bytes), key size (4 bytes), value size (4 bytes)
- 🔑 Key: Variable length string
- 📄 Value: Variable length string

## ⚠️ Limitations

- 🔒 FileDB is not designed for concurrent access. It's suitable for single-threaded applications or scenarios where external synchronization is applied.
- 💾 The entire key directory is kept in memory, which may not be suitable for very large datasets.
- 🗑️ Deleted keys are not removed from the file, potentially leading to file size growth over time.

## 🤝 Contributing

Contributions to FileDB are welcome! Please feel free to submit a Pull Request. 😊

## 📜 License

This project is licensed under the MIT License - see the LICENSE file for details.
