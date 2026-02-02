package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	INTERVALS_API_KEY = "INTERVALS_API_KEY"
	INTERVALS_ID      = "INTERVALS_ID"
)

type intervals struct{}

type Activity struct {
	ID                     string  `json:"id"`
	StartDateLocal         string  `json:"start_date_local"`
	Type                   string  `json:"type"`
	IcuIgnoreTime          bool    `json:"icu_ignore_time"`
	IcuPmCp                float64 `json:"icu_pm_cp"`
	IcuPmWPrime            float64 `json:"icu_pm_w_prime"`
	IcuPmPMax              float64 `json:"icu_pm_p_max"`
	IcuPmFtp               float64 `json:"icu_pm_ftp"`
	IcuPmFtpSecs           float64 `json:"icu_pm_ftp_secs"`
	IcuPmFtpWatts          float64 `json:"icu_pm_ftp_watts"`
	IcuIgnorePower         bool    `json:"icu_ignore_power"`
	IcuRollingCp           float64 `json:"icu_rolling_cp"`
	IcuRollingWPrime       float64 `json:"icu_rolling_w_prime"`
	IcuRollingPMax         float64 `json:"icu_rolling_p_max"`
	IcuRollingFtp          float64 `json:"icu_rolling_ftp"`
	IcuRollingFtpDelta     float64 `json:"icu_rolling_ftp_delta"`
	IcuTrainingLoad        float64 `json:"icu_training_load"`
	IcuAtl                 float64 `json:"icu_atl"`
	IcuCtl                 float64 `json:"icu_ctl"`
	SsPMax                 float64 `json:"ss_p_max"`
	SsWPrime               float64 `json:"ss_w_prime"`
	SsCp                   float64 `json:"ss_cp"`
	PairedEventID          float64 `json:"paired_event_id"`
	IcuFtp                 float64 `json:"icu_ftp"`
	IcuJoules              float64 `json:"icu_joules"`
	IcuRecordingTime       float64 `json:"icu_recording_time"`
	ElapsedTime            float64 `json:"elapsed_time"`
	IcuWeightedAvgWatts    float64 `json:"icu_weighted_avg_watts"`
	CarbsUsed              float64 `json:"carbs_used"`
	Name                   string  `json:"name"`
	Description            string  `json:"description"`
	StartDate              string  `json:"start_date"`
	Distance               float64 `json:"distance"`
	IcuDistance            float64 `json:"icu_distance"`
	MovingTime             float64 `json:"moving_time"`
	CoastingTime           float64 `json:"coasting_time"`
	TotalElevationGain     float64 `json:"total_elevation_gain"`
	TotalElevationLoss     float64 `json:"total_elevation_loss"`
	Timezone               string  `json:"timezone"`
	Trainer                bool    `json:"trainer"`
	SubType                string  `json:"sub_type"`
	Commute                bool    `json:"commute"`
	Race                   bool    `json:"race"`
	MaxSpeed               float64 `json:"max_speed"`
	AverageSpeed           float64 `json:"average_speed"`
	DeviceWatts            bool    `json:"device_watts"`
	HasHeartrate           bool    `json:"has_heartrate"`
	MaxHeartrate           float64 `json:"max_heartrate"`
	AverageHeartrate       float64 `json:"average_heartrate"`
	AverageCadence         float64 `json:"average_cadence"`
	Calories               float64 `json:"calories"`
	AverageTemp            float64 `json:"average_temp"`
	Mfloat64emp            float64 `json:"min_temp"`
	MaxTemp                float64 `json:"max_temp"`
	AvgLrBalance           float64 `json:"avg_lr_balance"`
	Gap                    float64 `json:"gap"`
	GapModel               string  `json:"gap_model"`
	UseElevationCorrection bool    `json:"use_elevation_correction"`
	Gear                   struct {
		ID       string  `json:"id"`
		Name     string  `json:"name"`
		Distance float64 `json:"distance"`
		Primary  bool    `json:"primary"`
	} `json:"gear"`
	PerceivedExertion      float64   `json:"perceived_exertion"`
	DeviceName             string    `json:"device_name"`
	PowerMeter             string    `json:"power_meter"`
	PowerMeterSerial       string    `json:"power_meter_serial"`
	PowerMeterBattery      string    `json:"power_meter_battery"`
	CrankLength            float64   `json:"crank_length"`
	ExternalID             string    `json:"external_id"`
	FileSportIndex         float64   `json:"file_sport_index"`
	FileType               string    `json:"file_type"`
	IcuAthleteID           string    `json:"icu_athlete_id"`
	Created                time.Time `json:"created"`
	IcuSyncDate            time.Time `json:"icu_sync_date"`
	Analyzed               time.Time `json:"analyzed"`
	IcuWPrime              float64   `json:"icu_w_prime"`
	PMax                   float64   `json:"p_max"`
	ThresholdPace          float64   `json:"threshold_pace"`
	IcuHrZones             []float64 `json:"icu_hr_zones"`
	PaceZones              []float64 `json:"pace_zones"`
	Lthr                   float64   `json:"lthr"`
	IcuRestingHr           float64   `json:"icu_resting_hr"`
	IcuWeight              float64   `json:"icu_weight"`
	IcuPowerZones          []float64 `json:"icu_power_zones"`
	IcuSweetSpotMin        float64   `json:"icu_sweet_spot_min"`
	IcuSweetSpotMax        float64   `json:"icu_sweet_spot_max"`
	IcuPowerSpikeThreshold float64   `json:"icu_power_spike_threshold"`
	Trimp                  float64   `json:"trimp"`
	IcuWarmupTime          float64   `json:"icu_warmup_time"`
	IcuCooldownTime        float64   `json:"icu_cooldown_time"`
	IcuChatID              float64   `json:"icu_chat_id"`
	IcuIgnoreHr            bool      `json:"icu_ignore_hr"`
	IgnoreVelocity         bool      `json:"ignore_velocity"`
	IgnorePace             bool      `json:"ignore_pace"`
	IgnoreParts            []struct {
		StartIndex float64 `json:"start_index"`
		EndIndex   float64 `json:"end_index"`
		Power      bool    `json:"power"`
		Pace       bool    `json:"pace"`
		Hr         bool    `json:"hr"`
	} `json:"ignore_parts"`
	IcuTrainingLoadData float64  `json:"icu_training_load_data"`
	IntervalSummary     []string `json:"float64erval_summary"`
	SkylineChartBytes   string   `json:"skyline_chart_bytes"`
	StreamTypes         []string `json:"stream_types"`
	HasWeather          bool     `json:"has_weather"`
	HasSegments         bool     `json:"has_segments"`
	PowerFieldNames     []string `json:"power_field_names"`
	PowerField          string   `json:"power_field"`
	IcuZoneTimes        []struct {
		ID   string  `json:"id"`
		Secs float64 `json:"secs"`
	} `json:"icu_zone_times"`
	IcuHrZoneTimes  []float64 `json:"icu_hr_zone_times"`
	PaceZoneTimes   []float64 `json:"pace_zone_times"`
	GapZoneTimes    []float64 `json:"gap_zone_times"`
	UseGapZoneTimes bool      `json:"use_gap_zone_times"`
	CustomZones     []struct {
		Code  string `json:"code"`
		Zones []struct {
			ID         string  `json:"id"`
			Start      float64 `json:"start"`
			End        float64 `json:"end"`
			StartValue float64 `json:"start_value"`
			EndValue   float64 `json:"end_value"`
			Secs       float64 `json:"secs"`
		} `json:"zones"`
	} `json:"custom_zones"`
	TizOrder          string  `json:"tiz_order"`
	PolarizationIndex float64 `json:"polarization_index"`
	IcuAchievements   []struct {
		ID        string  `json:"id"`
		Type      string  `json:"type"`
		Message   string  `json:"message"`
		Watts     float64 `json:"watts"`
		Secs      float64 `json:"secs"`
		Value     float64 `json:"value"`
		Distance  float64 `json:"distance"`
		Pace      float64 `json:"pace"`
		Pofloat64 struct {
			StartIndex float64 `json:"start_index"`
			EndIndex   float64 `json:"end_index"`
			Secs       float64 `json:"secs"`
			Value      float64 `json:"value"`
		} `json:"pofloat64"`
	} `json:"icu_achievements"`
	Icufloat64ervalsEdited bool    `json:"icu_float64ervals_edited"`
	Lockfloat64ervals      bool    `json:"lock_float64ervals"`
	IcuLapCount            float64 `json:"icu_lap_count"`
	IcuJoulesAboveFtp      float64 `json:"icu_joules_above_ftp"`
	IcuMaxWbalDepletion    float64 `json:"icu_max_wbal_depletion"`
	IcuHrr                 struct {
		StartIndex   float64 `json:"start_index"`
		EndIndex     float64 `json:"end_index"`
		StartTime    float64 `json:"start_time"`
		EndTime      float64 `json:"end_time"`
		StartBpm     float64 `json:"start_bpm"`
		EndBpm       float64 `json:"end_bpm"`
		AverageWatts float64 `json:"average_watts"`
		Hrr          float64 `json:"hrr"`
	} `json:"icu_hrr"`
	IcuSyncError       string   `json:"icu_sync_error"`
	IcuColor           string   `json:"icu_color"`
	IcuPowerHrZ2       float64  `json:"icu_power_hr_z2"`
	IcuPowerHrZ2Mins   float64  `json:"icu_power_hr_z2_mins"`
	IcuCadenceZ2       float64  `json:"icu_cadence_z2"`
	IcuRpe             float64  `json:"icu_rpe"`
	Feel               float64  `json:"feel"`
	KgLifted           float64  `json:"kg_lifted"`
	Decoupling         float64  `json:"decoupling"`
	IcuMedianTimeDelta float64  `json:"icu_median_time_delta"`
	P30SExponent       float64  `json:"p30s_exponent"`
	WorkoutShiftSecs   float64  `json:"workout_shift_secs"`
	StravaID           string   `json:"strava_id"`
	Lengths            float64  `json:"lengths"`
	PoolLength         float64  `json:"pool_length"`
	Compliance         float64  `json:"compliance"`
	CoachTick          float64  `json:"coach_tick"`
	Source             string   `json:"source"`
	OauthClientID      float64  `json:"oauth_client_id"`
	OauthClientName    string   `json:"oauth_client_name"`
	AverageAltitude    float64  `json:"average_altitude"`
	MinAltitude        float64  `json:"min_altitude"`
	MaxAltitude        float64  `json:"max_altitude"`
	PowerLoad          float64  `json:"power_load"`
	HrLoad             float64  `json:"hr_load"`
	PaceLoad           float64  `json:"pace_load"`
	HrLoadType         string   `json:"hr_load_type"`
	PaceLoadType       string   `json:"pace_load_type"`
	Tags               []string `json:"tags"`
	Attachments        []struct {
		ID       string `json:"id"`
		Filename string `json:"filename"`
		Mimetype string `json:"mimetype"`
		URL      string `json:"url"`
	} `json:"attachments"`
	RecordingStops      []float64 `json:"recording_stops"`
	AverageWeatherTemp  float64   `json:"average_weather_temp"`
	MinWeatherTemp      float64   `json:"min_weather_temp"`
	MaxWeatherTemp      float64   `json:"max_weather_temp"`
	AverageFeelsLike    float64   `json:"average_feels_like"`
	MinFeelsLike        float64   `json:"min_feels_like"`
	MaxFeelsLike        float64   `json:"max_feels_like"`
	AverageWindSpeed    float64   `json:"average_wind_speed"`
	AverageWindGust     float64   `json:"average_wind_gust"`
	PrevailingWindDeg   float64   `json:"prevailing_wind_deg"`
	HeadwindPercent     float64   `json:"headwind_percent"`
	TailwindPercent     float64   `json:"tailwind_percent"`
	AverageClouds       float64   `json:"average_clouds"`
	MaxRain             float64   `json:"max_rain"`
	MaxSnow             float64   `json:"max_snow"`
	CarbsIngested       float64   `json:"carbs_ingested"`
	RouteID             float64   `json:"route_id"`
	Pace                float64   `json:"pace"`
	AthleteMaxHr        float64   `json:"athlete_max_hr"`
	Group               string    `json:"group"`
	Icufloat64ensity    float64   `json:"icu_float64ensity"`
	IcuEfficiencyFactor float64   `json:"icu_efficiency_factor"`
	IcuPowerHr          float64   `json:"icu_power_hr"`
	SessionRpe          float64   `json:"session_rpe"`
	AverageStride       float64   `json:"average_stride"`
	IcuAverageWatts     float64   `json:"icu_average_watts"`
	IcuVariabilityIndex float64   `json:"icu_variability_index"`
	StrainScore         float64   `json:"strain_score"`
}

type Fitness struct {
	ID        string  `json:"id"`
	Ctl       float64 `json:"ctl"`
	Atl       float64 `json:"atl"`
	RampRate  float64 `json:"rampRate"`
	CtlLoad   float64 `json:"ctlLoad"`
	AtlLoad   float64 `json:"atlLoad"`
	SportInfo []struct {
		Type   string  `json:"type"`
		Eftp   float64 `json:"eftp"`
		WPrime float64 `json:"wPrime"`
		PMax   float64 `json:"pMax"`
	} `json:"sportInfo"`
	Updated                 time.Time `json:"updated"`
	Weight                  float64   `json:"weight"`
	RestingHR               float64   `json:"restingHR"`
	Hrv                     float64   `json:"hrv"`
	HrvSDNN                 float64   `json:"hrvSDNN"`
	MenstrualPhase          string    `json:"menstrualPhase"`
	MenstrualPhasePredicted string    `json:"menstrualPhasePredicted"`
	KcalConsumed            float64   `json:"kcalConsumed"`
	SleepSecs               float64   `json:"sleepSecs"`
	SleepScore              float64   `json:"sleepScore"`
	SleepQuality            float64   `json:"sleepQuality"`
	AvgSleepingHR           float64   `json:"avgSleepingHR"`
	Soreness                float64   `json:"soreness"`
	Fatigue                 float64   `json:"fatigue"`
	Stress                  float64   `json:"stress"`
	Mood                    float64   `json:"mood"`
	Motivation              float64   `json:"motivation"`
	Injury                  float64   `json:"injury"`
	SpO2                    float64   `json:"spO2"`
	Systolic                float64   `json:"systolic"`
	Diastolic               float64   `json:"diastolic"`
	Hydration               float64   `json:"hydration"`
	HydrationVolume         float64   `json:"hydrationVolume"`
	Readiness               float64   `json:"readiness"`
	BaevskySI               float64   `json:"baevskySI"`
	BloodGlucose            float64   `json:"bloodGlucose"`
	Lactate                 float64   `json:"lactate"`
	BodyFat                 float64   `json:"bodyFat"`
	Abdomen                 float64   `json:"abdomen"`
	Vo2Max                  float64   `json:"vo2max"`
	Comments                string    `json:"comments"`
	Steps                   float64   `json:"steps"`
	Respiration             float64   `json:"respiration"`
	Locked                  bool      `json:"locked"`
	TempWeight              bool      `json:"tempWeight"`
	TempRestingHR           bool      `json:"tempRestingHR"`
}

type DisplayData struct {
	// Fitness
	ID       string  `json:"id"`
	Ctl      float64 `json:"ctl"`
	Atl      float64 `json:"atl"`
	RampRate float64 `json:"rampRate"`

	// Activities
	Activities []MinimumActivity `json:"activities"`
}

type MinimumActivity struct {
	ID                  string  `json:"id"`
	Distance            float64 `json:"distance"`
	MovingTime          float64 `json:"moving_time"`
	StartDateLocal      string  `json:"start_date_local"`
	IcuAverageWatts     float64 `json:"icu_average_watts"`
	IcuWeightedAvgWatts float64 `json:"icu_weighted_avg_watts"`
	AverageHeartrate    float64 `json:"average_heartrate"`
	AvgLrBalance        float64 `json:"avg_lr_balance"`
	IcuRollingFtp       float64 `json:"icu_rolling_ftp"`
	Calories            float64 `json:"calories"`
}

func (i *intervals) makeRequest(url string, target any) error {
	apiKey, err := GetSecret(INTERVALS_API_KEY)
	if err != nil {
		return fmt.Errorf("failed to get API key: %w", err)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth("API_KEY", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}

func (i *intervals) GetActivities() ([]Activity, error) {
	athleteID, err := GetSecret(INTERVALS_ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get athlete ID: %w", err)
	}

	now := time.Now()
	oldest := now.AddDate(0, 0, -14).Format("2006-01-02")

	url := fmt.Sprintf("https://intervals.icu/api/v1/athlete/%s/activities?oldest=%s&limit=2", athleteID, oldest)

	var activities []Activity
	if err := i.makeRequest(url, &activities); err != nil {
		return nil, err
	}

	return activities, nil
}

func (i *intervals) GetFitness(date string) (*Fitness, error) {
	athleteID, err := GetSecret(INTERVALS_ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get athlete ID: %w", err)
	}

	url := fmt.Sprintf("https://intervals.icu/api/v1/athlete/%s/wellness/%s", athleteID, date)

	var fitness Fitness
	if err := i.makeRequest(url, &fitness); err != nil {
		return nil, err
	}

	return &fitness, nil
}

func (i *intervals) GetDisplayData(date string) (*DisplayData, error) {
	activities, err := i.GetActivities()
	if err != nil {
		return nil, err
	}

	fitness, err := i.GetFitness(date)
	if err != nil {
		return nil, err
	}

	displayData := &DisplayData{
		ID:       fitness.ID,
		Ctl:      fitness.Ctl,
		Atl:      fitness.Atl,
		RampRate: fitness.RampRate,
	}

	// Process activities: reverse order (oldest to newest) and format date
	// API returns newest first due to limit=2 on usually sorted chronological list
	count := min(len(activities), 2)

	for j := count - 1; j >= 0; j-- {
		a := activities[j]

		// Parse and format date: Mon 15:04
		formattedDate := a.StartDateLocal
		t, err := time.Parse("2006-01-02T15:04:05", a.StartDateLocal)
		if err == nil {
			formattedDate = t.Format("Mon 15:04")
		} else {
			// Try RFC3339 if the first one fails
			t, err = time.Parse(time.RFC3339, a.StartDateLocal)
			if err == nil {
				formattedDate = t.Format("Mon 15:04")
			}
		}

		displayData.Activities = append(displayData.Activities, MinimumActivity{
			ID:                  a.ID,
			StartDateLocal:      formattedDate,
			IcuAverageWatts:     a.IcuAverageWatts,
			IcuWeightedAvgWatts: a.IcuWeightedAvgWatts,
			AverageHeartrate:    a.AverageHeartrate,
			AvgLrBalance:        a.AvgLrBalance,
			IcuRollingFtp:       a.IcuRollingFtp,
			Calories:            a.Calories,
			Distance:            a.Distance,
			MovingTime:          a.MovingTime,
		})
	}

	return displayData, nil
}

// GetCachedIntervals retrieves cached intervals data from database
func GetCachedIntervals() (*DisplayData, time.Time, error) {
	var (
		ctl            float64
		atl            float64
		rampRate       float64
		fatigue        float64
		stress         float64
		activitiesJSON string
		lastUpdated    string
	)

	row := db.QueryRow(`
		SELECT ctl, atl, ramp_rate, activities_json, last_updated
		FROM intervals_cache
		WHERE id = 1
	`)

	err := row.Scan(&ctl, &atl, &rampRate, &fatigue, &stress, &activitiesJSON, &lastUpdated)
	if err != nil {
		// No cache exists
		return nil, time.Time{}, err
	}

	// Parse last updated time
	updatedAt, err := time.Parse(time.RFC3339, lastUpdated)
	if err != nil {
		return nil, time.Time{}, err
	}

	// Parse activities JSON
	var activities []MinimumActivity
	if err := json.Unmarshal([]byte(activitiesJSON), &activities); err != nil {
		return nil, time.Time{}, err
	}

	displayData := &DisplayData{
		Ctl:        ctl,
		Atl:        atl,
		RampRate:   rampRate,
		Activities: activities,
	}

	return displayData, updatedAt, nil
}

// SaveIntervalsCache stores intervals data in database
func SaveIntervalsCache(data *DisplayData) error {
	// Serialize activities to JSON
	activitiesJSON, err := json.Marshal(data.Activities)
	if err != nil {
		return err
	}

	now := time.Now().Format(time.RFC3339)

	// Use INSERT OR REPLACE to ensure only one row exists
	_, err = db.Exec(`
		INSERT OR REPLACE INTO intervals_cache (id, ctl, atl, ramp_rate, activities_json, last_updated)
		VALUES (1, ?, ?, ?, ?, ?, ?, ?)
	`, data.Ctl, data.Atl, data.RampRate, string(activitiesJSON), now)

	if err != nil {
		return err
	}

	log.Printf("Intervals cache updated at %s", now)
	return nil
}

func handleIntervals(w http.ResponseWriter, r *http.Request) {
	const cacheMaxAge = 8 * time.Hour

	// Try to get cached data
	cachedData, lastUpdated, err := GetCachedIntervals()
	if err == nil {
		// Check if cache is still valid (less than 8 hours old)
		age := time.Since(lastUpdated)
		if age < cacheMaxAge {
			log.Printf("Using cached intervals data (age: %v)", age.Round(time.Minute))
			w.Header().Set("X-Cache-Age", age.String())
			json.NewEncoder(w).Encode(cachedData)
			return
		}
		log.Printf("Cache expired (age: %v), fetching fresh data", age.Round(time.Minute))
	} else {
		log.Printf("No cache found, fetching fresh data")
	}

	// Cache is expired or doesn't exist, fetch fresh data
	i := &intervals{}
	now := time.Now()
	date := now.Format("2006-01-02")
	displayData, err := i.GetDisplayData(date)
	if err != nil {
		log.Printf("Failed to fetch intervals data: %v", err)
		// If we have cached data, return it even if expired
		if cachedData != nil {
			log.Printf("Returning stale cache due to API error")
			w.Header().Set("X-Cache-Stale", "true")
			json.NewEncoder(w).Encode(cachedData)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Save to cache
	log.Printf("Saving intervals data to db")
	if err := SaveIntervalsCache(displayData); err != nil {
		log.Printf("Failed to save intervals cache: %v", err)
		// Continue anyway, just log the error
	}

	log.Printf("Returning fresh intervals data")
	w.Header().Set("X-Cache-Fresh", "true")
	json.NewEncoder(w).Encode(displayData)
}
