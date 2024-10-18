package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"time"
)

type KeyInfo struct {
	timestamp uint32
	position  uint32
	totalSize uint32
}

type DiskStore struct {
	file     *os.File
	position int
	keyInfo  map[string]KeyInfo
}

const headerSize = 12
const defaultWhence = 0

func NewKeyEntry(timestamp uint32, position uint32, totalSize uint32) KeyInfo {
	return KeyInfo{timestamp, position, totalSize}
}

func encodeHeader(timestamp uint32, keySize uint32, valueSize uint32) []byte {
	header := make([]byte, headerSize)
	binary.LittleEndian.PutUint32(header[0:4], timestamp)
	binary.LittleEndian.PutUint32(header[4:8], keySize)
	binary.LittleEndian.PutUint32(header[8:12], valueSize)
	return header
}

func decodeHeader(header []byte) (uint32, uint32, uint32) {
	timestamp := binary.LittleEndian.Uint32(header[0:4])
	keySize := binary.LittleEndian.Uint32(header[4:8])
	valueSize := binary.LittleEndian.Uint32(header[8:12])
	return timestamp, keySize, valueSize
}

func encodeKV(timestamp uint32, key string, value string) (int, []byte) {
	header := encodeHeader(timestamp, uint32(len(key)), uint32(len(value)))
	data := append([]byte(key), []byte(value)...)
	return headerSize + len(data), append(header, data...)
}

func decodeKV(data []byte) (uint32, string, string) {
	timestamp, keySize, valueSize := decodeHeader(data[0:headerSize])
	key := string(data[headerSize : headerSize+keySize])
	value := string(data[headerSize+keySize : headerSize+keySize+valueSize])
	return timestamp, key, value
}

func isFileExist(fileName string) bool {
	if _, err := os.Stat(fileName); err == nil || errors.Is(err, fs.ErrExist) {
		return true
	}
	return false
}

func (d *DiskStore) initKeyDir(exitingFile string) {
	file, _ := os.Open(exitingFile)
	defer file.Close()
	for {
		header := make([]byte, headerSize)
		_, err := io.ReadFull(file, header)
		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}
		timestamp, keySize, valueSize := decodeHeader(header)
		key := make([]byte, keySize)
		value := make([]byte, valueSize)
		_, err = io.ReadFull(file, key)
		if err != nil {
			break
		}
		_, err = io.ReadFull(file, value)
		if err != nil {
			break
		}
		totalSize := headerSize + keySize + valueSize
		d.keyInfo[string(key)] = NewKeyEntry(timestamp, uint32(d.position), uint32(totalSize))
		d.position += int(totalSize)
		fmt.Printf("loaded key=%s, value=%s\n", key, value)
	}
}

func InitDiskStore(fileName string) (*DiskStore, error) {
	ds := &DiskStore{keyInfo: make(map[string]KeyInfo)}
	if isFileExist(fileName) {
		ds.initKeyDir(fileName)
	}

	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	ds.file = file
	return ds, nil
}

func (d *DiskStore) Get(key string) string {
	kInfo, ok := d.keyInfo[key]
	if !ok {
		return ""
	}
	d.file.Seek(int64(kInfo.position), defaultWhence)
	data := make([]byte, kInfo.totalSize)
	_, err := io.ReadFull(d.file, data)
	if err != nil {
		panic("read error")
	}
	_, _, value := decodeKV(data)
	return value
}

func (d *DiskStore) write(data []byte) {
	if _, err := d.file.Write(data); err != nil {
		panic(err)
	}
	if err := d.file.Sync(); err != nil {
		panic(err)
	}
}

func (d *DiskStore) Set(key string, value string) {
	timestamp := uint32(time.Now().Unix())
	size, data := encodeKV(timestamp, key, value)
	d.write(data)
	d.keyInfo[key] = NewKeyEntry(timestamp, uint32(d.position), uint32(size))
	d.position += size
}

func main() {
	fmt.Println("Hello, file DB!")
	ds, _ := InitDiskStore("./test.db")
	ds.Set("hello", "world")
	ds.Set("foo", "oro")
	res := ds.Get("foo")
	fmt.Println(res)
}
