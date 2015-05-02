package main

import "net/url"
import "golang.org/x/net/context"

func (client *Client) Disconnect(server, username, address, received, sent string) {
  httpClient := client.Client(context.TODO())
  data := url.Values{}
  data.Add("username", username)
  data.Add("address", address)
  data.Add("received", received)
  data.Add("sent", sent)

  httpClient.PostForm(client.CreateURL("server/" + server + "/disconnect/"), data)
}
