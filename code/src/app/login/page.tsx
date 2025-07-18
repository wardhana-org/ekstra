"use client";

import { signIn } from "next-auth/react";

export default function LoginPage() {
  return (
    <div className="min-h-screen flex items-center justify-center p-4">
      <div className="max-w-md w-full border-3 border-white rounded-xl shadow-lg p-8">
        <h2 className="text-2xl font-bold text-white mb-6 text-center">
          Ekstra
        </h2>
        <div className="flex justify-center">
          <button
            onClick={() => signIn("google", { redirectTo: "/user/home" })}
            className="px-4 py-2 border-2 flex gap-2 border-slate-200 rounded-lg text-slate-200 hover:border-slate-500 hover:text-slate-300 hover:shadow transition duration-150 hover:cursor-pointer"
          >
            <img
              className="w-6 h-6"
              src="https://www.svgrepo.com/show/475656/google-color.svg"
              loading="lazy"
              alt="google logo"
            />
            <span>Login with Google</span>
          </button>
        </div>
      </div>
    </div>
  );
}
