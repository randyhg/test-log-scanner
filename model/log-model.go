package model

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
)

var TimestampRegex = regexp.MustCompile(`(\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2})`)
var fileNameRegex = regexp.MustCompile(`(\w+\.go):(\d+)`)

type LogErrors struct {
	ID         int    `gorm:"primaryKey"`
	Message    string `gorm:"type:text"`
	StackTrace *string
	FailedAt   string
	Hash       *string
}

func Sha256(data string) string {
	match := fileNameRegex.FindStringSubmatch(data)
	if len(match) >= 3 {
		fileName := match[1]
		lineNumber := match[2]
		hashFormat := fmt.Sprintf("%s:%s", fileName, lineNumber)
		hash := sha256.Sum256([]byte(hashFormat))
		hashString := hex.EncodeToString(hash[:])
		return hashString
	}
	return ""
}
