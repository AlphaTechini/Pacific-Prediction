import { browser } from '$app/environment';
import { get, writable } from 'svelte/store';

import { fetchBackend } from '$lib/backend-api';

type GuestSessionStatus = 'idle' | 'loading' | 'ready' | 'error';

export interface SessionPlayer {
	id: string;
	displayName: string;
	expiresAt?: string;
}

interface GuestSessionState {
	status: GuestSessionStatus;
	player: SessionPlayer | null;
	error: string | null;
}

interface GuestSessionResponse {
	player_id: string;
	display_name: string;
	expires_at: string;
}

interface MeResponse {
	id: string;
	display_name: string;
}

const guestPlayerStorageKey = 'pacific_prediction_guest_player';
const storedInitialPlayer = loadStoredGuestPlayer();

const initialState: GuestSessionState = {
	status: storedInitialPlayer ? 'ready' : 'idle',
	player: storedInitialPlayer,
	error: null
};

export const guestSession = writable<GuestSessionState>(initialState);

let ensurePromise: Promise<SessionPlayer | null> | null = null;

export async function ensureGuestSession(): Promise<SessionPlayer | null> {
	const currentState = get(guestSession);

	if (currentState.status === 'ready' && currentState.player) {
		return currentState.player;
	}

	if (ensurePromise) {
		return ensurePromise;
	}

	guestSession.set({
		status: 'loading',
		player: currentState.player,
		error: null
	});

	ensurePromise = resolveGuestSession()
		.then((player) => {
			guestSession.set({
				status: 'ready',
				player,
				error: null
			});

			return player;
		})
		.catch((error: unknown) => {
			guestSession.set({
				status: 'error',
				player: null,
				error: toErrorMessage(error)
			});

			return null;
		})
		.finally(() => {
			ensurePromise = null;
		});

	return ensurePromise;
}

async function resolveGuestSession(): Promise<SessionPlayer> {
	const storedPlayer = loadStoredGuestPlayer();
	if (storedPlayer) {
		return storedPlayer;
	}

	const currentPlayer = await fetchCurrentPlayer();

	if (currentPlayer) {
		storeGuestPlayer(currentPlayer);
		return currentPlayer;
	}

	return createGuestSession();
}

async function fetchCurrentPlayer(): Promise<SessionPlayer | null> {
	const response = await fetchBackend('/api/v1/players/me');

	if (response.status === 401) {
		return null;
	}

	if (!response.ok) {
		throw new Error('Unable to check the current guest session.');
	}

	const payload = (await response.json()) as MeResponse;

	return {
		id: payload.id,
		displayName: payload.display_name
	};
}

async function createGuestSession(): Promise<SessionPlayer> {
	const response = await fetchBackend('/api/v1/players/guest', {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json'
		},
		body: JSON.stringify({})
	});

	if (!response.ok) {
		throw new Error('Unable to start a guest session.');
	}

	const payload = (await response.json()) as GuestSessionResponse;

	const player = {
		id: payload.player_id,
		displayName: payload.display_name,
		expiresAt: payload.expires_at
	};

	storeGuestPlayer(player);

	return player;
}

function loadStoredGuestPlayer(): SessionPlayer | null {
	if (!browser) {
		return null;
	}

	const rawPlayer = localStorage.getItem(guestPlayerStorageKey);
	if (!rawPlayer) {
		return null;
	}

	try {
		const parsed = JSON.parse(rawPlayer) as Partial<SessionPlayer>;
		if (typeof parsed.id !== 'string' || typeof parsed.displayName !== 'string') {
			localStorage.removeItem(guestPlayerStorageKey);
			return null;
		}

		return {
			id: parsed.id,
			displayName: parsed.displayName,
			expiresAt: typeof parsed.expiresAt === 'string' ? parsed.expiresAt : undefined
		};
	} catch {
		localStorage.removeItem(guestPlayerStorageKey);
		return null;
	}
}

function storeGuestPlayer(player: SessionPlayer): void {
	if (!browser) {
		return;
	}

	try {
		localStorage.setItem(
			guestPlayerStorageKey,
			JSON.stringify({
				id: player.id,
				displayName: player.displayName,
				expiresAt: player.expiresAt
			})
		);
	} catch {
		return;
	}
}

function toErrorMessage(error: unknown): string {
	if (error instanceof Error && error.message) {
		return error.message;
	}

	return 'Unable to start a guest session.';
}
