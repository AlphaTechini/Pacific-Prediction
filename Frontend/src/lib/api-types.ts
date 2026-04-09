export interface MarketResponse {
	id: string;
	title: string;
	symbol: string;
	market_type: string;
	condition_operator: string;
	threshold_value?: string;
	source_type: string;
	source_interval?: string;
	reference_value?: string;
	expiry_time: string;
	status: string;
	result?: string;
	settlement_value?: string;
	resolved_at?: string;
	resolution_reason?: string;
	created_by_player_id: string;
	created_at: string;
}

export interface ListMarketsResponse {
	active: MarketResponse[];
	resolved: MarketResponse[];
}

export interface BalanceResponse {
	player_id: string;
	available_balance: string;
	locked_balance: string;
	updated_at: string;
}

export interface PositionResponse {
	id: string;
	player_id: string;
	market_id: string;
	side: string;
	stake_amount: string;
	potential_payout: string;
	status: string;
	created_at: string;
	settled_at?: string;
}

export interface ListPositionsResponse {
	positions: PositionResponse[];
}
