package model

import (
	"bytes"
	"time"
	"unicode/utf8"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // the MySQL GORM dialect
	"golang.org/x/text/encoding/charmap"
)

var dbConnection *gorm.DB
var dbError error

func init() {
	dbConnection, dbError = gorm.Open("mysql", "pwc:pwc$1@/pwc_develop?charset=utf8mb4&parseTime=True&loc=Local")
}

func formatCurrentTime() string {
	return time.Now().Format("20060102150405")
}

func ConvertToUTF8(str string) string {
	var arr = []byte(str)

	if utf8.Valid(arr) {
		return str
	}

	encodedBytes, err := charmap.Windows1251.NewDecoder().Bytes(arr)
	if err == nil {
		if utf8.Valid(encodedBytes) {
			return string(encodedBytes)
		}
		arr = encodedBytes
	}

	var buf bytes.Buffer
	var r rune

	for _, b := range arr {
		switch b {
		case 0x80:
			r = 0x20AC
		case 0x82:
			r = 0x201A
		case 0x83:
			r = 0x0192
		case 0x84:
			r = 0x201E
		case 0x85:
			r = 0x2026
		case 0x86:
			r = 0x2020
		case 0x87:
			r = 0x2021
		case 0x88:
			r = 0x02C6
		case 0x89:
			r = 0x2030
		case 0x8A:
			r = 0x0160
		case 0x8B:
			r = 0x2039
		case 0x8C:
			r = 0x0152
		case 0x8E:
			r = 0x017D
		case 0x91:
			r = 0x2018
		case 0x92:
			r = 0x2019
		case 0x93:
			r = 0x201C
		case 0x94:
			r = 0x201D
		case 0x95:
			r = 0x2022
		case 0x96:
			r = 0x2013
		case 0x97:
			r = 0x2014
		case 0x98:
			r = 0x02DC
		case 0x99:
			r = 0x2122
		case 0x9A:
			r = 0x0161
		case 0x9B:
			r = 0x203A
		case 0x9C:
			r = 0x0153
		case 0x9E:
			r = 0x017E
		case 0x9F:
			r = 0x0178
		default:
			r = rune(b)
		}

		buf.WriteRune(r)
	}

	return string(buf.Bytes())
}
