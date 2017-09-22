package main

type Response struct {
	Errno int `json:"errno,omitempty" bson:"errno,omitempty"`
	Msg string `json:"msg,omitempty" bson:"msg,omitempty"`
}
