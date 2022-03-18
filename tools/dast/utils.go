package dast

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Create env files OWASP ZAP needs for correct authentication
// to the API or stores.
// These files will be created in a tmp file, that will later
// be copied to the Docker container.
// This function returns the path to the dir to be copied
// and a error, if there has been one
// The path should be removed once finished using, not in this
// function, otherwise, the caller cant access the files.
func PrepareContainerEnv(s *Settings) (str string, err error) {

	// Create a dir to write configs to copy
	f, err := os.MkdirTemp("", "talos-*")
	if err != nil {
		return "", errors.New("error: could not create temporal file")
	}

	// If we are running API SCAN (api)
	// Copy helper script files to new temp dir.
	// Create .env and auth_post_request.json files
	// To finish configuring api env.
	// Create new project structure
	err = CopyEmbedDir(helpers, f)
	if err != nil {
		return "", errors.New("Error: " + err.Error())
	}
	if strings.ToUpper(s.ScanType) == "api" {
		// Create json auth_post_request.json file
		if err := CreateAuthPostRequest(s, f+"/helpers"); err != nil {
			return "", err
		}
		// Create .env file
		if err := Createdotenv("Bearer $(/zap/helpers/getToken.sh)", s.ZAPConfigFile, f+"/helpers"); err != nil {
			return "", err
		}

		return f, nil
		// We are going to run a DAST scan against
		// website web store.
		// Create a .env file with required settings
	} else {
		if err := Createdotenv(s.Auth.Website.AccessToken, s.ZAPConfigFile, f+"/helpers"); err != nil {
			return "", err
		}
		return f, nil
	}
}

// Write both parameters to a file named .env located
// in d directory
func Createdotenv(at, zcf, dir string) (err error) {
	s := [3]string{
		"#AUTO-GENERATED FILE.\n#Refer to talos/tools/dast/createdotenv\n#to edit content\n",
		"\nexport ACCESS_TOKEN=\"" + at + "\"",
		"\nexport ZAP_CONFIG_FILE=\"" + zcf + "\"\n"}
	err = os.MkdirAll(dir, 0777)
	if err != nil {
		fmt.Println(err)
		return errors.New("error: Could not create " + dir + " dir")
	}
	f, err := os.Create(dir + "/.env")
	if err != nil {
		fmt.Println(err)
		return errors.New("error: Could not create " + f.Name() + " file")
	}

	for _, str := range s {
		_, err = f.WriteString(str)
		if err != nil {
			return errors.New("error: could not write string")
		}
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return errors.New("error: Could not close " + f.Name() + " file")
	}
	return
}

// Create a file called auth_post_request.json
// And write JSON Put body request in order to
// get authentication for the ZAP api analysis session
func CreateAuthPostRequest(s *Settings, dir string) (err error) {
	b, err := json.Marshal(s.Auth.Api)
	if err != nil {
		return errors.New("error: Could not marshal json")
	}
	err = ioutil.WriteFile(dir+"/auth_post_request.json", b, 0645)
	if err != nil {
		return errors.New("error: could not create " + dir + "/helpers/auth_post_request.json " + "file")
	}
	return
}

// Set website settings for DAST scan to work.
// The config.yaml/config.json parameters are:
// DAST.Auth.website.AccessToken
func SetwebsiteSettings(s *Settings, c *viper.Viper) (err error) {
	at := c.GetString("dast.auth.website.accesstoken")
	if at == "" {
		return errors.New("error: empty value for key accesstoken")
	}
	s.Auth.Website.AccessToken = at
	return
}

// Set api settings for DAST scan to work.
// The config.yaml/config.json parameters are:
// DAST.Auth.api.ClientID
// DAST.Auth.api.ClientSecret
// DAST.Auth.api.GrantType
// DAST.Auth.api.Username
// DAST.Auth.api.Password
func SetapiSettings(s *Settings, c *viper.Viper) (err error) {

	cid := c.GetString("dast.auth.api.clientid")
	if cid == "" {
		return errors.New("error: empty value for key clientid")
	}
	s.Auth.Api.ClientID = cid

	cs := c.GetString("dast.auth.api.clientsecret")
	if cs == "" {
		return errors.New("error: empty value for key clientsecret")
	}
	s.Auth.Api.ClientSecret = cs

	gt := c.GetString("dast.auth.api.granttype")
	if gt == "" {
		return errors.New("error: empty value for key granttype")
	}
	s.Auth.Api.GrantType = gt

	u := c.GetString("dast.auth.api.username")
	if u == "" {
		return errors.New("error: empty value for key username")
	}
	s.Auth.Api.Username = u

	p := c.GetString("dast.auth.api.password")
	if p == "" {
		return errors.New("error: empty value for key password")
	}
	s.Auth.Api.Password = p
	return
}
