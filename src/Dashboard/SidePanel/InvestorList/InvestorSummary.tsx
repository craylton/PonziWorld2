import type { Investor } from './Investor';
import './InvestorList.css';

interface InvestorSummaryProps {
    investor: Investor;
}

export default function InvestorSummary({ investor }: InvestorSummaryProps) {
    const formatCurrency = (amount: number) => {
        return new Intl.NumberFormat('en-GB', {
            style: 'currency',
            currency: 'GBP',
            minimumFractionDigits: 0,
            maximumFractionDigits: 0,
        }).format(amount);
    };

    return (
        <div className="investor-list__item">
            <div className="investor-list__name">{investor.name}</div>
            <div className="investor-list__amount">
                {formatCurrency(investor.amountInvested)}
            </div>
        </div>
    );
}
