package tools

import (
	"bytes"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/mail"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

const numbers = "0123456789"

const SpecialSymbols = "_=+-/|!@#$%^&*()"

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

func GetBasePath() string {
	log.Println(basepath)
	path := strings.ReplaceAll(basepath, "/tools", "")
	log.Println(path)
	return path
}

// image formats and magic numbers
var magicTable = map[string]string{
	"\xff\xd8\xff":      "image/jpeg",
	"\x89PNG\r\n\x1a\n": "image/png",
	"GIF87a":            "image/gif",
	"GIF89a":            "image/gif",
}

func EncryptTextSHA512(text string) string {

	bytes := []byte(text)
	enc := sha512.Sum512(bytes)
	stringenc := hex.EncodeToString(enc[:])
	return stringenc
}

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func RandomString(length int) string {
	return StringWithCharset(length, charset)
}

func RandomNumbers(length int) string {
	return StringWithCharset(length, numbers)
}

func AgeFromDate(date time.Time) int {
	now := time.Now()
	days := now.Sub(date).Hours() / 24
	years := days / 365
	return int(years)
}

func GetContentType(data []byte) {
	t := ""
	for magic, mime := range magicTable {
		byteString := string(data)
		if strings.HasPrefix(byteString, magic) {
			t = mime
		}
	}
	log.Println(t)
}

// mimeFromIncipit returns the mime type of an image file from its first few
// bytes or the empty string if the file does not look like a known file type
func mimeFromIncipit(incipit []byte) string {
	incipitStr := []byte(incipit)
	for magic, mime := range magicTable {
		byteString := string(incipitStr)
		if strings.HasPrefix(byteString, magic) {
			return mime
		}
	}
	return ""
}

func GetAmericanDateStringFrom(dateString string) string {
	parts := strings.Split(dateString, "/")
	newDateString := parts[2] + "-" + parts[1] + "-" + parts[0]
	return newDateString
}

func HeadersToString(headers http.Header) string {
	b := new(bytes.Buffer)
	for key, value := range headers {
		fmt.Fprintf(b, "%s=\"%s\"\n", key, value)
	}
	return b.String()
}

func RemoveSpecialCaracters(text string) string {
	result := strings.ReplaceAll(text, "(", "")
	result = strings.ReplaceAll(result, ")", "")
	result = strings.ReplaceAll(result, "-", "")
	result = strings.ReplaceAll(result, " ", "")
	result = strings.ReplaceAll(result, "+", "")
	result = strings.ReplaceAll(result, " ", "")
	return result
}

func DateStringToDate(dateString string) (date time.Time, err error) {
	layout := "2006-01-02"
	date, err = time.Parse(layout, dateString)
	if err != nil {
		return date, err
	}
	return date, nil
}

/************************************************
/**** MARK: PASSWORD SECTION ****/
/************************************************/

/* CheckPassword: Verifica se a senha é valida.
Caso não seja, retorna o critério desejado;
Caso seja, retorna uma string vazia; */
func CheckPassword(password string) string {
	message := "senha válida. A senha precisa ter 6 ou mais caracteres, letras maiúsculas, minúsculas e números."
	if len(password) < 6 {
		return message
	}
	if !strings.ContainsAny(password, numbers) {
		return message
	}
	return ""
}

func ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func ReplaceWithRegex(content string, regex string, replace string) string {
	var re = regexp.MustCompile(regex)
	s := re.ReplaceAllString(content, replace)
	return s
}

func MapToString(m map[string]string) string {
	b := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(b, "%s=\"%s\"\n", key, value)
	}
	return b.String()
}

func MapInterfaceToString(m map[string]interface{}) string {
	b := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(b, "%s=\"%s\"\n", key, value)
	}
	return b.String()
}
