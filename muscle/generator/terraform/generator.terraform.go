package generator

import (
	"errors"
	"fmt"
	"strings"
)

type GeneratorTerraform struct {
	// contains filtered or unexported fields
	Config map[string]string
}

var OpenstackCheckKeys = []string{
	"user_name",
	"tenant_name",
	"password",
	"auth_url",
	"region",
	"user_domain_name",
	"endpoint_type",
}

var OpenstackOverrideEndpoint = []string{
	"endpoint_override_identity",
	"endpoint_override_compute",
	"endpoint_override_network",
	"endpoint_override_image",
	"endpoint_override_object_store",
	"endpoint_override_placement",
	"endpoint_override_volume",
	"endpoint_override_metering",
	"endpoint_override_orchestration",
	"endpoint_override_baremetal",
	"endpoint_override_dns",
	"endpoint_override_shared_file_system",
}

var AWSCheckKeys = []string{
	"access_key",
	"secret_key",
	"region",
}

var AzureCheckKeys = []string{
	"client_id",
	"client_secret",
	"tenant_id",
	"subscription_id",
}

var GCPCheckKeys = []string{
	"project",
	"credentials",
}

func (g *GeneratorTerraform) CheckConfig() error {
	// CheckConfig
	if g.Config["provider"] == "openstack" {
		for _, key := range OpenstackCheckKeys {
			if _, ok := g.Config[key]; !ok {
				return errors.New("key not found : " + key)
			}
		}
		for _, key := range OpenstackOverrideEndpoint {
			if _, ok := g.Config[key]; !ok {
				fmt.Println("key not found : " + key)
			}
		}
	} else if g.Config["provider"] == "aws" {
		for _, key := range AWSCheckKeys {
			if _, ok := g.Config[key]; !ok {
				return errors.New("key not found : " + key)
			}
		}
	} else if g.Config["provider"] == "azure" {
		for _, key := range AzureCheckKeys {
			if _, ok := g.Config[key]; !ok {
				return errors.New("key not found : " + key)
			}
		}
	} else if g.Config["provider"] == "gcp" {
		for _, key := range GCPCheckKeys {
			if _, ok := g.Config[key]; !ok {
				return errors.New("key not found : " + key)
			}
		}
	} else {
		return errors.New("provider not found. you should check provider key")
	}
	return nil
}

func (g *GeneratorTerraform) Generate() (string, error) {
	if err := g.CheckConfig(); err != nil {
		return "", err
	}

	var tfConfig strings.Builder
	tfConfig.WriteString(g.generateCommonProvider())

	switch g.Config["provider"] {
	case "openstack":
		tfConfig.WriteString(g.generateOpenstackProvider())
	case "aws":
		tfConfig.WriteString(g.generateAWSProvider())
	case "azure":
		tfConfig.WriteString(g.generateAzureProvider())
	case "gcp":
		tfConfig.WriteString(g.generateGCPProvider())
	default:
		return "", errors.New("unsupported provider")
	}

	return tfConfig.String(), nil
}

func (g *GeneratorTerraform) generateOpenstackProvider() string {
	var tfConfig strings.Builder
	tfConfig.WriteString(fmt.Sprintf(`
provider "openstack" {
  user_name         = "%s"
  tenant_name       = "%s"
  password          = "%s"
  auth_url          = "%s"
  region            = "%s"
  user_domain_name  = "%s"
  endpoint_type     = "%s"
`,
		g.Config["user_name"],
		g.Config["tenant_name"],
		g.Config["password"],
		g.Config["auth_url"],
		g.Config["region"],
		g.Config["user_domain_name"],
		g.Config["endpoint_type"],
	))
	for _, key := range OpenstackOverrideEndpoint {
		if value, ok := g.Config[key]; ok {
			tfConfig.WriteString(fmt.Sprintf("  %s = \"%s\"\n", strings.TrimPrefix(key, "endpoint_override_"), value))
		}
	}
	tfConfig.WriteString("}\n")
	return tfConfig.String()
}

func (g *GeneratorTerraform) generateAWSProvider() string {
	var tfConfig strings.Builder

	tfConfig.WriteString(fmt.Sprintf(`
provider "aws" {
  access_key = "%s"
  secret_key = "%s"
  region     = "%s"
}
`,
		g.Config["access_key"],
		g.Config["secret_key"],
		g.Config["region"],
	))
	return tfConfig.String()

}

func (g *GeneratorTerraform) generateAzureProvider() string {
	var tfConfig strings.Builder

	tfConfig.WriteString(fmt.Sprintf(`
provider "azurerm" {
	  client_id       = "%s"
	  client_secret   = "%s"
	  tenant_id       = "%s"
	  subscription_id = "%s"
	}
`,
		g.Config["client_id"],
		g.Config["client_secret"],
		g.Config["tenant_id"],
		g.Config["subscription_id"],
	))
	return tfConfig.String()
}

func (g *GeneratorTerraform) generateGCPProvider() string {
	var tfConfig strings.Builder

	tfConfig.WriteString(fmt.Sprintf(`
provider "google" {
	  project     = "%s"
	  credentials = "%s"
	}
`,
		g.Config["project"],
		g.Config["credentials"],
	))
	return tfConfig.String()
}

func (g *GeneratorTerraform) generateCommonProvider() string {
	var tfConfig strings.Builder
	tfConfig.WriteString(`
terraform {
  required_version = ">= 1.0.0"
  required_providers {
`)

	switch g.Config["provider"] {
	case "openstack":
		tfConfig.WriteString(`
    openstack = {
      source  = "terraform-provider-openstack/openstack"
      version = "~> 1.42.0"
    }
`)
	case "aws":
		tfConfig.WriteString(`
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.0"
    }
`)
	case "azure":
		tfConfig.WriteString(`
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 2.0"
    }
`)
	case "gcp":
		tfConfig.WriteString(`
    google = {
      source  = "hashicorp/google"
      version = "~> 3.0"
    }
`)
	}
	tfConfig.WriteString(`
  }
}
`)
	return tfConfig.String()
}
