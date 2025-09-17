package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/justinas/alice"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"

	"github.com/joconcepts/sonnenbatterie-exporter/api"
)

const timeout = 15 * time.Second

var log = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).With().
	Timestamp().
	Logger()

type collector struct {
	api *api.Sonnenbatterie

	gridVoltage            *prometheus.Desc
	gridFrequency          *prometheus.Desc
	chargePercent          *prometheus.Desc
	usableChargePercent    *prometheus.Desc
	consumptionPower       *prometheus.Desc
	consumptionEnergy      *prometheus.Desc
	productionPower        *prometheus.Desc
	productionEnergy       *prometheus.Desc
	lastFullyCharged       *prometheus.Desc
	fullChargeCapacity     *prometheus.Desc
	remaningChargeCapacity *prometheus.Desc

	batteryCycleCount             *prometheus.Desc
	batteryMaximumCellTemperature *prometheus.Desc
	batteryMaximumCellVoltage     *prometheus.Desc
	batteryMaximumModuleCurrent   *prometheus.Desc
	batteryMaximumModuleDCVoltage *prometheus.Desc
	batteryMinimumCellTemperature *prometheus.Desc
	batteryMinimumCellVoltage     *prometheus.Desc
	batteryMinimumModuleCurrent   *prometheus.Desc
	batteryMinimumModuleDCVoltage *prometheus.Desc
	batteryRelativeStateOfCharge  *prometheus.Desc
	batteryRemainingCapacity      *prometheus.Desc
	batterySystemAlarm            *prometheus.Desc
	batterySystemCurrent          *prometheus.Desc
	batterySystemVoltage          *prometheus.Desc
	batterySystemDCVoltage        *prometheus.Desc
	batterySystemStatus           *prometheus.Desc
	batterySystemWarning          *prometheus.Desc
}

func newCollector(api *api.Sonnenbatterie) *collector {
	return &collector{
		api: api,
		gridVoltage: prometheus.NewDesc(
			"solar_battery_grid_voltage",
			"Solar battery Grid (AC) voltage",
			[]string{"phase"},
			nil,
		),
		gridFrequency: prometheus.NewDesc(
			"solar_battery_grid_frequency",
			"Solar battery Grid (AC) frequency in Hz",
			nil,
			nil,
		),
		chargePercent: prometheus.NewDesc(
			"solar_battery_charge_percent",
			"Solar battery charge in percent",
			nil,
			nil,
		),
		usableChargePercent: prometheus.NewDesc(
			"solar_battery_usable_charge_percent",
			"Solar battery usable charge in percent",
			nil,
			nil,
		),
		consumptionPower: prometheus.NewDesc(
			"solar_battery_consumption_power",
			"Solar battery consumption power in watts",
			[]string{"phase"},
			nil,
		),
		consumptionEnergy: prometheus.NewDesc(
			"solar_battery_consumption_energy_total",
			"Total consumption measured in kwH",
			nil,
			nil,
		),
		productionPower: prometheus.NewDesc(
			"solar_battery_production_power",
			"Solar battery production power in watts",
			[]string{"phase"},
			nil,
		),
		productionEnergy: prometheus.NewDesc(
			"solar_battery_production_energy_total",
			"Total production measured in kwH",
			nil,
			nil,
		),
		lastFullyCharged: prometheus.NewDesc(
			"solar_battery_last_fully_charged_unix_timestamp",
			"Timestamp of last full charge",
			nil,
			nil,
		),
		fullChargeCapacity: prometheus.NewDesc(
			"solar_battery_full_charge_capacity",
			"Full charge capacity in watt hours",
			nil,
			nil,
		),
		remaningChargeCapacity: prometheus.NewDesc(
			"solar_battery_remaining_charge_capacity",
			"Remaining charge capacity in watt hours",
			nil,
			nil,
		),

		batteryCycleCount: prometheus.NewDesc(
			"solar_battery_cycle_count",
			"Cycle count of battery module",
			nil,
			nil,
		),
		batteryMaximumCellTemperature: prometheus.NewDesc(
			"solar_battery_maximum_cell_temperature",
			"Maximum cell temperature of battery",
			nil,
			nil,
		),
		batteryMaximumCellVoltage: prometheus.NewDesc(
			"solar_battery_maximum_cell_voltage",
			"Maximum cell voltage of battery",
			nil,
			nil,
		),
		batteryMaximumModuleCurrent: prometheus.NewDesc(
			"solar_battery_maximum_module_current",
			"Maximum module current of battery",
			nil,
			nil,
		),
		batteryMaximumModuleDCVoltage: prometheus.NewDesc(
			"solar_battery_maximum_module_dc_voltage",
			"Maximum module DC voltage of battery",
			nil,
			nil,
		),
		batteryMinimumCellTemperature: prometheus.NewDesc(
			"solar_battery_minimum_cell_temperature",
			"Minimum cell temperature of battery",
			nil,
			nil,
		),
		batteryMinimumCellVoltage: prometheus.NewDesc(
			"solar_battery_minimum_cell_voltage",
			"Minimum cell voltage of battery",
			nil,
			nil,
		),
		batteryMinimumModuleCurrent: prometheus.NewDesc(
			"solar_battery_minimum_module_current",
			"Minimum module current of battery",
			nil,
			nil,
		),
		batteryMinimumModuleDCVoltage: prometheus.NewDesc(
			"solar_battery_minimum_module_dc_voltage",
			"Minimum module DC voltage of battery",
			nil,
			nil,
		),
		batteryRelativeStateOfCharge: prometheus.NewDesc(
			"solar_battery_relative_state_of_charge",
			"Relative state of charge of battery",
			nil,
			nil,
		),
		batteryRemainingCapacity: prometheus.NewDesc(
			"solar_battery_remaining_capacity",
			"Remaining capacity of battery",
			nil,
			nil,
		),
		batterySystemAlarm: prometheus.NewDesc(
			"solar_battery_system_alarm",
			"System alarm status of battery",
			nil,
			nil,
		),
		batterySystemCurrent: prometheus.NewDesc(
			"solar_battery_system_current",
			"System current of battery",
			nil,
			nil,
		),
		batterySystemVoltage: prometheus.NewDesc(
			"solar_battery_system_voltage",
			"System voltage of battery",
			nil,
			nil,
		),
		batterySystemDCVoltage: prometheus.NewDesc(
			"solar_battery_system_dc_voltage",
			"System DC voltage of battery",
			nil,
			nil,
		),
		batterySystemStatus: prometheus.NewDesc(
			"solar_battery_system_status",
			"System status of battery",
			nil,
			nil,
		),
		batterySystemWarning: prometheus.NewDesc(
			"solar_battery_system_warning",
			"System warning status of battery",
			nil,
			nil,
		),
	}
}

// Describe implements Collector.
func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.gridVoltage
	ch <- c.gridFrequency
	ch <- c.chargePercent
	ch <- c.usableChargePercent
	ch <- c.consumptionPower
	ch <- c.consumptionEnergy
	ch <- c.productionPower
	ch <- c.productionEnergy
}

func (c *collector) collectStatus(ch chan<- prometheus.Metric) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	status, err := c.api.GetStatus(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to get status")
		return
	}

	ch <- prometheus.MustNewConstMetric(c.gridVoltage, prometheus.GaugeValue, status.Uac, "")
	ch <- prometheus.MustNewConstMetric(c.gridFrequency, prometheus.GaugeValue, status.Fac)
	ch <- prometheus.MustNewConstMetric(c.chargePercent, prometheus.GaugeValue, float64(status.Rsoc))
	ch <- prometheus.MustNewConstMetric(c.usableChargePercent, prometheus.GaugeValue, float64(status.Usoc))
	ch <- prometheus.MustNewConstMetric(c.consumptionPower, prometheus.GaugeValue, float64(status.ConsumptionW), "")
	ch <- prometheus.MustNewConstMetric(c.productionPower, prometheus.GaugeValue, float64(status.ProductionW), "")
	ch <- prometheus.MustNewConstMetric(c.remaningChargeCapacity, prometheus.GaugeValue, float64(status.RemainingCapacityWh))
}

func (c *collector) collectPowerMeter(ch chan<- prometheus.Metric) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	production, consumption, err := c.api.GetPowerMeter(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to get power meter")
		return
	}

	ch <- prometheus.MustNewConstMetric(c.gridVoltage, prometheus.GaugeValue, consumption.VL1N, "L1")
	ch <- prometheus.MustNewConstMetric(c.gridVoltage, prometheus.GaugeValue, consumption.VL2N, "L2")
	ch <- prometheus.MustNewConstMetric(c.gridVoltage, prometheus.GaugeValue, consumption.VL3N, "L3")
	ch <- prometheus.MustNewConstMetric(c.gridVoltage, prometheus.GaugeValue, consumption.VL1L2, "L1-L2")
	ch <- prometheus.MustNewConstMetric(c.gridVoltage, prometheus.GaugeValue, consumption.VL2L3, "L2-L3")
	ch <- prometheus.MustNewConstMetric(c.gridVoltage, prometheus.GaugeValue, consumption.VL3L1, "L3-L1")

	ch <- prometheus.MustNewConstMetric(c.consumptionPower, prometheus.GaugeValue, consumption.WL1, "L1")
	ch <- prometheus.MustNewConstMetric(c.consumptionPower, prometheus.GaugeValue, consumption.WL2, "L2")
	ch <- prometheus.MustNewConstMetric(c.consumptionPower, prometheus.GaugeValue, consumption.WL3, "L3")
	ch <- prometheus.MustNewConstMetric(c.consumptionEnergy, prometheus.CounterValue, consumption.KwhImported)

	ch <- prometheus.MustNewConstMetric(c.productionPower, prometheus.GaugeValue, production.WL1, "L1")
	ch <- prometheus.MustNewConstMetric(c.productionPower, prometheus.GaugeValue, production.WL2, "L2")
	ch <- prometheus.MustNewConstMetric(c.productionPower, prometheus.GaugeValue, production.WL3, "L3")
	ch <- prometheus.MustNewConstMetric(c.productionEnergy, prometheus.CounterValue, production.KwhImported)
}

func (c *collector) collectLatestData(ch chan<- prometheus.Metric) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	latestData, err := c.api.GetLatestData(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to get latest data")
		return
	}

	ch <- prometheus.MustNewConstMetric(c.lastFullyCharged, prometheus.GaugeValue, (float64(time.Now().UnixNano())/1e9)-float64(latestData.IcStatus.SecondsSinceFullCharge))
	ch <- prometheus.MustNewConstMetric(c.fullChargeCapacity, prometheus.GaugeValue, float64(latestData.FullChargeCapacity))
}

func (c *collector) collectBatteryModuleData(ch chan<- prometheus.Metric) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	battery_module, err := c.api.GetBatteryModuleData(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to get status")
		return
	}

	ch <- prometheus.MustNewConstMetric(c.batteryCycleCount, prometheus.GaugeValue, battery_module.CycleCount)
	ch <- prometheus.MustNewConstMetric(c.batteryMaximumCellTemperature, prometheus.GaugeValue, battery_module.MaximumCellTemperature)
	ch <- prometheus.MustNewConstMetric(c.batteryMaximumCellVoltage, prometheus.GaugeValue, battery_module.MaximumCellVoltage)
	ch <- prometheus.MustNewConstMetric(c.batteryMaximumModuleCurrent, prometheus.GaugeValue, battery_module.MaximumModuleCurrent)
	ch <- prometheus.MustNewConstMetric(c.batteryMaximumModuleDCVoltage, prometheus.GaugeValue, battery_module.MaximumModuleDCVoltage)
	ch <- prometheus.MustNewConstMetric(c.batteryMinimumCellTemperature, prometheus.GaugeValue, battery_module.MinimumCellTemperature)
	ch <- prometheus.MustNewConstMetric(c.batteryMinimumCellVoltage, prometheus.GaugeValue, battery_module.MinimumCellVoltage)
	ch <- prometheus.MustNewConstMetric(c.batteryMinimumModuleCurrent, prometheus.GaugeValue, battery_module.MinimumModuleCurrent)
	ch <- prometheus.MustNewConstMetric(c.batteryMinimumModuleDCVoltage, prometheus.GaugeValue, battery_module.MinimumModuleDCVoltage)
	ch <- prometheus.MustNewConstMetric(c.batteryRelativeStateOfCharge, prometheus.GaugeValue, battery_module.RelativeStateOfCharge)
	ch <- prometheus.MustNewConstMetric(c.batteryRemainingCapacity, prometheus.GaugeValue, battery_module.RemainingCapacity)
	ch <- prometheus.MustNewConstMetric(c.batterySystemAlarm, prometheus.GaugeValue, battery_module.SystemAlarm)
	ch <- prometheus.MustNewConstMetric(c.batterySystemCurrent, prometheus.GaugeValue, battery_module.SystemCurrent)
	ch <- prometheus.MustNewConstMetric(c.batterySystemVoltage, prometheus.GaugeValue, battery_module.SystemVoltage)
	ch <- prometheus.MustNewConstMetric(c.batterySystemDCVoltage, prometheus.GaugeValue, battery_module.SystemDCVoltage)
	ch <- prometheus.MustNewConstMetric(c.batterySystemStatus, prometheus.GaugeValue, battery_module.SystemStatus)
	ch <- prometheus.MustNewConstMetric(c.batterySystemWarning, prometheus.GaugeValue, battery_module.SystemWarning)
}

func (c *collector) Collect(ch chan<- prometheus.Metric) {
	c.collectStatus(ch)
	if c.api.HasToken() {
		c.collectPowerMeter(ch)
		c.collectLatestData(ch)
		c.collectBatteryModuleData(ch)
	}
}

func run() error {

	var (
		addr        string
		metricsPath string
		url         string
		token       string
	)
	flag.StringVar(&addr, "listen-address", ":9110", "The address to listen on for HTTP requests.")
	flag.StringVar(&metricsPath, "metrics-path", "/metrics", "The path to mount the metrics endpoints.")
	flag.StringVar(&url, "sonnenbatterie-url", "", "URL for the Sonnenbattery storage battery.")
	flag.StringVar(&token, "sonnenbatterie-token", "", "Token for the Sonnenbattery storage battery API.")
	flag.Parse()

	if url == "" {
		return fmt.Errorf("no sonnenbatterie-url set")
	}
	// Take token from environment if not set
	if envToken := os.Getenv("SONNENBATTERIE_TOKEN"); token == "" && envToken != "" {
		token = envToken
	}

	// create sonnenbatterie collector
	a, err := api.NewSonnenbatterie(url, token)
	if err != nil {
		return err
	}

	coll := newCollector(a)

	reg := prometheus.NewRegistry()
	if err := reg.Register(coll); err != nil {
		return err
	}

	// go module build info.
	if err := reg.Register(collectors.NewBuildInfoCollector()); err != nil {
		return err
	}
	if err := reg.Register(collectors.NewGoCollector()); err != nil {
		return err
	}

	// Install the logger handler with default output on the console
	c := alice.New()
	c = c.Append(hlog.NewHandler(log))

	// Expose the registered metrics via HTTP.
	mux := http.NewServeMux()
	mux.Handle(metricsPath, promhttp.HandlerFor(
		reg,
		promhttp.HandlerOpts{
			// Opt into OpenMetrics to support exemplars.
			EnableOpenMetrics: true,
		},
	))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<html>
			<head><title>Sonnenbatterie Exporter</title></head>
			<body>
			<h1>Sonnenbatterie Exporter</h1>
			<p><a href="` + metricsPath + `">Metrics</a></p>
			</body>
			</html>`))
	})

	c = c.Append(hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		hlog.FromRequest(r).Info().
			Str("method", r.Method).
			Stringer("url", r.URL).
			Int("status", status).
			Int("size", size).
			Dur("duration", duration).
			Msg("")
	}))

	return http.ListenAndServe(addr, c.Then(mux))
}

func main() {
	if err := run(); err != nil {
		log.Fatal().Err(err).Msg("failed")
	}
}
