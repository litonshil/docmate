"use client";

import { useEffect, useState, useCallback } from "react";
import { useRouter } from "next/navigation";
import { useToast } from "@/components/Toast";
import { PrescriptionSettingReq, PrescriptionSettingResp } from "@/types/prescription-setting";

interface Chamber {
    id: number;
    name: string;
}

export default function PrescriptionSettingsPage() {
    const router = useRouter();
    const { success: successToast, error: errorToast } = useToast();

    const [chambers, setChambers] = useState<Chamber[]>([]);
    const [selectedChamberId, setSelectedChamberId] = useState<number | "">("");
    const [doctorProfile, setDoctorProfile] = useState<any>(null);
    const [loading, setLoading] = useState(true);
    const [saving, setSaving] = useState(false);

    const [formData, setFormData] = useState<PrescriptionSettingReq>({
        chamber_id: 0,
        header_left_bangla: "",
        header_right_english: "",
        footer_info_bangla: "",
        footer_info_english: "",
        template_type: "standard",
    });

    const fetchInitialData = useCallback(async () => {
        const token = localStorage.getItem('docmate_token');
        if (!token) {
            router.push('/login');
            return;
        }

        try {
            // 1. Fetch Doctor Profile
            const profileRes = await fetch('http://localhost:8081/v1/doctors/profile', {
                headers: { 'Authorization': `Bearer ${token}` }
            });
            const profileData = await profileRes.json();
            if (!profileRes.ok || !profileData.success) return;
            const doctor = profileData.data;
            setDoctorProfile(doctor);

            // 2. Fetch Chambers
            const chambersRes = await fetch(`http://localhost:8081/v1/doctors/${doctor.id}/chambers`, {
                headers: { 'Authorization': `Bearer ${token}` }
            });
            const chambersData = await chambersRes.json();
            if (chambersRes.ok && chambersData.success) {
                const chamberList = chambersData.data.records || [];
                setChambers(chamberList);
                if (chamberList.length > 0) {
                    setSelectedChamberId(chamberList[0].id);
                }
            }
        } catch (error) {
            console.error("Error fetching initial data:", error);
            errorToast("Failed to load chambers");
        } finally {
            setLoading(false);
        }
    }, [router, errorToast]);

    const fetchSettings = useCallback(async (chamberId: number) => {
        if (!doctorProfile) return;
        const token = localStorage.getItem('docmate_token');
        try {
            const res = await fetch(`http://localhost:8081/v1/doctors/${doctorProfile.id}/prescription-settings?chamber_id=${chamberId}`, {
                headers: { 'Authorization': `Bearer ${token}` }
            });
            const data = await res.json();
            if (res.ok && data.success && data.data) {
                const s: PrescriptionSettingResp = data.data;
                setFormData({
                    chamber_id: s.chamber_id,
                    header_left_bangla: s.header_left_bangla,
                    header_right_english: s.header_right_english,
                    footer_info_bangla: s.footer_info_bangla,
                    footer_info_english: s.footer_info_english,
                    template_type: s.template_type,
                });
            } else {
                // Reset form for new setting
                setFormData({
                    chamber_id: chamberId,
                    header_left_bangla: "",
                    header_right_english: "",
                    footer_info_bangla: "",
                    footer_info_english: "",
                    template_type: "standard",
                });
            }
        } catch (error) {
            console.error("Error fetching settings:", error);
            errorToast("Failed to fetch prescription settings");
        }
    }, [doctorProfile]);

    useEffect(() => {
        fetchInitialData();
    }, [fetchInitialData]);

    useEffect(() => {
        if (selectedChamberId !== "") {
            fetchSettings(Number(selectedChamberId));
        }
    }, [selectedChamberId, fetchSettings]);

    const handleSave = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!doctorProfile || selectedChamberId === "") return;

        setSaving(true);
        try {
            const token = localStorage.getItem('docmate_token');
            const res = await fetch(`http://localhost:8081/v1/doctors/${doctorProfile.id}/prescription-settings`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify({ ...formData, chamber_id: Number(selectedChamberId) })
            });

            const data = await res.json();
            if (res.ok && data.success) {
                successToast("Settings saved successfully");
            } else {
                errorToast(data.message || "Failed to save settings");
            }
        } catch (error) {
            console.error("Error saving settings:", error);
            errorToast("An error occurred while saving");
        } finally {
            setSaving(false);
        }
    };

    if (loading) {
        return <div className="p-8 flex items-center justify-center min-h-[400px]">
            <div className="w-12 h-12 border-4 border-primary border-t-transparent rounded-full animate-spin"></div>
        </div>;
    }

    return (
        <div className="mx-auto">
            <div className="grid grid-cols-1 gap-8">
                <section className="bg-card rounded-3xl border border-border shadow-sm p-8">
                    <form onSubmit={handleSave} className="space-y-8">
                        <div>
                            <label className="block text-sm font-bold text-slate-700 mb-2">Select Chamber</label>
                            <select
                                value={selectedChamberId}
                                onChange={(e) => setSelectedChamberId(e.target.value === "" ? "" : Number(e.target.value))}
                                className="w-full px-4 py-2.5 rounded-xl border border-slate-200 outline-none focus:ring-2 focus:ring-primary transition"
                            >
                                <option value="" disabled>Choose a chamber...</option>
                                {chambers.map(c => (
                                    <option key={c.id} value={c.id}>{c.name}</option>
                                ))}
                            </select>
                        </div>

                        {selectedChamberId !== "" && (
                            <>
                                <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
                                    <div className="space-y-2">
                                        <label className="text-sm font-bold text-slate-700 flex items-center gap-2">
                                            Header Left (Bangla) <span className="text-[10px] bg-red-50 text-red-500 px-2 py-0.5 rounded-full">বাংলা</span>
                                        </label>
                                        <textarea
                                            value={formData.header_left_bangla}
                                            onChange={(e) => setFormData({ ...formData, header_left_bangla: e.target.value })}
                                            className="w-full px-4 py-2.5 rounded-xl border border-slate-200 outline-none focus:ring-2 focus:ring-primary transition h-32 text-sm"
                                            placeholder="ডাাঃ আ.ফ.ম হেলাল উদ্দীীন..."
                                        ></textarea>
                                    </div>
                                    <div className="space-y-2">
                                        <label className="text-sm font-bold text-slate-700">Header Right (English)</label>
                                        <textarea
                                            value={formData.header_right_english}
                                            onChange={(e) => setFormData({ ...formData, header_right_english: e.target.value })}
                                            className="w-full px-4 py-2.5 rounded-xl border border-slate-200 outline-none focus:ring-2 focus:ring-primary transition h-32 text-sm"
                                            placeholder="DR. A. F. M Helal Uddin..."
                                        ></textarea>
                                    </div>
                                </div>

                                <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
                                    <div className="space-y-2">
                                        <label className="text-sm font-bold text-slate-700 flex items-center gap-2">
                                            Footer Info (Bangla) <span className="text-[10px] bg-red-50 text-red-500 px-2 py-0.5 rounded-full">বাংলা</span>
                                        </label>
                                        <textarea
                                            value={formData.footer_info_bangla}
                                            onChange={(e) => setFormData({ ...formData, footer_info_bangla: e.target.value })}
                                            className="w-full px-4 py-2.5 rounded-xl border border-slate-200 outline-none focus:ring-2 focus:ring-primary transition h-24 text-sm"
                                            placeholder="চেম্বারের বিস্তারিত..."
                                        ></textarea>
                                    </div>
                                    <div className="space-y-2">
                                        <label className="text-sm font-bold text-slate-700">Footer Info (English)</label>
                                        <textarea
                                            value={formData.footer_info_english}
                                            onChange={(e) => setFormData({ ...formData, footer_info_english: e.target.value })}
                                            className="w-full px-4 py-2.5 rounded-xl border border-slate-200 outline-none focus:ring-2 focus:ring-primary transition h-24 text-sm"
                                            placeholder="Visiting hours etc..."
                                        ></textarea>
                                    </div>
                                </div>



                                <div className="flex flex-col md:flex-row justify-between items-center gap-4 pt-8 border-t border-slate-100">
                                    <div className="flex items-center gap-4">
                                        <button
                                            type="button"
                                            onClick={() => {
                                                const previewData = {
                                                    ...formData,
                                                    signature_url: doctorProfile?.signature_url || ""
                                                };
                                                const data = encodeURIComponent(JSON.stringify(previewData));
                                                window.open(`/settings/prescription/preview?data=${data}`, '_blank');
                                            }}
                                            className="flex items-center gap-2 text-primary font-bold hover:bg-blue-50 px-6 py-3 rounded-xl transition border border-primary/20"
                                        >
                                            <span className="text-xl">👁️</span> Preview Prescription PDF
                                        </button>
                                        <p className="text-[10px] text-slate-400 max-w-[150px] leading-tight italic">
                                            Click to see how your headers and footers look on a real prescription.
                                        </p>
                                    </div>

                                    <button
                                        type="submit"
                                        disabled={saving}
                                        className="w-full md:w-auto bg-primary text-white px-10 py-3 rounded-xl font-bold medical-gradient shadow-lg hover:opacity-90 transition disabled:opacity-50 flex items-center justify-center gap-2"
                                    >
                                        {saving ? (
                                            <>
                                                <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
                                                Saving...
                                            </>
                                        ) : (
                                            "Save Template Settings"
                                        )}
                                    </button>
                                </div>
                            </>
                        )}

                        {selectedChamberId === "" && chambers.length === 0 && (
                            <div className="text-center py-10 bg-slate-50 rounded-2xl border-2 border-dashed border-slate-200">
                                <p className="text-slate-500 font-medium">Please add a chamber first to configure prescription settings.</p>
                                <button
                                    type="button"
                                    onClick={() => router.push('/chambers')}
                                    className="mt-4 text-primary font-bold hover:underline"
                                >
                                    Go to Chamber Management
                                </button>
                            </div>
                        )}
                    </form>
                </section>
            </div>
        </div>
    );
}
