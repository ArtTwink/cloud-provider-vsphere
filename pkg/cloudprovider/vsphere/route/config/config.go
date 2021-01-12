/*
 Copyright 2020 The Kubernetes Authors.

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package config

import (
	"fmt"
	"os"

	klog "k8s.io/klog/v2"
)

// FromEnv initializes the provided configuration object with values
// obtained from environment variables. If an environment variable is set
// for a property that's already initialized, the environment variable's value
// takes precedence.
func (cfg *RouteConfig) FromEnv() error {
	if v := os.Getenv("NSXT_ROUTER_PATH"); v != "" {
		cfg.RouterPath = v
	}
	return nil
}

/*
	TODO:
	When the INI based cloud-config is deprecated, the references to the
	INI based code (ie the call to ReadConfigINI) below should be deleted.
*/

// ReadRouteConfig parses vSphere cloud config file and stores it into VSphereConfig.
// Environment variables are also checked
func ReadRouteConfig(configData []byte) (*Config, error) {
	if len(configData) == 0 {
		return nil, fmt.Errorf("Invalid YAML/INI file")
	}

	cfg, err := ReadConfigYAML(configData)
	if err != nil {
		cfg, err = ReadConfigINI(configData)
		if err != nil {
			return nil, err
		}

		klog.Info("ReadConfig INI succeeded. Route INI-based cloud-config is deprecated and will be removed in 2.0. Please use YAML based cloud-config.")
	} else {
		klog.Info("ReadRouteConfig YAML succeeded")
	}

	// Env Vars should override config file entries if present
	if err := cfg.Route.FromEnv(); err != nil {
		return nil, err
	}
	if err := cfg.NSXT.FromEnv(); err != nil {
		return nil, err
	}

	klog.Info("Route Config initialized")
	return cfg, nil
}
