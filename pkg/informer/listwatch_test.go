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

package informer

import (
	"errors"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

const (
	listTotal        = "kube_events_exporter_list_total"
	listFailedTotal  = "kube_events_exporter_list_failed_total"
	watchTotal       = "kube_events_exporter_watch_total"
	watchFailedTotal = "kube_events_exporter_watch_failed_total"
)

type successListWatch struct{}

func (successListWatch) List(metav1.ListOptions) (runtime.Object, error) {
	return nil, nil
}
func (successListWatch) Watch(metav1.ListOptions) (watch.Interface, error) {
	return nil, nil
}

type errorListWatch struct{}

func (errorListWatch) List(metav1.ListOptions) (runtime.Object, error) {
	return nil, errors.New("")
}
func (errorListWatch) Watch(metav1.ListOptions) (watch.Interface, error) {
	return nil, errors.New("")
}

func TestInstrumentedListerWatcher(t *testing.T) {
	registry := prometheus.NewRegistry()
	metrics := NewListWatchMetrics(registry)

	successLW := NewInstrumentedListerWatcher(successListWatch{}, metrics)
	errorLW := NewInstrumentedListerWatcher(errorListWatch{}, metrics)

	testCases := []struct {
		desc       string
		lw         cache.ListerWatcher
		metricName string
		count      float64
		f          func()
	}{
		{
			desc:       "ListSuccess",
			lw:         successLW,
			metricName: listTotal,
			count:      1,
			f:          func() { _, _ = successLW.List(metav1.ListOptions{}) },
		},
		{
			desc:       "ListFailed",
			lw:         errorLW,
			metricName: listFailedTotal,
			count:      1,
			f:          func() { _, _ = errorLW.List(metav1.ListOptions{}) },
		},
		{
			desc:       "WatchSuccess",
			lw:         successLW,
			metricName: watchTotal,
			count:      1,
			f:          func() { _, _ = successLW.Watch(metav1.ListOptions{}) },
		},
		{
			desc:       "WatchFailed",
			lw:         errorLW,
			metricName: watchFailedTotal,
			count:      1,
			f:          func() { _, _ = errorLW.Watch(metav1.ListOptions{}) },
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			tc.f()

			mf, err := registry.Gather()
			if err != nil {
				t.Fatal(err)
			}

			for _, family := range mf {
				if *family.Name == tc.metricName {
					for _, metric := range family.Metric {
						if *metric.Counter.Value == tc.count {
							return
						}
					}
				}
			}
			t.Fatal("expected metric not found")
		})
	}
}
