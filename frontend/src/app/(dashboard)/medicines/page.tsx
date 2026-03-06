import Link from "next/link";

export default function MedicinesPage() {
    const medicines = [
        { name: 'Napa 500mg', generic: 'Paracetamol', form: 'Tablet', manufacturer: 'Beximco' },
        { name: 'Seclo 20mg', generic: 'Omeprazole', form: 'Capsule', manufacturer: 'Square' },
        { name: 'Fexo 120mg', generic: 'Fexofenadine', form: 'Tablet', manufacturer: 'ACG' },
        { name: 'Alatrol', generic: 'Cetirizine', form: 'Syrup', manufacturer: 'Square' },
        { name: 'Azithrocin 500', generic: 'Azithromycin', form: 'Tablet', manufacturer: 'Beximco' },
    ];

    return (
        <div className="p-8">
            <div className="flex justify-between items-center mb-10">
                <div>
                    <h1 className="text-3xl font-bold text-slate-900 tracking-tight">Medicine Library</h1>
                    <p className="text-slate-500">Maintain your frequently used medicine database</p>
                </div>
                <Link href="/medicines/new" className="bg-primary text-white px-6 py-2.5 rounded-xl font-bold medical-gradient shadow-lg">
                    + Add New Medicine
                </Link>
            </div>

            <div className="bg-card rounded-2xl border border-border shadow-sm overflow-hidden">
                <div className="p-6 border-b border-border bg-slate-50/50">
                    <input
                        type="text"
                        placeholder="Search by brand or generic name..."
                        className="w-full max-w-md px-4 py-2.5 rounded-xl border border-slate-200 focus:ring-2 focus:ring-primary outline-none transition bg-white"
                    />
                </div>
                <table className="w-full text-left">
                    <thead className="bg-slate-50 border-b border-border">
                        <tr>
                            <th className="px-8 py-4 font-bold text-slate-500 text-xs uppercase tracking-wider">Brand Name</th>
                            <th className="px-8 py-4 font-bold text-slate-500 text-xs uppercase tracking-wider">Generic Name</th>
                            <th className="px-8 py-4 font-bold text-slate-500 text-xs uppercase tracking-wider">Form</th>
                            <th className="px-8 py-4 font-bold text-slate-500 text-xs uppercase tracking-wider">Manufacturer</th>
                            <th className="px-8 py-4 font-bold text-slate-500 text-xs uppercase tracking-wider text-right">Actions</th>
                        </tr>
                    </thead>
                    <tbody className="divide-y divide-slate-100">
                        {medicines.map((m, i) => (
                            <tr key={i} className="hover:bg-slate-50/50 transition">
                                <td className="px-8 py-5">
                                    <div className="font-bold text-slate-900">{m.name}</div>
                                </td>
                                <td className="px-8 py-5 text-sm text-slate-600 font-medium">{m.generic}</td>
                                <td className="px-8 py-5 text-sm">
                                    <span className="px-3 py-1 rounded-full bg-slate-100 text-slate-600 text-xs font-bold uppercase">{m.form}</span>
                                </td>
                                <td className="px-8 py-5 text-sm text-slate-500">{m.manufacturer}</td>
                                <td className="px-8 py-5 text-right">
                                    <button className="text-slate-400 hover:text-primary transition font-bold mr-4">Edit</button>
                                    <button className="text-slate-300 hover:text-red-500 transition">Trash</button>
                                </td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>
        </div>
    );
}
