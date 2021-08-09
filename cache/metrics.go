package cache

var (
/*
	cpuTemp = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_temperature_celsius",
		Help: "Current temperature of the CPU.",
	})
	hdFailures = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "hd_errors_total",
			Help: "Number of hard-disk errors.",
		},
		[]string{"device"},
	)
*/
)

func init() {
	// Metrics have to be registered to be exposed:
	//prometheus.MustRegister(cpuTemp)
	//prometheus.MustRegister(hdFailures)
}
