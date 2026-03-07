export interface DoctorResp {
    id: number;
    user_id: number;
    email: string;
    full_name: string;
    degree: string[];
    specialization: string[];
    phone: string;
    bio: string;
    signature_url: string;
    created_at: string;
    updated_at: string | null;
}

export interface DoctorUpdateReq {
    id: number;
    full_name: string;
    degree: string[];
    specialization: string[];
    phone: string;
    bio: string;
    signature_url: string;
}
