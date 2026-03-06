"use client";

import React, { useState, useEffect, useRef } from 'react';

interface Medicine {
    id: number;
    brand_name: string;
    generic_name: string;
    form: string;
    strength: string;
}

interface MedicineAutocompleteProps {
    value: string;
    onChange: (value: string) => void;
    placeholder?: string;
}

export default function MedicineAutocomplete({ value, onChange, placeholder }: MedicineAutocompleteProps) {
    const [query, setQuery] = useState(value);
    const [results, setResults] = useState<Medicine[]>([]);
    const [isOpen, setIsOpen] = useState(false);
    const [loading, setLoading] = useState(false);
    const wrapperRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        setQuery(value);
    }, [value]);

    useEffect(() => {
        const handleClickOutside = (event: MouseEvent) => {
            if (wrapperRef.current && !wrapperRef.current.contains(event.target as Node)) {
                setIsOpen(false);
            }
        };
        document.addEventListener('mousedown', handleClickOutside);
        return () => document.removeEventListener('mousedown', handleClickOutside);
    }, []);

    useEffect(() => {
        const fetchMedicines = async () => {
            if (!query.trim()) {
                setResults([]);
                return;
            }
            setLoading(true);
            const token = localStorage.getItem('docmate_token');
            try {
                const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8081'}/v1/medicines?page=1&limit=10&search=${encodeURIComponent(query)}`, {
                    headers: { 'Authorization': `Bearer ${token}` }
                });
                const data = await res.json();
                if (data.success) {
                    setResults(data.data.records || []);
                }
            } catch (error) {
                console.error("Failed to fetch medicines", error);
            } finally {
                setLoading(false);
            }
        };

        const timer = setTimeout(() => {
            if (isOpen) {
                fetchMedicines();
            }
        }, 300);

        return () => clearTimeout(timer);
    }, [query, isOpen]);

    const handleSelect = (med: Medicine) => {
        // Construct the full display string for the prescription
        const formPrefix = med.form ? `${med.form}. ` : '';
        const fullName = `${formPrefix}${med.brand_name} ${med.strength}`.trim();

        setQuery(fullName);
        onChange(fullName);
        setIsOpen(false);
    };

    return (
        <div ref={wrapperRef} className="relative w-full">
            <input
                type="text"
                value={query}
                onChange={(e) => {
                    setQuery(e.target.value);
                    onChange(e.target.value); // Allow free text input anytime
                    setIsOpen(true);
                }}
                onFocus={() => {
                    if (query.trim()) setIsOpen(true);
                }}
                className="w-full px-4 py-2 rounded-lg border border-slate-200 focus:ring-2 focus:ring-primary outline-none bg-white font-medium"
                placeholder={placeholder || "Search medicine..."}
            />
            {isOpen && query.trim() && (
                <div className="absolute z-[100] w-full mt-1 bg-white border border-slate-200 rounded-lg shadow-lg max-h-60 overflow-y-auto">
                    {loading ? (
                        <div className="px-4 py-3 text-sm text-slate-500 font-medium">Searching...</div>
                    ) : results.length > 0 ? (
                        <ul className="py-1">
                            {results.map((med) => (
                                <li
                                    key={med.id}
                                    onClick={() => handleSelect(med)}
                                    className="px-4 py-2 cursor-pointer hover:bg-slate-50 border-b border-slate-50/50 last:border-0 transition"
                                >
                                    <div className="font-bold text-slate-900 text-[14px]">
                                        {med.brand_name} <span className="text-slate-500 font-normal">{med.strength}</span>
                                    </div>
                                    <div className="text-[12px] text-slate-400 mt-0.5">
                                        {med.generic_name} • {med.form}
                                    </div>
                                </li>
                            ))}
                        </ul>
                    ) : (
                        <div className="px-4 py-3 text-xs text-slate-500 font-medium">
                            No medicines found. You can still use the entered name.
                        </div>
                    )}
                </div>
            )}
        </div>
    );
}
