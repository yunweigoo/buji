import { createClient, type SupabaseClient } from '@supabase/supabase-js'

let client: SupabaseClient | undefined

/**
 * Browser-safe Supabase client for Auth, Realtime, or RLS-protected reads.
 * Never place the service-role key or database password in VITE_* variables.
 */
export function getSupabaseClient(): SupabaseClient | undefined {
  const url = import.meta.env.VITE_SUPABASE_URL
  const anonKey = import.meta.env.VITE_SUPABASE_ANON_KEY

  if (!url || !anonKey) {
    return undefined
  }
  client ??= createClient(url, anonKey)
  return client
}
