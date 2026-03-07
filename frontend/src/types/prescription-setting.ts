export interface PrescriptionSettingReq {
    chamber_id: number;
    header_left_bangla: string;
    header_right_english: string;
    footer_info_bangla: string;
    footer_info_english: string;
    template_type: 'standard' | 'modern';
}

export interface PrescriptionSettingResp {
    id: number;
    doctor_id: number;
    chamber_id: number;
    header_left_bangla: string;
    header_right_english: string;
    footer_info_bangla: string;
    footer_info_english: string;
    template_type: 'standard' | 'modern';
    created_at: string;
    updated_at: string;
}
