package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"ioa"
	"ioa/proto"
	"log"
	"net/http"
	"strings"
)

type ioaPlugin struct {
}

type Data struct {
}

type Config struct {
	JwtSecret  string `json:"jwtSecret"`
	ClaimsKeys []string `json:"claimsKeys"`
}

type RawConfig struct {
	JwtSecret  string `json:"jwtSecret"`
	ClaimsKeys string `json:"claimsKeys"`
}

func (c *Config) UnmarshalJSON(b []byte) error {
	rawConfig := RawConfig{}
	err := json.Unmarshal(b, &rawConfig)
	if err != nil {
		panic(err)
	}

	c.JwtSecret = rawConfig.JwtSecret
	c.ClaimsKeys = strings.Split(rawConfig.ClaimsKeys, ",")

	return nil
}

var name = "jwt"

func (s ioaPlugin) GetName() string {
	return name
}

func (s ioaPlugin) GetDescribe() string {
	return "jwt Authorization Bearer"
}

func (s ioaPlugin) GetConfigTemplate() proto.ConfigTpl {
	configTpl := proto.ConfigTpl{
		{Name: "jwtSecret", Desc: "jwt secret key", Required: true, FieldType: "string"},
		{Name: "claimsKeys", Desc: "fields in token, separated by , (e.g.:user_id,user_name)", Required: true, FieldType: "string"},
	}

	return configTpl
}

func (i ioaPlugin) InitApi(api *ioa.Api) error {
	err := i.InitApiConfig(api)
	if err != nil {
		return i.throwErr(err)
	}
	err = i.InitApiData(api)
	if err != nil {
		return i.throwErr(err)
	}

	return nil
}

func (i ioaPlugin) InitApiData(api *ioa.Api) error {
	return nil
}

func (i ioaPlugin) InitApiConfig(api *ioa.Api) error {
	var config Config
	json.Unmarshal(api.PluginRawConfig[name], &config)
	log.Println("this is config***********", config)
	api.PluginConfig[name] = config
	return nil
}

func (s ioaPlugin) Run(w http.ResponseWriter, r *http.Request, api *ioa.Api) error {
	config := api.PluginConfig[name].(Config)
	jwtSecret := config.JwtSecret
	claimsKeys := config.ClaimsKeys

	authorization := r.Header.Get("Authorization")
	if !strings.Contains(authorization, "Bearer") {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("error Header Authorization"))
		return nil
	}

	token := string([]byte(authorization)[7:])
	if token == "null" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("error Header Authorization"))
		return nil
	}

	claims, err := parseJwtToken(jwtSecret, token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("parse token err"))
		return nil
	}

	//todo claims.VerifyExpiresAt
	for _, claimsKey := range claimsKeys {
		r.Header.Add(claimsKey, claims[claimsKey].(string))
	}

	return nil
}

func (i ioaPlugin) throwErr(err error) error {
	return errors.New("plugin" + name + err.Error())
}

func parseJwtToken(key, token string) (jwt.MapClaims, error) {
	t, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(key), nil
	})

	if mc, ok := t.Claims.(jwt.MapClaims); ok {
		return mc, nil
	}
	return nil, fmt.Errorf("interface.(jwt.MapClaims) error")
}

var IoaPlugin ioaPlugin
