package engine

type BigFileConfig struct {
	MaxLines int `yaml:"maxLines"`
}

type LongFuncConfig struct {
	MaxLength int `yaml:"maxLength"`
}

type CycloConfig struct {
	IgnoreRegx string `yaml:"ignoreRegx"`
	Over       int    `yaml:"over"`
}

type CopyCheckConfig struct {
	Threshold  int    `yaml:"threshold"`
	IgnoreRegx string `yaml:"ignoreRegx"`
}

type LintersConfig struct {
	Enable []string `yaml:"enable"`
}

type SecurityConfig struct {
	Env []string `yaml:"env"`
}

type LintersSettingsConfig struct {
	Cyclo     CycloConfig     `yaml:"cyclo"`
	BigFile   BigFileConfig   `yaml:"bigFile"`
	LongFunc  LongFuncConfig  `yaml:"longFunc"`
	CopyCheck CopyCheckConfig `yaml:"copyCheck"`
	Security  SecurityConfig  `yaml:"security"`
}

type Config struct {
	IgnoreError     bool                  `yaml:"ignoreError"`
	ReportType      string                `yaml:"reportType"`
	Linters         LintersConfig         `yaml:"linters"`
	LintersSettings LintersSettingsConfig `yaml:"linters-settings"`
}

func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type config Config
	if err := unmarshal((*config)(c)); err != nil {
		return err
	}
	if len(c.ReportType) == 0 {
		c.ReportType = Console
	}
	if c.LintersSettings.BigFile.MaxLines == 0 {
		c.LintersSettings.BigFile.MaxLines = 800
	}
	if c.LintersSettings.LongFunc.MaxLength == 0 {
		c.LintersSettings.LongFunc.MaxLength = 80
	}

	if c.LintersSettings.CopyCheck.Threshold == 0 {
		c.LintersSettings.CopyCheck.Threshold = 50
	}

	if c.LintersSettings.Cyclo.Over == 0 {
		c.LintersSettings.Cyclo.Over = 15
	}

	return nil
}

func DefaultConfig() Config {
	return Config{ReportType: Console,
		LintersSettings: LintersSettingsConfig{
			BigFile:   BigFileConfig{MaxLines: 800},
			LongFunc:  LongFuncConfig{MaxLength: 80},
			CopyCheck: CopyCheckConfig{Threshold: 30},
			Security:  SecurityConfig{Env: []string{}},
			Cyclo:     CycloConfig{Over: 15}},
	}
}
