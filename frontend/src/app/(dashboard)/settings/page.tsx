"use client";

import { useEffect, useState, useRef } from 'react';
import { useRouter } from 'next/navigation';
import { useToast } from '@/components/Toast';
import { DoctorResp } from '@/types/doctor';

export default function SettingsPage() {
    const router = useRouter();
    const { success: successToast, error: errorToast } = useToast();
    const [template, setTemplate] = useState('modern');
    const [profile, setProfile] = useState<DoctorResp | null>(null);
    const [loading, setLoading] = useState(true);
    const [saving, setSaving] = useState(false);
    const [uploading, setUploading] = useState(false);
    const fileInputRef = useRef<HTMLInputElement>(null);

    const [doctorForm, setDoctorForm] = useState({
        full_name: '',
        degree: '',
        specialization: '',
        phone: '',
        bio: '',
        signature_url: ''
    });

    useEffect(() => {
        fetchProfile();
    }, []);

    const fetchProfile = async () => {
        const token = localStorage.getItem("docmate_token");
        if (!token) {
            router.push('/login');
            return;
        }

        try {
            const res = await fetch('http://localhost:8081/v1/doctors/profile', {
                headers: { 'Authorization': `Bearer ${token}` }
            });
            const data = await res.json();
            if (res.ok && data.success) {
                const d: DoctorResp = data.data;
                setProfile(d);
                setDoctorForm({
                    full_name: d.full_name || '',
                    degree: d.degree?.join(', ') || '',
                    specialization: d.specialization?.join(', ') || '',
                    phone: d.phone || '',
                    bio: d.bio || '',
                    signature_url: d.signature_url || ''
                });
            }
        } catch (err) {
            console.error("Failed to fetch profile", err);
            errorToast("Failed to load doctor profile");
        } finally {
            setLoading(false);
        }
    };

    const handleProfileSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!profile) return;

        setSaving(true);
        try {
            const token = localStorage.getItem("docmate_token");
            const res = await fetch(`http://localhost:8081/v1/doctors/${profile.id}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify({
                    id: profile.id,
                    full_name: doctorForm.full_name,
                    degree: doctorForm.degree.split(',').map(s => s.trim()),
                    specialization: doctorForm.specialization.split(',').map(s => s.trim()),
                    phone: doctorForm.phone,
                    bio: doctorForm.bio,
                    signature_url: doctorForm.signature_url
                })
            });

            const data = await res.json();
            if (res.ok && data.success) {
                successToast("Doctor profile updated successfully!");
                setProfile(data.data);
            } else {
                errorToast(data.message || "Failed to update profile");
            }
        } catch (err) {
            errorToast("An error occurred");
        } finally {
            setSaving(false);
        }
    };

    const handleSignatureUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.[0];
        if (!file) return;

        setUploading(true);
        const formData = new FormData();
        formData.append('file', file);

        try {
            const token = localStorage.getItem("docmate_token");
            const res = await fetch('http://localhost:8081/v1/upload', {
                method: 'POST',
                headers: { 'Authorization': `Bearer ${token}` },
                body: formData
            });

            const data = await res.json();
            if (res.ok && data.success) {
                setDoctorForm(prev => ({ ...prev, signature_url: data.data.url }));
                successToast("Signature uploaded!");
            } else {
                errorToast("Upload failed");
            }
        } catch (err) {
            errorToast("Upload error");
        } finally {
            setUploading(false);
        }
    };

    return (
        <div className="space-y-10">
            {/* Profile Card */}
            <section className="bg-card rounded-3xl border border-border shadow-sm p-8">
                <h2 className="text-xl font-bold text-slate-900 mb-6">Doctor Profile</h2>
                {loading ? (
                    <div className="flex justify-center py-10">
                        <div className="w-8 h-8 border-4 border-primary border-t-transparent rounded-full animate-spin"></div>
                    </div>
                ) : (
                    <form onSubmit={handleProfileSubmit}>
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
                            <div>
                                <label className="block text-xs font-bold text-slate-400 uppercase mb-2">Full Name</label>
                                <input
                                    type="text"
                                    className="w-full px-4 py-2.5 rounded-xl border border-slate-200 focus:ring-2 focus:ring-primary outline-none transition"
                                    value={doctorForm.full_name}
                                    onChange={(e) => setDoctorForm({ ...doctorForm, full_name: e.target.value })}
                                    placeholder="Dr. Faisal Ahmed"
                                />
                            </div>
                            <div>
                                <label className="block text-xs font-bold text-slate-400 uppercase mb-2">Degrees (comma separated)</label>
                                <input
                                    type="text"
                                    className="w-full px-4 py-2.5 rounded-xl border border-slate-200 focus:ring-2 focus:ring-primary outline-none transition"
                                    value={doctorForm.degree}
                                    onChange={(e) => setDoctorForm({ ...doctorForm, degree: e.target.value })}
                                    placeholder="MBBS, FCPS"
                                />
                            </div>
                            <div>
                                <label className="block text-xs font-bold text-slate-400 uppercase mb-2">Specialization</label>
                                <input
                                    type="text"
                                    className="w-full px-4 py-2.5 rounded-xl border border-slate-200 focus:ring-2 focus:ring-primary outline-none transition"
                                    value={doctorForm.specialization}
                                    onChange={(e) => setDoctorForm({ ...doctorForm, specialization: e.target.value })}
                                    placeholder="Medicine Specialist"
                                />
                            </div>
                            <div>
                                <label className="block text-xs font-bold text-slate-400 uppercase mb-2">Phone</label>
                                <input
                                    type="text"
                                    className="w-full px-4 py-2.5 rounded-xl border border-slate-200 focus:ring-2 focus:ring-primary outline-none transition"
                                    value={doctorForm.phone}
                                    onChange={(e) => setDoctorForm({ ...doctorForm, phone: e.target.value })}
                                    placeholder="+880..."
                                />
                            </div>
                            <div className="md:col-span-2">
                                <label className="block text-xs font-bold text-slate-400 uppercase mb-2">Bio</label>
                                <textarea
                                    className="w-full px-4 py-2.5 rounded-xl border border-slate-200 focus:ring-2 focus:ring-primary outline-none transition h-24"
                                    value={doctorForm.bio}
                                    onChange={(e) => setDoctorForm({ ...doctorForm, bio: e.target.value })}
                                    placeholder="Professional bio..."
                                ></textarea>
                            </div>
                            <div className="md:col-span-2">
                                <label className="block text-xs font-bold text-slate-400 uppercase mb-2">Signature</label>
                                <input
                                    type="file"
                                    ref={fileInputRef}
                                    className="hidden"
                                    accept="image/*"
                                    onChange={handleSignatureUpload}
                                />
                                <div
                                    onClick={() => fileInputRef.current?.click()}
                                    className={`border-2 border-dashed rounded-2xl p-6 text-center transition group cursor-pointer ${doctorForm.signature_url ? 'border-primary bg-blue-50/50' : 'border-slate-200 hover:border-primary'}`}
                                >
                                    {uploading ? (
                                        <div className="flex flex-col items-center">
                                            <div className="w-6 h-6 border-2 border-primary border-t-transparent rounded-full animate-spin mb-2"></div>
                                            <span className="text-sm font-bold text-primary italic underline">Uploading signature...</span>
                                        </div>
                                    ) : doctorForm.signature_url ? (
                                        <div className="flex flex-col items-center gap-2">
                                            <img src={doctorForm.signature_url} alt="Signature Preview" className="h-16 object-contain mix-blend-multiply" />
                                            <span className="text-xs font-bold text-primary italic underline">Click to change signature</span>
                                        </div>
                                    ) : (
                                        <div className="flex flex-col items-center">
                                            <span className="text-slate-400 group-hover:text-primary mb-1">Click to upload signature image</span>
                                            <span className="text-[10px] text-slate-300">PNG, JPG recommended</span>
                                        </div>
                                    )}
                                </div>
                            </div>
                        </div>
                        <div className="mt-8 flex justify-end">
                            <button
                                type="submit"
                                disabled={saving}
                                className="bg-primary text-white px-8 py-2.5 rounded-xl font-bold medical-gradient shadow-lg hover:opacity-90 transition disabled:opacity-50"
                            >
                                {saving ? "Saving..." : "Save Profile"}
                            </button>
                        </div>
                    </form>
                )}
            </section>

            {/* Prescription Template */}
            <section className="bg-card rounded-3xl border border-border shadow-sm p-8">
                <div className="flex justify-between items-center mb-6">
                    <h2 className="text-xl font-bold text-slate-900">Prescription Template</h2>
                    <button
                        onClick={() => router.push('/settings/prescription')}
                        className="text-sm font-bold text-primary hover:underline"
                    >
                        Advanced Settings →
                    </button>
                </div>
                <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                    {[
                        { id: 'classic', name: 'Classic Layout', color: 'bg-slate-100' },
                        { id: 'modern', name: 'Modern Minimal', color: 'bg-blue-100' },
                        { id: 'compact', name: 'Compact Medical', color: 'bg-teal-100' },
                    ].map((t) => (
                        <div
                            key={t.id}
                            onClick={() => setTemplate(t.id)}
                            className={`cursor-pointer rounded-2xl border-2 p-4 transition ${template === t.id ? 'border-primary bg-blue-50/30' : 'border-slate-100'}`}
                        >
                            <div className={`h-32 rounded-xl mb-4 ${t.color} flex items-center justify-center text-2xl`}>📄</div>
                            <div className="text-center">
                                <div className="font-bold text-slate-800 text-sm">{t.name}</div>
                                {template === t.id && <div className="text-[10px] font-bold text-primary uppercase mt-1">Default Template</div>}
                            </div>
                        </div>
                    ))}
                </div>
            </section>
        </div>
    );
}
