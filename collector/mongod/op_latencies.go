// Copyright 2017 Percona LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mongod

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	opLatencyMicrosecondsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: Namespace,
		Name:      "op_latency_sum",
		Help:      "sum microsecond latency measured internally to mongod, excluding networking layer",
	}, []string{"type"})
	opLatencyOpsCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: Namespace,
		Name:      "op_latency_ops_count",
		Help:      "ops counter incremented synchronously with op_latency_sum increases",
	}, []string{"type"})
)

type OpLatencyAndCount struct {
	Latency float64 `bson:"latency"`
	OpCount float64 `bson:"ops"`
}

type OpLatencyStats struct {
	Reads    OpLatencyAndCount `bson:"reads"`
	Writes   OpLatencyAndCount `bson:"writes"`
	Commands OpLatencyAndCount `bson:"commands"`
}

// Export exports the data to prometheus.
func (opLatencies *OpLatencyStats) Export(ch chan<- prometheus.Metric) {
	opLatencyMicrosecondsTotal.WithLabelValues("reads").Set(opLatencies.Reads.Latency)
	opLatencyOpsCounter.WithLabelValues("reads").Set(opLatencies.Reads.OpCount)
	opLatencyMicrosecondsTotal.WithLabelValues("writes").Set(opLatencies.Writes.Latency)
	opLatencyOpsCounter.WithLabelValues("writes").Set(opLatencies.Writes.OpCount)
	opLatencyMicrosecondsTotal.WithLabelValues("commands").Set(opLatencies.Commands.Latency)
	opLatencyOpsCounter.WithLabelValues("commands").Set(opLatencies.Commands.OpCount)

	opLatencyMicrosecondsTotal.Collect(ch)
	opLatencyOpsCounter.Collect(ch)
}

// Describe describes the metrics for prometheus
func (opLatencyGauges *OpLatencyStats) Describe(ch chan<- *prometheus.Desc) {
	opLatencyMicrosecondsTotal.Describe(ch)
	opLatencyOpsCounter.Describe(ch)
}
