"use client";

import Link from 'next/link';
import { useState, useEffect } from 'react';
import { useRouter, useParams } from 'next/navigation';
import { useToast } from '@/components/Toast';
import { PrescriptionReq, PrescriptionResp, Vitals, Medication } from '@/types/prescription';
import MedicineAutocomplete from '@/components/MedicineAutocomplete';

interface Patient {
    id: number;
    full_name: string;
}

interface Chamber {
    id: number;
    name: string;
}

export default function EditPrescription() {
    const router = useRouter();
    const params = useParams();
    const id = params.id as string;
    const { success, error: errorToast } = useToast();

    const [patients, setPatients] = useState<Patient[]>([]);
    const [chambers, setChambers] = useState<Chamber[]>([]);
    const [isLoading, setIsLoading] = useState(true);

    // Form State
    const [selectedPatient, setSelectedPatient] = useState<number>(0);
    const [selectedChamber, setSelectedChamber] = useState<number>(0);
    const [vitals, setVitals] = useState<Vitals>({
        weight_kg: undefined,
        blood_pressure: '',
        temperature_f: undefined,
        pulse_bpm: undefined,
    });

    const [complaints, setComplaints] = useState<string>('');
    const [diagnosis, setDiagnosis] = useState<string>('');
    const [investigations, setInvestigations] = useState<string>('');
    const [advice, setAdvice] = useState<string>('');
    const [followUpDate, setFollowUpDate] = useState<string>('');
    const [status, setStatus] = useState<'draft' | 'finalized'>('draft');

    const [medications, setMedications] = useState<Medication[]>([
        { medicine_name: '', dosage: '', frequency: '', duration: '', instructions: '' }
    ]);
    const [isSubmitting, setIsSubmitting] = useState(false);

    useEffect(() => {
        fetchInitialData();
    }, [id]);

    const fetchInitialData = async () => {
        const token = localStorage.getItem("docmate_token");
        if (!token) return;

        try {
            // Fetch Patients
            const patRes = await fetch(`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8081'}/v1/patients`, {
                headers: { 'Authorization': `Bearer ${token}` }
            });
            const patData = await patRes.json();
            if (patData.success) {
                setPatients(patData.data.records);
            }

            // Fetch Doctor Profile to get ID for chambers
            const docRes = await fetch(`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8081'}/v1/doctors/profile`, {
                headers: { 'Authorization': `Bearer ${token}` }
            });
            const docData = await docRes.json();

            if (docData.success) {
                const doctorId = docData.data.id;
                // Fetch Chambers
                const chamRes = await fetch(`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8081'}/v1/doctors/${doctorId}/chambers`, {
                    headers: { 'Authorization': `Bearer ${token}` }
                });
                const chamData = await chamRes.json();
                if (chamData.success) {
                    setChambers(chamData.data.records);
                }

                // Fetch existing prescription
                const pxRes = await fetch(`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8081'}/v1/prescriptions/${id}`, {
                    headers: { 'Authorization': `Bearer ${token}` }
                });
                const pxData = await pxRes.json();
                if (pxData.success) {
                    const px = pxData.data as PrescriptionResp;

                    if (px.status === 'finalized') {
                        errorToast("Finalized prescriptions cannot be edited");
                        router.push('/prescriptions');
                        return;
                    }

                    setSelectedPatient(px.patient_id);
                    setSelectedChamber(px.chamber_id);
                    setVitals(px.vitals);
                    setComplaints(px.chief_complaints.join(', '));
                    setDiagnosis(px.diagnosis.join(', '));
                    setInvestigations(px.investigations.join(', '));
                    setAdvice(px.advice);
                    setStatus(px.status);
                    if (px.follow_up_date) {
                        setFollowUpDate(new Date(px.follow_up_date).toISOString().split('T')[0]);
                    }
                    setMedications(px.medications.length > 0 ? px.medications : [{ medicine_name: '', dosage: '', frequency: '', duration: '', instructions: '' }]);
                }
            }
        } catch (error) {
            console.error("Failed to load data", error);
            errorToast("Failed to load prescription data");
        } finally {
            setIsLoading(false);
        }
    };

    const addMedication = () => {
        setMedications([...medications, { medicine_name: '', dosage: '', frequency: '', duration: '', instructions: '' }]);
    };

    const removeMedication = (index: number) => {
        setMedications(medications.filter((_, i) => i !== index));
    };

    const updateMedication = (index: number, field: keyof Medication, value: string) => {
        const newMeds = [...medications];
        newMeds[index] = { ...newMeds[index], [field]: value };
        setMedications(newMeds);
    };

    const handleUpdate = async (newStatus: 'draft' | 'finalized') => {
        if (!selectedPatient || !selectedChamber) {
            errorToast("Please select both a patient and a chamber");
            return;
        }

        setIsSubmitting(true);
        const token = localStorage.getItem("docmate_token");

        const payload: PrescriptionReq = {
            patient_id: selectedPatient,
            chamber_id: selectedChamber,
            vitals: vitals,
            chief_complaints: complaints.split(',').map((s: string) => s.trim()).filter(Boolean),
            diagnosis: diagnosis.split(',').map((s: string) => s.trim()).filter(Boolean),
            investigations: investigations.split(',').map((s: string) => s.trim()).filter(Boolean),
            advice: advice,
            status: newStatus,
            follow_up_date: followUpDate ? new Date(followUpDate).toISOString() : undefined,
            medications: medications.filter(m => m.medicine_name.trim() !== ''),
        };

        try {
            const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8081'}/v1/prescriptions/${id}`, {
                method: "PUT",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": `Bearer ${token}`
                },
                body: JSON.stringify(payload),
            });

            const data = await response.json();

            if (response.ok && data.success) {
                success(`Prescription ${newStatus === 'finalized' ? 'finalized' : 'updated'} successfully`);
                if (newStatus === 'finalized') {
                    router.push(`/prescriptions/${id}/print`);
                } else {
                    router.push('/prescriptions');
                }
            } else {
                errorToast(data.message || `Failed to ${newStatus} prescription`);
            }
        } catch (error) {
            console.error("Error updating prescription:", error);
            errorToast("An unexpected error occurred");
        } finally {
            setIsSubmitting(false);
        }
    };

    if (isLoading) {
        return <div className="p-8 text-center">Loading prescription...</div>;
    }

    return (
        <div className="p-8 max-w-5xl mx-auto">
            <div className="mb-8 flex justify-between items-start">
                <div>
                    <h1 className="text-3xl font-bold text-slate-900 tracking-tight">Edit Prescription</h1>
                    <p className="text-slate-500">Updating draft prescription</p>
                </div>
                <Link href="/prescriptions" className="text-sm font-bold text-slate-400 hover:text-slate-900 transition flex items-center gap-1">
                    ← Back to List
                </Link>
            </div>

            <div className="space-y-8">
                {/* Meta Section */}
                <section className="bg-card p-8 rounded-2xl border border-border shadow-sm">
                    <div className="grid grid-cols-1 sm:grid-cols-2 gap-6">
                        <div>
                            <label className="block text-xs font-bold text-slate-500 uppercase tracking-wider mb-2">Select Patient</label>
                            <select
                                value={selectedPatient}
                                onChange={(e) => setSelectedPatient(Number(e.target.value))}
                                className="w-full px-4 py-2 rounded-xl border border-slate-200 focus:ring-2 focus:ring-primary outline-none bg-white font-medium text-slate-700"
                            >
                                <option value={0} disabled>Select a patient...</option>
                                {patients.map(p => (
                                    <option key={p.id} value={p.id}>{p.full_name}</option>
                                ))}
                            </select>
                        </div>
                        <div>
                            <label className="block text-xs font-bold text-slate-500 uppercase tracking-wider mb-2">Select Chamber</label>
                            <select
                                value={selectedChamber}
                                onChange={(e) => setSelectedChamber(Number(e.target.value))}
                                className="w-full px-4 py-2 rounded-xl border border-slate-200 focus:ring-2 focus:ring-primary outline-none bg-white font-medium text-slate-700"
                            >
                                <option value={0} disabled>Select chamber...</option>
                                {chambers.map(c => (
                                    <option key={c.id} value={c.id}>{c.name}</option>
                                ))}
                            </select>
                        </div>
                    </div>
                </section>

                {/* Patient Vitals Section */}
                <section className="bg-card p-8 rounded-2xl border border-border shadow-sm">
                    <h2 className="text-lg font-bold text-slate-900 mb-6 flex items-center gap-2">
                        <span className="p-1.5 bg-blue-50 text-blue-600 rounded-lg">💓</span>
                        Patient Vitals
                    </h2>
                    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
                        <div>
                            <label className="block text-xs font-bold text-slate-500 uppercase tracking-wider mb-2">Weight (kg)</label>
                            <input
                                type="number"
                                value={vitals.weight_kg || ''}
                                onChange={(e) => setVitals({ ...vitals, weight_kg: parseFloat(e.target.value) || undefined })}
                                className="w-full px-4 py-2 rounded-xl border border-slate-200 focus:ring-2 focus:ring-primary outline-none"
                                placeholder="70"
                            />
                        </div>
                        <div>
                            <label className="block text-xs font-bold text-slate-500 uppercase tracking-wider mb-2">BP (mmHg)</label>
                            <input
                                type="text"
                                value={vitals.blood_pressure || ''}
                                onChange={(e) => setVitals({ ...vitals, blood_pressure: e.target.value })}
                                className="w-full px-4 py-2 rounded-xl border border-slate-200 focus:ring-2 focus:ring-primary outline-none"
                                placeholder="120/80"
                            />
                        </div>
                        <div>
                            <label className="block text-xs font-bold text-slate-500 uppercase tracking-wider mb-2">Temp (°F)</label>
                            <input
                                type="number" step="0.1"
                                value={vitals.temperature_f || ''}
                                onChange={(e) => setVitals({ ...vitals, temperature_f: parseFloat(e.target.value) || undefined })}
                                className="w-full px-4 py-2 rounded-xl border border-slate-200 focus:ring-2 focus:ring-primary outline-none"
                                placeholder="98.6"
                            />
                        </div>
                        <div>
                            <label className="block text-xs font-bold text-slate-500 uppercase tracking-wider mb-2">Pulse (bpm)</label>
                            <input
                                type="number"
                                value={vitals.pulse_bpm || ''}
                                onChange={(e) => setVitals({ ...vitals, pulse_bpm: parseInt(e.target.value) || undefined })}
                                className="w-full px-4 py-2 rounded-xl border border-slate-200 focus:ring-2 focus:ring-primary outline-none"
                                placeholder="72"
                            />
                        </div>
                    </div>
                </section>

                {/* Complaints & Diagnosis */}
                <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
                    <section className="bg-card p-8 rounded-2xl border border-border shadow-sm">
                        <h2 className="text-lg font-bold text-slate-900 mb-4">Chief Complaints</h2>
                        <textarea
                            value={complaints}
                            onChange={(e) => setComplaints(e.target.value)}
                            className="w-full h-32 px-4 py-2 rounded-xl border border-slate-200 focus:ring-2 focus:ring-primary outline-none resize-none"
                            placeholder="Enter patient complaints separated by commas (e.g., Fever, Dry Cough, Headache)"
                        ></textarea>
                    </section>
                    <section className="bg-card p-8 rounded-2xl border border-border shadow-sm">
                        <h2 className="text-lg font-bold text-slate-900 mb-4">Diagnosis</h2>
                        <textarea
                            value={diagnosis}
                            onChange={(e) => setDiagnosis(e.target.value)}
                            className="w-full h-32 px-4 py-2 rounded-xl border border-slate-200 focus:ring-2 focus:ring-primary outline-none resize-none"
                            placeholder="Enter diagnosis notes separated by commas..."
                        ></textarea>
                    </section>
                </div>

                {/* Medication Section */}
                <section className="bg-card p-8 rounded-2xl border border-border shadow-sm">
                    <div className="flex justify-between items-center mb-6">
                        <h2 className="text-lg font-bold text-slate-900 flex items-center gap-2">
                            <span className="p-1.5 bg-teal-50 text-teal-600 rounded-lg">💊</span>
                            Medications (Rx)
                        </h2>
                        <button onClick={addMedication} className="text-sm font-bold text-primary hover:underline">+ Add Medicine</button>
                    </div>

                    <div className="space-y-4">
                        {medications.map((med, index) => (
                            <div key={index} className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-12 gap-4 p-4 rounded-xl bg-slate-50 border border-slate-100 items-start">
                                <div className="lg:col-span-3">
                                    <label className="block text-xs font-bold text-slate-400 mb-1">Medicine Name <span className="text-red-400">*</span></label>
                                    <MedicineAutocomplete
                                        value={med.medicine_name}
                                        onChange={(val) => updateMedication(index, 'medicine_name', val)}
                                        placeholder="Tab. Napa 500mg"
                                    />
                                </div>
                                <div className="lg:col-span-2">
                                    <label className="block text-xs font-bold text-slate-400 mb-1">Dosage</label>
                                    <input
                                        type="text"
                                        value={med.dosage}
                                        onChange={(e) => updateMedication(index, 'dosage', e.target.value)}
                                        className="w-full px-4 py-2 rounded-lg border border-slate-200 focus:ring-2 focus:ring-primary outline-none bg-white"
                                        placeholder="1 tab"
                                    />
                                </div>
                                <div className="lg:col-span-2">
                                    <label className="block text-xs font-bold text-slate-400 mb-1">Frequency <span className="text-red-400">*</span></label>
                                    <input
                                        type="text"
                                        value={med.frequency}
                                        onChange={(e) => updateMedication(index, 'frequency', e.target.value)}
                                        className="w-full px-4 py-2 rounded-lg border border-slate-200 focus:ring-2 focus:ring-primary outline-none bg-white"
                                        placeholder="1+0+1"
                                    />
                                </div>
                                <div className="lg:col-span-2">
                                    <label className="block text-xs font-bold text-slate-400 mb-1">Duration <span className="text-red-400">*</span></label>
                                    <input
                                        type="text"
                                        value={med.duration}
                                        onChange={(e) => updateMedication(index, 'duration', e.target.value)}
                                        className="w-full px-4 py-2 rounded-lg border border-slate-200 focus:ring-2 focus:ring-primary outline-none bg-white"
                                        placeholder="5 days"
                                    />
                                </div>
                                <div className="lg:col-span-2">
                                    <label className="block text-xs font-bold text-slate-400 mb-1">Instructions</label>
                                    <input
                                        type="text"
                                        value={med.instructions || ''}
                                        onChange={(e) => updateMedication(index, 'instructions', e.target.value)}
                                        className="w-full px-4 py-2 rounded-lg border border-slate-200 focus:ring-2 focus:ring-primary outline-none bg-white font-bengali"
                                        placeholder="খাওয়ার পরে"
                                    />
                                </div>
                                <div className="flex justify-end lg:col-span-1 lg:mt-6">
                                    <button onClick={() => removeMedication(index)} className="p-2 text-slate-400 hover:text-red-500 transition">Trash</button>
                                </div>
                            </div>
                        ))}
                    </div>
                </section>

                {/* Investigations Section */}
                <section className="bg-card p-8 rounded-2xl border border-border shadow-sm">
                    <h2 className="text-lg font-bold text-slate-900 mb-4 flex items-center gap-2">
                        <span className="p-1.5 bg-purple-50 text-purple-600 rounded-lg">🔬</span>
                        Investigations
                    </h2>
                    <textarea
                        value={investigations}
                        onChange={(e) => setInvestigations(e.target.value)}
                        className="w-full h-24 px-4 py-2 rounded-xl border border-slate-200 focus:ring-2 focus:ring-primary outline-none resize-none"
                        placeholder="Enter recommended tests separated by commas (e.g., CBC, Chest X-Ray PA View)"
                    ></textarea>
                </section>

                {/* Advice & Footer */}
                <section className="bg-card p-8 rounded-2xl border border-border shadow-sm">
                    <h2 className="text-lg font-bold text-slate-900 mb-4">Advice & Recommendations</h2>
                    <textarea
                        value={advice}
                        onChange={(e) => setAdvice(e.target.value)}
                        className="w-full h-24 px-4 py-2 rounded-xl border border-slate-200 focus:ring-2 focus:ring-primary outline-none resize-none mb-6"
                        placeholder="Rest for 2 days, avoid cold water..."
                    ></textarea>

                    <div className="flex justify-between items-center bg-slate-50 -m-8 mt-8 p-8 rounded-b-2xl border-t border-border">
                        <div className="flex gap-4">
                            <button
                                onClick={() => handleUpdate('draft')}
                                disabled={isSubmitting}
                                className="px-6 py-2 rounded-xl border border-slate-200 font-bold text-slate-600 hover:bg-white transition disabled:opacity-50"
                            >
                                {isSubmitting ? 'Updating...' : 'Update Draft'}
                            </button>
                            <button
                                onClick={() => handleUpdate('finalized')}
                                disabled={isSubmitting}
                                className="px-6 py-2 rounded-xl bg-primary text-white font-bold medical-gradient shadow-lg disabled:opacity-50"
                            >
                                {isSubmitting ? 'Finalizing...' : 'Finalize & Print'}
                            </button>
                        </div>
                        <div className="text-right">
                            <p className="text-xs text-slate-400 font-bold uppercase tracking-widest">Next Follow-up</p>
                            <input
                                type="date"
                                value={followUpDate}
                                onChange={(e) => setFollowUpDate(e.target.value)}
                                className="bg-transparent font-bold text-slate-700 outline-none cursor-pointer"
                            />
                        </div>
                    </div>
                </section>
            </div>
        </div>
    );
}
