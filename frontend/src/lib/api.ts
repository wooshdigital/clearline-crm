/**
 * API client for the ClearLine CRM backend.
 * Reads the API base URL from environment variables.
 */

const API_BASE = import.meta.env.PUBLIC_API_URL || 'http://localhost:8080';

export interface Contact {
  id: string;
  owner_id: string | null;
  first_name: string;
  last_name: string;
  email: string | null;
  phone: string | null;
  company: string | null;
  title: string | null;
  status: 'lead' | 'prospect' | 'customer' | 'churned';
  notes: string | null;
  created_at: string;
  updated_at: string;
}

export interface Deal {
  id: string;
  contact_id: string | null;
  owner_id: string | null;
  title: string;
  value: number;
  stage: 'prospecting' | 'qualification' | 'proposal' | 'negotiation' | 'closed_won' | 'closed_lost';
  close_date: string | null;
  notes: string | null;
  contact_name: string | null;
  created_at: string;
  updated_at: string;
}

export interface Activity {
  id: string;
  user_id: string | null;
  contact_id: string | null;
  deal_id: string | null;
  type: 'call' | 'email' | 'meeting' | 'note' | 'task';
  subject: string;
  body: string | null;
  occurred_at: string;
  created_at: string;
  user_name: string | null;
}

export interface User {
  id: string;
  name: string;
  email: string;
  role: 'admin' | 'manager' | 'rep';
  created_at: string;
}

export interface LoginResponse {
  token: string;
  user: User;
}

async function apiFetch<T>(path: string, token?: string, options?: RequestInit): Promise<T> {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(token ? { Authorization: `Bearer ${token}` } : {}),
  };

  const res = await fetch(`${API_BASE}${path}`, {
    ...options,
    headers: { ...headers, ...(options?.headers as Record<string, string> || {}) },
  });

  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: 'Unknown error' }));
    throw new Error(err.error || `HTTP ${res.status}`);
  }

  if (res.status === 204) return undefined as T;
  return res.json();
}

export const api = {
  login(email: string, password: string): Promise<LoginResponse> {
    return apiFetch('/api/auth/login', undefined, {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    });
  },

  getContacts(token: string, params?: { status?: string; q?: string }): Promise<Contact[]> {
    const qs = new URLSearchParams();
    if (params?.status) qs.set('status', params.status);
    if (params?.q) qs.set('q', params.q);
    const query = qs.toString() ? `?${qs}` : '';
    return apiFetch(`/api/contacts${query}`, token);
  },

  getContact(token: string, id: string): Promise<Contact> {
    return apiFetch(`/api/contacts/${id}`, token);
  },

  getDeals(token: string, stage?: string): Promise<Deal[]> {
    const query = stage ? `?stage=${stage}` : '';
    return apiFetch(`/api/deals${query}`, token);
  },

  getActivities(token: string, params?: { contact_id?: string; deal_id?: string; type?: string }): Promise<Activity[]> {
    const qs = new URLSearchParams();
    if (params?.contact_id) qs.set('contact_id', params.contact_id);
    if (params?.deal_id) qs.set('deal_id', params.deal_id);
    if (params?.type) qs.set('type', params.type);
    const query = qs.toString() ? `?${qs}` : '';
    return apiFetch(`/api/activities${query}`, token);
  },

  getUsers(token: string): Promise<User[]> {
    return apiFetch('/api/users', token);
  },

  health(): Promise<{ status: string }> {
    return apiFetch('/api/health');
  },
};

