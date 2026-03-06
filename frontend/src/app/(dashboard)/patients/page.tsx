"use client";

import { useEffect, useState, useCallback } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useToast } from "@/components/Toast";

interface Patient {
    id: number;
    full_name: string;
    gender: string;
    age: number;
    phone: string;
    email: string;
    blood_group: string;
    created_at: string;
}

export default function PatientList() {
    const router = useRouter();
    const [patients, setPatients] = useState<Patient[]>([]);
    const [loading, setLoading] = useState(true);
    const { error: errorToast } = useToast();
    const [searchTerm, setSearchTerm] = useState("");
    const [pagination, setPagination] = useState({
        page: 1,
        limit: 10,
        total: 0,
        last_page: 1
    });

    const fetchPatients = useCallback(async (page = 1, search = "") => {
        setLoading(true);
        try {
            const token = localStorage.getItem('docmate_token');
            if (!token) {
                router.push('/login');
                return;
            }

            const queryParams = new URLSearchParams({
                page: page.toString(),
                limit: pagination.limit.toString(),
                name: search
            });

            const response = await fetch(`http://localhost:8081/v1/patients?${queryParams}`, {
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            });

            const data = await response.json();
            if (response.ok && data.success) {
                setPatients(data.data.records || []);
                setPagination({
                    page: data.data.page || 1,
                    limit: data.data.limit || 10,
                    total: data.data.total || 0,
                    last_page: data.data.last_page || 1
                });
            } else {
                console.error("Failed to fetch patients:", data.message);
                errorToast(data.message || "Failed to fetch patients");
            }
        } catch (error) {
            console.error("Error fetching patients:", error);
            errorToast("An error occurred while fetching patients");
        } finally {
            setLoading(false);
        }
    }, [pagination?.limit || 10, router]);

    useEffect(() => {
        fetchPatients(1, searchTerm);
    }, [searchTerm, fetchPatients]);

    const handleSearch = (e: React.ChangeEvent<HTMLInputElement>) => {
        setSearchTerm(e.target.value);
    };

    const handlePageChange = (newPage: number) => {
        if (newPage >= 1 && newPage <= pagination.last_page) {
            fetchPatients(newPage, searchTerm);
        }
    };

    const formatDate = (dateString: string) => {
        return new Date(dateString).toLocaleDateString('en-US', {
            month: 'short',
            day: 'numeric',
            year: 'numeric'
        });
    };

    return (
        <div className="p-8">
            <div className="flex justify-between items-center mb-10">
                <div>
                    <h1 className="text-3xl font-bold text-slate-900 tracking-tight">Patient Directory</h1>
                    <p className="text-slate-500">Manage and view your patient records</p>
                </div>
                <div className="flex gap-4">
                    <div className="relative">
                        <input
                            type="text"
                            placeholder="Search by name..."
                            className="pl-11 pr-4 py-2.5 w-80 rounded-xl border border-slate-200 focus:ring-2 focus:ring-primary outline-none transition shadow-sm"
                            value={searchTerm}
                            onChange={handleSearch}
                        />
                        <span className="absolute left-4 top-1/2 -translate-y-1/2 text-slate-400">🔍</span>
                    </div>
                    <Link href="/patients/new" className="bg-primary text-white px-6 py-2.5 rounded-xl font-bold medical-gradient shadow-lg">
                        + Add New Patient
                    </Link>
                </div>
            </div>

            <div className="bg-card rounded-2xl border border-border shadow-sm overflow-hidden">
                <table className="w-full text-left">
                    <thead className="bg-slate-50">
                        <tr>
                            <th className="px-8 py-4 font-bold text-slate-500 text-xs uppercase tracking-wider">Patient Details</th>
                            <th className="px-8 py-4 font-bold text-slate-500 text-xs uppercase tracking-wider">Contact</th>
                            <th className="px-8 py-4 font-bold text-slate-500 text-xs uppercase tracking-wider">Medical Summary</th>
                            <th className="px-8 py-4 font-bold text-slate-500 text-xs uppercase tracking-wider">Registered</th>
                            <th className="px-8 py-4 font-bold text-slate-500 text-xs uppercase tracking-wider">Actions</th>
                        </tr>
                    </thead>
                    <tbody className="divide-y divide-slate-100">
                        {loading ? (
                            <tr>
                                <td colSpan={5} className="px-8 py-20 text-center text-slate-400">
                                    <div className="flex flex-col items-center gap-2">
                                        <div className="w-8 h-8 border-4 border-primary border-t-transparent rounded-full animate-spin"></div>
                                        <span>Loading patients...</span>
                                    </div>
                                </td>
                            </tr>
                        ) : patients.length === 0 ? (
                            <tr>
                                <td colSpan={5} className="px-8 py-20 text-center text-slate-400">
                                    No patients found.
                                </td>
                            </tr>
                        ) : (
                            patients.map((p) => (
                                <tr key={p.id} className="hover:bg-slate-50/50 transition cursor-pointer">
                                    <td className="px-8 py-5">
                                        <div className="flex items-center gap-4">
                                            <div className="w-10 h-10 rounded-full bg-slate-100 flex items-center justify-center font-bold text-slate-400">
                                                {p.full_name.charAt(0)}
                                            </div>
                                            <div>
                                                <div className="font-bold text-slate-900">{p.full_name}</div>
                                                <div className="text-xs text-slate-500">{p.gender} • {p.age} years</div>
                                            </div>
                                        </div>
                                    </td>
                                    <td className="px-8 py-5">
                                        <div className="text-sm font-medium text-slate-700">{p.phone}</div>
                                        <div className="text-xs text-slate-400">Patient ID: #{p.id}</div>
                                    </td>
                                    <td className="px-8 py-5 text-sm">
                                        <div className="flex gap-2">
                                            <span className="px-2 py-0.5 rounded-lg bg-red-50 text-red-600 text-[10px] font-bold uppercase">{p.blood_group}</span>
                                            {/* Medical history preview or status placeholder */}
                                            <span className="px-2 py-0.5 rounded-lg bg-blue-50 text-blue-600 text-[10px] font-bold uppercase tracking-tight">Active</span>
                                        </div>
                                    </td>
                                    <td className="px-8 py-5 text-sm text-slate-500">{formatDate(p.created_at)}</td>
                                    <td className="px-8 py-5">
                                        <div className="flex items-center gap-3">
                                            <Link href="/prescriptions/new" className="text-primary font-bold text-sm tracking-tight hover:underline">Prescribe</Link>
                                            <span className="text-slate-200">|</span>
                                            <Link href={`/patients/${p.id}`} className="text-slate-400 hover:text-slate-900 transition tracking-tight">Profile</Link>
                                        </div>
                                    </td>
                                </tr>
                            ))
                        )}
                    </tbody>
                </table>
                <div className="p-6 border-t border-border flex justify-between items-center text-sm text-slate-500">
                    <span>
                        Showing {(pagination.page - 1) * pagination.limit + 1} to {Math.min(pagination.page * pagination.limit, pagination.total)} of {pagination.total} patients
                    </span>
                    <div className="flex gap-2">
                        <button
                            disabled={pagination.page === 1}
                            onClick={() => handlePageChange(pagination.page - 1)}
                            className="px-4 py-2 rounded-lg border border-border hover:bg-slate-50 transition disabled:opacity-50 disabled:cursor-not-allowed"
                        >
                            Previous
                        </button>
                        <button
                            disabled={pagination.page === pagination.last_page}
                            onClick={() => handlePageChange(pagination.page + 1)}
                            className="px-4 py-2 rounded-lg bg-primary text-white medical-gradient shadow-sm disabled:opacity-50 disabled:cursor-not-allowed"
                        >
                            Next
                        </button>
                    </div>
                </div>
            </div>
        </div>
    );
}
