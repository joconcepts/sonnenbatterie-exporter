# Sonnenbatterie

Uses the sonnenbatterie v2 API to expose its metrics

## Examples

```
# HELP solar_battery_charge_percent Solar battery charge in percent
# TYPE solar_battery_charge_percent gauge
solar_battery_charge_percent 7
# HELP solar_battery_consumption_energy_total Total consumption measured in kwH
# TYPE solar_battery_consumption_energy_total counter
solar_battery_consumption_energy_total 99.123
# HELP solar_battery_consumption_power Solar battery consumption power in watts
# TYPE solar_battery_consumption_power gauge
solar_battery_consumption_power{phase=""} 629
solar_battery_consumption_power{phase="L1"} 427.5
solar_battery_consumption_power{phase="L2"} 190.8000030517578
solar_battery_consumption_power{phase="L3"} 21
# HELP solar_battery_full_charge_capacity Full charge capacity in watt hours
# TYPE solar_battery_full_charge_capacity gauge
solar_battery_full_charge_capacity 15000
# HELP solar_battery_grid_frequency Solar battery Grid (AC) frequency in Hz
# TYPE solar_battery_grid_frequency gauge
solar_battery_grid_frequency 50.013
# HELP solar_battery_grid_voltage Solar battery Grid (AC) voltage
# TYPE solar_battery_grid_voltage gauge
solar_battery_grid_voltage{phase=""} 236
solar_battery_grid_voltage{phase="L1"} 236
solar_battery_grid_voltage{phase="L1-L2"} 410.5
solar_battery_grid_voltage{phase="L2"} 226
solar_battery_grid_voltage{phase="L2-L3"} 411.3
solar_battery_grid_voltage{phase="L3"} 236
solar_battery_grid_voltage{phase="L3-L1"} 408
# HELP solar_battery_last_fully_charged_unix_timestamp Timestamp of last full charge
# TYPE solar_battery_last_fully_charged_unix_timestamp gauge
solar_battery_last_fully_charged_unix_timestamp 1.735476305243205e+09
# HELP solar_battery_production_energy_total Total production measured in kwH
# TYPE solar_battery_production_energy_total counter
solar_battery_production_energy_total 12.38984375
# HELP solar_battery_production_power Solar battery production power in watts
# TYPE solar_battery_production_power gauge
solar_battery_production_power{phase=""} 0
solar_battery_production_power{phase="L1"} -7.599999904632568
solar_battery_production_power{phase="L2"} 0
solar_battery_production_power{phase="L3"} 0
# HELP solar_battery_remaining_charge_capacity Remaining charge capacity in watt hours
# TYPE solar_battery_remaining_charge_capacity gauge
solar_battery_remaining_charge_capacity 717
# HELP solar_battery_usable_charge_percent Solar battery usable charge in percent
# TYPE solar_battery_usable_charge_percent gauge
solar_battery_usable_charge_percent 0
```
