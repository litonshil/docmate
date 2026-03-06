import Link from "next/link";

export default function PrescriptionList() {
    const prescriptions = [
        { id: 'PR-7821', patient: 'John Doe', date: 'Mar 05, 2026', diagnosis: 'Hypertension', items: 3 },
        { id: 'PR-7815', patient: 'Sarah Wilson', date: 'Mar 04, 2026', diagnosis: 'Migraine', items: 2 },
        { id: 'PR-7802', patient: 'Michael Chen', date: 'Mar 02, 2026', diagnosis: 'Type 2 Diabetes', items: 5 },
        { id: 'PR-7798', patient: 'Emma Brown', date: 'Mar 01, 2026', diagnosis: 'Common Cold', items: 1 },
    ];

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
                        className="px-4 py-2.5 w-64 rounded-xl border border-slate-200 focus:ring-2 focus:ring-primary outline-none transition"
                    />
                    <Link href="/prescriptions/new" className="bg-primary text-white px-6 py-2.5 rounded-xl font-bold medical-gradient shadow-lg">
                        + Create New
                    </Link>
                </div>
            </div>

            <div className="grid grid-cols-1 gap-4">
                {prescriptions.map((px) => (
                    <div key={px.id} className="bg-card p-6 rounded-2xl border border-border shadow-sm hover:shadow-md transition flex items-center justify-between group">
                        <div className="flex items-center gap-6">
                            <div className="w-12 h-12 rounded-xl bg-slate-50 flex items-center justify-center text-slate-400 group-hover:bg-blue-50 group-hover:text-primary transition">
                                <span className="text-xl">📄</span>
                            </div>
                            <div>
                                <div className="flex items-center gap-2">
                                    <span className="text-xs font-bold text-slate-400 uppercase tracking-widest">{px.id}</span>
                                    <span className="text-slate-200">|</span>
                                    <span className="text-sm font-medium text-slate-500">{px.date}</span>
                                </div>
                                <h3 className="text-lg font-bold text-slate-900">{px.patient}</h3>
                                <p className="text-sm text-slate-600">Diagnosis: <span className="font-semibold text-primary">{px.diagnosis}</span> • {px.items} medicines</p>
                            </div>
                        </div>

                        <div className="flex items-center gap-4">
                            <button className="px-5 py-2 rounded-xl text-slate-600 font-bold hover:bg-slate-50 transition border border-transparent hover:border-slate-100">
                                Preview
                            </button>
                            <button className="px-5 py-2 rounded-xl bg-slate-900 text-white font-bold hover:bg-slate-800 transition shadow-sm">
                                Print PDF
                            </button>
                        </div>
                    </div>
                ))}
            </div>

            <div className="mt-8 text-center">
                <button className="text-primary font-bold hover:underline">View More History</button>
            </div>
        </div>
    );
}
