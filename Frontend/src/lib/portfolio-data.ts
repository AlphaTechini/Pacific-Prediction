import type { MarketResponse, PositionResponse } from '$lib/api-types';
import { fetchBackend } from '$lib/backend-api';

interface PortfolioData {
	activeMarkets: MarketResponse[];
	resolvedMarkets: MarketResponse[];
	positions: PositionResponse[];
	availableBalance: string;
	lockedBalance: string;
}

interface ListMarketsResponse {
	active: MarketResponse[];
	resolved: MarketResponse[];
}

interface BalanceResponse {
	available_balance: string;
	locked_balance: string;
}

interface ListPositionsResponse {
	positions: PositionResponse[];
}

export async function loadPortfolioData(): Promise<PortfolioData> {
	const [marketsResponse, balanceResponse, positionsResponse] = await Promise.all([
		fetchBackend('/api/v1/markets'),
		fetchBackend('/api/v1/players/me/balance'),
		fetchBackend('/api/v1/players/me/positions')
	]);

	if (!marketsResponse.ok) {
		throw new Error('Unable to load market data right now.');
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

	const markets = (await marketsResponse.json()) as ListMarketsResponse;
	const balance = (await balanceResponse.json()) as BalanceResponse;
	const positions = (await positionsResponse.json()) as ListPositionsResponse;

	return {
		activeMarkets: markets.active,
		resolvedMarkets: markets.resolved,
		positions: positions.positions,
		availableBalance: balance.available_balance,
		lockedBalance: balance.locked_balance
	};
}
