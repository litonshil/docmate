"use client";

import { useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useToast } from "@/components/Toast";

export default function LoginPage() {
    const router = useRouter();
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');

    const { success: successToast, error: errorToast } = useToast();

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        console.log('Login attempt:', { email, password });

        try {
            const response = await fetch('http://localhost:8081/v1/users/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ email, password }),
            });

            const data = await response.json();

            if (response.ok && data.success) {
                // Store token in localStorage
                if (typeof window !== 'undefined') {
                    localStorage.setItem('docmate_token', data.data.token);
                    successToast("Logged in successfully!");

                    const isProfileCompleted = data.data.user.is_profile_completed;

                    if (isProfileCompleted) {
                        router.push('/');
                    } else {
                        router.push('/profile-setup');
                    }
                }
            } else {
                errorToast(data.message || 'Login failed');
            }
        } catch (error) {
            console.error('Login error:', error);
            errorToast('An error occurred during login. Please ensure the backend is running.');
        }
    };

    return (
        <div className="min-h-screen flex items-center justify-center bg-background p-4">
            <div className="w-full max-w-md">
                <div className="text-center mb-10">
                    <div className="inline-block p-3 rounded-2xl bg-primary mb-4 medical-gradient">
                        <svg className="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19.428 15.428a2 2 0 00-1.022-.547l-2.387-.477a6 6 0 00-3.86.517l-.318.158a6 6 0 01-3.86.517L6.05 15.21a2 2 0 00-1.022.547l-2.387-.477a6 6 0 00-3.86.517l-.318.158a6 6 0 01-3.86.517L6.05 15.21a2 2 0 00-1.022.547l" />
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h.01M15 12h.01M9 16h.01M15 16h.01M12 12h.01M12 16h.01" />
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                        </svg>
                    </div>
                    <h1 className="text-4xl font-bold text-foreground">Doc-Mate</h1>
                    <p className="text-slate-500 mt-2">Welcome back, Doctor.</p>
                </div>

                <div className="bg-card p-8 rounded-2xl shadow-xl border border-border glass">
                    <form onSubmit={handleSubmit} className="space-y-6">
                        <div>
                            <label className="block text-sm font-medium text-slate-700 mb-2">Email Address</label>
                            <input
                                type="email"
                                required
                                className="w-full px-4 py-3 rounded-xl border border-slate-200 focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent transition"
                                placeholder="doctor@example.com"
                                value={email}
                                onChange={(e) => setEmail(e.target.value)}
                            />
                        </div>
                        <div>
                            <div className="flex justify-between items-center mb-2">
                                <label className="block text-sm font-medium text-slate-700">Password</label>
                                <a href="#" className="text-sm text-primary hover:underline">Forgot password?</a>
                            </div>
                            <input
                                type="password"
                                required
                                className="w-full px-4 py-3 rounded-xl border border-slate-200 focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent transition"
                                placeholder="••••••••"
                                value={password}
                                onChange={(e) => setPassword(e.target.value)}
                            />
                        </div>
                        <button
                            type="submit"
                            className="w-full py-3 px-4 rounded-xl bg-primary text-white font-semibold hover:bg-blue-700 transition medical-gradient shadow-lg"
                        >
                            Log In
                        </button>
                    </form>

                    <div className="mt-8 text-center">
                        <p className="text-slate-500">
                            Don't have an account?{' '}
                            <Link href="/register" className="text-primary font-semibold hover:underline">
                                Register as Doctor
                            </Link>
                        </p>
                    </div>
                </div>
            </div>
        </div>
    );
}
