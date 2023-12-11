package dto

import "github.com/Nerzal/gocloak/v13"

type FindFibonacciInput struct {
	N int `json:"n"`
}

type QueueMessageInput struct {
	M int `json:"m"`
}

type UserDTO struct {
	gocloak.User
	Roles interface{} `json:"roles"`
}

type MinUserDTO struct {
	Id   string `json:"id"`
	Name string `json:"name,omitempty"`
}

type ChatWindow struct {
	Id   string `json:"id"`
	Name string `json:"name,omitempty"`
}

type WsInMessageDTO struct {
	To          string `json:"to"`
	ChatWindow  string `json:"chatWindow"`
	Msg         string `json:"msg"`
	MessageType int    `json:"-"`
}

type WsOutMessageDTO struct {
	Sender     string `json:"sender"`
	ChatWindow string `json:"chatWindow"`
	Msg        string `json:"msg"`
}
