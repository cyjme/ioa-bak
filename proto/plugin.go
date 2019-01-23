package proto

type Plugin struct {
	Name      string    `json:"name"`
	Tags      []string  `json:"tags"`
	Describe  string    `json:"describe"`
	ConfigTpl ConfigTpl `json:"configTpl"`
}

type ConfigField struct {
	Name      string `json:"name"`
	Desc      string `json:"desc"`
	Required  bool   `json:"required"`
	FieldType string `json:"fieldType"`
}

type ConfigTpl []ConfigField
