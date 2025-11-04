package utils

import (
	"bytes"
	"errors"
	"fmt"
	"html"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/steambap/captcha"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecretKey = []byte(os.Getenv("JWT_SECRET_KEY"))
var magicLinkSecret = []byte(os.Getenv("MAGIC_LINK_SECRET"))


func GenerateCode() string {
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(1000000)
	return fmt.Sprintf("%06d", code)
}

func IsValidEmail(email string) bool {
	regex := `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`
	return regexMatch(email, regex)
}

func IsValidPassword(password string) bool {
	regex := `^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d]{8,}$`
	return regexMatch(password, regex)
}

func regexMatch(text, regex string) bool {
	matched, _ := regexp.MatchString(regex, text)
	return matched
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func IsNotXSS(data string) bool {
	xssPatterns := []string{
		`<script.*?>.*?</script>`, 
		`javascript:`,             
		`<iframe.*?>.*?</iframe>`, 
		`on\w+=['"][^'"]+['"]`,     
	}

	for _, pattern := range xssPatterns {
		if matched, _ := regexp.MatchString(pattern, data); matched {
			return false
		}
	}
	return true
}

func IsNotSQL(data string) bool {
	sqlKeywords := []string{
		`--`,      
		`DROP`,    
		`INSERT`,  
		`SELECT`,  
		`DELETE`,  
		`UPDATE`,  
		`OR 1=1`,   
		`AND 1=1`,  
		`UNION`,    
		`--`,      
		`#`,        
		`' OR '`,   
		`" OR "`,   
	}

	for _, keyword := range sqlKeywords {
		if strings.Contains(strings.ToUpper(data), keyword) {
			return false
		}
	}
	return true
}
// EscapeHTML escapira specijalne HTML znakove.
func EscapeHTML(data string) string {
	return html.EscapeString(data)
}

// UnescapeHTML uklanja escapirane HTML znakove.
func UnescapeHTML(data string) string {
	return html.UnescapeString(data)
}

func GenerateJWT(userID string, email string, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func RefreshJWT(oldToken string) (string, error) {
	claims, err := ValidateJWT(oldToken)
	if err != nil {
		return "", fmt.Errorf("invalid token for refresh: %w", err)
	}

	if time.Now().Unix() > int64(claims["exp"].(float64)) {
		return "", errors.New("token expired")
	}

	return GenerateJWT(claims["user_id"].(string), claims["email"].(string), claims["role"].(string))
}


func ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("token parsing error: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token or claims")
	}

	return claims, nil
}


// GenerateMagicLinkToken generates a JWT token for the magic link
func GenerateMagicLinkToken(email string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 240).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(magicLinkSecret)
}

// ValidateMagicLinkToken validates the magic link token and returns the email
func ValidateMagicLinkToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return magicLinkSecret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		email, _ := claims["email"].(string)
		return email, nil
	}

	return "", err
}

func GenerateCaptcha() ([]byte, error) {
	img, err := captcha.New(200, 100)	
	if err != nil {
		log.Printf("Error generating CAPTCHA: %v", err)
		return nil, err
	}

	// Konvertujemo CAPTCHA sliku u bajtove.
	var buf bytes.Buffer
	if err := img.WriteImage(&buf); err != nil {
		log.Printf("Error writing CAPTCHA to buffer: %v", err)
		return nil, err
	}

	return buf.Bytes(), nil
}