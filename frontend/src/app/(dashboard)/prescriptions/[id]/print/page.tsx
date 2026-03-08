"use client";

import { useEffect, useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { PrescriptionResp } from '@/types/prescription';
import { PrescriptionSettingResp } from '@/types/prescription-setting';
import { useToast } from '@/components/Toast';

// Interfaces for nested relationships
interface Patient { full_name: string; date_of_birth: string; gender: string; }
interface Chamber { name: string; address: string; fee: number; phone: string; }
interface Doctor { full_name: string; degree: any; specialization: any; phone: string; email: string; signature_url?: string; }

export default function PrintPrescription() {
    const params = useParams();
    const router = useRouter();
    const { error: errorToast } = useToast();

    const [prescription, setPrescription] = useState<PrescriptionResp | null>(null);
    const [patient, setPatient] = useState<Patient | null>(null);
    const [chamber, setChamber] = useState<Chamber | null>(null);
    const [doctor, setDoctor] = useState<Doctor | null>(null);
    const [settings, setSettings] = useState<PrescriptionSettingResp | null>(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        if (params.id) {
            fetchData(Number(params.id));
        }
    }, [params.id]);

    const fetchData = async (id: number) => {
        const token = localStorage.getItem("docmate_token");
        if (!token) return;

        try {
            // 1. Fetch Prescription
            const presRes = await fetch(`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8081'}/v1/prescriptions/${id}`, {
                headers: { 'Authorization': `Bearer ${token}` }
            });
            const presData = await presRes.json();

            if (!presData.success) {
                errorToast("Failed to load prescription");
                setLoading(false);
                return;
            }

            const p: PrescriptionResp = presData.data;
            setPrescription(p);

            // Fetch relations in parallel
            const [docRes, patRes, chamRes] = await Promise.all([
                fetch(`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8081'}/v1/doctors/${p.doctor_id}`, { headers: { 'Authorization': `Bearer ${token}` } }),
                fetch(`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8081'}/v1/patients/${p.patient_id}`, { headers: { 'Authorization': `Bearer ${token}` } }),
                fetch(`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8081'}/v1/doctors/${p.doctor_id}/chambers/${p.chamber_id}`, { headers: { 'Authorization': `Bearer ${token}` } })
            ]);

            const [docData, patData, chamData] = await Promise.all([docRes.json(), patRes.json(), chamRes.json()]);

            if (docData.success) setDoctor(docData.data);
            if (patData.success) setPatient(patData.data);
            if (chamData.success) setChamber(chamData.data);

            // 2. Fetch Prescription Settings
            const setRes = await fetch(`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8081'}/v1/doctors/${p.doctor_id}/prescription-settings?chamber_id=${p.chamber_id}`, {
                headers: { 'Authorization': `Bearer ${token}` }
            });
            const setData = await setRes.json();
            if (setData.success) setSettings(setData.data);

        } catch (error) {
            console.error("Error loading print data:", error);
            errorToast("An error occurred loading the prescription for printing");
        } finally {
            setLoading(false);
        }
    };

    const calculateAge = (dob: string) => {
        if (!dob) return 'N/A';
        const diff = Date.now() - new Date(dob).getTime();
        const age = new Date(diff);
        return Math.abs(age.getUTCFullYear() - 1970);
    };

    if (loading) {
        return <div className="flex h-screen items-center justify-center font-bold text-slate-500">Preparing Document...</div>;
    }

    if (!prescription || !doctor || !patient) {
        return (
            <div className="flex flex-col h-screen items-center justify-center">
                <p className="font-bold text-red-500 mb-4">Error generating prescription view.</p>
                <button onClick={() => router.back()} className="px-4 py-2 bg-slate-200 rounded-lg">Go Back</button>
            </div>
        );
    }

    // Helper functions to parse GORM datatypes.JSON from backend
    const parseJSONB = (val: any, fallback: string) => {
        if (!val) return fallback;
        if (typeof val === 'string') return val;
        if (Array.isArray(val) && val.length > 0) return val.join(', ');
        return fallback;
    };

    return (
        <div className="print-container bg-white text-black min-h-screen relative p-8 md:p-12 max-w-[210mm] mx-auto overflow-hidden text-sm">
            {/* Action Bar (Hidden when printing) */}
            <div className="no-print absolute top-4 right-4 flex gap-4">
                <button onClick={() => router.back()} className="px-4 py-2 border border-slate-300 rounded shadow-sm bg-white hover:bg-slate-50">Back</button>
                <button onClick={() => window.print()} className="px-4 py-2 bg-blue-600 text-white rounded shadow text-sm font-bold flex items-center gap-2">🖨️ Print PDF</button>
            </div>

            {/* Header section matches image */}
            <div className={`flex justify-between items-start border-b-2 border-green-700 pb-4 mb-4`}>
                {/* Left Side Doctor Info (Bangla/Custom) */}
                <div className="max-w-[50%] whitespace-pre-wrap">
                    {settings?.header_left_bangla ? (
                        <div className="text-[13px] leading-snug">
                            {settings.header_left_bangla.split('\n').map((line, i) => (
                                <p key={i} className={i === 0 ? "text-xl font-bold text-green-900 mb-1" : ""}>{line}</p>
                            ))}
                        </div>
                    ) : (
                        <>
                            <h1 className="text-2xl font-bold text-green-900 mb-1 leading-tight">ডাাঃ {doctor.full_name.split(' ').map(n => n.charAt(0).toUpperCase() + n.slice(1)).join(' ')}</h1>
                            <p className="text-sm font-semibold">{parseJSONB(doctor.degree, 'MBBS, FCPS')}</p>
                            <p className="text-[13px]">{parseJSONB(doctor.specialization, 'Medicine Specialist')}</p>
                            {chamber && (
                                <div className="mt-2 text-[12px] leading-tight text-gray-700">
                                    <p className="font-semibold text-red-600">Chamber:</p>
                                    <p>{chamber.name}</p>
                                    <p>{chamber.address}</p>
                                </div>
                            )}
                        </>
                    )}
                </div>

                {/* Right side English Doctor Info */}
                <div className="text-right max-w-[50%] whitespace-pre-wrap">
                    {settings?.header_right_english ? (
                        <div className="text-[13px] leading-snug">
                            {settings.header_right_english.split('\n').map((line, i) => (
                                <p key={i} className={i === 0 ? "text-lg font-bold uppercase mb-1" : ""}>{line}</p>
                            ))}
                        </div>
                    ) : (
                        <>
                            <h1 className="text-xl font-bold uppercase mb-1">DR. {doctor.full_name}</h1>
                            <p className="text-sm font-medium">{parseJSONB(doctor.degree, 'MBBS')}</p>
                            <p className="text-[13px] text-red-600 font-bold">{parseJSONB(doctor.specialization, '')}</p>
                            <p className="text-[12px] leading-tight mt-1">{chamber?.address || ''}</p>
                            <p className="text-[12px] leading-tight mt-1 font-semibold">{doctor.phone}</p>
                            <p className="text-[12px] leading-tight">{doctor.email}</p>
                        </>
                    )}
                </div>
            </div>

            {/* Patient Meta Row */}
            <div className="flex justify-between items-center border-b-2 border-black pb-2 mb-6 text-sm font-bold">
                <div className="flex-1">Name: <span className="font-normal">{patient.full_name}</span></div>
                <div className="flex-1 text-center">Age: <span className="font-normal">{calculateAge(patient.date_of_birth)}</span></div>
                <div className="flex-1 text-center">Date: <span className="font-normal">{new Date(prescription.created_at).toLocaleDateString('en-GB')}</span></div>
                <div className="flex-1 text-right">ID: <span className="font-normal">{prescription.id}</span></div>
            </div>

            {/* Main Content Area: 2 Columns */}
            <div className="flex gap-8 min-h-[600px]">
                {/* Left Column: Complaints, Investigations, Advice */}
                <div className="w-1/3 border-r-2 border-slate-200 pr-6 shrink-0 flex flex-col gap-6">
                    {prescription.vitals && Object.keys(prescription.vitals).length > 0 && (
                        <div>
                            <h3 className="font-bold mb-2">Vitals</h3>
                            <ul className="text-[13px] leading-relaxed">
                                {prescription.vitals.weight_kg && <li>Wt: {prescription.vitals.weight_kg} kg</li>}
                                {prescription.vitals.blood_pressure && <li>BP: {prescription.vitals.blood_pressure}</li>}
                                {prescription.vitals.pulse_bpm && <li>Pulse: {prescription.vitals.pulse_bpm} bpm</li>}
                                {prescription.vitals.temperature_f && <li>Temp: {prescription.vitals.temperature_f} °F</li>}
                            </ul>
                        </div>
                    )}

                    {prescription.chief_complaints && prescription.chief_complaints.length > 0 && (
                        <div>
                            <h3 className="font-bold mb-2">Chief Complaints</h3>
                            <ul className="text-[13px] leading-relaxed list-none space-y-1">
                                {prescription.chief_complaints.map((c, i) => (
                                    <li key={i}>- {c}</li>
                                ))}
                            </ul>
                        </div>
                    )}

                    {prescription.diagnosis && prescription.diagnosis.length > 0 && (
                        <div>
                            <h3 className="font-bold mb-2">Diagnosis</h3>
                            <ul className="text-[13px] leading-relaxed list-none space-y-1">
                                {prescription.diagnosis.map((d, i) => (
                                    <li key={i}>- {d}</li>
                                ))}
                            </ul>
                        </div>
                    )}

                    {prescription.investigations && prescription.investigations.length > 0 && (
                        <div>
                            <h3 className="font-bold mb-2">Investigation</h3>
                            <ul className="text-[13px] leading-relaxed list-none space-y-1">
                                {prescription.investigations.map((inv, i) => (
                                    <li key={i}>- {inv}</li>
                                ))}
                            </ul>
                        </div>
                    )}

                    {prescription.advice && (
                        <div className="mt-auto no-break">
                            <h3 className="font-bold mb-2">Advice</h3>
                            <p className="text-[13px] whitespace-pre-wrap">{prescription.advice}</p>
                        </div>
                    )}

                    {prescription.follow_up_date && (
                        <div className="mt-4 no-break">
                            <h3 className="font-bold mb-1 text-red-600">Next Follow-up</h3>
                            <p className="text-[13px]">{new Date(prescription.follow_up_date).toLocaleDateString()}</p>
                        </div>
                    )}
                </div>

                {/* Right Column: Medications (Rx) */}
                <div className="w-2/3 pl-2">
                    <div className="text-3xl font-bold mb-6 font-serif tracking-widest">Rx,</div>

                    <div className="space-y-8 no-break-inside">
                        {prescription.medications && prescription.medications.map((med, index) => (
                            <div key={index} className="flex gap-4 items-start break-inside-avoid">
                                <span className="font-bold text-lg">{index + 1}.</span>
                                <div className="flex-1">
                                    <h4 className="font-bold text-[15px] mb-2">{med.medicine_name}</h4>

                                    <div className="flex justify-between items-center text-[14px]">
                                        <div className="flex items-center gap-8">
                                            <span className="font-semibold tracking-[0.2em]">{med.frequency}</span>
                                            {med.dosage && <span>{med.dosage}</span>}
                                        </div>
                                        <span className="font-medium whitespace-nowrap min-w-[60px] text-right">{med.duration}</span>
                                    </div>

                                    {med.instructions && (
                                        <p className="text-[13px] mt-1 text-gray-800 italic">{med.instructions}</p>
                                    )}
                                    <div className="w-full border-b border-gray-300 mt-3 mix-blend-multiply opacity-50"></div>
                                </div>
                            </div>
                        ))}
                    </div>
                </div>
            </div>

            {/* Signature Area */}
            <div className="absolute bottom-12 right-12 text-right">
                {doctor.signature_url ? (
                    <div className="mb-1">
                        <img src={doctor.signature_url} alt="Doctor Signature" className="h-16 object-contain mix-blend-multiply ml-auto" />
                    </div>
                ) : (
                    <div className="font-custom-signature text-2xl mb-1">{doctor.full_name}</div>
                )}
                <div className="border-t border-black pt-1 px-4 text-xs font-bold uppercase">Signature</div>
            </div>

            {/* Footer Section */}
            {(settings?.footer_info_bangla || settings?.footer_info_english) && (
                <div className="absolute bottom-4 left-8 right-8 flex justify-between items-end border-t border-gray-200 pt-2 text-[10px] text-gray-500 italic">
                    <div className="max-w-[45%] whitespace-pre-wrap">{settings.footer_info_bangla}</div>
                    <div className="max-w-[45%] text-right whitespace-pre-wrap">{settings.footer_info_english}</div>
                </div>
            )}
        </div>
    );
}
