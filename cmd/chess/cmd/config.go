package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/dumbogo/chess/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(defaultConfigCmd)
	configCmd.AddCommand(viewConfigCmd)
	configCmd.AddCommand(alphaVersion)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Game configuration",
	Long:  "Set/Get game configuration",
}

var defaultConfigCmd = &cobra.Command{
	Use:   "default",
	Short: "Print default configuration",
	Long:  "Print default configuration with mandatory fields to play",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := config.NewClientConfiguration(config.WithDefaultBaseClientConfiguration())
		if err != nil {
			log.Fatalf("Unexpected error: %s", err)
		}
		str, err := c.Marshal()
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s", str)
	},
}

var alphaVersion = &cobra.Command{
	Use:   "alpha",
	Short: "Configure to play alpha version",
	Long:  "Alpha version is a playable version, configured to connect to a test server",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			homeDir string
			err     error
		)

		c, err := config.NewClientConfiguration(config.WithTestServer())
		if err != nil {
			log.Fatalf("Unexpected error: %s", err)
		}
		str, err := c.Marshal()
		if err != nil {
			panic(err)
		}

		if homeDir, err = os.UserHomeDir(); err != nil {
			panic("home directory not found")
		}

		// TODO: Once we have CA root  certificates provided and configured, this willbe deleted
		// once we have a stable environment
		fakeCert := `-----BEGIN CERTIFICATE-----
MIIC6TCCAo+gAwIBAgIUSNQUG/7271BeZEbrXlbKnase7dQwCgYIKoZIzj0EAwIw
gYMxCzAJBgNVBAYTAlVTMQ4wDAYDVQQIDAVUZXhhczEPMA0GA1UEBwwGQXVzdGlu
MQswCQYDVQQKDAJNUzELMAkGA1UECwwCREMxGDAWBgNVBAMMD3d3dy5jb250b3Nv
LmNvbTEfMB0GCSqGSIb3DQEJARYQdXNlckBjb250b3NvLmNvbTAeFw0yMjA1MTIw
NzE0NTRaFw0yMzA1MTIwNzE0NTRaMGwxCzAJBgNVBAYTAlVTMQ0wCwYDVQQIDARV
dGFoMQ0wCwYDVQQHDARMZWhpMRcwFQYDVQQKDA5mYWJyaWthbSwgSW5jLjELMAkG
A1UECwwCSVQxGTAXBgNVBAMMEHd3dy5mYWJyaWthbS5jb20wWTATBgcqhkjOPQIB
BggqhkjOPQMBBwNCAARZQ4DEC3OTSiIv4ZsbXEGsKrXxOTBcnmH9xKe7/HyOrTV7
4gSvjHKaiMaTpsTjSpCnHjt4Ad7TKEWqw7DXkb/Fo4H2MIHzMC0GA1UdEQQmMCSC
EHd3dy5mYWJyaWthbS5jb22CEHd3dy5mYWJyaWthbS5jb20wHQYDVR0OBBYEFEnk
ZWtMETesmc820XH8hk4gaHToMIGiBgNVHSMEgZowgZehgYmkgYYwgYMxCzAJBgNV
BAYTAlVTMQ4wDAYDVQQIDAVUZXhhczEPMA0GA1UEBwwGQXVzdGluMQswCQYDVQQK
DAJNUzELMAkGA1UECwwCREMxGDAWBgNVBAMMD3d3dy5jb250b3NvLmNvbTEfMB0G
CSqGSIb3DQEJARYQdXNlckBjb250b3NvLmNvbYIJAMzUFESSG/g/MAoGCCqGSM49
BAMCA0gAMEUCIQCxi060mk415y6v/oZAf3igwrEr6oLR3XGmrlDL4veKvwIgESE2
oSMVzGbTQIMDwzbVXVn6Ne0iWyRyvrf3lk5Gc0s=
-----END CERTIFICATE-----
`

		e := os.MkdirAll(homeDir+"/.chess/certs/x509/", os.ModePerm)
		if e != nil {
			panic(e)
		}

		err = ioutil.WriteFile(homeDir+"/.chess/certs/x509/client.crt", []byte(fakeCert), 0644)
		if err != nil {
			panic(err)
		}

		ioutil.WriteFile(homeDir+"/.chess/config.toml", str, 0644)

		fmt.Printf("All settings configured to play on alpha server!\n")
	},
}

var viewConfigCmd = &cobra.Command{
	Use:   "view",
	Short: "Show current configuration",
	Long:  "Print current configuration client chess game",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := config.LoadClientConfiguration()
		if err != nil {
			log.Fatalf("Unexpected Error %s", err)
		}
		str, err := c.Marshal()
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s", str)
	},
}

// TODO: chess config credentials, sub command to configure credentias
// Must investigate about how to fix the problem of using CA credentials in the project
