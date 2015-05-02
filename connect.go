package main

import "io/ioutil"
import "encoding/json"
import "net/url"
import "golang.org/x/net/context"

type ConnectResponse struct {
  Address string
  Type string
  Mask string
}

func ParseConnect(raw []byte) (ConnectResponse, error) {
  var result ConnectResponse
  err := json.Unmarshal(raw, &result)
  return result, err
}

func (client *Client) Connect(server, username, token, origin string) ([]byte, bool) {
  httpClient := client.Client(context.TODO())
  data := url.Values{}
  data.Add("username", username)
  data.Add("token", token)
  data.Add("origin", origin)

  resp, _ := httpClient.PostForm(client.CreateURL("server/" + server + "/connect/"), data)
  body, _ := ioutil.ReadAll(resp.Body)
  resp.Body.Close()

  if resp.StatusCode != 200 {
    return body, true
  }
  return body, false
}
