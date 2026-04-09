import type { PageLoad } from './$types';

import { loadLeaderboardPageData } from '$lib/leaderboard-data';

export const load: PageLoad = async ({ fetch }) => {
	return loadLeaderboardPageData(fetch);
};
