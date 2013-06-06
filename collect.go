package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"crypto/sha1"
	"io"
	"io/ioutil"
	"encoding/json"
	"regexp"
)

type Entry struct {
	collection string
	path string
	hash string
}

var hashRoot = "/mnt/us/documents/"
var hashRegex = regexp.MustCompile(`id="uuid_id">([\w-]+)</dc:identifier>`)

func listEntries(path string) (entries []Entry, e error) {
	walker := func (path string, info os.FileInfo, err error) (e error) {
		if err != nil {
			return err
		}

		if filepath.Ext(path) != ".mobi" {
			return
		}

		pathParts := strings.Split(path, `\`)

		var hash string
		if metadata, err := ioutil.ReadFile(filepath.Dir(path) + "/metadata.opf"); err == nil {
			hash = "#" + string(hashRegex.FindSubmatch(metadata)[1]) + "^EBOK"
		} else {
			relativePath := strings.Join(pathParts[2:], "/")
			sha := sha1.New()
			io.WriteString(sha, hashRoot + relativePath)
			hash = "*" + fmt.Sprintf("%x", sha.Sum(nil))
		}

		entry := Entry{pathParts[2] + "@en-US", path, hash}
		entries = append(entries, entry)
		return 
	}

	filepath.Walk(path, walker)
	return
}

type Collection struct {
	Items []string `json:"items"`
	LastAccess int `json:"lastAccess"`
}
type Collections map[string]Collection

func buildCollections(entries []Entry) Collections {
	c := make(Collections)

	for _, entry := range entries {
		collection, _ := c[entry.collection]
		collection.Items = append(collection.Items, entry.hash)
		c[entry.collection] = collection
	}

	return c
}


func main() {
	entries, _ := listEntries("E:/documents")
	collections := buildCollections(entries)
	encoded, _ := json.Marshal(collections)
	ioutil.WriteFile("collections.json", encoded, 0666)
}
