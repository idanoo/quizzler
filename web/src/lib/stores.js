import { writable } from 'svelte/store';

// Initialize from localStorage
const storedToken = localStorage.getItem('token');
const storedUser = JSON.parse(localStorage.getItem('user') || 'null');

export const token = writable(storedToken);
export const user = writable(storedUser);
export const isAuthenticated = writable(!!storedToken);

// Update isAuthenticated when token changes
token.subscribe((value) => {
  isAuthenticated.set(!!value);
});

export function logout() {
  token.set(null);
  user.set(null);
  localStorage.removeItem('token');
  localStorage.removeItem('user');
}

