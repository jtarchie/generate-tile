package manifest

type ReleaseStemcell struct {
	Name    string `validate:"required"`
	Version string `validate:"required"`
}

type Release struct {
	Name         string            `validate:"required"`
	Version      string            `validate:"required"`
	URL          string            `yaml:",omitempty"`
	SHA1         string            `yaml:"sha1,omitempty"`
	Stemcell     []ReleaseStemcell `yaml:",omitempty" validate:"dive"`
	ExportedFrom []ReleaseStemcell `yaml:"exported_from,omitempty" validate:"dive"`
}

type Stemcell struct {
	Alias   string `validate:"required"`
	OS      string `yaml:",omitempty" validate:"required_without=Name"`
	Version string `validate:"required"`
	Name    string `yaml:",omitempty" validate:"required_without=OS"`
}

type Features struct {
	ConvergeVariables    bool `yaml:"converge_variables"`
	RandomizeAZPlacement bool `yaml:"randomize_az_placement"`
	UseDNSAddresses      bool `yaml:"use_dns_addresses"`
	UseTmpfsConfig       bool `yaml:"use_tmpfs_config"`
	UseShortDNSAddresses bool `yaml:"use_short_dns_addresses"`
}

type Payload struct {
	Name      string     `validate:"required"`
	Features  Features   `yaml:",omitempty"`
	Releases  []Release  `validate:"required,dive"`
	Stemcells []Stemcell `validate:"required,dive"`
}
