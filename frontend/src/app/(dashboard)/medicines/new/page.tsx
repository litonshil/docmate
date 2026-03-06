import Link from "next/link";

export default function NewMedicinePage() {
    return (
        <div className="p-8 max-w-2xl mx-auto">
            <Link href="/medicines" className="text-primary font-bold mb-6 inline-block hover:underline">
                ← Back to Library
            </Link>
            <div className="bg-card p-8 rounded-3xl border border-border shadow-sm">
                <h1 className="text-3xl font-bold text-slate-900 mb-2">Add New Medicine</h1>
                <p className="text-slate-500 mb-8">Add a new medicine to your library for quick prescribing</p>

                <form className="space-y-6">
                    <div className="space-y-2">
                        <label className="text-sm font-bold text-slate-700">Brand Name</label>
                        <input type="text" className="w-full px-4 py-2.5 rounded-xl border border-slate-200 outline-none focus:ring-2 focus:ring-primary transition" placeholder="e.g. Napa" />
                    </div>
                    <div className="space-y-2">
                        <label className="text-sm font-bold text-slate-700">Generic Name</label>
                        <input type="text" className="w-full px-4 py-2.5 rounded-xl border border-slate-200 outline-none focus:ring-2 focus:ring-primary transition" placeholder="e.g. Paracetamol" />
                    </div>

                    <div className="grid grid-cols-2 gap-4">
                        <div className="space-y-2">
                            <label className="text-sm font-bold text-slate-700">Form</label>
                            <select className="w-full px-4 py-2.5 rounded-xl border border-slate-200 outline-none focus:ring-2 focus:ring-primary transition">
                                <option>Tablet</option>
                                <option>Capsule</option>
                                <option>Syrup</option>
                                <option>Injection</option>
                                <option>Ointment</option>
                            </select>
                        </div>
                        <div className="space-y-2">
                            <label className="text-sm font-bold text-slate-700">Strength</label>
                            <input type="text" className="w-full px-4 py-2.5 rounded-xl border border-slate-200 outline-none focus:ring-2 focus:ring-primary transition" placeholder="e.g. 500mg" />
                        </div>
                    </div>

                    <button type="button" className="w-full py-4 bg-primary text-white rounded-2xl font-bold medical-gradient shadow-lg hover:opacity-90 transition">
                        Add to Library
                    </button>
                </form>
            </div>
        </div>
    );
}
