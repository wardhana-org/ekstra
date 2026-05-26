"use client";

import { type SubmitEvent, useState } from "react";

type User = {
  id: number;
  email: string;
  username: string;
  status: string;
};

type UserResponse = {
  user: User;
};

type RequestState = "idle" | "loading" | "success" | "error";

async function readError(response: Response) {
  const data = (await response.json().catch(() => null)) as
    | { error?: string }
    | null;

  return data?.error ?? `Request failed with status ${response.status}`;
}

export function LoginTester() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [user, setUser] = useState<User | null>(null);
  const [state, setState] = useState<RequestState>("idle");
  const [message, setMessage] = useState("");

  async function checkSession() {
    setState("loading");
    setMessage("Checking current session...");

    const response = await fetch("/api/auth/me", {
      credentials: "include",
    });

    if (!response.ok) {
      setUser(null);
      setState("idle");
      setMessage("No active session.");
      return;
    }

    const data = (await response.json()) as UserResponse;
    setUser(data.user);
    setState("success");
    setMessage("Session is active.");
  }

  async function handleLogin(event: SubmitEvent<HTMLFormElement>) {
    event.preventDefault();
    setState("loading");
    setMessage("Logging in...");

    const response = await fetch("/api/auth/login", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      credentials: "include",
      body: JSON.stringify({
        email,
        password,
      }),
    });

    if (!response.ok) {
      setUser(null);
      setState("error");
      setMessage(await readError(response));
      return;
    }

    const data = (await response.json()) as UserResponse;
    setUser(data.user);
    setState("success");
    setMessage("Login successful.");
  }

  async function handleLogout() {
    setState("loading");
    setMessage("Logging out...");

    const response = await fetch("/api/auth/logout", {
      method: "POST",
      credentials: "include",
    });

    if (!response.ok) {
      setState("error");
      setMessage(await readError(response));
      return;
    }

    setUser(null);
    setPassword("");
    setState("idle");
    setMessage("Logged out.");
  }

  return (
    <section>
      <h2>Login</h2>

      <form onSubmit={handleLogin}>
        <p>
          <label htmlFor="email">Email</label>
          <br />
          <input
            id="email"
            name="email"
            type="email"
            autoComplete="email"
            value={email}
            onChange={(event) => setEmail(event.target.value)}
            required
          />
        </p>

        <p>
          <label htmlFor="password">Password</label>
          <br />
          <input
            id="password"
            name="password"
            type="password"
            autoComplete="current-password"
            value={password}
            onChange={(event) => setPassword(event.target.value)}
            required
          />
        </p>

        <button type="submit" disabled={state === "loading"}>
          Login
        </button>
        <button
          type="button"
          disabled={state === "loading"}
          onClick={() => void checkSession()}
        >
          Check Session
        </button>
        <button
          type="button"
          disabled={state === "loading" || !user}
          onClick={() => void handleLogout()}
        >
          Logout
        </button>
      </form>

      <p>Status: {message || "Idle."}</p>

      {user ? (
        <pre>{JSON.stringify(user, null, 2)}</pre>
      ) : (
        <p>No user loaded.</p>
      )}
    </section>
  );
}
