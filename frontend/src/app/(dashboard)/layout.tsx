"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";

export default function DashboardLayout({
    children,
}: {
    children: React.ReactNode;
}) {
    const router = useRouter();
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

                <nav className="flex-1 p-4 space-y-2 mt-4">
                    {[
                        { name: 'Dashboard', icon: '🏠', href: '/' },
                        { name: 'Patients', icon: '👤', href: '/patients' },
                        { name: 'Prescriptions', icon: '📋', href: '/prescriptions' },
                        { name: 'Medicines', icon: '💊', href: '/medicines' },
                        { name: 'Chambers', icon: '🏥', href: '/chambers' },
                        { name: 'Settings', icon: '⚙️', href: '/settings' },
                    ].map((item) => (
                        <Link
                            key={item.href}
                            href={item.href}
                            className="flex items-center gap-3 px-4 py-3 rounded-xl text-slate-600 hover:bg-slate-50 hover:text-primary transition font-medium"
                        >
                            <span>{item.icon}</span>
                            {item.name}
                        </Link>
                    ))}
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
