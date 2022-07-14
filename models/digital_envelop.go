package models

type DigitalEnvelope struct {
	Data   string `json:"data" form:"data"`
	Key    string `json:"key" form:"key"`
	Iv     string `json:"iv" form:"iv"`
	Digest string `json:"digest" form:"digest"`
}
