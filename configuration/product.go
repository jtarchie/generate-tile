package configuration

type Property struct {
	Value interface{}
}

type AvailabilityZone struct {
	Name string `validate:"required"`
}

type NetworkProperties struct {
	Network struct {
		Name string `validate:"required"`
	} `yaml:",omitempty" validate:"dive"`
	OtherAvailabilityZones []AvailabilityZone `yaml:"other_availability_zones,omitempty" validate:"dive"`
	ServiceNetwork         struct {
		Name string `validate:"required"`
	} `yaml:",omitempty" validate:"dive"`
	SingletonAvailabilityZone struct {
		Name string `validate:"required"`
	} `yaml:"singleton_availability_zone,omitempty" validate:"dive"`
}

// https://docs.pivotal.io/pivotalcf/2-5/opsman-api/#retrieving-resources-for-a-job
type ResourceConfig struct {
	Instances    int `validate:"gte=0"`
	InstanceType struct {
		ID string `validate:"required"`
	} `yaml:"instance_type,omitempty" validate:"dive"`
	PersistentDisk struct {
		SizeMB string `yaml:"size_mb" validate:"required"`
	} `yaml:"persistent_disk,omitempty" validate:"dive"`

	// AWS, Google, and Azure
	InternetConnected bool     `yaml:"internet_connected"`
	ElbNames          []string `yaml:"elb_names,omitempty"`

	// Vsphere
	NSXLbs []struct {
		EdgeName      string `yaml:"edge_name"`
		PoolName      string `yaml:"pool_name"`
		SecurityGroup string `yaml:"security_group"`
		Port          int    `yaml:"port" validate:"gt=0"`
	} `yaml:"nsx_lbs,omitempty" validate:"dive"`
	NSXSecurityGroups []string `yaml:"nsx_security_groups,omitempty"`

	// OpenStack
	FloatingIPs string `yaml:"floating_ips,omitempty" validate:""`

	AdditionalVMExtensions []string `yaml:"additional_vm_extensions,omitempty"`
	AdditionalNetworks     []struct {
		GUID string `validate:"required"`
	} `yaml:"additional_networks,omitempty" validate:"dive"`
	SwapAsPercentOfMemorySize int `yaml:"swap_as_percent_of_memory_size,omitempty" validate:"gte=0"`

	MaxInFlight interface{} `yaml:"max_in_flight,omitempty"`
}

type Product struct {
	Name              string                    `validate:"required"`
	NetworkProperties NetworkProperties         `yaml:"network-properties,omitempty" validate:"dive"`
	ProductProperties map[string]Property       `yaml:"product-properties,omitempty" validate:"dive"`
	ResourceConfig    map[string]ResourceConfig `yaml:"resource-config,omitempty" validate:"dive"`
}
