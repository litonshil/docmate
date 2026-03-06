"use client";

import React, { useEffect, useState, useCallback } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";

interface Medicine {
    id: number;
    brand_name: string;
    generic_name: string;
    form: string;
    strength: string;
    manufacturer: string;
    description: string;
    is_active: bool;
}

interface Pagination {
    page: number;
    limit: number;
    total: number;
    last_page: number;
}

export default function MedicinesPage() {
    const router = useRouter();
    const [medicines, setMedicines] = useState<Medicine[]>([]);
    const [loading, setLoading] = useState(true);
    const [searchTerm, setSearchTerm] = useState("");
    const [pagination, setPagination] = useState<Pagination>({
        page: 1,
        limit: 10,
        total: 0,
        last_page: 1
    });

    const fetchMedicines = useCallback(async (page: number, search: string) => {
        setLoading(true);
        const token = localStorage.getItem('docmate_token');
        try {
            const response = await fetch(`http://localhost:8081/v1/medicines?page=${page}&limit=10&search=${search}`, {
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            });
            const data = await response.json();
            if (response.ok && data.success) {
                setMedicines(data.data.records || []);
                setPagination(data.data.pagination);
            }
        } catch (error) {
            console.error("Error fetching medicines:", error);
        } finally {
            setLoading(false);
        }
    }, []);

    useEffect(() => {
        const delayDebounceFn = setTimeout(() => {
            fetchMedicines(1, searchTerm);
        }, 300);

        return () => clearTimeout(delayDebounceFn);
    }, [searchTerm, fetchMedicines]);

    const handleDelete = async (id: number) => {
        if (!confirm("Are you sure you want to delete this medicine?")) return;

        const token = localStorage.getItem('docmate_token');
        try {
            const response = await fetch(`http://localhost:8081/v1/medicines/${id}`, {
                method: 'DELETE',
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            });
            const data = await response.json();
            if (response.ok && data.success) {
                alert("Medicine deleted successfully");
                fetchMedicines(pagination.page, searchTerm);
            } else {
                alert(data.message || "Failed to delete medicine");
            }
        } catch (error) {
            console.error("Error deleting medicine:", error);
            alert("An error occurred while deleting the medicine");
        }
    };

    return (
        <div className="p-8">
            <div className="flex justify-between items-center mb-10">
                <div>
                    <h1 className="text-3xl font-bold text-slate-900 tracking-tight">Medicine Library</h1>
                    <p className="text-slate-500">Maintain your frequently used medicine database</p>
                </div>
                <Link href="/medicines/new" className="bg-primary text-white px-6 py-2.5 rounded-xl font-bold medical-gradient shadow-lg hover:opacity-90 transition">
                    + Add New Medicine
                </Link>
            </div>

            <div className="bg-card rounded-2xl border border-border shadow-sm overflow-hidden">
                <div className="p-6 border-b border-border bg-slate-50/50">
                    <input
                        type="text"
                        placeholder="Search by brand or generic name..."
                        value={searchTerm}
                        onChange={(e) => setSearchTerm(e.target.value)}
                        className="w-full max-w-md px-4 py-2.5 rounded-xl border border-slate-200 focus:ring-2 focus:ring-primary outline-none transition bg-white"
                    />
                </div>
                <div className="overflow-x-auto">
                    <table className="w-full text-left">
                        <thead className="bg-slate-50 border-b border-border">
                            <tr>
                                <th className="px-8 py-4 font-bold text-slate-500 text-xs uppercase tracking-wider">Medicine Info</th>
                                <th className="px-8 py-4 font-bold text-slate-500 text-xs uppercase tracking-wider">Generic Name</th>
                                <th className="px-8 py-4 font-bold text-slate-500 text-xs uppercase tracking-wider">Details</th>
                                <th className="px-8 py-4 font-bold text-slate-500 text-xs uppercase tracking-wider">Manufacturer</th>
                                <th className="px-8 py-4 font-bold text-slate-500 text-xs uppercase tracking-wider text-right">Actions</th>
                            </tr>
                        </thead>
                        <tbody className="divide-y divide-slate-100">
                            {loading && medicines.length === 0 ? (
                                <tr>
                                    <td colSpan={5} className="px-8 py-10 text-center text-slate-400">Loading medicines...</td>
                                </tr>
                            ) : medicines.length === 0 ? (
                                <tr>
                                    <td colSpan={5} className="px-8 py-10 text-center text-slate-400">No medicines found.</td>
                                </tr>
                            ) : (
                                medicines.map((m) => (
                                    <tr key={m.id} className="hover:bg-slate-50/50 transition">
                                        <td className="px-8 py-5">
                                            <div className="font-bold text-slate-900">{m.brand_name}</div>
                                            <div className="text-xs text-slate-400 uppercase font-bold">{m.form}</div>
                                        </td>
                                        <td className="px-8 py-5 text-sm text-slate-600 font-medium">{m.generic_name}</td>
                                        <td className="px-8 py-5 text-sm">
                                            <div className="text-slate-700 font-bold">{m.strength}</div>
                                        </td>
                                        <td className="px-8 py-5 text-sm text-slate-500">{m.manufacturer}</td>
                                        <td className="px-8 py-5 text-right">
                                            <button
                                                onClick={() => router.push(`/medicines/${m.id}/edit`)}
                                                className="text-slate-400 hover:text-primary transition font-bold mr-4"
                                            >
                                                Edit
                                            </button>
                                            <button
                                                onClick={() => handleDelete(m.id)}
                                                className="text-slate-300 hover:text-red-500 transition"
                                            >
                                                Trash
                                            </button>
                                        </td>
                                    </tr>
                                ))
                            )}
                        </tbody>
                    </table>
                </div>

                {pagination.last_page > 1 && (
                    <div className="p-6 border-t border-border bg-slate-50/30 flex justify-between items-center text-sm">
                        <div className="text-slate-500">
                            Showing page {pagination.page} of {pagination.last_page} ({pagination.total} medicines)
                        </div>
                        <div className="flex gap-2">
                            <button
                                disabled={pagination.page === 1}
                                onClick={() => fetchMedicines(pagination.page - 1, searchTerm)}
                                className="px-4 py-2 rounded-lg border border-slate-200 bg-white hover:bg-slate-50 disabled:opacity-50 transition font-bold"
                            >
                                Previous
                            </button>
                            <button
                                disabled={pagination.page === pagination.last_page}
                                onClick={() => fetchMedicines(pagination.page + 1, searchTerm)}
                                className="px-4 py-2 rounded-lg border border-slate-200 bg-white hover:bg-slate-50 disabled:opacity-50 transition font-bold"
                            >
                                Next
                            </button>
                        </div>
                    </div>
                )}
            </div>
        </div>
    );
}
