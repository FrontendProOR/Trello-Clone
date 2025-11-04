package utils

import (
	"errors"
	"fmt"
	"html"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecretKey = []byte("secret")

func GenerateCode() string {
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(1000000)
	return fmt.Sprintf("%06d", code)
}

func IsValidEmail(email string) bool {
	//regex for email validation check
	regex := `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`
	return regexMatch(email, regex)
}

func IsValidPassword(password string) bool {
	//regex for password validation check
	regex := `^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d]{8,}$`
	return regexMatch(password, regex)
}

func regexMatch(text, regex string) bool {
	matched, _ := regexp.MatchString(regex, text)
	return matched
}

// HashPassword hashes a password using bcrypt and returns the hashed password.
func HashPassword(password string) (string, error) {
	// Generate a bcrypt hash of the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// IsNotXSS checks for potential XSS (Cross-site Scripting) attack patterns
// A simple but important XSS check: looking for <script> or javascript: in input
func IsNotXSS(data string) bool {
	// A basic XSS attack check for script injections
	xssPatterns := []string{
		`<script.*?>.*?</script>`, // <script> tags
		`javascript:`,             // javascript: links
		`<iframe.*?>.*?</iframe>`, // iframe tag
		`<object.*?>.*?</object>`, // object tag
		`<embed.*?>.*?</embed>`,   // embed tag
		`<img.*?>`,                // img tag with potential src attribute injection
		`on\w+=['"][^'"]+['"]`,    // event handlers like onclick, onmouseover, etc.
	}

	for _, pattern := range xssPatterns {
		if matched, _ := regexp.MatchString(pattern, data); matched {
			return false
		}
	}
	return true
}

// IsNotSQL checks for SQL injection patterns by identifying common SQL keywords in inputs
func IsNotSQL(data string) bool {
	// List of SQL keywords that are commonly abused in injection attacks
	sqlKeywords := []string{
		`--`,      // SQL comment
		`DROP`,    // DROP table, column, etc.
		`INSERT`,  // INSERT statement
		`SELECT`,  // SELECT statement
		`DELETE`,  // DELETE statement
		`UPDATE`,  // UPDATE statement
		`OR 1=1`,  // Common injection pattern
		`AND 1=1`, // Common injection pattern
		`UNION`,   // UNION SELECT statement
		`--`,      // SQL comments
		`#`,       // Another form of SQL comment
		`' OR '`,  // Another SQL injection pattern
		`" OR "`,  // Another SQL injection pattern
	}

	for _, keyword := range sqlKeywords {
		if strings.Contains(strings.ToUpper(data), keyword) {
			return false
		}
	}
	return true
}

// GenerateJWT creates a JWT token for the user.
func GenerateJWT(userID string, email string) (string, error) {
	// Set the claims, which include the user's ID and email
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	}
	// Create a new token with the claims and sign it using the secret key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

// ValidateJWT validates the JWT token.
func ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	// Parse the token and validate it
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the token's signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecretKey, nil
	})
	if err != nil {
		return nil, err
	}

	// Check if the token is valid
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.New("invalid token")
	}
}
func EscapeHTML(data string) string {
	return html.EscapeString(data)
}

// UnescapeHTML uklanja escapirane HTML znakove.
func UnescapeHTML(data string) string {
	return html.UnescapeString(data)
}