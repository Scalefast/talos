package dast

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/spf13/viper"
	tester "github.com/scalefast/talos/tester"
)

func TestPrepareContainerEnvapi(t *testing.T) {
	nt := tester.NewTester(t)

	nt.Given("the necesary settings have been set correctly for a api DAST test")
	{
		// Create a new struct Settings
		s := new(Settings)
		s.ScanType = "api"
		s.Auth.api.ClientID = "ID"
		s.Auth.api.ClientSecret = "Secret"
		s.Auth.api.GrantType = "Password"
		s.Auth.api.Username = "Username"
		s.Auth.api.Password = "Password"
		nt.When("we prepare the api environment")
		{
			files, err := PrepareContainerEnv(s)
			var expectedDotEnvPath = files + "/helpers/.env"
			var expectedAPRPath = files + "/helpers/auth_post_request.json"
			nt.Then("we should have a tmp dir containing everything we need for api DAST tests")
			{
				// Check there is a file called .env, and a file called auth_post_request.json
				if !exists(expectedDotEnvPath) {
					nt.Error("file does not exist", expectedDotEnvPath, "Not there")
				}
				if !exists(expectedAPRPath) {
					nt.Error("file does not exist", expectedAPRPath, "Not there")
				}
			}
			nt.And("error shuld return \"nil\"")
			{
				if err != nil {
					nt.Error("expected \"nil\"", "\"nil\"", err.Error())
				}
			}
		}
	}
}

func TestPrepareContainerEnvwebsite(t *testing.T) {
	nt := tester.NewTester(t)

	nt.Given("We have to create necesary files for DAST website tests")
	{
		s := new(Settings)
		s.ScanType = "website"
		s.Auth.website.AccessToken = "basic token"
		nt.When("we prepare the website environment")
		{
			files, err := PrepareContainerEnv(s)
			var expectedDotEnvPath = files + "/helpers/.env"
			nt.Then("")
			{
				if !exists(expectedDotEnvPath) {
					nt.Error("file does not exist", expectedDotEnvPath, "Not there")
				}
			}
			nt.And("error shuld return \"nil\"")
			{
				if err != nil {
					nt.Error("expected \"nil\"", "\"nil\"", err.Error())
				}
			}
		}
	}
}

// This function checks we can create the Auth POST
// request correctly with the parameters given
func TestCreateAuthPostRequestGeneratesJSONFile(t *testing.T) {

	// Initialize tester, so we can use Given When Then standard
	nt := tester.NewTester(t)

	// Set the dir the file is going to be written to.
	var dir = os.TempDir()

	expectedFileContents := []byte("{\"client_id\":\"ID\",\"client_secret\":\"Secret\",\"grant_type\":\"Password\",\"username\":\"Username\",\"password\":\"Password\"}")

	nt.Given("the need to generate a auth_post_request.json file in the temp directory")
	{
		nt.When("supplied with valid values")
		{
			var s = new(Settings)
			s.Auth.api.ClientID = "ID"
			s.Auth.api.ClientSecret = "Secret"
			s.Auth.api.GrantType = "Password"
			s.Auth.api.Username = "Username"
			s.Auth.api.Password = "Password"
			nt.Then("the file should be correctly generated")
			{
				err := CreateAuthPostRequest(s, dir)
				f, _ := ioutil.ReadFile(os.TempDir() + "/auth_post_request.json")
				if !bytes.Equal(f, expectedFileContents) {
					t.Fatalf("Error: Contents of files are not the same\nFile:\t\t\"%s\"\nExpected: \t\"%s\"", string(f), string(expectedFileContents))
				}
				nt.And("the error message should be: \"nil\"")
				{
					if err != nil {
						t.Fatalf("Error message is \"%s\" instead of \"nil\"", err.Error())
					}
				}
			}
		}
	}
}

func TestSetwebsiteSettingIsEmpty(t *testing.T) {

	nt := tester.NewTester(t)

	expectedErr := "error: empty value for key accesstoken"

	nt.Given("The need to store user settings in a settings struct")
	{
		nt.When("The user does not specify the necessary website settings")
		{
			// Create empty Configurations
			s := &Settings{}
			c := &viper.Viper{}
			nt.Then("The function should return a error")
			{
				// Method to test
				err := SetwebsiteSettings(s, c)
				// Check if condition is met
				if err.Error() != expectedErr {
					nt.Error("Should have returned an error", expectedErr, err.Error())
				}
			}
		}
	}
}

func TestSetwebsiteSettingCompletes(t *testing.T) {

	nt := tester.NewTester(t)

	//Initialization of Auth settings
	expectedValue := "basic token"
	//Set the dir the file is going to be written to.

	nt.Given("the need to store website configs in the Settings struct")
	{
		nt.When("the user has specified all the website settings needed")
		{
			s := &Settings{}
			c := viper.New()
			c.Set("dast.auth.website.accesstoken", expectedValue)
			nt.Then("the website Settings struct should be correctly initialized")
			{
				// Method to test
				err := SetwebsiteSettings(s, c)
				// Check if condition is met
				if strings.Compare(s.Auth.website.AccessToken, expectedValue) != 0 {
					nt.Error("website Settings struct has not been initialized correctly", expectedValue, s.Auth.website.AccessToken)
				}
				nt.And("the return value should be \"nil\"")
				if err != nil {
					nt.Error("Should have returned \"nil\" error", "nil", err.Error())
				}
			}
		}
	}
}

func TestCreatedotenvWorking(t *testing.T) {

	nt := tester.NewTester(t)
	var at = "basic token"
	var zcf = "/location"

	expectedFileContents := []byte("#AUTO-GENERATED FILE.\n#Refer to talos/tools/dast/createdotenv\n#to edit content\n\nexport ACCESS_TOKEN=\"" + at + "\"\nexport ZAP_CONFIG_FILE=\"" + zcf + "\"\n")
	//Initialization of Auth settings

	//Set the dir the file is going to be written to.

	nt.Given("the need to generate a .env file in the temp directory")
	{
		nt.When("supplied with valid values")
		{
			// Configuration
			nt.Then("the file is successfully generated")
			{
				// Method to test
				err := Createdotenv(at, zcf, os.TempDir())
				// Check if condition is met
				f, _ := ioutil.ReadFile(os.TempDir() + "/.env")
				if !bytes.Equal(f, expectedFileContents) {
					nt.Error("not nil", string(expectedFileContents), string(f))
				}
				nt.And("the error message should be: \"nil\" ")
				{
					if err != nil {
						nt.Error("not nil", "\"nil\"", err.Error())
					}
				}
			}
		}
	}
}
func TestSetapiSettingCompletes(t *testing.T) {

	nt := tester.NewTester(t)

	//Initialization of Auth settings
	var expectedClientID = "ID"
	var expectedClientSecret = "Secret"
	var expectedGrantType = "Password"
	var expectedUsername = "Username"
	var expectedPassword = "Password"
	//Set the dir the file is going to be written to.

	nt.Given("the need to store website configs in the Settings struct")
	{
		nt.When("the user has specified all the website settings needed")
		{
			s := new(Settings)
			c := viper.New()
			c.Set("dast.auth.api.ClientID", expectedClientID)
			c.Set("dast.auth.api.ClientSecret", expectedClientSecret)
			c.Set("dast.auth.api.GrantType", expectedGrantType)
			c.Set("dast.auth.api.Username", expectedUsername)
			c.Set("dast.auth.api.Password", expectedPassword)
			nt.Then("the website Settings struct should be correctly initialized")
			{
				// Method to test
				err := SetapiSettings(s, c)
				// Check if condition is met
				if strings.Compare(s.Auth.api.ClientID, expectedClientID) != 0 {
					nt.Error("api ClientID has not been set correctly", expectedClientID, s.Auth.api.ClientID)
				}
				if strings.Compare(s.Auth.api.ClientSecret, expectedClientSecret) != 0 {
					nt.Error("api ClientSecret has not been set correctly", expectedClientSecret, s.Auth.api.ClientSecret)
				}
				if strings.Compare(s.Auth.api.GrantType, expectedGrantType) != 0 {
					nt.Error("api GrantType has not been set correctly", expectedGrantType, s.Auth.api.GrantType)
				}
				if strings.Compare(s.Auth.api.Username, expectedUsername) != 0 {
					nt.Error("api Username has not been set correctly", expectedUsername, s.Auth.api.Username)
				}
				if strings.Compare(s.Auth.api.Password, expectedPassword) != 0 {
					nt.Error("api Password has not been set correctly", expectedPassword, s.Auth.api.Password)
				}
				nt.And("the return value should be \"nil\"")
				if err != nil {
					nt.Error("Should have returned \"nil\" error", "nil", err.Error())
				}
			}
		}
	}
}
