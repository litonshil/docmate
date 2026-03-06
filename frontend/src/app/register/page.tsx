"use client";

import { useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useToast } from "@/components/Toast";

export default function RegisterPage() {
    const router = useRouter();
    const [formData, setFormData] = useState({
        username: '', // This will map to user_name in backend
        email: '',
        password: '',
    });

    const { success: successToast, error: errorToast } = useToast();

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        console.log('Registration attempt (Users table):', formData);

        try {
            const response = await fetch('http://localhost:8081/v1/users/register', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    user_name: formData.username, // Using the field as full name but mapping to user_name per backend
                    email: formData.email,
                    password: formData.password,
                    role: 'doctor',
                }),
            });

            const data = await response.json();

            if (response.ok && data.success) {
                successToast('Registration successful! Please log in.');
                router.push('/login');
            } else {
                errorToast(data.message || 'Registration failed');
            }
        } catch (error) {
            console.error('Registration error:', error);
            errorToast('An error occurred during registration. Please ensure the backend is running.');
        }
    };

    return (
        <div className="min-h-screen flex items-center justify-center bg-background p-4 py-12">
            <div className="w-full max-w-md">
                <div className="text-center mb-10">
                    <div className="inline-block p-3 rounded-2xl bg-primary mb-4 medical-gradient">
                        <svg className="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z" />
                        </svg>
                    </div>
                    <h1 className="text-4xl font-extrabold text-foreground tracking-tight">Create Account</h1>
                    <p className="text-slate-500 mt-2">Join Doc-Mate community</p>
                </div>

                <div className="bg-card p-10 rounded-3xl shadow-2xl border border-border glass relative overflow-hidden">
                    <form onSubmit={handleSubmit} className="space-y-6">
                        <div>
                            <label className="block text-sm font-semibold text-slate-700 mb-2">Full Name</label>
                            <input
                                type="text"
                                required
                                className="w-full px-4 py-3 rounded-xl border border-slate-200 focus:ring-2 focus:ring-primary outline-none transition"
                                placeholder="Dr. John Doe"
                                value={formData.username}
                                onChange={(e) => setFormData({ ...formData, username: e.target.value })}
                            />
                        </div>
                        <div>
                            <label className="block text-sm font-semibold text-slate-700 mb-2">Email Address</label>
                            <input
                                type="email"
                                required
                                className="w-full px-4 py-3 rounded-xl border border-slate-200 focus:ring-2 focus:ring-primary outline-none transition"
                                placeholder="doctor@example.com"
                                value={formData.email}
                                onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                            />
                        </div>
                        <div>
                            <label className="block text-sm font-semibold text-slate-700 mb-2">Password</label>
                            <input
                                type="password"
                                required
                                className="w-full px-4 py-3 rounded-xl border border-slate-200 focus:ring-2 focus:ring-primary outline-none transition"
                                placeholder="••••••••"
                                value={formData.password}
                                onChange={(e) => setFormData({ ...formData, password: e.target.value })}
                            />
                        </div>

                        <button
                            type="submit"
                            className="w-full py-4 px-6 rounded-2xl bg-primary text-white font-bold medical-gradient shadow-lg hover:shadow-xl transition-all"
                        >
                            Register
                        </button>
                    </form>

                    <div className="mt-8 text-center border-t border-slate-100 pt-6">
                        <p className="text-slate-500">
                            Already have an account?{' '}
                            <Link href="/login" className="text-primary font-bold hover:underline">
                                Log In
                            </Link>
                        </p>
                    </div>
                </div>
            </div>
        </div>
    );
}
