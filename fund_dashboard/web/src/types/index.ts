export interface Fund {
  code: string;
  name: string;
  type: string;
  created_at?: string;
  updated_at?: string;
}

export interface Transaction {
  id: number;
  fund_code: string;
  type: 'buy' | 'sell';
  date: string;
  shares: number;
  price: number;
  amount: number;
  fee: number;
  created_at?: string;
  updated_at?: string;
}

export interface FundPosition {
  fund_code: string;
  fund_name: string;
  total_shares: number;
  total_cost: number;
  average_cost: number;
  current_nav: number;
  current_value: number;
  daily_pnl: number;
  daily_pnl_rate: number;
  total_pnl: number;
  total_pnl_rate: number;
  nav_date: string;
  est_nav_date: string;
}

export interface PortfolioSummary {
  total_cost: number;
  total_value: number;
  daily_pnl: number;
  daily_pnl_rate: number;
  total_pnl: number;
  total_pnl_rate: number;
  positions: FundPosition[];
}
