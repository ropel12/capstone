package helper

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	depedency "github.com/education-hub/BE/config/dependency"
	"github.com/golang-jwt/jwt"
	"github.com/mojocn/base64Captcha"
	"golang.org/x/crypto/bcrypt"
)

func GetUid(token *jwt.Token) int {
	parse := token.Claims.(jwt.MapClaims)
	id := int(parse["id"].(float64))

	return id
}
func GetRole(token *jwt.Token) string {
	parse := token.Claims.(jwt.MapClaims)
	return parse["role"].(string)
}
func GetStatus(token *jwt.Token) string {
	parse := token.Claims.(jwt.MapClaims)
	return parse["verified"].(string)
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func VerifyPassword(passhash string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(passhash), []byte(password))
}

func GenerateJWT(id int, role string, is_verified string, dp depedency.Depend) string {
	var informasi = jwt.MapClaims{}
	informasi["id"] = id
	informasi["role"] = role
	informasi["verified"] = is_verified
	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS256, informasi)
	resultToken, err := rawToken.SignedString([]byte(dp.Config.JwtSecret))
	if err != nil {
		log.Println("generate jwt error ", err.Error())
		return ""
	}
	return resultToken
}

func GenerateEndTime(timee string, duration float32) string {
	t, err := time.Parse("2006-01-02 15:04:05", strings.Replace(timee, "T", " ", 1))
	if err != nil {
		log.Printf("error when generate endtime : %v", err)
		return ""
	}
	minute := duration * 60
	return t.Add(time.Minute * time.Duration(int(minute))).Format("2006-01-02 15:04:05")
}
func GenerateExpiretime(timee string, duration int) string {
	t, err := time.Parse("2006-01-02 15:04:05", timee)
	if err != nil {
		return ""
	}
	return t.Add(time.Minute * time.Duration(duration)).Format("2006-01-02 15:04:05")
}
func GenerateInvoice(eventid int, userid int) string {
	rand.Seed(time.Now().UnixNano())

	randomNum := rand.Intn(9999) + 1000

	return fmt.Sprintf("INV-%d%d%d", userid, eventid, randomNum)

}

func GenerateCaptcha() (string, string, error) {
	DriverString := &base64Captcha.DriverString{
		Height:          60,
		Width:           240,
		ShowLineOptions: 0,
		NoiseCount:      0,
		Source:          "1234567890qwertyuioplkjhgfdsazxcvbnm",
		Length:          7,
		Fonts:           []string{"wqy-microhei.ttc"},
	}
	var driver base64Captcha.Driver
	driver = DriverString.ConvertFonts()
	c := base64Captcha.NewCaptcha(driver, base64Captcha.DefaultMemStore)
	id, b64s, err := c.Generate()
	return id, b64s, err
}

func VerifyCaptcha(captcha string, value string) bool {
	var store = base64Captcha.DefaultMemStore
	return store.Verify(captcha, value, true)
}
