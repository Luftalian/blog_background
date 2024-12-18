package model

import (
	"fmt"
)

// ContactForm はフロントエンドから送信されるデータの構造体です。
type ContactForm struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
}

func SendEmail(form ContactForm) error {
	return fmt.Errorf("Not implemented")
}
