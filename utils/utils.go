package utils

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"math"
	"math/rand"
	"net/http"
	"net/mail"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/smartpet/websocket/constant"
	"github.com/smartpet/websocket/models"
)

func ConvertByteToString(v []byte) (data [][]string, err error) {

	reader := csv.NewReader(bytes.NewBuffer(v))
	for {
		line, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				return data, err

			}
		}

		data = append(data, line)

	}

}

func RoundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func GetRequestID(r *http.Request, partyCode string) string {
	// if X-RequestId is not set in request, reqID will be "". In that case, formatted request id will be set.
	reqID := r.Header.Get(constant.RequestIDHeader)
	if reqID == "" {
		currentTime := time.Now().Format(time.RFC3339)
		reqID = fmt.Sprintf("%v_%v", currentTime, partyCode)
	}
	return reqID
}

func JSONErrorResponder(r *http.Request, w http.ResponseWriter, httpCode int, reqID, partycode, description string, reqStartTime time.Time, err error) {
	log.WithFields(log.Fields{
		"reqID":      reqID,
		"statusCode": httpCode,
		"clientID":   partycode,
		"latency":    time.Since(reqStartTime).Milliseconds(),
		"publicIP":   r.RemoteAddr,
		"method":     r.Method,
		"url_path":   r.URL.Path,
		"error":      err,
	}).Errorln(description)

	w.WriteHeader(httpCode)

	resp := models.Response{
		StatusCode:        httpCode,
		StatusDescription: http.StatusText(httpCode),
		Description:       err.Error(),
	}
	o, _ := json.Marshal(resp)
	w.Write(o)

}

func JSONSuccessResponder(ctx *gin.Context, httpCode int, reqID, partycode, description string, reqStartTime time.Time, response interface{}) {
	log.WithFields(log.Fields{
		"reqID":      reqID,
		"statusCode": httpCode,
		"clientID":   partycode,
		"latency":    time.Since(reqStartTime).Milliseconds(),
		"publicIP":   ctx.ClientIP(),
		"method":     ctx.Request.Method,
		"url_path":   ctx.Request.URL.Path,
	}).Infoln(description)

	ctx.JSON(httpCode, models.Response{
		StatusCode:        httpCode,
		StatusDescription: http.StatusText(httpCode),
		Description:       description,
		Response:          response,
	})
}

func IsNumeric(s string) bool {
	val, err := strconv.ParseFloat(s, 64)
	if val <= 0 {
		return false
	}
	return err == nil
}

func IsAlphanumeric(s string) bool {
	re := regexp.MustCompile("^[a-zA-Z0-9_]*$")
	return re.MatchString(s)
}

func IsNumericAndValidInt(s string) bool {
	val, err := strconv.Atoi(s)
	log.Printf("val=%v\n", val)
	if err != nil {
		log.Printf("err=%v\n", err)
		return false
	}

	if val <= 0 {
		return false
	}

	return true
}

func IsValidateUserId(s string) bool {
	_, err := mail.ParseAddress(s)
	if err != nil {
		return IsNumericAndValidInt(s)
	}
	return true
}

func IsUserIdIsEmail(s string) bool {
	_, err := mail.ParseAddress(s)
	return err == nil
}

func IsBlank(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func GenerateRandomOtp() string {
	rand.Seed(time.Now().UnixNano())
	otp := rand.Intn(900000) + 100000
	otpStr := fmt.Sprintf("%06d", otp)

	return otpStr
}

func GetClientIPByHeaders(req *http.Request) (ip string, err error) {

	// Client could be behid a Proxy, so Try Request Headers (X-Forwarder)
	ipSlice := []string{}

	ipSlice = append(ipSlice, req.Header.Get("X-Forwarded-For"))
	ipSlice = append(ipSlice, req.Header.Get("x-forwarded-for"))
	ipSlice = append(ipSlice, req.Header.Get("X-FORWARDED-FOR"))

	for _, v := range ipSlice {
		log.Printf("debug: client request header check gives ip: %v", v)
		if v != "" {
			return v, nil
		}
	}
	err = fmt.Errorf("Could not find clients IP address from the Request Headers")

	return "", err

}

func ValidateDates(startDate, endDate string) error {
	layout := "2006-01-02"

	startTime, err := time.Parse(layout, startDate)
	if err != nil {
		return fmt.Errorf("Invalid startDate format: %v", err)
	}

	endTime, err := time.Parse(layout, endDate)
	if err != nil {
		return fmt.Errorf("Invalid endDate format: %v", err)
	}

	// Compare the dates
	if !startTime.Before(endTime) && !startTime.Equal(endTime) {
		return fmt.Errorf("startDate must be before or equal to endDate")
	}

	return nil
}

func IsValidatePan(pan string) bool {
	re := regexp.MustCompile("[A-Z]{5}[0-9]{4}[A-Z]{1}")
	return re.MatchString(pan)
}

func GetUrlWithoutQueryParams(url string) string {
	queryIndex := strings.Index(url, "?")
	var urlWithoutQuery string
	if queryIndex != -1 {
		urlWithoutQuery = url[:queryIndex]
	} else {
		urlWithoutQuery = url
	}
	return urlWithoutQuery
}

func ConvertStringCammel(str string) string {
	words := strings.Fields(str)
	camelCase := ""
	for _, word := range words {
		camelCase += " " + strings.ToUpper(word[0:1]) + strings.ToLower(word[1:])
	}
	return strings.TrimLeft(camelCase, " ")
}

func StrLoop(str string, length int) string {
	var mask string
	for i := 1; i <= length; i++ {
		mask += str
	}
	return mask
}

func MaskEmail(i string) string {

	l := len([]rune(i))
	if l == 0 {
		return ""
	}

	tmp := strings.Split(i, "@")
	if len(tmp) < 2 {
		return ""
	}

	addr := tmp[0]
	domain := tmp[1]

	addr = Overlay(addr, 2)

	tmpd := strings.Split(domain, ".")
	if len(tmpd) < 2 {
		return ""
	}
	domain = fmt.Sprintf("%s.%s", Overlay(tmpd[0], 2), tmpd[1])

	return fmt.Sprintf("%s@%s", addr, domain)

}

func MaskMobile(i string) string {

	l := len([]rune(i))
	if l == 0 {
		return ""
	}

	return Overlay(i, 4)

}

func Overlay(str string, cnt int) (overlayed string) {
	maskcnt := cnt / 2
	s := "*"
	if len([]rune(str)) > cnt {
		maskcnt = len([]rune(str)) - cnt
	}

	overlayed = ""
	overlayed += StrLoop(s, maskcnt)
	overlayed += string(str[maskcnt:])
	return overlayed
}
func CheckXSSAttack(input map[string]interface{}) bool {
	r, _ := regexp.Compile(`<[^>]*>`)

	for k := range input {
		aValue := reflect.ValueOf(input[k])
		fmt.Printf("%s - %s", k, aValue.Kind())
		fmt.Println("")
		if aValue.Kind() != reflect.String {
			continue
		}
		if matched := r.MatchString(input[k].(string)); matched {
			return true
		}
	}
	return false
}
func GetFileFromBase64(docuemnt, fileext string) ([]byte, string, error) {

	dec, err := base64.StdEncoding.DecodeString(docuemnt)
	if err != nil {
		return nil, "", err
	}

	fileName := fmt.Sprintf("%s.%s", uuid.New().String(), fileext)
	f, err := os.Create(fileName)
	if err != nil {
		return nil, "", err
	}

	defer f.Close()

	if _, err := f.Write(dec); err != nil {
		return nil, "", err
	}
	if err := f.Sync(); err != nil {
		return nil, "", err
	}

	file, err := os.Open(f.Name())
	if err != nil {
		return nil, "", err
	}
	defer file.Close()
	// Get the file size
	stat, err := file.Stat()
	if err != nil {
		return nil, "", err
	}

	bs := make([]byte, stat.Size())
	_, err = bufio.NewReader(file).Read(bs)
	if err != nil && err != io.EOF {
		return nil, "", err
	}

	return bs, f.Name(), nil
}

func GetEmailContent(data map[string]interface{}, cgTemplatePath string) (string, error) {

	templ, err := template.ParseFiles(cgTemplatePath)
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML template file %s", err)
	}

	var body bytes.Buffer
	if err = templ.Execute(&body, data); err != nil {
		return "", fmt.Errorf("failed in adding data to HTML template %s", err)
	}

	emailContent := (body).String()
	return emailContent, nil

}
