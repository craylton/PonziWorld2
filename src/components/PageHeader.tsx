import './PageHeader.css';
import { useCurrentDay } from '../contexts/CurrentDayContext';

interface PageHeaderProps {
    title: string;
}

export default function PageHeader({ title }: PageHeaderProps) {
    const { currentDay } = useCurrentDay();

    return (
        <header className="page-header">
            {currentDay !== null && (
                <div className="page-header__day">Day {currentDay}</div>
            )}
            <h1 className="page-header__title">{title}</h1>
        </header>
    );
}
