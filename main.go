package main

import "os"
import "io/ioutil"
import "fmt"
import "net"
import "gopkg.in/alecthomas/kingpin.v1"

var (
  app = kingpin.New("magpie-backend", "OpenVPN helper for Magpie service")
  config = app.Flag("config", "Configuration file").Required().ExistingFile()

  credentials = app.Command("credentials", "Check user credentials")
  credentialsServer = credentials.Arg("server", "Server instance").Required().String()

  connect = app.Command("connect", "Get connection configuration")
  connectServer = connect.Arg("server", "Server instance").Required().String()
  connectOutfile = connect.Arg("outfile", "Configuration output file").Required().String()

  disconnect = app.Command("disconnect", "Send disconnection signal")
  disconnectServer = disconnect.Arg("server", "Server instance").Required().String()

  client Client
)

func Credentials(server string) {
  username := os.Getenv("username")
  token := os.Getenv("password")
  origin := os.Getenv("untrusted_ip")

  body, err := client.Connect(server, username, token, origin)
  if err {
    os.Exit(1)
  }

  ioutil.WriteFile("/tmp/magpie-" + server + "-" + username, body, os.ModePerm)
}

func Connect(server, outfile string) {
  username := os.Getenv("common_name")
  filename := "/tmp/magpie-" + server + "-" + username
  raw, _ := ioutil.ReadFile(filename)
  data, _ := ParseConnect(raw)
  fmt.Println(data)

  file, _ := os.Create(outfile)
  defer file.Close()

  file.WriteString("ifconfig-push ")
  file.WriteString(data.Address + " ")

  if data.Type == "TUN" {
    ip := net.ParseIP(data.Address)
    ip[len(ip)-1] += 1
    file.WriteString(ip.String())
  } else {
    file.WriteString(data.Mask)
  }

  file.WriteString("\npush \"redirect-gateway\"")
}

func Disconnect(server string) {
  username := os.Getenv("common_name")
  address := os.Getenv("ifconfig_pool_remote_ip")
  received := os.Getenv("bytes_received")
  sent := os.Getenv("bytes_sent")

  client.Disconnect(server, username, address, received, sent)
}

func main() {
  command, err := app.Parse(os.Args[1:])

  if *config != "" {
    client.ConfigureFromIni(*config)
  }

  switch kingpin.MustParse(command, err) {
    case credentials.FullCommand():
      Credentials(*credentialsServer)
      break

    case connect.FullCommand():
      Connect(*connectServer, *connectOutfile)
      break

    case disconnect.FullCommand():
      Disconnect(*disconnectServer)
      break

    default:
      app.Usage(os.Stdout)
  }
}
