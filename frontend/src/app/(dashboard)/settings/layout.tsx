"use client";

import { usePathname } from "next/navigation";

export default function SettingsLayout({
    children,
}: {
    children: React.ReactNode;
}) {
    const pathname = usePathname();

    const isPrescription = pathname === '/settings/prescription';

    const title = isPrescription ? "Prescription Settings" : "Doctor Profile";
    const description = isPrescription
        ? "Customize your prescription header, footer, and templates for different chambers."
        : "Update your professional information, qualifications, and digital signature.";

    return (
        <div className="p-8">
            <div className="mb-8">
                <h1 className="text-3xl font-bold text-slate-900 tracking-tight">{title}</h1>
                <p className="text-slate-500 font-medium">{description}</p>
            </div>

            <div className="max-w-4xl">
                {children}
            </div>
        </div>
    );
}
