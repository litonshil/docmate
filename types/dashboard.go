package types

type DashboardSummaryResp struct {
	TotalPatients      int               `json:"total_patients"`
	TodayVisits        int               `json:"today_visits"`
	TotalPrescriptions int               `json:"total_prescriptions"`
	ActiveMedicines    int               `json:"active_medicines"`
	TotalDoctors       int               `json:"total_doctors,omitempty"`
	TodayCollection    float64           `json:"today_collection"`
	TotalCollection    float64           `json:"total_collection"`
	RecentPatients     []PatientSummary  `json:"recent_patients"`
	TodaySchedule      []ScheduleSummary `json:"today_schedule"`
}

type PatientSummary struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Gender    string `json:"gender"`
	Age       int16  `json:"age"`        // Example: "45",
	LastVisit string `json:"last_visit"` // Example: "2 hours ago", "Yesterday"
}

type ScheduleSummary struct {
	ID             int    `json:"id"`
	PrescriptionID int    `json:"prescription_id,omitempty"`
	PatientID      int    `json:"patient_id"`
	PatientName    string `json:"patient_name"`
	Time           string `json:"time"`   // Example: "09:00 AM"
	Type           string `json:"type"`   // Example: "Checkup", "Follow-up"
	Status         string `json:"status"` // pending, confirmed, etc.
}
