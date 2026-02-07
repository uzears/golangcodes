"use client";

import { useEffect, useMemo, useState } from "react";

const API_BASE =
  process.env.NEXT_PUBLIC_API_BASE_URL?.replace(/\/$/, "") ||
  "http://localhost:8080";

type Mode = "login" | "register";

type Profile = {
  id: string;
  email: string;
};

export default function Home() {
  const [mode, setMode] = useState<Mode>("login");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [status, setStatus] = useState("");
  const [token, setToken] = useState("");
  const [profile, setProfile] = useState<Profile | null>(null);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    const saved = window.localStorage.getItem("research_token");
    if (saved) {
      setToken(saved);
    }
  }, []);

  useEffect(() => {
    if (token) {
      window.localStorage.setItem("research_token", token);
    }
  }, [token]);

  const endpoint = useMemo(() => {
    return mode === "login" ? "/auth/login" : "/auth/register";
  }, [mode]);

  async function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();
    setStatus("");
    setProfile(null);
    setLoading(true);

    try {
      const res = await fetch(`${API_BASE}${endpoint}`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ email, password }),
      });

      const data = await res.json().catch(() => ({}));
      if (!res.ok) {
        setStatus(data?.error || "Request failed");
        return;
      }

      if (mode === "login") {
        setToken(data?.token || "");
        setStatus("Logged in. Token saved.");
      } else {
        setStatus(`Registered. Your user id is ${data?.id ?? "unknown"}.`);
      }
    } catch (err) {
      setStatus("Network error. Is the API running?");
    } finally {
      setLoading(false);
    }
  }

  async function handleMe() {
    setStatus("");
    setProfile(null);

    if (!token) {
      setStatus("No token available. Login first.");
      return;
    }

    try {
      const res = await fetch(`${API_BASE}/me`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      const data = await res.json().catch(() => ({}));
      if (!res.ok) {
        setStatus(data?.error || "Failed to fetch profile");
        return;
      }

      setProfile({ id: data.id, email: data.email });
      setStatus("Profile loaded.");
    } catch {
      setStatus("Network error. Is the API running?");
    }
  }

  function clearToken() {
    window.localStorage.removeItem("research_token");
    setToken("");
    setProfile(null);
    setStatus("Token cleared.");
  }

  if (token) {
    return (
      <div className="min-h-screen bg-slate-950 text-slate-100">
        <div className="relative isolate overflow-hidden">
          <div className="pointer-events-none absolute inset-0 -z-10">
            <div className="absolute -top-32 right-0 h-96 w-96 rounded-full bg-[radial-gradient(circle_at_top,_rgba(56,189,248,0.25),_transparent_60%)]" />
            <div className="absolute -bottom-48 left-0 h-[36rem] w-[36rem] rounded-full bg-[radial-gradient(circle_at_top,_rgba(16,185,129,0.2),_transparent_60%)]" />
            <div className="absolute inset-0 bg-[linear-gradient(120deg,_rgba(15,23,42,0.75),_rgba(2,6,23,0.95))]" />
            <div className="absolute inset-0 opacity-[0.12] mix-blend-soft-light [background-image:radial-gradient(#ffffff_1px,transparent_1px)] [background-size:20px_20px]" />
          </div>

          <div className="mx-auto flex min-h-screen w-full max-w-6xl flex-col gap-10 px-6 py-16">
            <header className="flex flex-wrap items-center justify-between gap-4">
              <div>
                <p className="text-xs uppercase tracking-[0.2em] text-emerald-200">
                  Dashboard
                </p>
                <h1 className="font-display text-3xl text-white md:text-4xl">
                  Research APP Console
                </h1>
              </div>
              <button
                type="button"
                onClick={clearToken}
                className="rounded-full border border-emerald-300/60 px-4 py-2 text-xs text-emerald-100 transition hover:border-emerald-200 hover:text-white"
              >
                Logout
              </button>
            </header>

            <div className="grid gap-6 md:grid-cols-[1.4fr_1fr]">
              <section className="rounded-3xl border border-emerald-400/30 bg-emerald-400/10 p-6 shadow-[0_30px_80px_-50px_rgba(15,23,42,0.9)] backdrop-blur">
                <h2 className="font-display text-2xl text-white">
                  Session Token
                </h2>
                <p className="mt-2 text-sm text-emerald-100/80">
                  Stored locally. Use it to access protected endpoints.
                </p>
                <div className="mt-4 rounded-2xl border border-emerald-400/30 bg-slate-950/40 p-4 text-xs text-emerald-100">
                  <p className="break-all">{token}</p>
                </div>
              </section>

              <section className="rounded-3xl border border-slate-800/80 bg-slate-900/70 p-6 shadow-[0_30px_80px_-50px_rgba(15,23,42,0.9)] backdrop-blur">
                <h2 className="font-display text-2xl text-white">
                  User Profile
                </h2>
                <p className="mt-2 text-sm text-slate-300">
                  Fetch your current user details from `/me`.
                </p>
                <button
                  type="button"
                  onClick={handleMe}
                  className="mt-4 w-full rounded-xl border border-emerald-400/60 px-4 py-3 text-sm font-semibold text-emerald-200 transition hover:border-emerald-300 hover:text-emerald-100"
                >
                  Fetch /me
                </button>
                {status && (
                  <div className="mt-4 rounded-xl border border-slate-800 bg-slate-950/40 px-4 py-3 text-xs text-slate-200">
                    {status}
                  </div>
                )}
                {profile && (
                  <div className="mt-4 rounded-xl border border-emerald-400/40 bg-emerald-400/10 px-4 py-3 text-xs text-emerald-100">
                    <div>ID: {profile.id}</div>
                    <div>Email: {profile.email}</div>
                  </div>
                )}
              </section>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-slate-950 text-slate-100">
      <div className="relative isolate overflow-hidden">
        <div className="pointer-events-none absolute inset-0 -z-10">
          <div className="absolute -top-32 right-0 h-96 w-96 rounded-full bg-[radial-gradient(circle_at_top,_rgba(56,189,248,0.25),_transparent_60%)]" />
          <div className="absolute -bottom-48 left-0 h-[36rem] w-[36rem] rounded-full bg-[radial-gradient(circle_at_top,_rgba(16,185,129,0.2),_transparent_60%)]" />
          <div className="absolute inset-0 bg-[linear-gradient(120deg,_rgba(15,23,42,0.75),_rgba(2,6,23,0.95))]" />
          <div className="absolute inset-0 opacity-[0.12] mix-blend-soft-light [background-image:radial-gradient(#ffffff_1px,transparent_1px)] [background-size:20px_20px]" />
        </div>

        <div className="mx-auto flex min-h-screen w-full max-w-6xl flex-col gap-10 px-6 py-16 md:flex-row md:items-center md:gap-14">
          <section className="flex-1 space-y-6">
            <div className="inline-flex items-center gap-3 rounded-full border border-emerald-400/30 bg-emerald-400/10 px-4 py-1 text-xs uppercase tracking-[0.2em] text-emerald-200">
              Research API
              <span className="h-1.5 w-1.5 rounded-full bg-emerald-300" />
              JWT Auth Lab
            </div>
            <h1 className="font-display text-4xl leading-tight text-slate-50 md:text-6xl">
              Stocks research at your fingertips by experts.
            </h1>
            <p className="max-w-xl text-base text-slate-300 md:text-lg">
              Subscribe to our newsletter to get the latest updates and insights on stock.
              Guidance and research delivered straight to your inbox, helping you make informed investment decisions.
            </p>

            <div className="grid gap-4 rounded-3xl border border-slate-800/80 bg-slate-900/60 p-6 backdrop-blur">
              <div className="flex flex-col gap-2 text-sm text-slate-300">
                <span className="text-xs uppercase tracking-[0.2em] text-slate-400">
                  API Base URL
                </span>
                <div className="rounded-xl border border-slate-800 bg-slate-950/60 px-4 py-3 font-mono text-xs text-slate-200">
                  {API_BASE}
                </div>
              </div>
              <div className="flex flex-wrap gap-3 text-xs text-slate-400">
                <span className="rounded-full border border-slate-800 px-3 py-1">
                  POST /auth/register
                </span>
                <span className="rounded-full border border-slate-800 px-3 py-1">
                  POST /auth/login
                </span>
                <span className="rounded-full border border-slate-800 px-3 py-1">
                  GET /me
                </span>
              </div>
            </div>
          </section>

          <section className="w-full max-w-md">
            {!token ? (
              <div className="rounded-3xl border border-slate-800/80 bg-slate-900/70 p-6 shadow-[0_30px_80px_-50px_rgba(15,23,42,0.9)] backdrop-blur">
                <div className="flex items-center justify-between gap-3">
                  <h2 className="font-display text-2xl text-white">Auth Console</h2>
                  <div className="flex rounded-full border border-slate-800 bg-slate-950/40 p-1 text-xs">
                    {(["login", "register"] as Mode[]).map((item) => (
                      <button
                        key={item}
                        type="button"
                        onClick={() => setMode(item)}
                        className={`rounded-full px-4 py-1.5 capitalize transition ${
                          mode === item
                            ? "bg-emerald-400 text-slate-950"
                            : "text-slate-400 hover:text-slate-200"
                        }`}
                      >
                        {item}
                      </button>
                    ))}
                  </div>
                </div>

                <form className="mt-6 space-y-4" onSubmit={handleSubmit}>
                  <label className="block text-sm text-slate-300">
                    Email
                    <input
                      type="email"
                      required
                      value={email}
                      onChange={(e) => setEmail(e.target.value)}
                      placeholder="you@company.com"
                      className="mt-2 w-full rounded-xl border border-slate-800 bg-slate-950/70 px-4 py-3 text-sm text-slate-100 outline-none ring-emerald-400/40 focus:ring-2"
                    />
                  </label>
                  <label className="block text-sm text-slate-300">
                    Password
                    <input
                      type="password"
                      required
                      value={password}
                      onChange={(e) => setPassword(e.target.value)}
                      placeholder="••••••••"
                      className="mt-2 w-full rounded-xl border border-slate-800 bg-slate-950/70 px-4 py-3 text-sm text-slate-100 outline-none ring-emerald-400/40 focus:ring-2"
                    />
                  </label>
                  <button
                    type="submit"
                    disabled={loading}
                    className="w-full rounded-xl bg-emerald-400 px-4 py-3 text-sm font-semibold text-slate-950 transition hover:bg-emerald-300 disabled:cursor-not-allowed disabled:opacity-70"
                  >
                    {loading
                      ? "Working..."
                      : mode === "login"
                        ? "Login"
                        : "Create account"}
                  </button>
                </form>

                {status && (
                  <div className="mt-4 rounded-xl border border-slate-800 bg-slate-950/40 px-4 py-3 text-xs text-slate-200">
                    {status}
                  </div>
                )}
              </div>
            ) : (
              <div className="rounded-3xl border border-emerald-400/30 bg-emerald-400/10 p-6 shadow-[0_30px_80px_-50px_rgba(15,23,42,0.9)] backdrop-blur">
                <div className="flex items-start justify-between">
                  <div>
                    <p className="text-xs uppercase tracking-[0.2em] text-emerald-200">
                      Dashboard
                    </p>
                    <h2 className="font-display text-2xl text-white">
                      Welcome back
                    </h2>
                  </div>
                  <button
                    type="button"
                    onClick={clearToken}
                    className="text-xs text-emerald-200 hover:text-emerald-100"
                  >
                    Logout
                  </button>
                </div>

                <div className="mt-6 space-y-4 rounded-2xl border border-emerald-400/30 bg-slate-950/40 p-4 text-xs text-emerald-100">
                  <div>
                    <span className="text-[10px] uppercase tracking-[0.2em] text-emerald-200">
                      Token
                    </span>
                    <p className="mt-2 break-all text-emerald-100">{token}</p>
                  </div>
                  {profile ? (
                    <div>
                      <span className="text-[10px] uppercase tracking-[0.2em] text-emerald-200">
                        Profile
                      </span>
                      <p className="mt-2">ID: {profile.id}</p>
                      <p>Email: {profile.email}</p>
                    </div>
                  ) : (
                    <p className="text-emerald-100/80">
                      Fetch your profile to see the user id.
                    </p>
                  )}
                </div>

                <div className="mt-4 flex flex-col gap-3">
                  <button
                    type="button"
                    onClick={handleMe}
                    className="w-full rounded-xl border border-emerald-300 px-4 py-3 text-sm font-semibold text-emerald-50 transition hover:border-emerald-200 hover:text-white"
                  >
                    Fetch /me
                  </button>

                  {status && (
                    <div className="rounded-xl border border-emerald-300/40 bg-slate-950/40 px-4 py-3 text-xs text-emerald-100">
                      {status}
                    </div>
                  )}
                </div>
              </div>
            )}
          </section>
        </div>
      </div>
    </div>
  );
}
