export function useAnalytics() {
  const formatMoney = (value: number) => {
    return new Intl.NumberFormat('zh-CN', {
      style: 'currency',
      currency: 'CNY',
      minimumFractionDigits: 2,
    }).format(value);
  };

  const formatRate = (value: number) => {
    return (value * 100).toFixed(2) + '%';
  };

  const getPnlClass = (value: number) => {
    if (value > 0) return 'text-red'; // In China, red is usually positive
    if (value < 0) return 'text-green'; // and green is negative
    return '';
  };

  const getMatrixPnlClass = (value: number) => {
    // For Matrix terminal theme, let's keep Matrix green as positive, and red as negative
    // The instructions say "loss red #ff4444", so loss is red, profit is green
    if (value > 0) return 'text-green';
    if (value < 0) return 'text-red';
    return '';
  };

  return { formatMoney, formatRate, getPnlClass, getMatrixPnlClass };
}
