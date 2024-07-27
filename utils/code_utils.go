package utils

import (
	"bytes"
	"encoding/base64"
	"github.com/google/uuid"
	"strings"
)

func GenerateAuthCode(clientId string, userId string) (string, error) {
	buf := bytes.NewBufferString(clientId)
	buf.WriteString(userId)
	token := uuid.NewMD5(uuid.Must(uuid.NewRandom()), buf.Bytes())
	code := base64.URLEncoding.EncodeToString([]byte(token.String()))
	code = strings.ToUpper(strings.TrimRight(code, "="))
	return code, nil
}
