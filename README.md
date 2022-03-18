# Talos

Is a project whose aim is to unite every security check in a single,
easy-to-use tool.

Developers will use it by configuring a `config.yaml` file,
and running a single command, like so: `talos run dast`.

## Development set-up

Talos depends on [mage](https://github.com/magefile/mage) to be built.


Talos runs security checks on developed applications or features.

The currently developed modules include DAST and CSA.
DAST is available for API and for basic store

modules are configured through a configuration file. This file
(config.yaml) contains the global configuration for every tool available for
talos.

Below is a sample configuration file in YAML format:

```
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
                                                      # Below fields will be used to generate authentication/oauth post body
      ClientID: "XXXXXXXXXX"                          # Required; Will be used to get authorization from API
      ClientSecret: "XXXXXXXXXX"                      # Required; Will be used to get authorization from API
      GrantType: "password"                           # Required; Will be used to get authorization from API
      Username: "XXXXXXXXXXXXXXXXXXXX"                # Required; Will be used to get authorization from API
      Password: "XXXXXXXXXXXXXXXXXXXX"                # Required; Will be used to get authorization from API
```

CSA stands for Container Security Analysis.
It's goal is to analyze the layers a specific container Image is composed of
and subsequently provide a report of the vulnerabilities that the whole image
contain.

**CSA** Configuration explanation: The Image parameter will be sent to the Docker
client, it will look for a container with that name.  If the container is not
found locally, it will be looked for in the default registries, like Dockerhub.

DAST stands for Dynamic Application Security Test.
It's goal is to run functional tests on the application, to evaluate how it
performs against common security vulnerabilities.
These include SQLi, XML, XSS amongst others.

**DAST** configuration explained: The dast tool spinns-up a container and attaches
it to the network the ImageName container is on.  Then, the ZAP server
starts performing the automated security analysis configured in the
ZAPConfigFileName file.
To obtain a specific ZAP automation application scan configuration file, please
contact the security team at Scalefast.

The ZAP config file is a sensitive file, as it contains infirmation on
authentication and authorization, so it should NOT be commited.


