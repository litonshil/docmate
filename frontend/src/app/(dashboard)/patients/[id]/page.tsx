"use client";

import { useEffect, useState, use } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useToast } from "@/components/Toast";

interface Patient {
    id: number;
    full_name: string;
    gender: string;
    age: number;
    phone: string;
    email: string;
    blood_group: string;
    allergies: string[];
    medical_history: string;
    created_at: string;
}

export default function PatientDetail({ params: paramsPromise }: { params: Promise<{ id: string }> }) {
    const params = use(paramsPromise);
    const router = useRouter();
    const { error: errorToast } = useToast();
    const [patient, setPatient] = useState<Patient | null>(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchPatient = async () => {
            setLoading(true);
            try {
                const token = localStorage.getItem('docmate_token');
                if (!token) {
                    router.push('/login');
                    return;
                }

                const response = await fetch(`http://localhost:8081/v1/patients/${params.id}`, {
                    headers: {
                        'Authorization': `Bearer ${token}`
                    }
                });

                const data = await response.json();
                if (response.ok && data.success) {
                    setPatient(data.data);
                } else {
                    console.error("Failed to fetch patient:", data.message);
                    errorToast(data.message || "Failed to fetch patient details");
                }
            } catch (error) {
                console.error("Error fetching patient:", error);
                errorToast("An error occurred while loading patient profile");
            } finally {
                setLoading(false);
            }
        };

        if (params.id) {
            fetchPatient();
        }
    }, [params.id, router]);

    if (loading) {
        return (
            <div className="p-8 flex flex-col items-center justify-center min-h-[400px] gap-4">
                <div className="w-12 h-12 border-4 border-primary border-t-transparent rounded-full animate-spin"></div>
                <p className="text-slate-500 font-medium">Loading patient profile...</p>
            </div>
        );
    }

    if (!patient) {
        return (
            <div className="p-8 text-center min-h-[400px] flex flex-col items-center justify-center">
                <h2 className="text-2xl font-bold text-slate-900 mb-2">Patient Not Found</h2>
                <p className="text-slate-500 mb-6">The patient record you are looking for does not exist or you don't have access.</p>
                <Link href="/patients" className="text-primary font-bold hover:underline">Back to Directory</Link>
            </div>
        );
    }

    const registrationDate = new Date(patient.created_at).toLocaleDateString('en-US', {
        month: 'short',
        day: 'numeric',
        year: 'numeric'
    });

    return (
        <div className="p-8">
            <div className="flex justify-between items-start mb-10">
                <div className="flex items-center gap-6">
                    <div className="w-20 h-20 rounded-3xl bg-primary medical-gradient flex items-center justify-center text-white text-3xl font-bold shadow-xl">
                        {patient.full_name.charAt(0)}
                    </div>
                    <div>
                        <div className="flex items-center gap-3">
                            <h1 className="text-3xl font-extrabold text-slate-900 tracking-tight">{patient.full_name}</h1>
                            <span className="px-3 py-1 rounded-full bg-blue-50 text-blue-600 text-xs font-bold uppercase tracking-wider">Patient</span>
                        </div>
                        <p className="text-slate-500 font-medium mt-1">ID: #P-{patient.id} • Registered {registrationDate}</p>
                    </div>
                </div>
                <div className="flex gap-4">
                    <button className="px-6 py-2.5 rounded-xl border border-slate-200 font-bold text-slate-600 hover:bg-slate-50 transition">Edit Info</button>
                    <button className="px-6 py-2.5 rounded-xl bg-primary text-white font-bold medical-gradient shadow-lg">+ New Prescription</button>
                </div>
            </div>

            <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
                {/* Left Column: Basic Info & Health Summary */}
                <div className="space-y-8">
                    <section className="bg-card rounded-3xl border border-border shadow-sm p-8">
                        <h3 className="text-sm font-bold text-slate-400 uppercase tracking-widest mb-6">Contact & Personal</h3>
                        <div className="space-y-5">
                            <div className="flex justify-between">
                                <span className="text-slate-500 text-sm">Age/Gender</span>
                                <span className="text-slate-900 font-bold">{patient.age} / {patient.gender}</span>
                            </div>
                            <div className="flex justify-between">
                                <span className="text-slate-500 text-sm">Phone</span>
                                <span className="text-slate-900 font-bold">{patient.phone || 'N/A'}</span>
                            </div>
                            <div className="flex justify-between">
                                <span className="text-slate-500 text-sm">Blood Group</span>
                                <span className="text-red-600 font-extrabold">{patient.blood_group || 'N/A'}</span>
                            </div>
                            <div className="flex justify-between">
                                <span className="text-slate-500 text-sm">Email</span>
                                <span className="text-slate-900 font-bold text-xs">{patient.email || 'N/A'}</span>
                            </div>
                        </div>
                    </section>

                    <section className="bg-card rounded-3xl border border-border shadow-sm p-8">
                        <h3 className="text-sm font-bold text-slate-400 uppercase tracking-widest mb-6">Medical Summary</h3>
                        <div className="space-y-6">
                            <div>
                                <div className="text-xs font-bold text-slate-400 uppercase mb-2">Allergies</div>
                                <div className="flex flex-wrap gap-2">
                                    {patient.allergies && patient.allergies.length > 0 ? (
                                        patient.allergies.map((allergy, idx) => (
                                            <span key={idx} className="px-3 py-1 bg-red-50 text-red-600 rounded-lg text-xs font-bold">{allergy}</span>
                                        ))
                                    ) : (
                                        <span className="text-slate-400 text-xs italic">No known allergies</span>
                                    )}
                                </div>
                            </div>
                            <div>
                                <div className="text-xs font-bold text-slate-400 uppercase mb-2">History</div>
                                <p className="text-sm text-slate-600 leading-relaxed">
                                    {patient.medical_history || 'No medical history recorded.'}
                                </p>
                            </div>
                        </div>
                    </section>
                </div>

                {/* Right Column: Visit History */}
                <div className="lg:col-span-2 space-y-8">
                    <section className="bg-card rounded-3xl border border-border shadow-sm p-8">
                        <h3 className="text-xl font-bold text-slate-900 mb-8">Visit History</h3>
                        <div className="p-10 text-center border-2 border-dashed border-slate-100 rounded-2xl">
                            <p className="text-slate-400 font-medium">No visit history available yet.</p>
                            <button className="text-primary font-bold mt-2 hover:underline">Create first prescription</button>
                        </div>
                    </section>
                </div>
            </div>
        </div>
    );
}
