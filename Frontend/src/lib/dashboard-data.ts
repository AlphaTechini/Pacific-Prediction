import type {
	BalanceResponse,
	ListMarketsResponse,
	ListPositionsResponse,
	MarketResponse,
	PositionResponse
} from '$lib/api-types';
import { fetchBackend } from '$lib/backend-api';

export interface DashboardData {
	activeMarkets: MarketResponse[];
	resolvedMarkets: MarketResponse[];
	balance: BalanceResponse;
	positions: PositionResponse[];
}

export async function loadDashboardData(): Promise<DashboardData> {
	const [marketsResponse, balanceResponse, positionsResponse] = await Promise.all([
		fetchBackend('/api/v1/markets'),
		fetchBackend('/api/v1/players/me/balance'),
		fetchBackend('/api/v1/players/me/positions')
	]);

	if (!marketsResponse.ok) {
		throw new Error('Unable to load markets right now.');
	}

	if (balanceResponse.status === 401 || positionsResponse.status === 401) {
		throw new Error('Your guest session is not ready yet.');
	}

	if (!balanceResponse.ok) {
		throw new Error('Unable to load your balance right now.');
	}

	if (!positionsResponse.ok) {
		throw new Error('Unable to load your positions right now.');
	}

	const marketsPayload = (await marketsResponse.json()) as ListMarketsResponse;
	const balancePayload = (await balanceResponse.json()) as BalanceResponse;
	const positionsPayload = (await positionsResponse.json()) as ListPositionsResponse;

	return {
		activeMarkets: marketsPayload.active,
		resolvedMarkets: marketsPayload.resolved,
		balance: balancePayload,
		positions: positionsPayload.positions
	};
}
