import type { LeaderboardResponse } from '$lib/api-types';

export interface LeaderboardPageData {
	leaderboard: LeaderboardResponse | null;
	error: string | null;
}

export async function loadLeaderboardPageData(
	fetchFn: typeof globalThis.fetch
): Promise<LeaderboardPageData> {
	try {
		const response = await fetchFn('/api/v1/leaderboard');
		if (!response.ok) {
			return {
				leaderboard: null,
				error: 'Unable to load the leaderboard right now.'
			};
		}

		return {
			leaderboard: (await response.json()) as LeaderboardResponse,
			error: null
		};
	} catch {
		return {
			leaderboard: null,
			error: 'Unable to reach the leaderboard right now.'
		};
	}
}
