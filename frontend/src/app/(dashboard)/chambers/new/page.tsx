"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { useToast } from "@/components/Toast";

interface VisitingSlot {
    start_time: string;
    end_time: string;
}

interface VisitingDay {
    day: string;
    slots: VisitingSlot[];
}

export default function NewChamberPage() {
    const router = useRouter();
    const [doctorProfile, setDoctorProfile] = useState<any>(null);
    const [formData, setFormData] = useState({
        name: '',
        address: '',
        area: '',
        city: 'Dhaka',
        country: 'Bangladesh',
        phone: '',
        email: '',
        fee: '',
        follow_up_fee: '',
    });
    const [visitingHours, setVisitingHours] = useState<VisitingDay[]>([]);
    const { success: successToast, error: errorToast } = useToast();
    const [isSubmitting, setIsSubmitting] = useState(false);

    const days = ['Saturday', 'Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday'];

    const toggleDay = (day: string) => {
        setVisitingHours(prev => {
            const exists = prev.find(vh => vh.day === day);
            if (exists) {
                return prev.filter(vh => vh.day !== day);
            } else {
                return [...prev, { day, slots: [{ start_time: '09:00', end_time: '13:00' }] }];
            }
        });
    };

    const addSlot = (day: string) => {
        setVisitingHours(prev => prev.map(vh =>
            vh.day === day
                ? { ...vh, slots: [...vh.slots, { start_time: '17:00', end_time: '21:00' }] }
                : vh
        ));
    };

    const removeSlot = (day: string, index: number) => {
        setVisitingHours(prev => prev.map(vh =>
            vh.day === day
                ? { ...vh, slots: vh.slots.filter((_, i) => i !== index) }
                : vh
        ));
    };

    const updateSlot = (day: string, index: number, field: keyof VisitingSlot, value: string) => {
        setVisitingHours(prev => prev.map(vh =>
            vh.day === day
                ? {
                    ...vh,
                    slots: vh.slots.map((slot, i) => i === index ? { ...slot, [field]: value } : slot)
                }
                : vh
        ));
    };

    useEffect(() => {
        const fetchProfile = async () => {
            const token = localStorage.getItem('docmate_token');
            if (!token) return;

            try {
                const response = await fetch('http://localhost:8081/v1/doctors/profile', {
                    headers: { 'Authorization': `Bearer ${token}` }
                });
                const data = await response.json();
                if (response.ok && data.success) {
                    setDoctorProfile(data.data);
                }
            } catch (error) {
                console.error('Error fetching doctor profile:', error);
                errorToast('Failed to load doctor profile');
            }
        };
        fetchProfile();
    }, []);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!doctorProfile) {
            errorToast('Doctor profile not loaded. Please try again.');
            return;
        }

        if (visitingHours.length === 0) {
            errorToast('Please select at least one day for your schedule.');
            return;
        }

        setIsSubmitting(true);
        const token = localStorage.getItem('docmate_token');

        try {
            const response = await fetch(`http://localhost:8081/v1/doctors/${doctorProfile.id}/chambers`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`,
                },
                body: JSON.stringify({
                    name: formData.name,
                    address: formData.address,
                    area: formData.area,
                    city: formData.city,
                    country: formData.country,
                    phone: formData.phone,
                    email: formData.email || undefined,
                    fee: parseFloat(formData.fee),
                    visiting_hours: visitingHours
                }),
            });

            const data = await response.json();

            if (response.ok && data.success) {
                successToast('Chamber registered successfully!');
                router.push('/chambers');
            } else {
                errorToast(data.message || 'Failed to register chamber');
            }
        } catch (error) {
            console.error('Error creating chamber:', error);
            errorToast('An error occurred while saving chamber information.');
        } finally {
            setIsSubmitting(false);
        }
    };

    return (
        <div className="p-8 max-w-2xl mx-auto">
            <Link href="/chambers" className="text-primary font-bold mb-6 inline-block hover:underline">
                ← Back to Management
            </Link>
            <div className="bg-card p-8 rounded-3xl border border-border shadow-sm">
                <h1 className="text-3xl font-bold text-slate-900 mb-2">Add New Chamber</h1>
                <p className="text-slate-500 mb-8">Setup a new consultation location and schedule</p>

                <form onSubmit={handleSubmit} className="space-y-8">
                    <section className="space-y-4">
                        <h3 className="text-sm font-bold text-slate-400 uppercase tracking-widest">Basic Information</h3>
                        <div className="space-y-4">
                            <div className="space-y-2">
                                <label className="text-sm font-bold text-slate-700">Chamber Name</label>
                                <input
                                    type="text"
                                    required
                                    value={formData.name}
                                    onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                                    className="w-full px-4 py-2.5 rounded-xl border border-slate-200 outline-none focus:ring-2 focus:ring-primary transition"
                                    placeholder="e.g. City Medical Center"
                                />
                            </div>
                            <div className="space-y-2">
                                <label className="text-sm font-bold text-slate-700">Address</label>
                                <textarea
                                    required
                                    value={formData.address}
                                    onChange={(e) => setFormData({ ...formData, address: e.target.value })}
                                    className="w-full px-4 py-2.5 rounded-xl border border-slate-200 outline-none focus:ring-2 focus:ring-primary transition h-24"
                                    placeholder="Full address..."
                                ></textarea>
                            </div>
                            <div className="grid grid-cols-2 gap-4">
                                <div className="space-y-2">
                                    <label className="text-sm font-bold text-slate-700">Consultation Fee (BDT)</label>
                                    <input
                                        type="number"
                                        required
                                        value={formData.fee}
                                        onChange={(e) => setFormData({ ...formData, fee: e.target.value })}
                                        className="w-full px-4 py-2.5 rounded-xl border border-slate-200 outline-none focus:ring-2 focus:ring-primary transition"
                                        placeholder="e.g. 1000"
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
                                        placeholder="e.g. 01800-000000"
                                    />
                                </div>
                            </div>
                        </div>
                    </section>

                    <section className="space-y-4">
                        <h3 className="text-sm font-bold text-slate-400 uppercase tracking-widest">Weekly Schedule</h3>
                        <div className="space-y-6">
                            <div className="flex flex-wrap gap-2">
                                {days.map(day => (
                                    <button
                                        key={day}
                                        type="button"
                                        onClick={() => toggleDay(day)}
                                        className={`px-4 py-2 rounded-xl border transition text-sm font-bold ${visitingHours.find(vh => vh.day === day)
                                            ? 'bg-primary text-white border-primary shadow-md'
                                            : 'border-slate-200 text-slate-500 hover:border-primary hover:text-primary hover:bg-slate-50'
                                            }`}
                                    >
                                        {day.substring(0, 3)}
                                    </button>
                                ))}
                            </div>

                            <div className="space-y-4">
                                {visitingHours.length === 0 && (
                                    <p className="text-slate-400 text-sm italic py-4 text-center border-2 border-dashed border-slate-100 rounded-2xl">
                                        Select practice days above to add time slots
                                    </p>
                                )}
                                {visitingHours.map((vh) => (
                                    <div key={vh.day} className="p-5 bg-slate-50 rounded-2xl border border-slate-100 space-y-4">
                                        <div className="flex justify-between items-center">
                                            <span className="font-bold text-slate-900">{vh.day}</span>
                                            <button
                                                type="button"
                                                onClick={() => addSlot(vh.day)}
                                                className="text-xs font-bold text-primary hover:underline"
                                            >
                                                + Add Slot
                                            </button>
                                        </div>
                                        <div className="space-y-3">
                                            {vh.slots.map((slot, index) => (
                                                <div key={index} className="flex items-center gap-3">
                                                    <input
                                                        type="time"
                                                        value={slot.start_time}
                                                        onChange={(e) => updateSlot(vh.day, index, 'start_time', e.target.value)}
                                                        className="flex-1 px-3 py-1.5 rounded-lg border border-slate-200 text-sm focus:ring-2 focus:ring-primary outline-none"
                                                    />
                                                    <span className="text-slate-400 text-xs">to</span>
                                                    <input
                                                        type="time"
                                                        value={slot.end_time}
                                                        onChange={(e) => updateSlot(vh.day, index, 'end_time', e.target.value)}
                                                        className="flex-1 px-3 py-1.5 rounded-lg border border-slate-200 text-sm focus:ring-2 focus:ring-primary outline-none"
                                                    />
                                                    {vh.slots.length > 1 && (
                                                        <button
                                                            type="button"
                                                            onClick={() => removeSlot(vh.day, index)}
                                                            className="text-slate-400 hover:text-red-500 p-1"
                                                        >
                                                            ✕
                                                        </button>
                                                    )}
                                                </div>
                                            ))}
                                        </div>
                                    </div>
                                ))}
                            </div>
                        </div>
                    </section>

                    <button
                        type="submit"
                        disabled={isSubmitting || !doctorProfile}
                        className="w-full py-4 bg-primary text-white rounded-2xl font-bold medical-gradient shadow-lg hover:opacity-90 transition disabled:opacity-50"
                    >
                        {isSubmitting ? 'Registering...' : doctorProfile ? 'Register Chamber' : 'Loading Profile...'}
                    </button>
                </form>
            </div>
        </div>
    );
}
