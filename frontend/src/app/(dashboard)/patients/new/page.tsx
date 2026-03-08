"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { useToast } from "@/components/Toast";

export default function NewPatientPage() {
    const router = useRouter();
    const [formData, setFormData] = useState({
        full_name: '',
        phone: '',
        age: '',
        gender: 'male',
        blood_group: 'A+',
    });
    const { success: successToast, error: errorToast } = useToast();
    const [isSubmitting, setIsSubmitting] = useState(false);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setIsSubmitting(true);

        const token = localStorage.getItem('docmate_token');
        if (!token) {
            errorToast('Unauthorized: Please log in again.');
            router.push('/login');
            return;
        }

        try {
            const response = await fetch('http://localhost:8081/v1/patients', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`,
                },
                body: JSON.stringify({
                    full_name: formData.full_name,
                    phone: formData.phone,
                    age: parseInt(formData.age),
                    gender: formData.gender.toLowerCase(),
                    blood_group: formData.blood_group,
                }),
            });

            const data = await response.json();

            if (response.ok && data.success) {
                successToast('Patient created successfully!');
                router.push('/patients');
            } else {
                errorToast(data.message || 'Failed to create patient');
            }
        } catch (error) {
            console.error('Error creating patient:', error);
            errorToast('An error occurred while saving patient information.');
        } finally {
            setIsSubmitting(false);
        }
    };

    return (
        <div className="p-8 max-w-2xl mx-auto">
            <Link href="/patients" className="text-primary font-bold mb-6 inline-block hover:underline">
                ← Back to List
            </Link>
            <div className="bg-card p-8 rounded-3xl border border-border shadow-sm">
                <h1 className="text-3xl font-bold text-slate-900 mb-2">Add New Patient</h1>
                <p className="text-slate-500 mb-8">Register a new patient to your directory</p>

                <form onSubmit={handleSubmit} className="space-y-6">
                    <div className="grid grid-cols-2 gap-4">
                        <div className="space-y-2">
                            <label className="text-sm font-bold text-slate-700">Full Name</label>
                            <input
                                type="text"
                                required
                                value={formData.full_name}
                                onChange={(e) => setFormData({ ...formData, full_name: e.target.value })}
                                className="w-full px-4 py-2.5 rounded-xl border border-slate-200 outline-none focus:ring-2 focus:ring-primary transition"
                                placeholder="e.g. John Doe"
                            />
                        </div>
                        <div className="space-y-2">
                            <label className="text-sm font-bold text-slate-700">Phone Number</label>
                            <input
                                type="text"
                                required
                                value={formData.phone}
                                onChange={(e) => setFormData({ ...formData, phone: e.target.value })}
                                className="w-full px-4 py-2.5 rounded-xl border border-slate-200 outline-none focus:ring-2 focus:ring-primary transition"
                                placeholder="e.g. 01700-000000"
                            />
                        </div>
                    </div>

                    <div className="grid grid-cols-3 gap-4">
                        <div className="space-y-2">
                            <label className="text-sm font-bold text-slate-700">Age</label>
                            <input
                                type="number"
                                required
                                value={formData.age}
                                onChange={(e) => setFormData({ ...formData, age: e.target.value })}
                                className="w-full px-4 py-2.5 rounded-xl border border-slate-200 outline-none focus:ring-2 focus:ring-primary transition"
                            />
                        </div>
                        <div className="space-y-2">
                            <label className="text-sm font-bold text-slate-700">Gender</label>
                            <select
                                value={formData.gender}
                                onChange={(e) => setFormData({ ...formData, gender: e.target.value })}
                                className="w-full px-4 py-2.5 rounded-xl border border-slate-200 outline-none focus:ring-2 focus:ring-primary transition"
                            >
                                <option value="male">Male</option>
                                <option value="female">Female</option>
                                <option value="other">Other</option>
                            </select>
                        </div>
                        <div className="space-y-2">
                            <label className="text-sm font-bold text-slate-700">Blood Group</label>
                            <select
                                value={formData.blood_group}
                                onChange={(e) => setFormData({ ...formData, blood_group: e.target.value })}
                                className="w-full px-4 py-2.5 rounded-xl border border-slate-200 outline-none focus:ring-2 focus:ring-primary transition"
                            >
                                <option>A+</option><option>A-</option><option>B+</option><option>B-</option>
                                <option>O+</option><option>O-</option><option>AB+</option><option>AB-</option>
                            </select>
                        </div>
                    </div>

                    <button
                        type="submit"
                        disabled={isSubmitting}
                        className="w-full py-4 bg-primary text-white rounded-2xl font-bold medical-gradient shadow-lg hover:opacity-90 transition disabled:opacity-50"
                    >
                        {isSubmitting ? 'Saving...' : 'Save Patient Information'}
                    </button>
                </form>
            </div>
        </div>
    );
}
