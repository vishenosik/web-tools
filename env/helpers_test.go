package env

import (
	"testing"
)

type srvs struct {
	AuthenticationService string `env:"AUTH" default:"auth" desc:"Authentication service setting"`
}

type test_config struct {
	// Add your configuration fields here
	// Example:
	DatabaseHost string `env:"DB_HOST" default:"localhost" desc:"database host"`
	Serv         srvs
}

func Test_ConfigInfoTags(t *testing.T) {

	// fmt.Println(string(genEnvConfig(test_config{})))

}
