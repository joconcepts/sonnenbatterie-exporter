package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
)

type Sonnenbatterie struct {
	baseURL url.URL
	token   string
	Client  *http.Client
}

func NewSonnenbatterie(urlString, token string) (*Sonnenbatterie, error) {
	u, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}
	u.Path = filepath.Join(u.Path, "api/v2/")
	return &Sonnenbatterie{
		baseURL: *u,
		Client:  http.DefaultClient,
		token:   token,
	}, nil
}

func (f *Sonnenbatterie) HasToken() bool {
	return f.token != ""
}

func (f *Sonnenbatterie) newRequest(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "sonnenbatterie-exporter")
	if f.HasToken() {
		req.Header.Set("Auth-Token", f.token)
	}
	return req, nil

}

// see https://jlunz.github.io/homeassistant/#/api/getApiV2Status
type Status struct {
	// All AC output of apparent power in VA
	ApparentOutput int `json:"Apparent_output"`
	// 	Backup-buffer in percentage that is set on the system.
	BackupBuffer string `json:"BackupBuffer"`
	// Boolean that indicates the charge status. True if charging
	BatteryCharging bool `json:"BatteryCharging"`
	// Boolean that indicates the discharge status. True if discharging
	BatteryDischarging bool `json:"BatteryDischarging"`
	// House consumption in watts, average over the last 60s
	ConsumptionAvg int `json:"Consumption_Avg"`
	// House consumption in watts, direct measurement
	ConsumptionW int `json:"Consumption_W"`
	// AC frequency in hertz
	Fac float64 `json:"Fac"`
	// Boolean that indicates the energy flow at the installation site. True if battery feeds the consumption
	FlowConsumptionBattery bool `json:"FlowConsumptionBattery"`
	//Boolean that indicates the energy flow at the installation site. True if grid feeds the consumption
	FlowConsumptionGrid bool `json:"FlowConsumptionGrid"`
	// Boolean that indicates the energy flow at the installation site. True if production feeds the consumption
	FlowConsumptionProduction bool `json:"FlowConsumptionProduction"`
	// Boolean that indicates the energy flow at the installation site. True if battery is charging from grid
	FlowGridBattery bool `json:"FlowGridBattery"`
	// Boolean that indicates the energy flow at the installation site. True if production is charging the battery
	FlowProductionBattery bool `json:"FlowProductionBattery"`
	// Boolean that indicates the energy flow at the installation site. True if production feeds into the grid
	FlowProductionGrid bool `json:"FlowProductionGrid"`
	// Grid Feed in negative is consumption and positive is feed in
	GridFeedInW float64 `json:"GridFeedIn_W"`
	// System is installed or not
	IsSystemInstalled int `json:"IsSystemInstalled"`
	// Operating mode that is set on the system:
	// * 1: Manual charging or discharging via API
	// * 2: Automatic Self Consumption. Default
	OperatingMode string `json:"OperatingMode"`
	// AC Power greater than ZERO is discharging Inverter AC Power less than ZERO is charging
	PacTotalW int `json:"Pac_total_W"`
	// PV production in watts
	ProductionW int `json:"Production_W"`
	// Relative state of charge
	Rsoc int `json:"RSOC"`
	// Remaining capacity based on RSOC
	RemainingCapacityWh int `json:"RemainingCapacity_Wh"`
	// Output of apparent power in VA on Phase 1-3
	Sac1 int `json:"Sac1"`
	Sac2 int `json:"Sac2"`
	Sac3 int `json:"Sac3"`
	// String that indicates if the system is connected to the grid (“OnGrid”) or disconnected (“OffGrid”)
	SystemStatus string `json:"SystemStatus"`
	// Local system time
	Timestamp string `json:"Timestamp"`
	// User state of charge
	Usoc int `json:"USOC"`
	// AC voltage in volts
	Uac float64 `json:"Uac"`
	// Battery voltage in volts
	Ubat float64 `json:"Ubat"`
	// Boolean that indicates the discharge status. True if no discharge allowed, based on battery maintenance
	DischargeNotAllowed bool `json:"dischargeNotAllowed"`
	// Boolean that indicates the autostart setting of the generator.
	GeneratorAutostart bool `json:"generator_autostart"`
}

func (f *Sonnenbatterie) GetStatus(ctx context.Context) (*Status, error) {
	u := f.baseURL
	u.Path = filepath.Join(u.Path, "status")
	req, err := f.newRequest(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := f.Client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected http status: %s", resp.Status)
	}
	defer resp.Body.Close()

	var status Status
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return nil, fmt.Errorf("error parsing status: %w", err)
	}
	return &status, nil
}

type PowerMeter struct {
	AL1         float64 `json:"a_l1"`
	AL2         float64 `json:"a_l2"`
	AL3         float64 `json:"a_l3"`
	ATotal      float64 `json:"a_total"`
	Channel     int     `json:"channel"`
	Deviceid    int     `json:"deviceid"`
	Direction   string  `json:"direction"`
	Error       int     `json:"error"`
	Frequency   float64 `json:"frequency"`
	KwhExported float64 `json:"kwh_exported"`
	KwhImported float64 `json:"kwh_imported"`
	VL1L2       float64 `json:"v_l1_l2"`
	VL1N        float64 `json:"v_l1_n"`
	VL2L3       float64 `json:"v_l2_l3"`
	VL2N        float64 `json:"v_l2_n"`
	VL3L1       float64 `json:"v_l3_l1"`
	VL3N        float64 `json:"v_l3_n"`
	VaTotal     float64 `json:"va_total"`
	VarTotal    float64 `json:"var_total"`
	WL1         float64 `json:"w_l1"`
	WL2         float64 `json:"w_l2"`
	WL3         float64 `json:"w_l3"`
	WTotal      float64 `json:"w_total"`
}

// Gets the latest power-meter measurements (Read API)
func (f *Sonnenbatterie) GetPowerMeter(ctx context.Context) (production *PowerMeter, consumption *PowerMeter, err error) {
	u := f.baseURL
	u.Path = filepath.Join(u.Path, "powermeter")
	req, err := f.newRequest(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	resp, err := f.Client.Do(req)
	if err != nil {
		return nil, nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("unexpected http status: %s", resp.Status)
	}
	defer resp.Body.Close()

	var status []PowerMeter
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return nil, nil, fmt.Errorf("error parsing powermeters: %w", err)
	}

	for i := range status {
		if status[i].Direction == "consumption" {
			consumption = &status[i]
		} else if status[i].Direction == "production" {
			production = &status[i]
		}
	}

	if consumption == nil {
		return nil, nil, fmt.Errorf("no consumption powermeter found")
	}
	if production == nil {
		return nil, nil, fmt.Errorf("no production powermeter found")
	}

	return production, consumption, nil
}

type LatestData struct {
	FullChargeCapacity int `json:"FullChargeCapacity"`
	IcStatus           struct {
		SecondsSinceFullCharge int `json:"secondssincefullcharge"`
	} `json:"ic_status"`
}

// Gets latest data for this sonnenBatterie (Read API)
func (f *Sonnenbatterie) GetLatestData(ctx context.Context) (*LatestData, error) {
	u := f.baseURL
	u.Path = filepath.Join(u.Path, "latestdata")
	req, err := f.newRequest(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := f.Client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected http status: %s", resp.Status)
	}
	defer resp.Body.Close()

	var status LatestData
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return nil, fmt.Errorf("error parsing powermeters: %w", err)
	}

	return &status, nil
}

type BatteryModuleData struct {
	CycleCount             float64 `json:"cyclecount"`
	FullChargeCapacity     float64 `json:"fullchargecapacity"`
	MaximumCellTemperature float64 `json:"maximumcelltemperature"`
	MaximumCellVoltage     float64 `json:"maximumcellvoltage"`
	MaximumModuleCurrent   float64 `json:"maximummodulecurrent"`
	MaximumModuleDCVoltage float64 `json:"maximummoduledcvoltage"`
	MinimumCellTemperature float64 `json:"minimumcelltemperature"`
	MinimumCellVoltage     float64 `json:"minimumcellvoltage"`
	MinimumModuleCurrent   float64 `json:"minimummodulecurrent"`
	MinimumModuleDCVoltage float64 `json:"minimummoduledcvoltage"`
	RelativeStateOfCharge  float64 `json:"relativestateofcharge"`
	RemainingCapacity      float64 `json:"remainingcapacity"`
	SystemAlarm            float64 `json:"systemalarm"`
	SystemCurrent          float64 `json:"systemcurrent"`
	SystemVoltage          float64 `json:"systemvoltage"`
	SystemDCVoltage        float64 `json:"systemdcvoltage"`
	SystemStatus           float64 `json:"systemstatus"`
	SystemWarning          float64 `json:"systemwarning"`
}

// Gets battery module data for this sonnenBatterie (Read API)
func (f *Sonnenbatterie) GetBatteryModuleData(ctx context.Context) (*BatteryModuleData, error) {
	u := f.baseURL
	u.Path = filepath.Join(u.Path, "battery")
	req, err := f.newRequest(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := f.Client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected http status: %s", resp.Status)
	}
	defer resp.Body.Close()

	var battery_module BatteryModuleData
	if err := json.NewDecoder(resp.Body).Decode(&battery_module); err != nil {
		return nil, fmt.Errorf("error parsing battery module: %w", err)
	}

	return &battery_module, nil
}
