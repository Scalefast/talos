package dast

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// Settings holds the settings used for state handling
// We can determine the type of analysis we are executing
// And get data from the user, like Authorization data
type Settings struct {
	ctx                      context.Context                      `json:"-" yaml:"-"`
	cli                      *client.Client                       `json:"-" yaml:"-"`
	targetContainer          *types.Container                     `json:"-" yaml:"-"`
	ImageName                string                               `json:"image_name" yaml:"ImageName"`
	TargetContainerNetworkID string                               `json:"network" yaml:"Network"`
	ZAPContainer             container.ContainerCreateCreatedBody `json:"-" yaml:"-"`
	ZAPConfigFile            string                               `json:"zap_config_file" yaml:"ZAPConfigFileName"`
	OpenApiConfigFile        string                               `json:"open_api_config_file" yaml:"OpenAPIConfigFileName"`
	ScanType                 string                               `json:"scan_type" yaml:"ScanType"`
	Auth                     struct {
		Website struct {
			AccessToken string
		}
		Api struct {
			ClientID     string `json:"client_id" yaml:"ClientID"`
			ClientSecret string `json:"client_secret" yaml:"ClientSecret"`
			GrantType    string `json:"grant_type" yaml:"GrantType"`
			Username     string `json:"username" yaml:"Username"`
			Password     string `json:"password" yaml:"Password"`
		}
	}
}

type DAST struct {
	Env  Env    `yaml:"env"`
	Jobs []Jobs `yaml:"jobs"`
}
type PollAdditionalHeaders struct {
	Header string `yaml:"header"`
	Value  string `yaml:"value"`
}
type Verification struct {
	Method                string                  `yaml:"method"`
	LoggedInRegex         string                  `yaml:"loggedInRegex"`
	LoggedOutRegex        string                  `yaml:"loggedOutRegex"`
	PollFrequency         int                     `yaml:"pollFrequency"`
	PollUnits             string                  `yaml:"pollUnits"`
	PollURL               string                  `yaml:"pollUrl"`
	PollPostData          string                  `yaml:"pollPostData"`
	PollAdditionalHeaders []PollAdditionalHeaders `yaml:"pollAdditionalHeaders"`
}
type Authentication struct {
	Method       string       `yaml:"method"`
	Parameters   Parameters   `yaml:"parameters"`
	Verification Verification `yaml:"verification"`
}
type SessionManagement struct {
	Method       string `yaml:"method"`
	Script       string `yaml:"script"`
	ScriptEngine string `yaml:"scriptEngine"`
}
type Users struct {
	Name     string `yaml:"name"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}
type Contexts struct {
	Name              string            `yaml:"name"`
	Urls              []string          `yaml:"urls"`
	IncludePaths      []string          `yaml:"includePaths"`
	ExcludePaths      []string          `yaml:"excludePaths"`
	Authentication    Authentication    `yaml:"authentication"`
	SessionManagement SessionManagement `yaml:"sessionManagement"`
	Users             []Users           `yaml:"users"`
}
type Env struct {
	Contexts   []Contexts  `yaml:"contexts"`
	Vars       interface{} `yaml:"vars"`
	Parameters Parameters  `yaml:"parameters"`
}
type AlertFilters struct {
	RuleID         int    `yaml:"ruleId"`
	NewRisk        string `yaml:"newRisk"`
	Context        string `yaml:"context"`
	URL            string `yaml:"url"`
	URLRegex       int    `yaml:"urlRegex"`
	Parameter      string `yaml:"parameter"`
	ParameterRegex int    `yaml:"parameterRegex"`
	Attack         string `yaml:"attack"`
	AttackRegex    int    `yaml:"attackRegex"`
	Evidence       string `yaml:"evidence"`
	EvidenceRegex  int    `yaml:"evidenceRegex"`
}
type Requests struct {
	URL          string `yaml:"url"`
	Name         string `yaml:"name"`
	Method       string `yaml:"method"`
	Data         string `yaml:"data"`
	ResponseCode int    `yaml:"responseCode"`
}
type Parameters struct {
	Hostname                     string `yaml:"hostname"`
	Port                         int    `yaml:"port"`
	Realm                        string `yaml:"realm"`
	FailOnError                  bool   `yaml:"failOnError"`
	FailOnWarning                bool   `yaml:"failOnWarning"`
	ProgressToStdout             bool   `yaml:"progressToStdout"`
	DeleteGlobalAlerts           bool   `yaml:"deleteGlobalAlerts"`
	UpdateAddOns                 bool   `yaml:"updateAddOns"`
	MaxAlertsPerRule             int    `yaml:"maxAlertsPerRule"`
	ScanOnlyInScope              bool   `yaml:"scanOnlyInScope"`
	MaxBodySizeInBytesToScan     int    `yaml:"maxBodySizeInBytesToScan"`
	EnableTags                   bool   `yaml:"enableTags"`
	Endpoint                     string `yaml:"endpoint"`
	SchemaURL                    string `yaml:"schemaUrl"`
	SchemaFile                   string `yaml:"schemaFile"`
	MaxQueryDepth                int    `yaml:"maxQueryDepth"`
	LenientMaxQueryDepthEnabled  bool   `yaml:"lenientMaxQueryDepthEnabled"`
	MaxAdditionalQueryDepth      int    `yaml:"maxAdditionalQueryDepth"`
	MaxArgsDepth                 int    `yaml:"maxArgsDepth"`
	OptionalArgsEnabled          bool   `yaml:"optionalArgsEnabled"`
	ArgsType                     string `yaml:"argsType"`
	QuerySplitType               string `yaml:"querySplitType"`
	RequestMethod                string `yaml:"requestMethod"`
	APIFile                      string `yaml:"apiFile"`
	APIURL                       string `yaml:"apiUrl"`
	TargetURL                    string `yaml:"targetUrl"`
	WsdlFile                     string `yaml:"wsdlFile"`
	WsdlURL                      string `yaml:"wsdlUrl"`
	Context                      string `yaml:"context"`
	User                         string `yaml:"user"`
	URL                          string `yaml:"url"`
	MaxDuration                  int    `yaml:"maxDuration"`
	MaxDepth                     int    `yaml:"maxDepth"`
	MaxChildren                  int    `yaml:"maxChildren"`
	AcceptCookies                bool   `yaml:"acceptCookies"`
	HandleODataParametersVisited bool   `yaml:"handleODataParametersVisited"`
	HandleParameters             string `yaml:"handleParameters"`
	MaxParseSizeBytes            int    `yaml:"maxParseSizeBytes"`
	ParseComments                bool   `yaml:"parseComments"`
	ParseGit                     bool   `yaml:"parseGit"`
	ParseRobotsTxt               bool   `yaml:"parseRobotsTxt"`
	ParseSitemapXML              bool   `yaml:"parseSitemapXml"`
	ParseSVNEntries              bool   `yaml:"parseSVNEntries"`
	PostForm                     bool   `yaml:"postForm"`
	ProcessForm                  bool   `yaml:"processForm"`
	RequestWaitTime              int    `yaml:"requestWaitTime"`
	SendRefererHeader            bool   `yaml:"sendRefererHeader"`
	ThreadCount                  int    `yaml:"threadCount"`
	UserAgent                    string `yaml:"userAgent"`
	MaxCrawlDepth                int    `yaml:"maxCrawlDepth"`
	NumberOfBrowsers             int    `yaml:"numberOfBrowsers"`
	BrowserID                    string `yaml:"browserId"`
	ClickDefaultElems            bool   `yaml:"clickDefaultElems"`
	ClickElemsOnce               bool   `yaml:"clickElemsOnce"`
	EventWait                    int    `yaml:"eventWait"`
	MaxCrawlStates               int    `yaml:"maxCrawlStates"`
	RandomInputs                 bool   `yaml:"randomInputs"`
	ReloadWait                   int    `yaml:"reloadWait"`
	Time                         string `yaml:"time"`
	FileName                     string `yaml:"fileName"`
	Policy                       string `yaml:"policy"`
	MaxRuleDurationInMins        int    `yaml:"maxRuleDurationInMins"`
	MaxScanDurationInMins        int    `yaml:"maxScanDurationInMins"`
	AddQueryParam                bool   `yaml:"addQueryParam"`
	DefaultPolicy                string `yaml:"defaultPolicy"`
	DelayInMs                    int    `yaml:"delayInMs"`
	HandleAntiCSRFTokens         bool   `yaml:"handleAntiCSRFTokens"`
	InjectPluginIDInHeader       bool   `yaml:"injectPluginIdInHeader"`
	ScanHeadersAllRequests       bool   `yaml:"scanHeadersAllRequests"`
	ThreadPerHost                int    `yaml:"threadPerHost"`
	Format                       string `yaml:"format"`
	SummaryFile                  string `yaml:"summaryFile"`
	Template                     string `yaml:"template"`
	Theme                        string `yaml:"theme"`
	ReportDir                    string `yaml:"reportDir"`
	ReportFile                   string `yaml:"reportFile"`
	ReportTitle                  string `yaml:"reportTitle"`
	ReportDescription            string `yaml:"reportDescription"`
	DisplayReport                bool   `yaml:"displayReport"`
}
type Tests struct {
	Name      string `yaml:"name"`
	Type      string `yaml:"type"`
	Statistic string `yaml:"statistic"`
	Operator  string `yaml:"operator"`
	Value     int    `yaml:"value"`
	OnFail    string `yaml:"onFail"`
}
type Rules struct {
	ID        int    `yaml:"id"`
	Name      string `yaml:"name"`
	Strength  string `yaml:"strength"`
	Threshold string `yaml:"threshold"`
}
type PolicyDefinition struct {
	DefaultStrength  string  `yaml:"defaultStrength"`
	DefaultThreshold string  `yaml:"defaultThreshold"`
	Rules            []Rules `yaml:"rules"`
}
type Jobs struct {
	Type             string           `yaml:"type"`
	Parameters       Parameters       `yaml:"parameters,omitempty"`
	Install          interface{}      `yaml:"install,omitempty"`
	Uninstall        interface{}      `yaml:"uninstall,omitempty"`
	AlertFilters     []AlertFilters   `yaml:"alertFilters,omitempty"`
	Rules            []Rules          `yaml:"rules,omitempty"`
	Requests         []Requests       `yaml:"requests,omitempty"`
	Tests            []Tests          `yaml:"tests,omitempty"`
	PolicyDefinition PolicyDefinition `yaml:"policyDefinition,omitempty"`
	Risks            []string         `yaml:"risks,omitempty"`
	Confidences      []string         `yaml:"confidences,omitempty"`
	Sections         interface{}      `yaml:"sections,omitempty"`
}
