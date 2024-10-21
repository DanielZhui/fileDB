package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
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

func (d *DiskStore) initKeyDir(exitingFile string) error {
	file, err := os.Open(exitingFile)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()
	for {
		header := make([]byte, headerSize)
		_, err := io.ReadFull(file, header)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read header: %w", err)
		}
		timestamp, keySize, valueSize := decodeHeader(header)
		key := make([]byte, keySize)
		value := make([]byte, valueSize)
		_, err = io.ReadFull(file, key)
		if err != nil {
			return fmt.Errorf("failed to read key: %w", err)
		}
		_, err = io.ReadFull(file, value)
		if err != nil {
			return fmt.Errorf("failed to read value: %w", err)
		}
		totalSize := headerSize + keySize + valueSize
		d.keyInfo[string(key)] = NewKeyEntry(timestamp, uint32(d.position), uint32(totalSize))
		d.position += int(totalSize)
		fmt.Printf("loaded key=%s, value=%s\n", key, value)
	}
	return nil
}

func InitDiskStore(fileName string) (*DiskStore, error) {
	ds := &DiskStore{keyInfo: make(map[string]KeyInfo)}
	if isFileExist(fileName) {
		ds.initKeyDir(fileName)
	}

	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		file.Close()
		return nil, err
	}
	ds.file = file
	return ds, nil
}

func (d *DiskStore) Get(key string) (string, error) {
	kInfo, ok := d.keyInfo[key]
	if !ok {
		return "", fmt.Errorf("key not found: %s", key)
	}
	_, err := d.file.Seek(int64(kInfo.position), defaultWhence)
	if err != nil {
		return "", fmt.Errorf("seek error: %w", err)
	}
	data := make([]byte, kInfo.totalSize)
	_, err = io.ReadFull(d.file, data)
	if err != nil {
		return "", fmt.Errorf("read error: %w", err)
	}
	_, _, value := decodeKV(data)
	return value, nil
}

func (d *DiskStore) write(data []byte) error {
	if _, err := d.file.Write(data); err != nil {
		return err
	}
	if err := d.file.Sync(); err != nil {
		return err
	}
	return nil
}

func (d *DiskStore) Set(key string, value string) error {
	timestamp := uint32(time.Now().Unix())
	size, data := encodeKV(timestamp, key, value)
	if err := d.write(data); err != nil {
		return fmt.Errorf("write error: %w", err)
	}
	d.keyInfo[key] = NewKeyEntry(timestamp, uint32(d.position), uint32(size))
	return nil
}

func (d *DiskStore) Delete(key string) error {
	if _, ok := d.keyInfo[key]; !ok {
		return fmt.Errorf("key not found: %s", key)
	}
	delete(d.keyInfo, key)

	return nil
}

func (d *DiskStore) Update() {
	//
}

func (d *DiskStore) List() {
	//
}

func main() {
	fmt.Println("Hello, file DB!")
	ds, err := InitDiskStore("./test.db")
	if err != nil {
		log.Printf("Failed to initialize disk store: %v", err)
	}
	err = ds.Set("hello", "world")
	if err != nil {
		log.Printf("Failed to set 'hello': %v", err)
	}
	err = ds.Set("foo", "oro")
	if err != nil {
		log.Printf("Failed to set 'foo': %v", err)
	}
	res, err := ds.Get("hello")
	if err != nil {
		log.Printf("Failed to get 'hello': %v", err)
	} else {
		fmt.Println(res)
	}
	ds.Delete("foo")
}
