package gokku

// GokkuConfig represents the configuration of Gokku
type GokkuConfig struct {
	AppDirectory string
}

//
var Config = NewDefaultGokkuConfig()

func NewDefaultGokkuConfig() GokkuConfig {
	return GokkuConfig{
		AppDirectory: "/var/lib/gokku",
	}
}