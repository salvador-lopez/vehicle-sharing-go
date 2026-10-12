//go:build unit

package cmd

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

const (
	configPath    = "/config"
	configName    = "inventory"
	configFile    = "/config/inventory.yaml"
	configYAML    = "service: inventory\n"
	serviceKey    = "service"
	serviceName   = "inventory"
	serviceEnvVar = "CMD_ROOT_TEST_SERVICE"
	serviceEnvKey = "cmd_root_test_service"
)

type rootSuite struct {
	suite.Suite
	fs  afero.Fs
	cfg *Config
}

func TestRootSuite(t *testing.T) {
	suite.Run(t, new(rootSuite))
}

func (s *rootSuite) SetupTest() {
	s.fs = afero.NewMemMapFs()
	viper.SetFs(s.fs)
	s.cfg = NewConfigYaml(configPath, configName)
}

func (s *rootSuite) TearDownTest() {
	viper.Reset()
}

func (s *rootSuite) TestNewConfigYaml() {
	s.Require().Equal(&Config{
		path:     configPath,
		fileType: "yaml",
		name:     configName,
	}, s.cfg)
}

func (s *rootSuite) TestInitConfigReadsYamlFile() {
	s.Require().NoError(s.fs.MkdirAll(configPath, 0755))
	s.Require().NoError(afero.WriteFile(s.fs, configFile, []byte(configYAML), 0644))

	InitConfig(s.cfg)

	s.Require().Equal(serviceName, viper.GetString(serviceKey))
	s.Require().Equal(configFile, viper.ConfigFileUsed())
}

func (s *rootSuite) TestInitConfigReadsEnvironmentWhenConfigFileIsMissing() {
	s.T().Setenv(serviceEnvVar, serviceName)

	InitConfig(s.cfg)

	s.Require().Equal(serviceName, viper.GetString(serviceEnvKey))
	s.Require().Empty(viper.ConfigFileUsed())
}
