package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

func encodeBasicAuth(user string, passwd string) string {
	data := []byte(user + ":" + passwd)
	str := base64.StdEncoding.EncodeToString(data)
	authHdr := "Basic " + str
	return authHdr
}

func validToken(t string) bool {
	// Parsing to get expiration only, no need to validate signature
	token, _ := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		return nil, nil
	})
	if token == nil {
		return false
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}
	e := claims["exp"]
	if e == nil {
		return false
	}
	exp := int64(e.(float64))
	if exp < time.Now().Unix() {
		return false
	}
	return true
}

func readConfig() (config, error) {
	// Get current user's home directory
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	dirPath := filepath.Join(usr.HomeDir, ".idaas")
	configPath := filepath.Join(dirPath, "config.json")
	content, _ := ioutil.ReadFile(configPath)
	var c config
	err = json.Unmarshal(content, &c)
	if err != nil {
		return c, err
	}
	return c, nil
}

func writeConfig(c config) error {
	// Get current user's home directory
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	dirPath := filepath.Join(usr.HomeDir, ".idaas")
	configPath := filepath.Join(dirPath, "config.json")

	// Create IDAAS config directory if doesn't exist
	_, err = os.Stat(dirPath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dirPath, 0755)
		if err != nil {
			return err
		}
	}
	// Write config as JSON
	content, err := json.MarshalIndent(c, "", "   ")
	err = ioutil.WriteFile(configPath, content, 0644)
	if err != nil {
		return err
	}
	return nil
}

func refreshToken(c config) (config, error) {

	body := config{
		UserID:     c.UserID,
		CustomerID: c.CustomerID,
		Rtoken:     c.Rtoken}
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return c, err
	}

	// TODO: Move below values to config file or constants?
	u := url.URL{Scheme: "http", Host: "127.0.0.1:3000", Path: "api/v1/auth/refresh"}
	client := &http.Client{}
	req, err := http.NewRequest("POST", u.String(), bytes.NewReader(bodyJSON))
	if err != nil {
		return c, err
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return c, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return c, err
	}
	var res returnMsg
	err = json.Unmarshal(b, &res)
	if err != nil {
		return c, err
	}
	if res.Status.Code != 200 {
		return c, errors.New(res.Status.Message)
	}
	c.Atoken = res.Data[0].AccessToken
	c.Rtoken = res.Data[0].RefreshToken

	return c, nil
}

func userLogin(c *config, pw string) error {
	body := login{CustomerID: c.CustomerID, IP: "127.0.0.1"}
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return err
	}
	// TODO: Move below values to config file or constants?
	u := url.URL{Scheme: "http", Host: "127.0.0.1:3000", Path: "api/v1/auth/token"}
	client := &http.Client{}
	req, err := http.NewRequest("POST", u.String(), bytes.NewReader(bodyJSON))
	if err != nil {
		return err
	}
	auth := encodeBasicAuth(c.UserID, pw)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", auth)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var res returnMsg
	err = json.Unmarshal(b, &res)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	if res.Status.Code != 200 {
		log.Println(res.Status.Message)
		os.Exit(1)
	}
	c.Atoken = res.Data[0].AccessToken
	c.Rtoken = res.Data[0].RefreshToken

	return nil
}
