package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type TokenUtils struct{}

func NewTokenUtils() *TokenUtils {
	return &TokenUtils{}
}

func (tokenUtils *TokenUtils) Token(clientId string, userId string, createdAt time.Time) (string, error) {
	buf := bytes.NewBufferString(clientId)
	buf.WriteString(userId)
	buf.WriteString(strconv.FormatInt(createdAt.UnixNano(), 10))

	tokenUUID := uuid.NewMD5(uuid.Must(uuid.NewRandom()), buf.Bytes())
	tokenHash := sha256.Sum256([]byte(tokenUUID.String()))
	token := base64.URLEncoding.EncodeToString(tokenHash[:])
	token = strings.ToUpper(strings.TrimRight(token, "="))
	return token, nil
}
