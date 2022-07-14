package config

import (
	"encoding/json"
	"log"
	"os"
)

type Configuration struct {
	Environment      string       `json:"environment"`
	Host             string       `json:"host"`
	AdminURL         string       `json:"admin_url"`
	ApiPort          string       `json:"api-port"`
	ApiName          string       `json:"api-name"`
	DatabaseIp       string       `json:"database-ip"`
	DatabasePort     string       `json:"database-port"`
	DatabaseName     string       `json:"database-name"`
	DatabaseUsername string       `json:"database-username"`
	DatabasePassword string       `json:"database-password"`
	Email1           string       `json:"email1"`
	EmailPassword1   string       `json:"email1-password"`
	EmailDomain1     string       `json:"email1-domain"`
	EmailServer1     string       `json:"email1-server"`
	EmailPort1       string       `json:"email1-port"`
	Email2           string       `json:"email2"`
	EmailPassword2   string       `json:"email2-password"`
	EmailDomain2     string       `json:"email2-domain"`
	EmailServer2     string       `json:"email2-server"`
	EmailPort2       string       `json:"email2-port"`
	Site             string       `json:"site"`
	Database         string       `json:"database"`
	LogPath          string       `json:"logPath"`
	PetsInfoPath     string       `json:"petsInfoPath"`
	Features         Features     `json:"features"`
	Integrations     Integrations `json:"integrations"`
}

type Features struct {
	UserShop   bool `json:"userShop"`
	PetModule  bool `json:"petModule"`
	Service    bool `json:"service"`
	Providers  bool `json:"providers"`
	Products   bool `json:"products"`
	Highlights bool `json:"highlights"`
	Dashboard  bool `json:"dashboard"`
	Messages   bool `json:"messages"`
}

type Integrations struct {
	OneSingalKey   string `json:"oneSignalKey"`
	LalamoveKey    string `json:"lalamoveKey"`
	LalamoveSecret string `json:"lalamoveSecret"`
	LalamoveURL    string `json:"lalamoveURL"`
	PagarmeKey     string `json:"pagarmeKey"`
	PagarmeURL     string `json:"PagarmeURL"`
}

/* ****************************************************
**	Exported functions
** ***************************************************/

func Get(path string) Configuration {
	// Tries to open the configuration file
	file, err := os.Open(path)
	if err != nil {
		// If it can't open the configuration file
		log.Println("Arquivo de configuração inválido.")
		log.Fatalln(err)
	}

	// Tries to decode the json file
	decoder := json.NewDecoder(file)
	conf := Configuration{}
	err = decoder.Decode(&conf)
	if err != nil {
		// If it can't decode the configuration file
		log.Println("Error decoding the configuration file...")
		log.Fatalln(err)
	}

	return conf
}
