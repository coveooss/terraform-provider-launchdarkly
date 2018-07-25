package launchdarkly

type JsonEnvironment struct {
	Name  string `json:"name"`
	Key   string `json:"key"`
	Color string `json:"color"`
}

type JsonProject struct {
	Name         string            `json:"name"`
	Key          string            `json:"key"`
	Environments []JsonEnvironment `json:"environments"`
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
	Tags             []string                      `json:"tags"`
	CustomProperties map[string]JsonCustomProperty `json:"customProperties"`
}
