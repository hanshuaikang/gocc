package engine

type BigFileConfig struct {
	MaxLines int `yaml:"maxLines"`
}

type LongFuncConfig struct {
	MaxLength int `yaml:"maxLength"`
}

type CycloConfig struct {
	IgnoreRegx string `yaml:"ignoreRegx"`
}

type Config struct {
	ReportType string         `yaml:"reportType"`
	Cyclo      CycloConfig    `yaml:"cyclo"`
	BigFile    BigFileConfig  `yaml:"bigFile"`
	LongFunc   LongFuncConfig `yaml:"longFunc"`
}

func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type config Config
	if err := unmarshal((*config)(c)); err != nil {
		return err
	}
	if len(c.ReportType) == 0 {
		c.ReportType = Console
	}
	if c.BigFile.MaxLines == 0 {
		c.BigFile.MaxLines = 800
	}
	if c.LongFunc.MaxLength == 0 {
		c.LongFunc.MaxLength = 80
	}

	return nil
}

func DefaultConfig() Config {
	return Config{ReportType: Console, BigFile: BigFileConfig{MaxLines: 800}, LongFunc: LongFuncConfig{MaxLength: 80}}
}
