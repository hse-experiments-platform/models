package core

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

type ModelConfig struct {
	ModelName            string           `yaml:"model_name"`
	ModelDescription     string           `yaml:"model_description"`
	Problem              string           `yaml:"problem"`
	TrainMetrics         []string         `yaml:"train_metrics"`
	TrainHyperparameters []Hyperparameter `yaml:"train_hyperparameters"`
}

type Hyperparameter struct {
	Name         string `yaml:"name"`
	Description  string `yaml:"description"`
	Type         string `yaml:"type"`
	DefaultValue string `yaml:"default_value"`
}

func (c *ModelConfig) Validate() error {
	if c.ModelName == "" {
		return errors.New("model_name cannot be empty")
	}

	if c.ModelDescription == "" {
		return errors.New("model_description cannot be empty")
	}

	if c.Problem != "regression" && c.Problem != "classification" && c.Problem != "clusterization" {
		return fmt.Errorf("problem must be one of 'regression', 'classification', or 'clusterization', got '%s'", c.Problem)
	}

	if len(c.TrainMetrics) == 0 {
		return errors.New("train_metrics cannot be empty")
	}

	for _, metric := range c.TrainMetrics {
		if metric == "" {
			return errors.New("train_metrics cannot contain empty strings")
		}
	}

	if len(c.TrainHyperparameters) == 0 {
		return errors.New("train_hyperparameters cannot be empty")
	}

	for _, hyperparameter := range c.TrainHyperparameters {
		if hyperparameter.Name == "" {
			return errors.New("hyperparameter name cannot be empty")
		}

		if hyperparameter.Description == "" {
			return errors.New("hyperparameter description cannot be empty")
		}

		if hyperparameter.Type != "int" && hyperparameter.Type != "float" && hyperparameter.Type != "bool" && hyperparameter.Type != "string" {
			return fmt.Errorf("hyperparameter type must be one of 'int', 'float', 'bool', or 'string', got '%s'", hyperparameter.Type)
		}
	}

	return nil
}

func (c *ModelConfig) Hash() (string, error) {
	// Convert the ModelConfig to JSON
	jsonData, err := json.Marshal(c)
	if err != nil {
		return "", err
	}

	// Compute the SHA-256 hash
	hash := sha256.Sum256(jsonData)

	// Convert the hash to a hexadecimal string
	return hex.EncodeToString(hash[:]), nil
}

func ParseConfig(data []byte) (*ModelConfig, error) {
	c := &struct {
		Model *ModelConfig `yaml:"model"`
	}{}
	err := yaml.Unmarshal(data, c)
	if err != nil {
		return nil, err
	}
	return c.Model, nil
}

func ParseConfigFromSubfolders(dir string) ([]*ModelConfig, error) {
	var configs []*ModelConfig

	// Read the directory
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	// Iterate over the files in the directory
	for _, file := range files {
		// If the file is a directory, check for a config.yaml file
		if file.IsDir() {
			configPath := filepath.Join(dir, file.Name(), "config.yaml")
			if _, err = os.Stat(configPath); err == nil {
				// If the config.yaml file exists, read and parse it
				data, err := os.ReadFile(configPath)
				if err != nil {
					log.Error().Err(err).Msg("failed to read config")
					continue
				}

				config, err := ParseConfig(data)
				if err != nil {
					log.Error().Err(err).Msg("failed to parse config")
					continue
				}

				if err = config.Validate(); err != nil {
					log.Error().Err(err).Msg("invalid config")
					continue
				}

				configs = append(configs, config)
			}
		}
	}

	return configs, nil
}
