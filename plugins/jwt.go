package main

import (
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"ioa"
	"ioa/proto"
	"net/http"
	"strings"
)

type Plugin struct {
	ioa.BasePlugin
}

type Data struct {
}

type Config struct {
	JwtSecret  string   `json:"jwtSecret"`
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

func (i Plugin) GetName() string {
	return name
}

func (i Plugin) GetDescribe() string {
	return "jwt Authorization Bearer"
}

func (i Plugin) GetConfigTemplate() proto.ConfigTpl {
	configTpl := proto.ConfigTpl{
		{Name: "jwtSecret", Desc: "jwt secret key", Required: true, FieldType: "string"},
		{Name: "claimsKeys", Desc: "fields in token, separated by , (e.g.:user_id,user_name)", Required: true, FieldType: "string"},
	}

	return configTpl
}

func (i Plugin) InitApi(api *ioa.Api) error {
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

func (i Plugin) InitApiData(api *ioa.Api) error {
	return nil
}

func (i Plugin) InitApiConfig(api *ioa.Api) error {
	var config Config
	json.Unmarshal(api.PluginRawConfig[name], &config)
	i.Logger().Debug("this is config***********", config)
	api.PluginConfig[name] = config
	return nil
}

func (i Plugin) Run(ctx ioa.Context) error {
	config := ctx.Api.PluginConfig[name].(Config)
	jwtSecret := config.JwtSecret
	claimsKeys := config.ClaimsKeys

	authorization := ctx.Request.Header.Get("Authorization")
	if !strings.Contains(authorization, "Bearer") {
		ctx.ResponseWriter.WriteHeader(http.StatusUnauthorized)
		ctx.ResponseWriter.Write([]byte("error Header Authorization"))
		return nil
	}

	token := string([]byte(authorization)[7:])
	if token == "null" {
		ctx.ResponseWriter.WriteHeader(http.StatusUnauthorized)
		ctx.ResponseWriter.Write([]byte("error Header Authorization"))
		return nil
	}

	claims, err := parseJwtToken(jwtSecret, token)
	if err != nil {
		ctx.ResponseWriter.WriteHeader(http.StatusUnauthorized)
		ctx.ResponseWriter.Write([]byte("parse token err"))
		return nil
	}

	//todo claims.VerifyExpiresAt
	for _, claimsKey := range claimsKeys {
		ctx.Request.Header.Add(claimsKey, claims[claimsKey].(string))
	}

	return nil
}

func (i Plugin) throwErr(err error) error {
	return errors.New("plugin" + name + err.Error())
}

func parseJwtToken(key, token string) (jwt.MapClaims, error) {
	t, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(key), nil
	})

	if mc, ok := t.Claims.(jwt.MapClaims); ok {
		return mc, nil
	}
	return nil, errors.New("interface.(jwt.MapClaims) error")
}

var ExportPlugin Plugin
