import { token, user, logout } from './stores.js';
import { get } from 'svelte/store';

export async function api(endpoint, options = {}) {
  const currentToken = get(token);
  const headers = {
    'Content-Type': 'application/json',
    ...(currentToken && { Authorization: `Bearer ${currentToken}` }),
  };

  const response = await fetch(`/api${endpoint}`, {
    ...options,
    headers: { ...headers, ...options.headers },
    body: options.body ? JSON.stringify(options.body) : undefined,
  });

  if (response.status === 204) {
    return null;
  }

  const data = await response.json();

  if (response.status === 401) {
    // Only logout if we had a token (session expired)
    if (currentToken) {
      logout();
    }
    throw new Error(data.error || 'Unauthorized');
  }

  if (!response.ok) {
    throw new Error(data.error || 'Something went wrong');
  }

  return data;
}

export async function login(email, password) {
  const data = await api('/login', {
    method: 'POST',
    body: { email, password },
  });
  token.set(data.token);
  user.set(data.user);
  localStorage.setItem('token', data.token);
  localStorage.setItem('user', JSON.stringify(data.user));
  return data;
}

export async function register(email, password) {
  const data = await api('/register', {
    method: 'POST',
    body: { email, password },
  });
  token.set(data.token);
  user.set(data.user);
  localStorage.setItem('token', data.token);
  localStorage.setItem('user', JSON.stringify(data.user));
  return data;
}

// Decks
export const getDecks = () => api('/decks');
export const getDeck = (id) => api(`/decks/${id}`);
export const getPublicDecks = () => api('/decks/public');
export const getPublicDecksBrowse = async () => {
  const response = await fetch('/api/public-decks');
  if (!response.ok) throw new Error('Failed to fetch public decks');
  return response.json();
};
export const getPublicDeck = async (id) => {
  const response = await fetch(`/api/public-decks/${id}`);
  if (!response.ok) throw new Error('Deck not found');
  return response.json();
};
export const getPublicDeckCards = async (id) => {
  const response = await fetch(`/api/public-decks/${id}/cards`);
  if (!response.ok) throw new Error('Failed to fetch cards');
  return response.json();
};
export const createDeck = (name, description, isPublic = false) =>
  api('/decks', { method: 'POST', body: { name, description, public: isPublic } });
export const updateDeck = (id, name, description, isPublic = false) =>
  api(`/decks/${id}`, { method: 'PUT', body: { name, description, public: isPublic } });
export const deleteDeck = (id) => api(`/decks/${id}`, { method: 'DELETE' });

// Cards
export const getCards = (deckId) => api(`/decks/${deckId}/cards`);
export const getCard = (id) => api(`/cards/${id}`);
export const createCard = (deckId, front, back) =>
  api(`/decks/${deckId}/cards`, { method: 'POST', body: { front, back } });
export const updateCard = (id, front, back) =>
  api(`/cards/${id}`, { method: 'PUT', body: { front, back } });
export const deleteCard = (id) => api(`/cards/${id}`, { method: 'DELETE' });
export const importCards = (deckId, cards) =>
  api(`/decks/${deckId}/cards/import`, { method: 'POST', body: { cards } });

