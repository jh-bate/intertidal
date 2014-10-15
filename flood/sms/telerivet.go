package sms

import ()

type (
	TelerivetClient struct {
		config TelerivetConfig
	}
	TelerivetConfig struct {
		UserId     string `json:"-"`
		AccountSid string `json:"accountSid"`
		AuthToken  string `json:"authToken"`
	}
)
