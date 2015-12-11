package ebi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

type Entry struct {
	Hash string
}

type Imap struct {
	Name    string
	Entries map[string]Entry
}

var NullImap = Imap{}

//opens index file and finds value associated with "key"
func IndexResolve(indexFile string, key string) (string, error) {

	imap, err := ReadIndex(indexFile)
	if err != nil {
		return "", err
	}

	//Find key in map
	value, exists := imap.Entries[key]
	if !exists {
		return "", fmt.Errorf("Index does not contain entry for key")
	}

	return value.Hash, nil

}

func ReadIndexFile(indexPath string) (Imap, error) {

	if _, err := os.Stat(indexPath); err != nil {
		return NullImap, fmt.Errorf("Index file does not exist")
	}

	indexData, err := ioutil.ReadFile(indexPath)
	if err != nil {
		return NullImap, fmt.Errorf("Unable to read index file: %s", indexPath)
	}

	var imap Imap //imap map[string]Entry
	if err := json.Unmarshal(indexData, &imap); err != nil {
		return NullImap, fmt.Errorf("Failed to read index")
	}

	return imap, nil
}

func ReadIndex(indexFile string) (Imap, error) {
	indexPath := path.Join(Index, indexFile)

	ret, err := ReadIndexFile(indexPath)
	if err != nil {
		return NullImap, err
	}

	return ret, nil
}

func (imap *Imap) SetKey(key string, value string) error {
	if value == "" {
		delete(imap.Entries, key)
	} else {
		var entry Entry
		entry.Hash = value
		imap.Entries[key] = entry
	}

	return nil
}

func WriteIndex(imap Imap) error {
	indexData, err := json.Marshal(imap)
	if err != nil {
		return fmt.Errorf("Unable to Marshal index to json")
	}

	indexPath := path.Join(Index, imap.Name)

	err = ioutil.WriteFile(indexPath, indexData, 0644)
	if err != nil {
		return fmt.Errorf("Unable to write Index file")
	}

	return nil
}

//Convenience for opening, editing and writing an index file
func AddEntry(indexFile string, key string, value string) error {
	imap, err := ReadIndex(indexFile)
	if err != nil {
		return err
	}

	err = imap.SetKey(key, value)
	if err != nil {
		return err
	}

	err = WriteIndex(imap)
	if err != nil {
		return err
	}

	return nil
}

//Creating an empty index
func NewIndex(indexFile string) error {
	indexPath := path.Join(Index, indexFile)

	if _, err := os.Stat(indexPath); err == nil {
		return fmt.Errorf("Index file already exists")
	}

	var imap Imap
	imap.Name = indexFile
	imap.Entries = make(map[string]Entry)

	err := WriteIndex(imap)
	if err != nil {
		return err
	}

	return nil
}
