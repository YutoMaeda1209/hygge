package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type discordIdentify struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

func fetchDiscordIdentify(client *http.Client) (*discordIdentify, error) {
	resp, err := client.Get("https://discord.com/api/users/@me")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, errors.New(string(body))
	}

	var user discordIdentify
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}
