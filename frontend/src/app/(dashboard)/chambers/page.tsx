"use client";

import { useEffect, useState, useCallback } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useToast } from "@/components/Toast";

interface VisitingSlot {
    start_time: string;
    end_time: string;
}

interface VisitingDay {
    day: string;
    slots: VisitingSlot[];
}

interface Chamber {
    id: number;
    doctor_id: number;
    name: string;
    address: string;
    phone: string;
    fee: number;
    visiting_hours: VisitingDay[];
    is_active: boolean;
}

export default function ChambersPage() {
    const router = useRouter();
    const [chambers, setChambers] = useState<Chamber[]>([]);
    const [loading, setLoading] = useState(true);
    const { success: successToast, error: errorToast } = useToast();

    const fetchChambers = useCallback(async () => {
        setLoading(true);
        try {
            const token = localStorage.getItem('docmate_token');
            if (!token) {
                router.push('/login');
                return;
            }

            const profileRes = await fetch('http://localhost:8081/v1/doctors/profile', {
                headers: { 'Authorization': `Bearer ${token}` }
            });
            const profileData = await profileRes.json();
            if (!profileRes.ok || !profileData.success) return;
            const doctorId = profileData.data.id;

            const chambersRes = await fetch(`http://localhost:8081/v1/doctors/${doctorId}/chambers`, {
                headers: { 'Authorization': `Bearer ${token}` }
            });
            const chambersData = await chambersRes.json();

            if (chambersRes.ok && chambersData.success) {
                setChambers(chambersData.data.records || []);
            }
        } catch (error) {
            console.error("Error fetching chambers:", error);
        } finally {
            setLoading(false);
        }
    }, [router]);

    const handleDelete = async (chamberId: number) => {
        if (!confirm('Are you sure you want to delete this chamber?')) return;

        try {
            const token = localStorage.getItem('docmate_token');
            const profileRes = await fetch('http://localhost:8081/v1/doctors/profile', {
                headers: { 'Authorization': `Bearer ${token}` }
            });
            const profileData = await profileRes.json();
            const doctorId = profileData.data.id;

            const res = await fetch(`http://localhost:8081/v1/doctors/${doctorId}/chambers/${chamberId}`, {
                method: 'DELETE',
                headers: { 'Authorization': `Bearer ${token}` }
            });

            if (res.ok) {
                successToast('Chamber deleted successfully');
                fetchChambers();
            } else {
                const data = await res.json();
                errorToast(data.message || 'Failed to delete chamber');
            }
        } catch (error) {
            console.error('Error deleting chamber:', error);
            errorToast('An error occurred');
        }
    };

    useEffect(() => {
        fetchChambers();
    }, [fetchChambers]);

    const renderSchedule = (visiting_hours: VisitingDay[]) => {
        if (!visiting_hours || visiting_hours.length === 0) return (
            <div className="text-slate-400 italic">No schedule set</div>
        );

        return (
            <div className="space-y-1">
                {visiting_hours.map((vh, i) => (
                    <div key={i} className="flex items-start gap-2">
                        <span className="font-bold text-slate-700 min-w-[32px]">{vh.day.substring(0, 3)}:</span>
                        <div className="flex flex-wrap gap-x-2 gap-y-0.5">
                            {vh.slots.map((slot, j) => (
                                <span key={j} className="bg-slate-50 px-2 py-0.5 rounded text-xs border border-slate-100">
                                    {slot.start_time} - {slot.end_time}
                                </span>
                            ))}
                        </div>
                    </div>
                ))}
            </div>
        );
    };

    return (
        <div className="p-8">
            <div className="flex justify-between items-center mb-10">
                <div>
                    <h1 className="text-3xl font-bold text-slate-900 tracking-tight">Chamber Management</h1>
                    <p className="text-slate-500">Manage your consultation locations and schedules</p>
                </div>
                <Link href="/chambers/new" className="bg-primary text-white px-6 py-2.5 rounded-xl font-bold medical-gradient shadow-lg">
                    + Add New Chamber
                </Link>
            </div>

            <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
                {loading ? (
                    <div className="lg:col-span-2 flex flex-col items-center justify-center py-20 gap-4">
                        <div className="w-12 h-12 border-4 border-primary border-t-transparent rounded-full animate-spin"></div>
                        <p className="text-slate-500 font-medium">Loading chambers...</p>
                    </div>
                ) : chambers.length === 0 ? (
                    <div className="lg:col-span-2 text-center py-20 bg-slate-50 rounded-3xl border-2 border-dashed border-slate-200">
                        <p className="text-slate-400 font-medium mb-4">No chambers found. Add your first practice location!</p>
                        <Link href="/chambers/new" className="text-primary font-bold hover:underline">+ Add New Chamber</Link>
                    </div>
                ) : (
                    chambers.map((c) => (
                        <div key={c.id} className="bg-card rounded-3xl border border-border shadow-sm p-8 hover:shadow-md transition group h-full flex flex-col">
                            <div className="flex justify-between items-start mb-6">
                                <div className="p-3 bg-blue-50 text-primary rounded-2xl group-hover:bg-primary group-hover:text-white transition">
                                    <span className="text-2xl">🏥</span>
                                </div>
                                <div className="flex gap-2">
                                    <Link href={`/chambers/${c.id}/edit`} className="p-2 text-slate-400 hover:text-slate-900 transition">Edit</Link>
                                    <button onClick={() => handleDelete(c.id)} className="p-2 text-slate-400 hover:text-red-500 transition">Delete</button>
                                </div>
                            </div>

                            <h2 className="text-xl font-bold text-slate-900 mb-2">{c.name}</h2>
                            <p className="text-slate-500 text-sm mb-6">{c.address}</p>

                            <div className="space-y-4 mt-auto">
                                <div className="flex items-center gap-3 text-sm font-medium text-slate-600">
                                    <span className="w-5 text-center">📞</span>
                                    {c.phone || 'No phone provided'}
                                </div>
                                <div className="flex items-center gap-3 text-sm font-medium text-slate-600">
                                    <span className="w-5 text-center">💰</span>
                                    Consultation Fee: <span className="text-slate-900 font-bold">{c.fee} BDT</span>
                                </div>
                                <div className="flex items-start gap-3 text-sm font-medium text-slate-600">
                                    <span className="w-5 text-center mt-1">📅</span>
                                    {renderSchedule(c.visiting_hours)}
                                </div>
                            </div>

                            <div className="mt-8 pt-6 border-t border-slate-50">
                                <span className={`inline-block px-3 py-1 rounded-full text-xs font-bold uppercase ${c.is_active ? 'bg-green-50 text-green-600' : 'bg-slate-50 text-slate-400'}`}>
                                    {c.is_active ? 'Active' : 'Inactive'}
                                </span>
                            </div>
                        </div>
                    ))
                )}

                {!loading && (
                    <Link href="/chambers/new" className="h-full min-h-[300px] border-2 border-dashed border-slate-200 rounded-3xl flex flex-col items-center justify-center gap-4 text-slate-400 hover:border-primary hover:text-primary transition group">
                        <span className="text-4xl group-hover:scale-110 transition">+</span>
                        <span className="font-bold">Add Another Location</span>
                    </Link>
                )}
            </div>
        </div>
    );
}
