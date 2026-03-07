"use client";

import { useState } from 'react';
import { useRouter } from 'next/navigation';

export default function SettingsPage() {
    const router = useRouter();
    const [template, setTemplate] = useState('modern');

    return (
        <div className="p-8 max-w-4xl">
            <div className="mb-10">
                <h1 className="text-3xl font-bold text-slate-900 tracking-tight">Account Settings</h1>
                <p className="text-slate-500">Configure your profile and prescription preferences</p>
            </div>

            <div className="space-y-10">
                {/* Profile Card */}
                <section className="bg-card rounded-3xl border border-border shadow-sm p-8">
                    <h2 className="text-xl font-bold text-slate-900 mb-6">Doctor Profile</h2>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
                        <div>
                            <label className="block text-xs font-bold text-slate-400 uppercase mb-2">Full Name</label>
                            <input type="text" className="w-full px-4 py-2.5 rounded-xl border border-slate-200" defaultValue="Dr. Faisal" />
                        </div>
                        <div>
                            <label className="block text-xs font-bold text-slate-400 uppercase mb-2">Designation</label>
                            <input type="text" className="w-full px-4 py-2.5 rounded-xl border border-slate-200" defaultValue="Cardiologist" />
                        </div>
                        <div className="md:col-span-2">
                            <label className="block text-xs font-bold text-slate-400 uppercase mb-2">Signature</label>
                            <div className="border-2 border-dashed border-slate-200 rounded-2xl p-10 text-center hover:border-primary transition group cursor-pointer">
                                <span className="text-slate-400 group-hover:text-primary">Click to upload or drag & drop image</span>
                            </div>
                        </div>
                    </div>
                    <div className="mt-8 flex justify-end">
                        <button className="bg-primary text-white px-8 py-2.5 rounded-xl font-bold medical-gradient shadow-lg">Save Profile</button>
                    </div>
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
        </div>
    );
}
