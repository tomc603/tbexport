/*
 *    Copyright 2018 Tom Cameron
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 *
 */

package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

type Manifest struct {
	Revision int `xml:"revision,attr"`
	Notes []NoteEntry `xml:"note"`
}

type NoteEntry struct {
	Id string `xml:"id,attr"`
	Revision int `xml:"rev,attr"`
}

func (s Manifest) String() string {
	return fmt.Sprintf("Revision: %d\nNotes: %s\n", s.Revision, s.Notes)
}

func (n NoteEntry) String() string {
	return fmt.Sprintf("Revision: %d, Id: %s\n", n.Revision, n.Id)
}

func readManifest(path string) (*Manifest, error) {
	var xmlData Manifest

	manifestFile, readErr := ioutil.ReadFile(path)
	if readErr != nil {
		log.Printf("Error opening %s. %s", path, readErr)
		return nil, readErr
	}

	unmarshallErr := xml.Unmarshal(manifestFile, &xmlData)
	if unmarshallErr != nil {
		log.Printf("Error unmarshalling %s. %s.\n", path, unmarshallErr)
		return nil, unmarshallErr
	}
	return &xmlData, nil
}

func processManifest(path string, manifest *Manifest) {
	for _, note := range manifest.Notes {
		notePath := filepath.Join(path, strconv.Itoa(note.Revision), note.Id) + ".note"
		fmt.Printf("Note: %s\n", notePath)
	}
}

func walker(path string, info os.FileInfo, err error) error {
	if err != nil {
		log.Printf("Error accessing path \"%s\".\n", path)
		log.Printf("%s\n", err)
		return err
	}

	switch mode := info.Mode(); {
	case mode.IsRegular():
		readManifest(path)
	}
	return nil
}

func main() {
	var outPath string
	var inPath string
	var saveRevisions bool

	flag.StringVar(&outPath,"out", "", "Output path for converted notes")
	flag.StringVar(&inPath,"in", "", "Source path for Tomboy notes")
	flag.BoolVar(&saveRevisions,"revisions", false, "Export all note revisions")
	flag.Parse()

	log.Printf("Output path: %s\n", outPath)
	log.Printf("Save revisions: %v\n", saveRevisions)

	// TODO: Add a feature to walk the directory structure looking for orphaned note files
	//filepath.Walk(inPath, walker)

	manifest, manifestErr := readManifest(filepath.Join(inPath, "manifest.xml"))
	if manifestErr != nil {
		os.Exit(1)
	}
	processManifest(filepath.Join(inPath, "0"), manifest)
}
