//go:build esp32c3

package main

import (
	"device/esp"
	"runtime/interrupt"
	"time"
)

// deepSleepAvailable tells if the chip has the deep sleep calls.
const deepSleepAvailable = true

const (
	// See RESET_REASON_CORE_DEEP_SLEEP in esp-idf
	// components/soc/esp32c3/include/soc/reset_reasons.h.
	resetCauseDeepSleep = 5

	// See RTC_TIMER_TRIG_EN in esp-idf components/soc/esp32c3/include/soc/rtc.h.
	wakeupTimer = 1 << 3

	// 1024 slow clock cycles take about 7.5 ms to count.
	calibrationCycles = 1024
	xtalHz            = 40_000_000

	// See SOC_CLK_RC_SLOW_FREQ_APPROX in esp-idf
	// components/soc/esp32c3/include/soc/clk_tree_defs.h.
	defaultSlowHz = 136_000
)

// slowHz is the measured frequency of the RTC slow clock.
var slowHz uint64

// wokeFromDeepSleep tells if the chip started again after a deep sleep.
func wokeFromDeepSleep() bool {
	return esp.RTC_CNTL.GetRESET_STATE_RESET_CAUSE_PROCPU() == resetCauseDeepSleep
}

// rtcTime returns the time since power on. The RTC timer keeps counting in
// deep sleep.
func rtcTime() time.Duration {
	if slowHz == 0 {
		slowHz = calibrateSlowClock()
	}
	ticks := rtcTicks()
	sec := ticks / slowHz
	rest := ticks % slowHz
	return time.Duration(sec)*time.Second + time.Duration(rest*uint64(time.Second)/slowHz)
}

// rtcTicks reads the 48 bit RTC timer.
func rtcTicks() uint64 {
	esp.RTC_CNTL.SetTIME_UPDATE(1)
	return uint64(esp.RTC_CNTL.GetTIME_HIGH0_TIMER_VALUE0_HIGH())<<32 | uint64(esp.RTC_CNTL.GetTIME_LOW0())
}

// calibrateSlowClock measures the RTC slow clock against the crystal. See
// rtc_clk_cal_internal in esp-idf components/esp_hw_support/port/esp32c3/rtc_time.c.
func calibrateSlowClock() uint64 {
	t := esp.TIMG0
	t.SetRTCCALICFG_RTC_CALI_CLK_SEL(0)
	t.SetRTCCALICFG_RTC_CALI_START_CYCLING(0)
	t.SetRTCCALICFG_RTC_CALI_MAX(calibrationCycles)
	t.SetRTCCALICFG2_RTC_CALI_TIMEOUT_THRES(calibrationCycles << 10)
	t.SetRTCCALICFG_RTC_CALI_START(0)
	t.SetRTCCALICFG_RTC_CALI_START(1)

	var value uint32
	for {
		if t.GetRTCCALICFG_RTC_CALI_RDY() != 0 {
			value = t.GetRTCCALICFG1_RTC_CALI_VALUE()
			break
		}
		if t.GetRTCCALICFG2_RTC_CALI_TIMEOUT() != 0 {
			break
		}
	}
	t.SetRTCCALICFG_RTC_CALI_START(0)

	if value == 0 {
		return defaultSlowHz
	}
	return uint64(xtalHz) * calibrationCycles / uint64(value)
}

// loadSleepState reads the state that saveSleepState kept. Only ESP-IDF uses
// STORE2, see esp-idf components/esp_rom/include/esp32c3/rom/rtc.h.
func loadSleepState() uint32 {
	return esp.RTC_CNTL.STORE2.Get()
}

// saveSleepState keeps the state in a register that deep sleep does not clear.
func saveSleepState(state uint32) {
	esp.RTC_CNTL.STORE2.Set(state)
}

// deepSleep stops the chip for d, then the chip starts again from reset. See
// rtc_sleep_init in esp-idf components/esp_hw_support/port/esp32c3/rtc_sleep.c.
func deepSleep(d time.Duration) {
	if slowHz == 0 {
		slowHz = calibrateSlowClock()
	}
	interrupt.Disable()

	rtc := esp.RTC_CNTL

	// Reset forces these domains on. See rtc_init in esp-idf
	// components/esp_hw_support/port/esp32c3/rtc_init.c.
	rtc.DIG_PWC.ClearBits(esp.RTC_CNTL_DIG_PWC_DG_WRAP_FORCE_PU |
		esp.RTC_CNTL_DIG_PWC_WIFI_FORCE_PU |
		esp.RTC_CNTL_DIG_PWC_BT_FORCE_PU |
		esp.RTC_CNTL_DIG_PWC_CPU_TOP_FORCE_PU |
		esp.RTC_CNTL_DIG_PWC_DG_PERI_FORCE_PU)
	rtc.DIG_ISO.ClearBits(esp.RTC_CNTL_DIG_ISO_DG_WRAP_FORCE_NOISO |
		esp.RTC_CNTL_DIG_ISO_WIFI_FORCE_NOISO |
		esp.RTC_CNTL_DIG_ISO_BT_FORCE_NOISO |
		esp.RTC_CNTL_DIG_ISO_CPU_TOP_FORCE_NOISO |
		esp.RTC_CNTL_DIG_ISO_DG_PERI_FORCE_NOISO |
		esp.RTC_CNTL_DIG_ISO_WIFI_FORCE_ISO |
		esp.RTC_CNTL_DIG_ISO_BT_FORCE_ISO)
	rtc.DIG_PWC.SetBits(esp.RTC_CNTL_DIG_PWC_DG_WRAP_PD_EN |
		esp.RTC_CNTL_DIG_PWC_WIFI_PD_EN |
		esp.RTC_CNTL_DIG_PWC_BT_PD_EN |
		esp.RTC_CNTL_DIG_PWC_CPU_TOP_PD_EN |
		esp.RTC_CNTL_DIG_PWC_DG_PERI_PD_EN)

	rtc.SetANA_CONF_CKGEN_I2C_PU(0)
	rtc.SetANA_CONF_PLL_I2C_PU(0)
	rtc.SetANA_CONF_RFRX_PBUS_PU(0)
	rtc.SetANA_CONF_TXRF_I2C_PU(0)
	rtc.SetANA_CONF_SAR_I2C_PU(0)
	rtc.SetOPTIONS0_BB_I2C_FORCE_PU(0)
	rtc.SetOPTIONS0_XTL_FORCE_PU(0)
	rtc.SetCLK_CONF_XTAL_GLOBAL_FORCE_NOGATING(0)
	rtc.SetCLK_CONF_CK8M_FORCE_PU(0)
	rtc.SetCLK_CONF_CK8M_FORCE_NOGATING(0)
	rtc.SetSENSOR_CTRL_FORCE_XPD_SAR(0)

	// See rtc_sleep_get_default_config in rtc_sleep.c.
	rtc.SetBIAS_CONF_DBG_ATTEN_DEEP_SLP(15)
	rtc.SetBIAS_CONF_BIAS_SLEEP_DEEP_SLP(1)
	rtc.SetBIAS_CONF_PD_CUR_DEEP_SLP(1)

	ticks := uint64(d/time.Microsecond) * slowHz / 1_000_000
	alarm := rtcTicks() + ticks
	rtc.SetSLP_TIMER0(uint32(alarm))
	rtc.SetSLP_TIMER1_SLP_VAL_HI(uint32(alarm >> 32))
	rtc.SetINT_CLR_RTC_MAIN_TIMER_INT_CLR(1)
	rtc.SetSLP_TIMER1_MAIN_TIMER_ALARM_EN(1)

	rtc.SetWAKEUP_STATE_WAKEUP_ENA(wakeupTimer)
	rtc.SLP_REJECT_CONF.Set(0)
	rtc.SetINT_CLR_RTC_SLP_REJECT_INT_CLR(1)
	rtc.SetINT_CLR_RTC_SLP_WAKEUP_INT_CLR(1)
	rtc.SetSTATE0_SLEEP_EN(1)

	for {
	}
}
