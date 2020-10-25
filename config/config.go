package config

type App struct {
	Host              string `yaml:"host"`
	Topic             string `yaml:"topic"`
	DescriptorSetPath string `yaml:"descriptor_set_path"`
	Schema            string `yaml:"schema"`
}
