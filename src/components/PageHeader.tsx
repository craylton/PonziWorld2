import './PageHeader.css';
import { useCurrentDayContext } from '../contexts/useCurrentDayContext';

interface PageHeaderProps {
    title: string;
}

export default function PageHeader({ title }: PageHeaderProps) {
    const { currentDay } = useCurrentDayContext();

    return (
        <header className="page-header">
            {currentDay !== null && (
                <div className="page-header__day">Day {currentDay}</div>
            )}
            <h1 className="page-header__title">{title}</h1>
        </header>
    );
}
