// License information

/*
DAST is a tool created by OWASP that enables security experts to perform tests
on the security of the application being tested.

The way we are running it in TALOS is using the automation framework approach.
This approach is currently being developed by the OWASP team, but has enough
implemented so that it's use is also possible.

To run the scan, we have to generate a configuration file specifying the ZAP
program will run.

By default, the scans that the tools will be able to run are:

	* FastApplicationScan
	* CompleteApplicationScan
	* APIScan

Below is the full configuration TALOS uses.

CSA:                                                  # Required; To use talos like: "talos run csa"
  Image: "golang:latest"                              # Required; Can be a local image or a image from a container registry
DAST:                                                 # Required; To use talos like: "talos run dast"
  ScanType: "API"                                     # Required; one of Website or API
  Network: "development-network"                      # Optional; (Change if network is other than bridge)
  ImageName: "perf"                                   # Required; Container name of the target app
  ZAPConfigFileDir: "/path/to/<user>/talos/"          # Optional; Location of ZAPConfigFileName (absolute path)
  ZAPConfigFileName: "api-dev-automation.yaml"        # Required; Name of config file
  OpenApiConfigFileDir: "/path/to/<user>/talos/"      # Optional; Location of OpenApiConfigFileName (absolute path)
  OpenApiConfigFileName: "OpenAPISchema.json"         # Required; Name of config file
  Auth:                                               # Required; Configures Auth for Website or API, so both are not required
    Website:                                          # Required; If using talos to perform DAST analysis on store
      AccessToken: "basic XXXXXXX"                    # Required; AccessToken stands for Authorization header, sent in every request
    API:                                              # Required; If usnig talos to perform a API-DAST analysis
doc                                                      # Below fields will be used to generate authentication/oauth post body
      ClientID: "XXXXXXXXXX"                          # Required; Will be used to get authorization from API
      ClientSecret: "XXXXXXXXXX"                      # Required; Will be used to get authorization from API
      GrantType: "password"                           # Required; Will be used to get authorization from API
      Username: "XXXXXXXXXXXXXXXXXXXX"                # Required; Will be used to get authorization from API
      Password: "XXXXXXXXXXXXXXXXXXXX"                # Required; Will be used to get authorization from API
The `ScanType` parameter, can be website, or api.
It decides if Talos is going to do a website web scan, or
a api API scan.

The `Network` parameter, connect the Docker OWASP ZAP container to the
specified network.
If the network setting is not set, the OWASP ZAP Docker container will
use the default bridge network

The `ImageName` parameter ensures the Docker container we need to test
the given site will be running

The `ZAPConfigFileDir` parameter allows the user to specify the exact
dir the config file is located at.
If this field is empty or even absent, talos implies the file is located
where talos is run.

The `ZAPConfigFileName` parameter states the name of the file ZAP uses for the
analysis. This file is obtained through the OWASP ZAP Program, under Automation
Framework. If the `ZAPConfigFileDir` parameter is not set, then it is assumed
that the `ZAPConfigFileName` file is located in the same dir as talos is being
executed.

The `OpenApiConfigFileDir` has the same functionality as the `ZAPConfigFileDir`

The `OpenApiConfigFileName` has the same functionality as the `ZAPConfigFileName`

The `Auth` parameter sets the authentication for DAST on stores or API.
It has two sub-sections, website is for web stores.
For now, this parameter authenticates the user for the initial login, but does not
provide authentication for the store, only to it's access.
Store authentication is something in development.

The `ClientID`, `ClientSecret`, `GrantType`, `Username` and `Password` parameters
are all necesary, because they form the body of the Post request that provides
the ZAP framework with the Authorization: Bearer <token> header sent in every request.


The configuration needed for a DAST scan varies on the scan you want to run.
If you are using TALOS to run a DAST scan on the web store, then the following fields are used:

DAST:                                                 # Required; To use talos like: "talos run dast"
  ScanType: "website"                                 # Required; one of website, api
  Network: "api-development-environment_code-network" # Optional; (Change if network is other than bridge)
  ImageName: "perf"                                   # Required; Container name of the target app
  ZAPConfigFileDir: "/path/to/<user>/talos/"          # Optional; Location of ZAPConfigFileName (absolute path)
  ZAPConfigFileName: "api-automation.yaml"            # Required; Name of config file
  Auth:                                               # Required; Configures Auth for website or api, so both are not required
    website:                                          # Required; If using talos to perform DAST analysis on store
      AccessToken: "Basic XXXXXXXXX"                  # Required; AccessToken stands for Authorization header, sent in every request

If the DAST analysis is meant to be for api API, then the required fields are:

DAST:                                                 # Required; To use talos like: "talos run dast"
  ScanType: "api"                                     # Required; one of website, api
  Network: "api-development-environment_code-network" # Optional; (Change if network is other than bridge)
  ImageName: "perf"                                   # Required; Container name of the target app
  ZAPConfigFileDir: "/path/to/<user>/talos/"     # Optional; Location of ZAPConfigFileName (absolute path)
  ZAPConfigFileName: "api-api-dev-automation.yaml"    # Required; Name of config file
  OpenApiConfigFileDir: "/path/to/<user>/talos/" # Optional; Location of OpenApiConfigFileName (absolute path)
  OpenApiConfigFileName: "OpenAPISchema.json"         # Required; Name of config file
  Auth:                                               # Required; Configures Auth for website or api, so both are not required
    api:                                              # Required; If usnig talos to perform a API-DAST analysis
                                                      # Below fields will be used to generate authentication/oauth post body
      ClientID: "XXXXXXXXXX"                          # Required; Will be used to get authorization from API
      ClientSecret: "XXXXXXXXXXX"                     # Required; Will be used to get authorization from API
      GrantType: "password"                           # Required; Will be used to get authorization from API
      Username: "XXXXXXXXXXXXXXXXXXXX"                # Required; Will be used to get authorization from API
      Password: "XXXXXXXXXXXXXXXXXXX"                 # Required; Will be used to get authorization from API

*/
package dast
