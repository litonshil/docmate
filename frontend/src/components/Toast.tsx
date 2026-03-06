"use client";

import React, { createContext, useContext, useState, useCallback, ReactNode } from "react";

type ToastType = "success" | "error" | "info" | "warning";

interface Toast {
    id: number;
    message: string;
    type: ToastType;
}

interface ToastContextType {
    toast: (message: string, type?: ToastType) => void;
    success: (message: string) => void;
    error: (message: string) => void;
}

const ToastContext = createContext<ToastContextType | undefined>(undefined);

export const useToast = () => {
    const context = useContext(ToastContext);
    if (!context) {
        throw new Error("useToast must be used within a ToastProvider");
    }
    return context;
};

export const ToastProvider = ({ children }: { children: ReactNode }) => {
    const [toasts, setToasts] = useState<Toast[]>([]);

    const removeToast = useCallback((id: number) => {
        setToasts((prev) => prev.filter((t) => t.id !== id));
    }, []);

    const toast = useCallback((message: string, type: ToastType = "info") => {
        const id = Date.now();
        setToasts((prev) => [...prev, { id, message, type }]);
        setTimeout(() => removeToast(id), 5000);
    }, [removeToast]);

    const success = useCallback((message: string) => toast(message, "success"), [toast]);
    const error = useCallback((message: string) => toast(message, "error"), [toast]);

    return (
        <ToastContext.Provider value={{ toast, success, error }}>
            {children}
            <div className="fixed bottom-5 right-5 z-50 flex flex-col gap-3 pointer-events-none">
                {toasts.map((t) => (
                    <div
                        key={t.id}
                        className={`pointer-events-auto px-6 py-4 rounded-2xl shadow-xl text-white font-bold animate-slide-in flex items-center gap-3 transition-all ${t.type === "success" ? "bg-emerald-500" :
                                t.type === "error" ? "bg-rose-500" :
                                    t.type === "warning" ? "bg-amber-500" :
                                        "bg-slate-800"
                            }`}
                    >
                        <span>
                            {t.type === "success" && "✅"}
                            {t.type === "error" && "❌"}
                            {t.type === "warning" && "⚠️"}
                            {t.type === "info" && "ℹ️"}
                        </span>
                        {t.message}
                        <button
                            onClick={() => removeToast(t.id)}
                            className="ml-4 opacity-70 hover:opacity-100 transition"
                        >
                            ✕
                        </button>
                    </div>
                ))}
            </div>
        </ToastContext.Provider>
    );
};
