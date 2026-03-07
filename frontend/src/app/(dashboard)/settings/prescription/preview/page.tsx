"use client";

import { useEffect, useState, Suspense } from 'react';
import { useSearchParams } from 'next/navigation';
import { PrescriptionSettingResp } from '@/types/prescription-setting';

function PreviewContent() {
    const searchParams = useSearchParams();
    const [settings, setSettings] = useState<Partial<PrescriptionSettingResp>>({});

    useEffect(() => {
        const data = searchParams.get('data');
        if (data) {
            try {
                setSettings(JSON.parse(decodeURIComponent(data)));
            } catch (e) {
                console.error("Failed to parse preview data", e);
            }
        }
    }, [searchParams]);

    // Dummy data for preview
    const dummyDoctor = {
        full_name: "Dr. Faisal Ahmed",
        degree: "MBBS, FCPS (Medicine)",
        specialization: "Internal Medicine Specialist",
        phone: "+880 1234 567890",
        email: "faisal.ahmed@example.com"
    };

    const dummyPatient = {
        full_name: "John Doe",
        age: 35,
        gender: "Male"
    };

    const dummyPrescription = {
        id: "PRE-12345",
        created_at: new Date().toISOString(),
        vitals: { weight_kg: 72, blood_pressure: "120/80", pulse_bpm: 72, temperature_f: 98.6 },
        chief_complaints: ["Persistent cough for 3 days", "Mild fever", "Headache"],
        diagnosis: ["Acute Bronchitis"],
        investigations: ["Chest X-Ray", "CBC"],
        medications: [
            { medicine_name: "Tab. Azithromycin 500mg", frequency: "1+0+0", dosage: "1 tab", duration: "5 days", instructions: "After meal" },
            { medicine_name: "Syr. Adryl", frequency: "2 tsp", dosage: "3 times daily", duration: "7 days", instructions: "If cough persists" },
            { medicine_name: "Tab. Paracetamol 500mg", frequency: "1+1+1", dosage: "1 tab", duration: "3 days", instructions: "If fever > 100°F" }
        ],
        advice: "Drink plenty of warm water.\nAvoid cold drinks and ice cream.\nComplete the full course of antibiotics."
    };

    return (
        <div className="print-container bg-white text-black min-h-screen relative p-8 md:p-12 max-w-[210mm] mx-auto shadow-2xl my-8 text-sm overflow-hidden">
            <div className="no-print absolute top-4 right-4 flex gap-4">
                <button onClick={() => window.close()} className="px-4 py-2 border border-slate-300 rounded shadow-sm bg-white hover:bg-slate-50 text-xs font-bold">Close Preview</button>
                <button onClick={() => window.print()} className="px-4 py-2 bg-blue-600 text-white rounded shadow text-xs font-bold flex items-center gap-2">🖨️ Print Test</button>
            </div>

            <div className={`flex justify-between items-start border-b-2 border-green-700 pb-4 mb-4 ${settings.template_type === 'modern' ? 'border-blue-600' : ''}`}>
                <div className="max-w-[50%] whitespace-pre-wrap">
                    {settings.header_left_bangla ? (
                        <div className="text-[13px] leading-snug">
                            {settings.header_left_bangla.split('\n').map((line, i) => (
                                <p key={i} className={i === 0 ? "text-xl font-bold text-green-900 mb-1" : ""}>{line}</p>
                            ))}
                        </div>
                    ) : (
                        <div className="opacity-40 italic">Bangla Header Placeholder</div>
                    )}
                </div>

                <div className="text-right max-w-[50%] whitespace-pre-wrap">
                    {settings.header_right_english ? (
                        <div className="text-[13px] leading-snug">
                            {settings.header_right_english.split('\n').map((line, i) => (
                                <p key={i} className={i === 0 ? "text-lg font-bold uppercase mb-1" : ""}>{line}</p>
                            ))}
                        </div>
                    ) : (
                        <div className="opacity-40 italic">English Header Placeholder</div>
                    )}
                </div>
            </div>

            <div className="flex justify-between items-center border-b-2 border-black pb-2 mb-6 text-sm font-bold">
                <div className="flex-1">Name: <span className="font-normal">{dummyPatient.full_name}</span></div>
                <div className="flex-1 text-center">Age: <span className="font-normal">{dummyPatient.age}</span></div>
                <div className="flex-1 text-center">Date: <span className="font-normal">{new Date().toLocaleDateString('en-GB')}</span></div>
                <div className="flex-1 text-right">ID: <span className="font-normal">{dummyPrescription.id}</span></div>
            </div>

            <div className="flex gap-8 min-h-[600px]">
                <div className="w-1/3 border-r-2 border-slate-200 pr-6 shrink-0 flex flex-col gap-6">
                    <div>
                        <h3 className="font-bold mb-2">Vitals</h3>
                        <ul className="text-[13px] leading-relaxed">
                            <li>Wt: {dummyPrescription.vitals.weight_kg} kg</li>
                            <li>BP: {dummyPrescription.vitals.blood_pressure}</li>
                            <li>Pulse: {dummyPrescription.vitals.pulse_bpm} bpm</li>
                        </ul>
                    </div>
                    <div>
                        <h3 className="font-bold mb-2">Chief Complaints</h3>
                        <ul className="text-[13px] leading-relaxed list-none space-y-1">
                            {dummyPrescription.chief_complaints.map((c, i) => <li key={i}>- {c}</li>)}
                        </ul>
                    </div>
                </div>

                <div className="w-2/3 pl-2">
                    <div className="text-3xl font-bold mb-6 font-serif tracking-widest">Rx,</div>
                    <div className="space-y-8">
                        {dummyPrescription.medications.map((med, index) => (
                            <div key={index} className="flex gap-4 items-start border-b border-gray-100 pb-4">
                                <span className="font-bold text-lg">{index + 1}.</span>
                                <div className="flex-1">
                                    <h4 className="font-bold text-[15px] mb-2">{med.medicine_name}</h4>
                                    <div className="flex justify-between items-center text-[14px]">
                                        <div className="flex items-center gap-8">
                                            <span className="font-semibold tracking-[0.2em]">{med.frequency}</span>
                                            <span>{med.dosage}</span>
                                        </div>
                                        <span className="font-medium">{med.duration}</span>
                                    </div>
                                    <p className="text-[13px] mt-1 text-gray-600 italic">{med.instructions}</p>
                                </div>
                            </div>
                        ))}
                    </div>
                </div>
            </div>

            <div className="absolute bottom-12 right-12 text-right">
                {settings.signature_url ? (
                    <div className="mb-1">
                        <img src={settings.signature_url} alt="Doctor Signature" className="h-16 object-contain mix-blend-multiply ml-auto" />
                    </div>
                ) : (
                    <div className="text-xl font-bold mb-1 opacity-20">[Signature Placeholder]</div>
                )}
                <div className="border-t border-black pt-1 px-4 text-xs font-bold uppercase">Signature</div>
            </div>

            {(settings.footer_info_bangla || settings.footer_info_english) && (
                <div className="absolute bottom-4 left-8 right-8 flex justify-between items-end border-t border-gray-200 pt-2 text-[10px] text-gray-500 italic">
                    <div className="max-w-[45%] whitespace-pre-wrap">{settings.footer_info_bangla}</div>
                    <div className="max-w-[45%] text-right whitespace-pre-wrap">{settings.footer_info_english}</div>
                </div>
            )}
        </div>
    );
}

export default function PreviewPage() {
    return (
        <div className="bg-slate-100 min-h-screen p-4 md:p-12">
            <Suspense fallback={<div className="flex h-screen items-center justify-center font-bold text-slate-500">Loading Preview...</div>}>
                <PreviewContent />
            </Suspense>
        </div>
    );
}
