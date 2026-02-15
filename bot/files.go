package bot

import (
	"os"
	"log"
)

const (
	ICONS = "./assets/icons/"

	attachment = "attachment://"
)

type Image struct {
	Name string
	Type string
}

func (img Image) FullPath() string {
	return img.Type + img.Name
}

func (img Image) Attach() string {
	return attachment + img.Name
}

func (img Image) File() *os.File {
	file, err := os.Open(img.FullPath())
	if (err != nil) {
		log.Printf("Failed to open image: %s\n%v", img.Name, err)
		return nil
	}
	return file
}

