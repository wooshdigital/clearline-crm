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