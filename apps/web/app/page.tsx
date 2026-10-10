"use client";

import { FormEvent, useEffect, useState } from "react";

type Workflow = {
  id: string;
  brief: string;
  status: string;
  created_at: string;
};

const apiURL = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

export default function Home() {
  const [brief, setBrief] = useState("");
  const [workflows, setWorkflows] = useState<Workflow[]>([]);
  const [error, setError] = useState("");
  const [submitting, setSubmitting] = useState(false);

  async function loadWorkflows() {
    const response = await fetch(`${apiURL}/api/v1/workflows`);
    if (!response.ok) {
      throw new Error("Could not load workflows.");
    }
    setWorkflows(await response.json());
  }

  useEffect(() => {
    loadWorkflows().catch((loadError: Error) => setError(loadError.message));
  }, []);

  async function submitWorkflow(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setError("");
    setSubmitting(true);
    try {
      const response = await fetch(`${apiURL}/api/v1/workflows`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ brief }),
      });
      if (!response.ok) {
        throw new Error("Could not submit workflow.");
      }
      setBrief("");
      await loadWorkflows();
    } catch (submitError) {
      setError((submitError as Error).message);
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <main className="shell">
      <section className="hero">
        <p className="eyebrow">Research Workflow Platform</p>
        <h1>Turn a research question into a trackable workflow.</h1>
        <p className="intro">
          Submit a brief, send it through the Go API and RabbitMQ, and monitor its
          progress from one simple dashboard.
        </p>
      </section>

      <section className="panel">
        <h2>Start a workflow</h2>
        <form onSubmit={submitWorkflow}>
          <label htmlFor="brief">Research brief</label>
          <textarea
            id="brief"
            value={brief}
            onChange={(event) => setBrief(event.target.value)}
            placeholder="Example: Compare recent approaches to reproducible scientific software."
            required
            rows={5}
          />
          <button type="submit" disabled={submitting}>
            {submitting ? "Submitting..." : "Submit brief"}
          </button>
        </form>
        {error && <p className="error">{error}</p>}
      </section>

      <section className="panel">
        <div className="section-heading">
          <h2>Recent workflows</h2>
          <span>{workflows.length} total</span>
        </div>
        {workflows.length === 0 ? (
          <p className="empty">No workflows yet. Submit the first research brief.</p>
        ) : (
          <ul className="workflow-list">
            {workflows.map((workflow) => (
              <li key={workflow.id}>
                <div>
                  <strong>{workflow.brief}</strong>
                  <small>{new Date(workflow.created_at).toLocaleString()}</small>
                </div>
                <span className="status">{workflow.status}</span>
              </li>
            ))}
          </ul>
        )}
      </section>
    </main>
  );
}
