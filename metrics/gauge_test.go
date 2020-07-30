package metrics_test

import (
	"testing"

	"git.code.oa.com/trpc-go/trpc-go/metrics"
)

func Test_gauge_Set(t *testing.T) {
	metrics.RegisterMetricsSink(&metrics.ConsoleSink{})
	type fields struct {
		name string
	}
	type args struct {
		v float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"gauge-cpu.avgload", fields{"cpu.avgload"}, args{0.75}},
		{"gauge-mem.avgload", fields{"mem.avgload"}, args{0.80}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := metrics.Gauge(tt.fields.name)
			g.Set(tt.args.v)
		})
	}
}
