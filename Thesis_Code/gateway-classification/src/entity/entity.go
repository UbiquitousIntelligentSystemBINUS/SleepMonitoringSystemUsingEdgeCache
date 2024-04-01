package entity

import (
	"time"
)

// Doctor represents the Doctor table
type Doctor struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	FirstName        string    `json:"first_name,omitempty"`
	LastName         string    `json:"last_name,omitempty"`
	Gender           string    `json:"gender,omitempty"`
	Email            string    `json:"email,omitempty"`
	PhoneNumber      string    `json:"phone_number,omitempty"`
	BirthRate        time.Time `json:"birth_rate,omitempty"`
	RegistrationDate time.Time `json:"registration_date,omitempty"`
	Address          string    `json:"address,omitempty"`
}

// Patient represents the Patient table
type Patient struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	FirstName        string    `json:"first_name,omitempty"`
	LastName         string    `json:"last_name,omitempty"`
	Gender           string    `json:"gender,omitempty"`
	Email            string    `json:"email,omitempty"`
	PhoneNumber      string    `json:"phone_number,omitempty"`
	BirthDate        time.Time `json:"birth_date,omitempty"`
	RegistrationDate time.Time `json:"registration_date,omitempty"`
	Address          string    `json:"address,omitempty"`
	SensorToken      string    `json:"sensor_token,omitempty"`
	StartSleepTime   time.Time `json:"start_sleep_time,omitempty"`
	EndSleepTime     time.Time `json:"end_sleep_time,omitempty"`
}

// SleepData represents the Sleep Data table
type SleepData struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	PatientID         uint      `json:"patient_id,omitempty"`
	FirstECGID        uint      `json:"first_ecg_id,omitempty"`
	FirstSleepStageID uint      `json:"first_sleep_stage_id,omitempty"`
	SleepQualityID    uint      `json:"sleep_quality_id,omitempty"`
	FirstInputTime    time.Time `json:"first_input_time,omitempty"`
	LastInputTime     time.Time `json:"last_input_time,omitempty"`
}

// ECG represents the ECG table
type ECG struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ReferenceID uint      `json:"reference_id,omitempty"`
	Value       float64   `json:"value,omitempty"`
	InputTime   time.Time `json:"input_time,omitempty"`
}

// SleepStage represents the Sleep Stage table
type SleepStage struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	ReferenceID uint   `json:"reference_id,omitempty"`
	Value       string `json:"value,omitempty"`
	Method      string `json:"method,omitempty"`
}

// SleepQuality represents the Sleep Quality table
type SleepQuality struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	Value              string    `json:"value,omitempty"`
	InBedDuration      time.Time `json:"in_bed_duration,omitempty"`
	BeginToSleepTime   time.Time `json:"begin_to_sleep_time,omitempty"`
	AwakeFromSleepTime time.Time `json:"awake_from_sleep_time,omitempty"`
	SleepEfficiency    float64   `json:"sleep_efficiency,omitempty"`
	LightSleepDuration time.Time `json:"light_sleep_duration,omitempty"`
	DeepSleepDuration  time.Time `json:"deep_sleep_duration,omitempty"`
	REMDuration        time.Time `json:"rem_duration,omitempty"`
	AwakeDuration      time.Time `json:"awake_duration,omitempty"`
	InputTime          time.Time `json:"input_time,omitempty"`
}
