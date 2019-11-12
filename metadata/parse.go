package metadata

import (
	"regexp"
)

type PropertyInput struct {
	Description            string
	Label                  string
	Placeholder            string
	PropertyInputs         []PropertyInput `yaml:"property_inputs" validate:"dive"`
	Reference              string          `validate:"property-exists"`
	SelectorPropertyInputs []PropertyInput `yaml:"selector_property_inputs" validate:"dive"`
}

type Verifier struct {
	Name       string `validate:"required"`
	Properties map[string]interface{}
}

type FormType struct {
	Description    string
	Label          string `validate:"required"`
	Markdown       string
	Name           string          `validate:"required"`
	PropertyInputs []PropertyInput `yaml:"property_inputs" validate:"required,dive"`
	Verifiers      []Verifier      `validate:"dive"`
}

type NamedManifest struct {
	Name     string `validate:"required"`
	Manifest string `validate:"required"`
}

type OptionTemplate struct {
	Name               string              `validate:"required"`
	SelectValue        string              `yaml:"select_value" validate:"required"`
	NamedManifests     []NamedManifest     `yaml:"named_manifests" validate:"dive"`
	PropertyBlueprints []PropertyBlueprint `yaml:"property_blueprints,omitempty" validate:"dive"`
}

type Option struct {
	Name  string `validate:"required"`
	Label string
}

type PropertyBlueprint struct {
	Configurable       bool
	Constraints        []Constraints       `yaml:",omitempty"`
	Default            interface{}         `yaml:"default,omitempty"`
	FreeOnDeploy       bool                `yaml:"freeze_on_deploy"`
	Name               string              `validate:"required"`
	NamedManifests     []NamedManifest     `yaml:"named_manifests" validate:"dive"`
	Optional           bool                `yaml:",omitempty"`
	OptionTemplates    []OptionTemplate    `yaml:"option_templates" validate:"dive"`
	Options            []Option            `validate:"dive"`
	PropertyBlueprints []PropertyBlueprint `yaml:"property_blueprints,omitempty" validate:"dive"`
	Type               string              `validate:"required,oneof=boolean ca_certificate collection disk_type_dropdown domain dropdown_select email http_url integer ip_address ip_ranges ldap_url multi_select_options network_address network_address_list port rsa_cert_credentials rsa_pkey_credentials salted_credentials secret selector service_network_az_multi_select service_network_az_single_select simple_credentials smtp_authentication stemcell_selector string_list string text uuid vm_type_dropdown wildcard_domain"`
}

type Template struct {
	Consumes string `yaml:",omitempty"`
	Name     string `validate:"required"`
	Provides string `yaml:",omitempty"`
	Release  string `validate:"required"`
	Manifest string
}

type ResourceDefinition struct {
	Configurable bool
	Constraints  Constraints `yaml:",omitempty"`
	Default      interface{}
	Label        string `validate:"required"`
	Name         string `validate:"required"`
	Type         string `validate:"required"`
}

type Constraints struct {
	Max                int  `yaml:",omitempty"`
	Min                int  `yaml:",omitempty"`
	MaxOnlyBeOddOrZero bool `yaml:"may_only_be_odd_or_zero,omitempty"`
	MaxOnlyIncrease    bool `yaml:"may_only_increase,omitempty"`
	Modulo             int  `yaml:",omitempty"`
	PowerOfTwo         bool `yaml:"power_of_two,omitempty"`
	ZeroOrMin          int  `yaml:"zero_or_min,omitempty"`
}

type ZeroIf struct {
	PropertyReference string   `yaml:"property_reference"`
	PropertyValues    []string `yaml:"property_values"`
}

type InstanceDefinition struct {
	Configurable bool
	Constraints  Constraints `yaml:",omitempty"`
	Default      int
	Label        string
	Name         string
	Type         string
	ZeroIf       ZeroIf `yaml:"zero_if,omitempty"`
}

type JobType struct {
	Description         string
	InstanceDefinition  InstanceDefinition `yaml:"instance_definition" validate:"required"`
	Label               string
	Manifest            string               `yaml:",omitempty"`
	MaxInFlight         interface{}          `yaml:"max_in_flight" validate:"required"`
	Name                string               `validate:"required"`
	PropertyBlueprints  []PropertyBlueprint  `yaml:"property_blueprints,omitempty" validate:"dive"`
	ResourceDefinitions []ResourceDefinition `yaml:"resource_definitions" validate:"required,dive"`
	ResourceLabel       string               `yaml:"resource_label" validate:"required"`
	Serial              bool
	SingleAZOnly        bool       `yaml:"single_az_only"`
	Templates           []Template `validate:"required,dive"`
	UseStemcell         string     `yaml:"use_stemcell,omitempty"`
}

type StemcellCriteria struct {
	EnablePatchSecurityUpdates bool   `yaml:"enable_patch_security_updates"`
	OS                         string `validate:"required"`
	Version                    string `validate:"required"`
}

type Release struct {
	File    string `validate:"required"`
	Name    string `validate:"required"`
	Version string `validate:"required"`
	SHA1    string `yaml:"sha1"`
}

type InstallTimeVerifier struct {
	Name       string `validate:"required"`
	Ignorable  bool
	Properties map[string]interface{}
}

type RuntimeConfig struct {
	Name          string `validate:"required"`
	RuntimeConfig string `yaml:"runtime_config" validate:"required"`
}

type Variable struct {
	Name    string                 `validate:"required"`
	Type    string                 `validate:"required"`
	Options map[string]interface{} `yaml:",omitempty"`
}

// documented: https://docs.pivotal.io/tiledev/2-7/tile-errands.html
type Errand struct {
	Name          string `validate:"required"`
	Colocated     bool
	RunDefault    bool     `yaml:"run_default"` //default is true
	Instances     []string `yaml:",omitempty"`
	Label         string   `validate:"required"`
	Description   string   `yaml:",omitempty"`
	ImpactWarning string   `yaml:"impact_warning,omitempty"`
}

type ProductVersion struct {
	Name     string `validate:"required"`
	Version  string `validate:"required"`
	Optional bool   // not documented
}

// documented: https://docs.pivotal.io/tiledev/2-7/property-template-references.html
type Payload struct {
	AdditionalStemcellsCriteria []StemcellCriteria `yaml:"additional_stemcells_criteria" validate:"dive"`
	Description                 string
	FormTypes                   []FormType `yaml:"form_types" validate:"dive"`
	IconImage                   string     `yaml:"icon_image" validate:"required"`
	JobTypes                    []JobType  `yaml:"job_types" validate:"dive"`
	Label                       string
	InstallTimeVerifiers        []InstallTimeVerifier `yaml:"install_time_verifiers" validate:"dive"`
	MetadataVersion             string                `yaml:"metadata_version"`
	MinimumVersionForUpgrade    string                `yaml:"minimum_version_for_upgrade" validate:"required"`
	Name                        string                `validate:"required"`
	OpsmanagerSyslog            bool                  `yaml:"opsmanager_syslog"`
	PivnetFilenameRegex         string                `yaml:"pivnet_filename_regex,omitempty"`
	ProductVersion              string                `yaml:"product_version" validate:"required"`
	PropertyBlueprints          []PropertyBlueprint   `yaml:"property_blueprints" validate:"dive"`
	Rank                        int
	Releases                    []Release        `validate:"required,dive"`
	RequiresProductVersions     []ProductVersion `yaml:"requires_product_versions" validate:"dive"`
	RuntimeConfigs              []RuntimeConfig  `yaml:"runtime_configs" validate:"dive"`
	ServiceBroker               bool             `yaml:"service_broker"`
	StemcellCriteria            StemcellCriteria `yaml:"stemcell_criteria" validate:"required,dive"`
	Variables                   []Variable       `validate:"dive"`
	PostDeployErrands           []Errand         `yaml:"post_deploy_errands,omitempty" validate:"dive"`
	PreDeleteErrands            []Errand         `yaml:"pre_delete_errands,omitempty" validate:"dive"`
}

var metadataFile = regexp.MustCompile(`metadata\/.*\.yml`)
