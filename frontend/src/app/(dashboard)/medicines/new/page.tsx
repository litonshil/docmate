"use client";

import React, { useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useToast } from "@/components/Toast";

export default function NewMedicinePage() {
    const router = useRouter();
    const { success: successToast, error: errorToast } = useToast();
    const [isSubmitting, setIsSubmitting] = useState(false);
    const [formData, setFormData] = useState({
        brand_name: '',
        generic_name: '',
        form: 'tablet',
        strength: '',
        manufacturer: '',
        description: ''
    });

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setIsSubmitting(true);
        const token = localStorage.getItem('docmate_token');

        try {
            const response = await fetch('http://localhost:8081/v1/medicines', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify(formData)
            });

            const data = await response.json();
            if (response.ok && data.success) {
                successToast("Medicine added successfully");
                router.push('/medicines');
            } else {
                errorToast(data.message || "Failed to add medicine");
            }
        } catch (error) {
            console.error("Error adding medicine:", error);
            errorToast("An error occurred while adding the medicine");
        } finally {
            setIsSubmitting(false);
        }
    };

    return (
        <div className="p-8 max-w-2xl mx-auto">
            <Link href="/medicines" className="text-primary font-bold mb-6 inline-block hover:underline">
                ← Back to Library
            </Link>
            <div className="bg-card p-8 rounded-3xl border border-border shadow-sm">
                <h1 className="text-3xl font-bold text-slate-900 mb-2">Add New Medicine</h1>
                <p className="text-slate-500 mb-8">Add a new medicine to your library for quick prescribing</p>

                <form onSubmit={handleSubmit} className="space-y-6">
                    <div className="space-y-2">
                        <label className="text-sm font-bold text-slate-700">Brand Name</label>
                        <input
                            type="text"
                            required
                            value={formData.brand_name}
                            onChange={(e) => setFormData({ ...formData, brand_name: e.target.value })}
                            className="w-full px-4 py-2.5 rounded-xl border border-slate-200 outline-none focus:ring-2 focus:ring-primary transition"
                            placeholder="e.g. Napa"
                        />
                    </div>
                    <div className="space-y-2">
                        <label className="text-sm font-bold text-slate-700">Generic Name</label>
                        <input
                            type="text"
                            required
                            value={formData.generic_name}
                            onChange={(e) => setFormData({ ...formData, generic_name: e.target.value })}
                            className="w-full px-4 py-2.5 rounded-xl border border-slate-200 outline-none focus:ring-2 focus:ring-primary transition"
                            placeholder="e.g. Paracetamol"
                        />
                    </div>

                    <div className="grid grid-cols-2 gap-4">
                        <div className="space-y-2">
                            <label className="text-sm font-bold text-slate-700">Form</label>
                            <select
                                value={formData.form}
                                onChange={(e) => setFormData({ ...formData, form: e.target.value })}
                                className="w-full px-4 py-2.5 rounded-xl border border-slate-200 outline-none focus:ring-2 focus:ring-primary transition"
                            >
                                <option value="tablet">Tablet</option>
                                <option value="capsule">Capsule</option>
                                <option value="syrup">Syrup</option>
                                <option value="suspension">Suspension</option>
                                <option value="injection">Injection</option>
                                <option value="inhaler">Inhaler</option>
                                <option value="drops">Drops</option>
                                <option value="cream">Cream</option>
                                <option value="ointment">Ointment</option>
                                <option value="gel">Gel</option>
                                <option value="patch">Patch</option>
                                <option value="suppository">Suppository</option>
                                <option value="powder">Powder</option>
                                <option value="sachet">Sachet</option>
                                <option value="other">Other</option>
                            </select>
                        </div>
                        <div className="space-y-2">
                            <label className="text-sm font-bold text-slate-700">Strength</label>
                            <input
                                type="text"
                                value={formData.strength}
                                onChange={(e) => setFormData({ ...formData, strength: e.target.value })}
                                className="w-full px-4 py-2.5 rounded-xl border border-slate-200 outline-none focus:ring-2 focus:ring-primary transition"
                                placeholder="e.g. 500mg"
                            />
                        </div>
                    </div>

                    <div className="space-y-2">
                        <label className="text-sm font-bold text-slate-700">Manufacturer</label>
                        <input
                            type="text"
                            value={formData.manufacturer}
                            onChange={(e) => setFormData({ ...formData, manufacturer: e.target.value })}
                            className="w-full px-4 py-2.5 rounded-xl border border-slate-200 outline-none focus:ring-2 focus:ring-primary transition"
                            placeholder="e.g. Beximco Pharmaceuticals"
                        />
                    </div>

                    <div className="space-y-2">
                        <label className="text-sm font-bold text-slate-700">Description</label>
                        <textarea
                            value={formData.description}
                            onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                            className="w-full px-4 py-2.5 rounded-xl border border-slate-200 outline-none focus:ring-2 focus:ring-primary transition h-24"
                            placeholder="Additional details about the medicine..."
                        ></textarea>
                    </div>

                    <button
                        type="submit"
                        disabled={isSubmitting}
                        className="w-full py-4 bg-primary text-white rounded-2xl font-bold medical-gradient shadow-lg hover:opacity-90 transition disabled:opacity-50"
                    >
                        {isSubmitting ? "Adding..." : "Add to Library"}
                    </button>
                </form>
            </div>
        </div>
    );
}
