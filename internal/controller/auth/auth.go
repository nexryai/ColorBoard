package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	logger "github.com/nexryai/ColorBoard/internal/logger"
	"github.com/nexryai/ColorBoard/internal/service"
)

type FirebaseClaims struct {
	Name     string `json:"name"`
	Picture  string `json:"picture"`
	Iss      string `json:"iss"`
	Aud      string `json:"aud"`
	AuthTime int64  `json:"auth_time"`
	UserId   string `json:"user_id"`
	Sub      string `json:"sub"`
	Iat      int64  `json:"iat"`
	Exp      int64  `json:"exp"`
	Email    string `json:"email"`
}

type RegisterSessionReq struct {
	Token string `form:"token" json:"token"`
}

var (
	log         = logger.GetLogger("Auth")
	googleCerts = getGoogleCerts()
	ErrTokenIsNotValid = errors.New("token is not valid")
	ErrTokenIsExpired = errors.New("token is expired")
	ErrUntrustedKey = errors.New("untrusted key")
)

func getGoogleCerts() *map[string]string {
	log = logger.GetLogger("Init")
	log.Info("Registering google certifications for authentication...")

	// Googleの公開鍵を取得
	resp, err := http.Get("https://www.googleapis.com/robot/v1/metadata/x509/securetoken@system.gserviceaccount.com")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var data map[string]interface{}
	err = json.Unmarshal([]byte(body), &data)
	if err != nil {
		panic(err)
	}

	keys := make(map[string]string)
	for k, v := range data {
		if str, ok := v.(string); ok {
			log.Info("Registered certification: ", k)
			keys[k] = str
		} else {
			panic(fmt.Errorf("value for key %s is not a string", k))
		}
	}

	return &keys
}

func parseFirebaseJWT(tokenString string) (*FirebaseClaims, error) {
	// JWTのヘッダを解析し署名に用いられている鍵を取得
	parts := strings.Split(tokenString, ".")

	// decode the header
	headerJson, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, err
	}

	var header map[string]interface{}
	err = json.Unmarshal(headerJson, &header)
	if err != nil {
		return nil, err
	}

	kid, ok := header["kid"].(string)
	if (!ok) {
		return nil, ErrUntrustedKey
	}

	certString, ok := (*googleCerts)[kid]
	if (!ok) {
		return nil, ErrUntrustedKey
	}

	block, _ := pem.Decode([]byte(certString))
	if block == nil {
		return nil, err
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaPublicKey := cert.PublicKey.(*rsa.PublicKey)

	// 署名を検証
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if rsaPublicKey == nil {
			return nil, errors.New("public key is nil")
		}

		return rsaPublicKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		firebaseClaims := &FirebaseClaims{}
        for key, value := range claims {
            switch key {
            case "name":
                if str, ok := value.(string); ok {
                    firebaseClaims.Name = str
                }
            case "picture":
                if str, ok := value.(string); ok {
                    firebaseClaims.Picture = str
                }
            case "iss":
                if str, ok := value.(string); ok {
                    firebaseClaims.Iss = str
                }
            case "aud":
                if str, ok := value.(string); ok {
                    firebaseClaims.Aud = str
                }
            case "auth_time":
                if f, ok := value.(float64); ok {
                    firebaseClaims.AuthTime = int64(f)
                }
            case "user_id":
                if str, ok := value.(string); ok {
                    firebaseClaims.UserId = str
                }
            case "sub":
                if str, ok := value.(string); ok {
                    firebaseClaims.Sub = str
                }
            case "iat":
                if f, ok := value.(float64); ok {
                    firebaseClaims.Iat = int64(f)
                }
            case "exp":
                if f, ok := value.(float64); ok {
                    firebaseClaims.Exp = int64(f)
                }
            case "email":
                if str, ok := value.(string); ok {
                    firebaseClaims.Email = str
                }
            }
        }

		if time.Unix(firebaseClaims.Exp, 0).Before(time.Now()) {
			return nil, ErrTokenIsExpired
		} else {
			return firebaseClaims, nil
		}
	} else {
		return nil, ErrTokenIsNotValid
	}
}

func ConfigFirebaseAuthRouter(router *gin.Engine, userService service.IUserService) {
	router.POST("/auth/register-session", func(ctx *gin.Context) {
		// クライアントがFirebase Authから受け取ったセッショントークンをPOSTしてくるので、正しければCookieに載せて返す
		var req RegisterSessionReq
		err := ctx.ShouldBindBodyWithJSON(&req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		claims, err := parseFirebaseJWT(req.Token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"ident": claims.Sub})
	})
}
