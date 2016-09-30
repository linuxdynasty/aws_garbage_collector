package shared

import (
	"log"
	"os"

	"github.com/asdine/storm"
)

func PrepareDb(filePath string) *storm.DB {
	if _, err := os.Stat(filePath); err == nil {
		if err := os.Remove(filePath); err != nil {
			log.Fatal(err)
		} else {
			log.Printf("DB %s removed", filePath)
		}
	}
	db, _ := storm.Open(filePath, storm.AutoIncrement())
	return db
}
