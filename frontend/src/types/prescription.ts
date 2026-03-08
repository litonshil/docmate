export interface Vitals {
    weight_kg?: number;
    height_cm?: number;
    blood_pressure?: string;
    temperature_f?: number;
    pulse_bpm?: number;
    spo2_percent?: number;
}

export interface Medication {
    medicine_id?: number | null;
    medicine_name: string;
    generic_name?: string;
    form?: string;
    strength?: string;
    dosage: string;
    frequency: string;
    timing?: string;
    duration: string;
    instructions?: string;
    sort_order?: number;
}

export interface PrescriptionReq {
    patient_id: number;
    chamber_id: number;
    vitals: Vitals;
    chief_complaints: string[];
    diagnosis: string[];
    medications: Medication[];
    investigations: string[];
    advice: string;
    status: 'draft' | 'finalized';
    follow_up_date?: string; // ISO date string or null
}

export interface PrescriptionResp {
    id: number;
    doctor_id: number;
    patient_id: number;
    patient_name: string;
    chamber_id: number;
    vitals: Vitals;
    chief_complaints: string[];
    diagnosis: string[];
    medications: Medication[];
    investigations: string[];
    advice: string;
    status: 'draft' | 'finalized';
    follow_up_date?: string;
    created_at: string;
    updated_at: string;
}

export interface PaginatedPrescriptionResp {
    pagination: {
        page: number;
        limit: number;
        total: number;
        last_page: number;
    };
    records: PrescriptionResp[];
}
