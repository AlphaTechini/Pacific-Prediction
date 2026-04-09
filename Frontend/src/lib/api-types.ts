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

export interface MarketCreateContextSymbolResponse {
	symbol: string;
	tick_size: string;
	min_tick: string;
	max_tick: string;
	lot_size: string;
	min_order_size: string;
	max_order_size: string;
	max_leverage: number;
	isolated_only: boolean;
	mark_price?: string;
	oracle_price?: string;
	funding_rate?: string;
	next_funding_rate?: string;
	open_interest?: string;
	volume_24h?: string;
	updated_at?: string;
}

export interface MarketValidationModelResponse {
	market_type: string;
	source_type: string;
	allowed_operators: string[];
	requires_threshold: boolean;
	requires_interval: boolean;
	allowed_intervals?: string[];
}

export interface MarketCreateContextResponse {
	symbols: MarketCreateContextSymbolResponse[];
	validation_models: MarketValidationModelResponse[];
}

export interface CreateMarketRequest {
	title: string;
	symbol: string;
	market_type: string;
	condition_operator: string;
	creator_side: string;
	creator_stake_amount: string;
	threshold_value: string;
	source_type: string;
	source_interval: string;
	reference_value: string;
	expiry_time?: string;
}

export interface CreateMarketResponse extends MarketResponse {}

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

export interface CreatePositionRequest {
	side: string;
	stake_amount: string;
}

export interface CreatePositionResponse extends PositionResponse {}

export interface LeaderboardOverviewResponse {
	total_predictions: number;
	resolved_predictions: number;
	active_predictors: number;
	active_creators: number;
	average_win_rate: string;
}

export interface PredictorLeaderboardEntryResponse {
	rank: number;
	player_id: string;
	display_name: string;
	resolved_positions: number;
	won_positions: number;
	lost_positions: number;
	win_rate: string;
	net_profit: string;
	total_staked: string;
}

export interface CreatorLeaderboardEntryResponse {
	rank: number;
	player_id: string;
	display_name: string;
	created_markets: number;
	resolved_markets: number;
	total_positions: number;
	unique_participants: number;
	total_staked_on_markets: string;
}

export interface StreakLeaderboardEntryResponse {
	rank: number;
	player_id: string;
	display_name: string;
	current_win_streak: number;
	longest_win_streak: number;
	resolved_positions: number;
	win_rate: string;
	net_profit: string;
}

export interface ActivityLeaderboardEntryResponse {
	rank: number;
	player_id: string;
	display_name: string;
	total_positions: number;
	open_positions: number;
	resolved_positions: number;
	created_markets: number;
	total_staked: string;
}

export interface LeaderboardResponse {
	generated_at: string;
	overview: LeaderboardOverviewResponse;
	top_predictors: PredictorLeaderboardEntryResponse[];
	top_creators: CreatorLeaderboardEntryResponse[];
	best_streaks: StreakLeaderboardEntryResponse[];
	most_active: ActivityLeaderboardEntryResponse[];
}
