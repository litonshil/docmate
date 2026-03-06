import Link from "next/link";

export default function Dashboard() {
    return (
        <div className="p-8">
            <div className="flex justify-between items-center mb-8">
                <div>
                    <h1 className="text-3xl font-bold text-slate-900 tracking-tight">Dashboard Overview</h1>
                    <p className="text-slate-500">Welcome back, Dr. Smith</p>
                </div>
                <Link href="/prescriptions/new" className="bg-primary text-white px-6 py-2 rounded-xl font-semibold medical-gradient shadow-lg">
                    + New Prescription
                </Link>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
                {[
                    { label: 'Total Patients', value: '1,284', icon: '👤', trend: '+12% this month' },
                    { label: 'Today Visits', value: '24', icon: '📅', trend: '+4 since morning' },
                    { label: 'Prescriptions', value: '458', icon: '📋', trend: '+25% this week' },
                    { label: 'Medicines', value: '82', icon: '💊', trend: 'Active list' },
                ].map((stat, i) => (
                    <div key={i} className="bg-card p-6 rounded-2xl border border-border shadow-sm hover:shadow-md transition">
                        <div className="flex justify-between items-start mb-4">
                            <span className="text-2xl">{stat.icon}</span>
                            <span className="text-xs font-semibold px-2 py-1 bg-green-50 text-green-600 rounded-full">
                                {stat.trend}
                            </span>
                        </div>
                        <h3 className="text-slate-500 text-sm font-medium">{stat.label}</h3>
                        <p className="text-2xl font-bold text-slate-900 mt-1">{stat.value}</p>
                    </div>
                ))}
            </div>

            <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
                <div className="lg:col-span-2 bg-card rounded-2xl border border-border shadow-sm p-6">
                    <h2 className="text-xl font-bold text-slate-900 mb-6">Recent Patients</h2>
                    <div className="overflow-x-auto">
                        <table className="w-full text-left">
                            <thead>
                                <tr className="border-b border-slate-100">
                                    <th className="pb-4 font-semibold text-slate-600 text-sm">PATIENT</th>
                                    <th className="pb-4 font-semibold text-slate-600 text-sm">GENDER</th>
                                    <th className="pb-4 font-semibold text-slate-600 text-sm">AGE</th>
                                    <th className="pb-4 font-semibold text-slate-600 text-sm">LAST VISIT</th>
                                    <th className="pb-4 font-semibold text-slate-600 text-sm">ACTION</th>
                                </tr>
                            </thead>
                            <tbody className="divide-y divide-slate-50">
                                {[
                                    { name: 'John Doe', gender: 'Male', age: '45', lastVisit: '2 hours ago' },
                                    { name: 'Sarah Wilson', gender: 'Female', age: '28', lastVisit: 'Yesterday' },
                                    { name: 'Michael Chen', gender: 'Male', age: '52', lastVisit: '2 days ago' },
                                    { name: 'Emma Brown', gender: 'Female', age: '34', lastVisit: '3 days ago' },
                                ].map((patient, i) => (
                                    <tr key={i} className="group hover:bg-slate-50 transition">
                                        <td className="py-4">
                                            <div className="font-semibold text-slate-900">{patient.name}</div>
                                        </td>
                                        <td className="py-4 text-slate-600">{patient.gender}</td>
                                        <td className="py-4 text-slate-600">{patient.age}</td>
                                        <td className="py-4 text-slate-600 text-sm">{patient.lastVisit}</td>
                                        <td className="py-4">
                                            <Link href="/patients/102" className="text-primary font-medium hover:underline">View History</Link>
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                </div>

                <div className="bg-card rounded-2xl border border-border shadow-sm p-6">
                    <h2 className="text-xl font-bold text-slate-900 mb-6">Today's Schedule</h2>
                    <div className="space-y-6">
                        {[
                            { time: '09:00 AM', patient: 'Robert Fox', type: 'Checkup' },
                            { time: '10:30 AM', patient: 'Jane Cooper', type: 'Consultation' },
                            { time: '12:00 PM', patient: 'Cody Fisher', type: 'Follow-up' },
                            { time: '02:30 PM', patient: 'Esther Howard', type: 'Special' },
                        ].map((slot, i) => (
                            <div key={i} className="flex items-center gap-4">
                                <div className="w-20 text-xs font-bold text-slate-400 uppercase tracking-wider">{slot.time}</div>
                                <div className="flex-1 p-3 rounded-xl bg-slate-50 border border-slate-100">
                                    <div className="font-semibold text-slate-900 text-sm">{slot.patient}</div>
                                    <div className="text-xs text-slate-500">{slot.type}</div>
                                </div>
                            </div>
                        ))}
                        <Link href="/patients" className="w-full py-3 rounded-xl border-2 border-dashed border-slate-200 text-slate-400 font-medium hover:border-primary hover:text-primary transition flex justify-center">
                            View All Schedule
                        </Link>
                    </div>
                </div>
            </div>
        </div>
    );
}
