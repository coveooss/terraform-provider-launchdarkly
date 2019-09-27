package launchdarkly

type JsonEnvironment struct {
	Name      string `json:"name"`
	Key       string `json:"key"`
	Color     string `json:"color"`
	ApiKey    string `json:"apiKey"`
	MobileKey string `json:"mobileKey"`
}

type JsonProject struct {
	Name         string            `json:"name"`
	Key          string            `json:"key"`
	Environments []JsonEnvironment `json:"environments"`
}

type JsonVariations struct {
	Value string `json:"value"`
	Name  string `json:"name"`
}

type JsonCustomProperty struct {
	Name  string   `json:"name"`
	Value []string `json:"value"`
}

type JsonFeatureFlag struct {
	Name             string                        `json:"name"`
	Key              string                        `json:"key"`
	Description      string                        `json:"description"`
	Temporary        bool                          `json:"temporary"`
	IncludeInSnippet bool                          `json:"includeInSnippet"`
	Variations       []JsonVariations              `json:"variations"`
	Tags             []string                      `json:"tags"`
	CustomProperties map[string]JsonCustomProperty `json:"customProperties"`
}
