import Link from 'next/link';
import { useState } from 'react';

export default function NewPrescription() {
    const [medications, setMedications] = useState([{ name: '', dose: '', freq: '', duration: '', instructions: '' }]);

    const addMedication = () => {
        setMedications([...medications, { name: '', dose: '', freq: '', duration: '', instructions: '' }]);
    };

    const removeMedication = (index: number) => {
        setMedications(medications.filter((_, i) => i !== index));
    };

    return (
        <div className="p-8 max-w-5xl mx-auto">
            <div className="mb-8 flex justify-between items-start">
                <div>
                    <h1 className="text-3xl font-bold text-slate-900 tracking-tight">Create Prescription</h1>
                    <p className="text-slate-500">Generating digital prescription for Patient #1024</p>
                </div>
                <Link href="/dashboard/prescriptions" className="text-sm font-bold text-slate-400 hover:text-slate-900 transition flex items-center gap-1">
                    ← Back to List
                </Link>
            </div>

            <div className="space-y-8">
                {/* Patient Vitals Section */}
                <section className="bg-card p-8 rounded-2xl border border-border shadow-sm">
                    <h2 className="text-lg font-bold text-slate-900 mb-6 flex items-center gap-2">
                        <span className="p-1.5 bg-blue-50 text-blue-600 rounded-lg">💓</span>
                        Patient Vitals
                    </h2>
                    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
                        <div>
                            <label className="block text-xs font-bold text-slate-500 uppercase tracking-wider mb-2">Weight (kg)</label>
                            <input type="text" className="w-full px-4 py-2 rounded-xl border border-slate-200 focus:ring-2 focus:ring-primary outline-none" placeholder="70" />
                        </div>
                        <div>
                            <label className="block text-xs font-bold text-slate-500 uppercase tracking-wider mb-2">BP (mmHg)</label>
                            <input type="text" className="w-full px-4 py-2 rounded-xl border border-slate-200 focus:ring-2 focus:ring-primary outline-none" placeholder="120/80" />
                        </div>
                        <div>
                            <label className="block text-xs font-bold text-slate-500 uppercase tracking-wider mb-2">Temp (°F)</label>
                            <input type="text" className="w-full px-4 py-2 rounded-xl border border-slate-200 focus:ring-2 focus:ring-primary outline-none" placeholder="98.6" />
                        </div>
                        <div>
                            <label className="block text-xs font-bold text-slate-500 uppercase tracking-wider mb-2">Pulse (bpm)</label>
                            <input type="text" className="w-full px-4 py-2 rounded-xl border border-slate-200 focus:ring-2 focus:ring-primary outline-none" placeholder="72" />
                        </div>
                    </div>
                </section>

                {/* Complaints & Diagnosis */}
                <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
                    <section className="bg-card p-8 rounded-2xl border border-border shadow-sm">
                        <h2 className="text-lg font-bold text-slate-900 mb-4">Chief Complaints</h2>
                        <textarea
                            className="w-full h-32 px-4 py-2 rounded-xl border border-slate-200 focus:ring-2 focus:ring-primary outline-none resize-none"
                            placeholder="Enter patient complaints separated by commas..."
                        ></textarea>
                    </section>
                    <section className="bg-card p-8 rounded-2xl border border-border shadow-sm">
                        <h2 className="text-lg font-bold text-slate-900 mb-4">Diagnosis</h2>
                        <textarea
                            className="w-full h-32 px-4 py-2 rounded-xl border border-slate-200 focus:ring-2 focus:ring-primary outline-none resize-none"
                            placeholder="Enter diagnosis notes..."
                        ></textarea>
                    </section>
                </div>

                {/* Medication Section */}
                <section className="bg-card p-8 rounded-2xl border border-border shadow-sm">
                    <div className="flex justify-between items-center mb-6">
                        <h2 className="text-lg font-bold text-slate-900 flex items-center gap-2">
                            <span className="p-1.5 bg-teal-50 text-teal-600 rounded-lg">💊</span>
                            Medications
                        </h2>
                        <button onClick={addMedication} className="text-sm font-bold text-primary hover:underline">+ Add Medicine</button>
                    </div>

                    <div className="space-y-4">
                        {medications.map((med, index) => (
                            <div key={index} className="grid grid-cols-1 lg:grid-cols-6 gap-4 p-4 rounded-xl bg-slate-50 border border-slate-100 items-end">
                                <div className="lg:col-span-2">
                                    <label className="block text-xs font-bold text-slate-400 mb-1">Medicine Name</label>
                                    <input type="text" className="w-full px-4 py-2 rounded-lg border border-slate-200 focus:ring-2 focus:ring-primary outline-none bg-white font-medium" placeholder="Napa 500mg" />
                                </div>
                                <div>
                                    <label className="block text-xs font-bold text-slate-400 mb-1">Dosage</label>
                                    <input type="text" className="w-full px-4 py-2 rounded-lg border border-slate-200 focus:ring-2 focus:ring-primary outline-none bg-white" placeholder="1 tab" />
                                </div>
                                <div>
                                    <label className="block text-xs font-bold text-slate-400 mb-1">Frequency</label>
                                    <input type="text" className="w-full px-4 py-2 rounded-lg border border-slate-200 focus:ring-2 focus:ring-primary outline-none bg-white" placeholder="1+0+1" />
                                </div>
                                <div>
                                    <label className="block text-xs font-bold text-slate-400 mb-1">Duration</label>
                                    <input type="text" className="w-full px-4 py-2 rounded-lg border border-slate-200 focus:ring-2 focus:ring-primary outline-none bg-white" placeholder="5 days" />
                                </div>
                                <div className="flex justify-end">
                                    <button onClick={() => removeMedication(index)} className="p-2 text-slate-400 hover:text-red-500 transition">Trash</button>
                                </div>
                            </div>
                        ))}
                    </div>
                </section>

                {/* Advice & Footer */}
                <section className="bg-card p-8 rounded-2xl border border-border shadow-sm">
                    <h2 className="text-lg font-bold text-slate-900 mb-4">Advice & Recommendations</h2>
                    <textarea
                        className="w-full h-24 px-4 py-2 rounded-xl border border-slate-200 focus:ring-2 focus:ring-primary outline-none resize-none mb-6"
                        placeholder="Rest for 2 days, avoid cold water..."
                    ></textarea>

                    <div className="flex justify-between items-center bg-slate-50 -m-8 mt-8 p-8 rounded-b-2xl border-t border-border">
                        <div className="flex gap-4">
                            <button className="px-6 py-2 rounded-xl border border-slate-200 font-bold text-slate-600 hover:bg-white transition">Save Draft</button>
                            <button className="px-6 py-2 rounded-xl bg-primary text-white font-bold medical-gradient shadow-lg">Finalize & Print</button>
                        </div>
                        <div className="text-right">
                            <p className="text-xs text-slate-400 font-bold uppercase tracking-widest">Next Follow-up</p>
                            <input type="date" className="bg-transparent font-bold text-slate-700 outline-none" />
                        </div>
                    </div>
                </section>
            </div>
        </div>
    );
}
