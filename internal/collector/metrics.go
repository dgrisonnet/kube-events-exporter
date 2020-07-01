/*
Copyright 2020 Red Hat, Inc. All rights reserved.

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

package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rhobs/kube-events-exporter/pkg/informer"
)

type exporterMetrics struct {
	listWatchMetrics *informer.ListWatchMetrics
}

func newExporterMetrics(exporterRegistry *prometheus.Registry) *exporterMetrics {
	return &exporterMetrics{
		listWatchMetrics: informer.NewListWatchMetrics(exporterRegistry),
	}
}
