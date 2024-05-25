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
	Cyclo    CycloConfig    `yaml:"cyclo"`
	BigFile  BigFileConfig  `yaml:"bigFile"`
	LongFunc LongFuncConfig `yaml:"longFunc"`
}

func DefaultConfig() Config {
	return Config{BigFile: BigFileConfig{MaxLines: 800}, LongFunc: LongFuncConfig{MaxLength: 10}}
}
