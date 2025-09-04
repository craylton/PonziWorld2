import type { Investor } from './Investor';
import './InvestorList.css';
import { formatCurrencyFromString } from '../../../utils/currency';

interface InvestorSummaryProps {
    investor: Investor;
}

export default function InvestorSummary({ investor }: InvestorSummaryProps) {
    return (
        <div className="investor-list__item">
            <div className="investor-list__name">{investor.name}</div>
            <div className="investor-list__amount">
                {formatCurrencyFromString(investor.amountInvested)}
            </div>
        </div>
    );
}
