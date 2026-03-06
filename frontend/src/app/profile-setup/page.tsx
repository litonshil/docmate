"use client";

import { useState } from 'react';
import { useRouter } from 'next/navigation';

export default function ProfileSetupPage() {
    const router = useRouter();
    const [formData, setFormData] = useState({
        fullName: '',
        degree: '',
        specialization: '',
        phone: '',
        bio: '',
    });

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        console.log('Profile Setup (Doctors table):', formData);

        try {
            const token = typeof window !== 'undefined' ? localStorage.getItem('docmate_token') : '';
            if (!token) {
                alert('Authentication token missing. Please log in again.');
                router.push('/login');
                return;
            }

            const response = await fetch('http://localhost:8081/v1/doctors', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`,
                },
                body: JSON.stringify({
                    full_name: formData.fullName,
                    degree: [formData.degree],
                    specialization: [formData.specialization],
                    phone: formData.phone,
                    bio: formData.bio,
                    signature_url: '', // Default empty for now
                }),
            });

            const data = await response.json();

            if (response.ok && data.success) {
                // Simulation: Persist profile completion for local UI logic
                if (typeof window !== 'undefined') {
                    localStorage.setItem('doctorProfileComplete', 'true');
                }
                alert('Profile setup successful!');
                router.push('/');
            } else {
                alert(data.message || 'Profile setup failed');
            }
        } catch (error) {
            console.error('Profile setup error:', error);
            alert('An error occurred. Please ensure the backend is running.');
        }
    };

    const specializations = [
        'General Practice', 'Cardiology', 'Neurology', 'Gastroenterology',
        'Pulmonology', 'Pediatrics', 'Gynecology', 'Dermatology', 'ENT', 'Dentistry'
    ];

    return (
        <div className="min-h-screen bg-slate-50 flex items-center justify-center p-6 py-12">
            <div className="w-full max-w-2xl bg-white rounded-[2rem] shadow-xl border border-slate-100 overflow-hidden glass">
                <div className="bg-primary p-12 text-center medical-gradient">
                    <h1 className="text-3xl font-bold text-white mb-2">Complete Your Profile</h1>
                    <p className="text-blue-100 italic">Tell us more about your professional medical background.</p>
                </div>

                <form onSubmit={handleSubmit} className="p-10 space-y-8">
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                        <div className="md:col-span-2">
                            <label className="block text-xs font-bold text-slate-400 uppercase tracking-widest mb-2">Full Name</label>
                            <input
                                type="text"
                                required
                                className="w-full px-5 py-3 rounded-2xl border border-slate-200 focus:ring-2 focus:ring-primary outline-none transition font-medium"
                                placeholder="Dr. Faisal Ahmed"
                                value={formData.fullName}
                                onChange={(e) => setFormData({ ...formData, fullName: e.target.value })}
                            />
                        </div>

                        <div>
                            <label className="block text-xs font-bold text-slate-400 uppercase tracking-widest mb-2">Primary Degree</label>
                            <select
                                required
                                className="w-full px-5 py-3 rounded-2xl border border-slate-200 focus:ring-2 focus:ring-primary outline-none bg-white transition font-medium"
                                value={formData.degree}
                                onChange={(e) => setFormData({ ...formData, degree: e.target.value })}
                            >
                                <option value="">Select Degree</option>
                                <option value="MBBS">MBBS</option>
                                <option value="BDS">BDS</option>
                                <option value="MD">MD</option>
                                <option value="FCPS">FCPS</option>
                            </select>
                        </div>

                        <div>
                            <label className="block text-xs font-bold text-slate-400 uppercase tracking-widest mb-2">Specialization</label>
                            <select
                                required
                                className="w-full px-5 py-3 rounded-2xl border border-slate-200 focus:ring-2 focus:ring-primary outline-none bg-white transition font-medium"
                                value={formData.specialization}
                                onChange={(e) => setFormData({ ...formData, specialization: e.target.value })}
                            >
                                <option value="">Select Specialization</option>
                                {specializations.map(s => (
                                    <option key={s} value={s}>{s}</option>
                                ))}
                                <option value="other">Other</option>
                            </select>
                        </div>

                        <div>
                            <label className="block text-xs font-bold text-slate-400 uppercase tracking-widest mb-2">Phone Number</label>
                            <input
                                type="tel"
                                required
                                className="w-full px-5 py-3 rounded-2xl border border-slate-200 focus:ring-2 focus:ring-primary outline-none transition font-medium"
                                placeholder="+880 1XXX XXXXXX"
                                value={formData.phone}
                                onChange={(e) => setFormData({ ...formData, phone: e.target.value })}
                            />
                        </div>

                        <div className="md:col-span-2">
                            <label className="block text-xs font-bold text-slate-400 uppercase tracking-widest mb-2">Professional Bio</label>
                            <textarea
                                className="w-full px-5 py-3 rounded-2xl border border-slate-200 focus:ring-2 focus:ring-primary outline-none transition h-32 resize-none"
                                placeholder="Tell your patients about your experience..."
                                value={formData.bio}
                                onChange={(e) => setFormData({ ...formData, bio: e.target.value })}
                            ></textarea>
                        </div>
                    </div>

                    <button
                        type="submit"
                        className="w-full py-4 rounded-2xl bg-primary text-white font-bold medical-gradient shadow-lg hover:shadow-xl transition-all text-lg"
                    >
                        Finish Setup & Go to Dashboard
                    </button>
                </form>
            </div>
        </div>
    );
}
