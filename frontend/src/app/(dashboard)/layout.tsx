"use client";

import React, { useEffect, useState } from "react";
import { useRouter, usePathname } from "next/navigation";
import Link from "next/link";

export default function DashboardLayout({
    children,
}: {
    children: React.ReactNode;
}) {
    const router = useRouter();
    const pathname = usePathname();
    const [isLoading, setIsLoading] = useState(true);

    useEffect(() => {
        // Check for token in localStorage
        const token = typeof window !== 'undefined' ? localStorage.getItem('docmate_token') : null;
        if (!token) {
            router.push('/login');
        } else {
            setIsLoading(false);
        }
    }, [router]);

    const handleLogout = () => {
        localStorage.removeItem('docmate_token');
        router.push('/login');
    };

    const navItems = [
        { name: 'Dashboard', icon: '🏠', href: '/' },
        { name: 'Patients', icon: '👤', href: '/patients' },
        { name: 'Prescriptions', icon: '📋', href: '/prescriptions' },
        { name: 'Medicines', icon: '💊', href: '/medicines' },
        { name: 'Chambers', icon: '🏥', href: '/chambers' },
    ];

    const settingsItems = [
        { name: 'Profile', href: '/settings' },
        { name: 'Prescription', href: '/settings/prescription' },
    ];

    const isSettingsActive = pathname.startsWith('/settings');
    const [isSettingsOpen, setIsSettingsOpen] = useState(isSettingsActive);

    // Update isSettingsOpen when pathname changes to ensure it's open when active
    useEffect(() => {
        if (isSettingsActive) {
            setIsSettingsOpen(true);
        }
    }, [isSettingsActive]);

    if (isLoading) {
        return (
            <div className="min-h-screen flex items-center justify-center bg-background">
                <div className="flex flex-col items-center gap-4">
                    <div className="w-12 h-12 border-4 border-primary border-t-transparent rounded-full animate-spin"></div>
                    <p className="text-slate-500 font-medium">Authenticating...</p>
                </div>
            </div>
        );
    }

    return (
        <div className="min-h-screen flex bg-background">
            {/* Sidebar */}
            <aside className="w-64 bg-card border-r border-border flex flex-col hidden md:flex sticky top-0 h-screen">
                <div className="p-6 border-b border-border flex items-center gap-3">
                    <div className="w-8 h-8 rounded-lg bg-primary medical-gradient flex items-center justify-center text-white">
                        <span className="font-bold">D</span>
                    </div>
                    <span className="text-xl font-extrabold text-slate-900 tracking-tight">Doc-Mate</span>
                </div>

                <nav className="flex-1 p-4 space-y-1 mt-4 overflow-y-auto">
                    {navItems.map((item) => (
                        <Link
                            key={item.href}
                            href={item.href}
                            className={`flex items-center gap-3 px-4 py-3 rounded-xl transition font-medium ${pathname === item.href
                                ? 'bg-slate-50 text-primary'
                                : 'text-slate-600 hover:bg-slate-50 hover:text-primary'
                                }`}
                        >
                            <span>{item.icon}</span>
                            {item.name}
                        </Link>
                    ))}

                    <div className="space-y-1">
                        <button
                            onClick={() => setIsSettingsOpen(!isSettingsOpen)}
                            className={`w-full flex items-center justify-between px-4 py-3 rounded-xl font-medium transition ${isSettingsActive ? 'text-primary' : 'text-slate-600'
                                } hover:bg-slate-50`}
                        >
                            <div className="flex items-center gap-3">
                                <span>⚙️</span>
                                Settings
                            </div>
                            <svg
                                className={`w-4 h-4 transition-transform duration-200 ${isSettingsOpen ? 'rotate-180' : ''}`}
                                fill="none"
                                stroke="currentColor"
                                viewBox="0 0 24 24"
                                xmlns="http://www.w3.org/2000/svg"
                            >
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M19 9l-7 7-7-7"></path>
                            </svg>
                        </button>

                        {isSettingsOpen && (
                            <div className="pl-11 space-y-1 animate-in fade-in slide-in-from-top-1 duration-200">
                                {settingsItems.map((child) => (
                                    <Link
                                        key={child.href}
                                        href={child.href}
                                        className={`block py-2 text-sm transition font-medium ${pathname === child.href
                                            ? 'text-primary'
                                            : 'text-slate-400 hover:text-primary'
                                            }`}
                                    >
                                        {child.name}
                                    </Link>
                                ))}
                            </div>
                        )}
                    </div>
                </nav>

                <div className="p-6 border-t border-border mt-auto">
                    <div className="bg-slate-50 p-4 rounded-2xl flex items-center gap-3">
                        <div className="w-10 h-10 rounded-full bg-slate-200 mask-avatar"></div>
                        <div>
                            <div className="text-sm font-bold text-slate-900">Doctor</div>
                            <div className="text-xs text-slate-500">Medical Professional</div>
                        </div>
                    </div>
                </div>
            </aside>

            {/* Main Content */}
            <div className="flex-1 flex flex-col">
                {/* Top Header */}
                <header className="h-16 bg-card border-b border-border flex items-center justify-between px-8 sticky top-0 z-10 glass">
                    <div className="text-slate-400 text-sm">March 6, 2026 • 12:12 PM</div>
                    <div className="flex items-center gap-4">
                        <button className="p-2 text-slate-400 hover:text-primary transition">
                            <span className="text-xl">🔔</span>
                        </button>
                        <div className="h-8 w-px bg-border mx-2"></div>
                        <button
                            onClick={handleLogout}
                            className="text-sm font-semibold text-slate-600 hover:text-red-500 transition"
                        >
                            Logout
                        </button>
                    </div>
                </header>

                <main className="flex-1 overflow-y-auto">
                    {children}
                </main>
            </div>
        </div>
    );
}
