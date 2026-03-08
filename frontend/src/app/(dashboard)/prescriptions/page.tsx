"use client";

import Link from "next/link";
import { useState, useEffect } from "react";
import { PaginatedPrescriptionResp, PrescriptionResp } from "@/types/prescription";
import { useToast } from "@/components/Toast";

export default function PrescriptionList() {
    const [prescriptions, setPrescriptions] = useState<PrescriptionResp[]>([]);
    const [isLoading, setIsLoading] = useState(true);
    const [searchTerm, setSearchTerm] = useState("");
    const [debouncedSearch, setDebouncedSearch] = useState("");
    const { error: errorToast } = useToast();

    useEffect(() => {
        const timer = setTimeout(() => {
            setDebouncedSearch(searchTerm);
        }, 500);
        return () => clearTimeout(timer);
    }, [searchTerm]);

    useEffect(() => {
        fetchPrescriptions();
    }, [debouncedSearch]);

    const fetchPrescriptions = async () => {
        const token = localStorage.getItem("docmate_token");
        if (!token) return;

        try {
            const url = new URL(`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8081'}/v1/prescriptions`);
            if (debouncedSearch) {
                url.searchParams.append("search", debouncedSearch);
            }

            const response = await fetch(url.toString(), {
                headers: { 'Authorization': `Bearer ${token}` }
            });
            const data = await response.json();
            if (data.success) {
                setPrescriptions(data.data.records);
            } else {
                errorToast(data.message || "Failed to fetch prescriptions");
            }
        } catch (error) {
            console.error("Error fetching prescriptions:", error);
            errorToast("An unexpected error occurred");
        } finally {
            setIsLoading(false);
        }
    };

    if (isLoading && !debouncedSearch) {
        return <div className="p-8 text-center text-slate-500">Loading prescriptions...</div>;
    }

    return (
        <div className="p-8">
            <div className="flex justify-between items-center mb-10">
                <div>
                    <h1 className="text-3xl font-bold text-slate-900 tracking-tight">Prescriptions</h1>
                    <p className="text-slate-500">History of all digital prescriptions generated</p>
                </div>
                <div className="flex gap-4">
                    <input
                        type="text"
                        placeholder="Search by ID or Patient..."
                        value={searchTerm}
                        onChange={(e) => setSearchTerm(e.target.value)}
                        className="px-4 py-2.5 w-64 rounded-xl border border-slate-200 focus:ring-2 focus:ring-primary outline-none transition"
                    />
                    <Link href="/prescriptions/new" className="bg-primary text-white px-6 py-2.5 rounded-xl font-bold medical-gradient shadow-lg">
                        + Create New
                    </Link>
                </div>
            </div>

            <div className="grid grid-cols-1 gap-4">
                {prescriptions.length === 0 ? (
                    <div className="text-center py-20 bg-slate-50 rounded-2xl border border-dashed border-slate-200">
                        <p className="text-slate-400 font-medium">No prescriptions found</p>
                    </div>
                ) : (
                    prescriptions.map((px) => (
                        <div key={px.id} className="bg-card p-6 rounded-2xl border border-border shadow-sm hover:shadow-md transition flex items-center justify-between group">
                            <div className="flex items-center gap-6">
                                <div className="w-12 h-12 rounded-xl bg-slate-50 flex items-center justify-center text-slate-400 group-hover:bg-blue-50 group-hover:text-primary transition">
                                    <span className="text-xl">📄</span>
                                </div>
                                <div>
                                    <div className="flex items-center gap-2">
                                        <span className="text-xs font-bold text-slate-400 uppercase tracking-widest">PR-{px.id}</span>
                                        <span className="text-slate-200">|</span>
                                        <span className="text-sm font-medium text-slate-500">{new Date(px.created_at).toLocaleDateString()}</span>
                                        {px.status === 'draft' ? (
                                            <span className="px-2 py-0.5 rounded-full text-[10px] font-bold bg-amber-50 text-amber-600 border border-amber-100 uppercase tracking-tighter">Draft</span>
                                        ) : (
                                            <span className="px-2 py-0.5 rounded-full text-[10px] font-bold bg-emerald-50 text-emerald-600 border border-emerald-100 uppercase tracking-tighter">Finalized</span>
                                        )}
                                    </div>
                                    <h3 className="text-lg font-bold text-slate-900">{px.patient_name}</h3>
                                    <p className="text-sm text-slate-600">{px.medications.length} medicines</p>
                                </div>
                            </div>

                            <div className="flex items-center gap-4">
                                {px.status === 'draft' ? (
                                    <Link
                                        href={`/prescriptions/${px.id}/edit`}
                                        className="px-5 py-2 rounded-xl text-primary font-bold hover:bg-blue-50 transition border border-transparent"
                                    >
                                        Edit Draft
                                    </Link>
                                ) : (
                                    <Link
                                        href={`/prescriptions/${px.id}/print`}
                                        className="px-5 py-2 rounded-xl text-slate-600 font-bold hover:bg-slate-50 transition border border-transparent hover:border-slate-100"
                                    >
                                        Preview
                                    </Link>
                                )}
                                <Link
                                    href={`/prescriptions/${px.id}/print`}
                                    className="px-5 py-2 rounded-xl bg-slate-900 text-white font-bold hover:bg-slate-800 transition shadow-sm"
                                >
                                    Print PDF
                                </Link>
                            </div>
                        </div>
                    ))
                )}
            </div>

            <div className="mt-8 text-center">
                <button className="text-primary font-bold hover:underline">View More History</button>
            </div>
        </div>
    );
}
