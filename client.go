package main

import "golang.org/x/oauth2/clientcredentials"
import "github.com/vaughan0/go-ini"

type Client struct {
  clientcredentials.Config
  Endpoint string
}

func (self *Client) Configure(id string, secret string, tokenUrl string, endpoint string) {
  self.Config.Scopes = []string{"service"}
  self.Config.ClientID = id
  self.Config.ClientSecret = secret
  self.Config.TokenURL = tokenUrl
  self.Endpoint = endpoint
}

func (self *Client) ConfigureFromIni(filename string) {
  file := ini.File{}
  file.LoadFile(filename)
  id, _ := file.Get("oauth", "client_id")
  secret, _ := file.Get("oauth", "client_secret")
  url, _ := file.Get("oauth", "token_endpoint")
  endpoint, _ := file.Get("api", "endpoint")
  self.Configure(id, secret, url, endpoint)
}

func (self *Client) CreateURL(dest string) (url string) {
  url = self.Endpoint + dest
  return
}
